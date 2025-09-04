package normdb

import (
	"database/sql"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mrizkifadil26/medix/internal/db"
)

type Movie struct {
	File      db.File
	Subtitles []db.File
}

type MovieGroup struct {
	Group      string // top-level group folder
	Collection string // optional
	Folder     string // movie folder
	Movie      *Movie

	Icon       *db.File // .ico in same group folder
	DesktopIni *db.File // desktop.ini in same group folder
}

type Episode struct {
	File      db.File
	Subtitles []db.File
}

type Season struct {
	Season   string
	Episodes map[string]*Episode // key: episode filename or epNumber
}

type Show struct {
	Group      string
	Show       string
	Seasons    map[string]*Season
	Icon       *db.File // optional .ico file
	DesktopIni *db.File // optional desktop.ini file
}

type FileType string

const (
	FileTypeVideo     FileType = "video"
	FileTypeSubtitle  FileType = "subtitle"
	FileTypeThumbnail FileType = "thumbnail"
	FileTypeOther     FileType = "other"
)

// Matches "Title (2020)"
var moviePattern = regexp.MustCompile(`^(?P<title>.+) \((?P<year>\d{4})\)$`)

// Matches season folder like "S01"
var seasonPattern = regexp.MustCompile(`(?i)^S(\d{1,2})$`)

func handleMovieFile(movieGroups map[string]*MovieGroup, group string, parts []string, f db.File) {
	var collection, movieFolder string

	if moviePattern.MatchString(parts[1]) {
		movieFolder = parts[1]
	} else if len(parts) >= 3 && moviePattern.MatchString(parts[2]) {
		collection = parts[1]
		movieFolder = parts[2]
	} else {
		// might be desktop.ini / .ico in group folder
		if strings.EqualFold(f.Filename, "desktop.ini") {
			// assign to all movieGroups in this group
			for _, mg := range movieGroups {
				if mg.Group == group && mg.DesktopIni == nil {
					mg.DesktopIni = &f
				}
			}
		}
		if strings.EqualFold(filepath.Ext(f.Filename), ".ico") {
			for _, mg := range movieGroups {
				if mg.Group == group && mg.Icon == nil {
					mg.Icon = &f
				}
			}
		}
		return
	}

	key := group + "/" + collection + "/" + movieFolder
	mg, ok := movieGroups[key]
	if !ok {
		mg = &MovieGroup{
			Group:      group,
			Collection: collection,
			Folder:     movieFolder,
			Movie:      &Movie{},
		}
		movieGroups[key] = mg
	}

	switch f.Type {
	case "subtitle":
		mg.Movie.Subtitles = append(mg.Movie.Subtitles, f)
	default:
		mg.Movie.File = f
	}
}

func handleTVFile(shows map[string]*Show, group string, parts []string, f db.File) {
	if len(parts) < 2 {
		return
	}

	showName := parts[1]
	show, ok := shows[group+"/"+showName]
	if !ok {
		show = &Show{
			Group:   group,
			Show:    showName,
			Seasons: map[string]*Season{},
		}
		shows[group+"/"+showName] = show
	}

	// root-level files in show folder
	if len(parts) == 2 {
		name := strings.ToLower(f.Filename)
		ext := strings.ToLower(filepath.Ext(name))
		switch {
		case name == "desktop.ini":
			if show.DesktopIni == nil {
				show.DesktopIni = &f
			}
		case ext == ".ico":
			if show.Icon == nil {
				show.Icon = &f
			}
		}
		return
	}

	// must be under a season
	if len(parts) < 3 {
		return
	}

	seasonName := parts[2]
	season, ok := show.Seasons[seasonName]
	if !ok {
		season = &Season{
			Season:   seasonName,
			Episodes: map[string]*Episode{},
		}
		show.Seasons[seasonName] = season
	}

	// episode file
	if len(parts) >= 4 {
		epFileName := parts[len(parts)-1]
		ep, ok := season.Episodes[epFileName]
		if !ok {
			ep = &Episode{}
			season.Episodes[epFileName] = ep
		}

		switch f.Type {
		case "subtitle":
			ep.Subtitles = append(ep.Subtitles, f)
		default: // assume video
			ep.File = f
		}
	}
}

