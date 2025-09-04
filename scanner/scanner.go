package scanner

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mrizkifadil26/medix/internal/db"
)

type FileType string

const (
	FileTypeVideo     FileType = "video"
	FileTypeSubtitle  FileType = "subtitle"
	FileTypeThumbnail FileType = "thumbnail"
	FileTypeOther     FileType = "other"
)

// ScanDirectory walks the root directory and inserts/updates files into DB.
func ScanDirectory(dbConn *sql.DB, root, parentKey string) error {
	var files []db.File

	spinner := NewSpinner("Processing")
	defer spinner.Finish()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't read a file, just skip it
			return nil
		}

		if d.IsDir() {
			return nil // skip directories
		}

		info, err := d.Info()
		if err != nil {
			return nil // skip broken files
		}

		rel, _ := filepath.Rel(root, path)

		ext := filepath.Ext(info.Name())
		ftype, normExt := detectFileType(ext)

		fp, err := makeFingerprintWithPartial(path, info.Size(), info.ModTime())
		if err != nil {
			// fallback: just metadata-based
			fp = makeFingerprint(info.Size(), info.ModTime())
		}

		files = append(files, db.File{
			AbsolutePath: path,
			RelativePath: rel,
			Filename:     info.Name(),
			Extension:    normExt,
			Type:         string(ftype), // you could improve this later (video, audio, etc.)
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			ParentType:   parentKey,
			ParentID:     nil,
			Fingerprint:  fp,
			CreatedAt:    time.Now(),
		})

		spinner.Increment(path)
		return nil
	})

	if err != nil {
		return err
	}

	// Batch upsert into DB
	return db.BatchUpsertFiles(dbConn, files)
}

func detectFileType(ext string) (FileType, string) {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	switch ext {
	case "mp4", "mkv", "avi", "mov", "wmv":
		return FileTypeVideo, ext
	case "srt", "sub", "vtt", "ass":
		return FileTypeSubtitle, ext
	case "jpg", "jpeg", "png", "webp", "ico", "ini":
		return FileTypeThumbnail, ext
	default:
		return FileTypeOther, ext
	}
}

func makeFingerprintWithPartial(path string, size int64, modTime time.Time) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 1024*1024) // 1MB
	n, _ := f.Read(buf)

	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%d", size)))
	h.Write([]byte(modTime.UTC().Format(time.RFC3339Nano)))
	h.Write(buf[:n])

	return hex.EncodeToString(h.Sum(nil)), nil
}

func makeFingerprint(size int64, modTime time.Time) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%d", size)))
	h.Write([]byte(modTime.UTC().Format(time.RFC3339Nano)))
	return hex.EncodeToString(h.Sum(nil))
}