// GroupFiles scans all files and organizes them into movies and shows.
func GroupFiles(dbConn *sql.DB, mediaType string) (map[string]*MovieGroup, map[string]*Show, error) {
	files, err := db.ListFiles(dbConn)
	if err != nil {
		return nil, nil, err
	}

	movieGroups := map[string]*MovieGroup{}
	tvShows := map[string]*Show{}

	for _, f := range files {
		parts := strings.Split(filepath.ToSlash(f.RelativePath), "/")
		if len(parts) < 2 {
			continue
		}
		group := parts[0]

		switch mediaType {
		case "movie":
			handleMovieFile(movieGroups, group, parts, f)
		case "tv":
			handleTVFile(tvShows, group, parts, f)
		}
	}

	return movieGroups, tvShows, nil
}

func ToMovieWithExtras(group MovieGroup, mediaID int64) db.MovieWithExtras {
	abs := group.Movie.File.AbsolutePath // assuming db.File has Path (absolute)
	rel := group.Movie.File.RelativePath
	source := deriveSource(abs, rel)

	m := db.Movie{
		Title:      extractTitle(group.Folder), // you’ll define how to parse title
		Year:       extractYear(group.Folder),  // parse year from folder name
		MediaID:    mediaID,
		Source:     source,
		Group:      group.Group,
		Collection: group.Collection,
	}

	var thumb *db.MediaThumbnail
	if group.Icon != nil || group.DesktopIni != nil {
		thumb = &db.MediaThumbnail{}
		if group.Icon != nil {
			thumb.IconID = &group.Icon.ID
		}
		if group.DesktopIni != nil {
			thumb.DesktopIniID = &group.DesktopIni.ID
		}
	}

	var subs []db.MediaSubtitle
	for _, f := range group.Movie.Subtitles {
		subs = append(subs, db.MediaSubtitle{
			MediaType:  "movie",
			SubtitleID: f.ID,                   // adjust
			Language:   detectLang(f.Filename), // implement simple lang guesser
		})
	}

	return db.MovieWithExtras{
		Movie:     m,
		Thumbnail: thumb,
		Subtitles: subs,
	}
}

func NormalizeMovies(dbConn *sql.DB, movieGroups []MovieGroup, source string) error {
	var batch []db.MovieWithExtras

	for _, mg := range movieGroups {
		// `mg.Movie.File.ID` should exist after file scan
		if mg.Movie == nil || mg.Movie.File.ID == 0 {
			continue
		}

		movieWithExtras := ToMovieWithExtras(mg, mg.Movie.File.ID)
		batch = append(batch, movieWithExtras)
	}

	if len(batch) == 0 {
		return nil
	}

	return db.BatchUpsertMovies(dbConn, batch)
}

func deriveSource(absPath, relPath string) string {
	if strings.HasSuffix(absPath, relPath) {
		return strings.TrimSuffix(absPath, relPath)
	}

	return filepath.Dir(absPath[:len(absPath)-len(relPath)])
}

var yearRe = regexp.MustCompile(`(?i)[\(\[\.\s_-](19|20)\d{2}[\)\]\.\s_-]?`)

func extractTitle(folder string) string {
	// remove year if present
	title := yearRe.ReplaceAllString(folder, "")
	// replace dots/underscores with spaces
	title = strings.ReplaceAll(title, ".", " ")
	title = strings.ReplaceAll(title, "_", " ")
	// trim whitespace
	return strings.TrimSpace(title)
}

func extractYear(folder string) int {
	m := yearRe.FindString(folder)
	if m == "" {
		return 0 // no year detected
	}
	m = strings.Trim(m, " ()[]-_.") // clean up
	year, _ := strconv.Atoi(m)
	return year
}

func detectLang(name string) string {
	name = strings.ToLower(name)

	// strip Pahe.in or other tags
	name = strings.ReplaceAll(name, "pahe.in", "")

	if strings.Contains(name, ".en.") || strings.Contains(name, "_en") {
		return "en"
	}
	if strings.Contains(name, ".id.") || strings.Contains(name, "_id") {
		return "id"
	}
	return "unknown"
}
