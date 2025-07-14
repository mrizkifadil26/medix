package webgen

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/logger"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

var htmlMinifier = func() *minify.M {
	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})

	return minifier
}()

var DryRun = false

func RenderTemplate(files []string, outPath string, data any) {
	// Resolve all template file paths
	tmpl, err := template.ParseFiles(files...)
	Must(err)

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "base", data)
	Must(err)

	minified, err := htmlMinifier.Bytes("text/html", buf.Bytes())
	Must(err)

	if DryRun {
		logger.Dry("[DRY-RUN] Skipped writing file: " + outPath)
		return
	}

	// Create target directory if missing
	err = os.MkdirAll(filepath.Dir(outPath), 0755)
	Must(err)

	f, err := os.Create(outPath)
	Must(err)
	defer f.Close()

	_, err = f.Write(minified)
	Must(err)

	logger.Done("📄 Rendered → " + outPath)
}

func RenderStaticPages() {
	pages := []struct {
		Files   []string
		OutPath string
		Data    map[string]any
	}{
		{
			Files:   []string{"templates/layouts/base.html", "templates/pages/index.html"},
			OutPath: "dist/index.html",
			Data:    nil,
		},
		{
			Files:   []string{"templates/layouts/base.html", "templates/pages/about.html"},
			OutPath: "dist/about.html",
			Data:    nil,
		},
	}

	for _, page := range pages {
		RenderTemplate(page.Files, page.OutPath, page.Data)
	}
}

func RenderDataPage(jsonFile, title, outFile string) {
	raw, err := os.ReadFile(filepath.Join("data", jsonFile))
	Must(err)

	var parsed any
	Must(json.Unmarshal(raw, &parsed))

	data := map[string]any{
		"Title": title,
		"Type":  strings.TrimSuffix(jsonFile, ".json"),
		"Data":  template.JS(string(raw)), // inline JS object in <script>
	}

	RenderTemplate(
		[]string{"templates/layouts/base.html", "templates/pages/list.html"},
		filepath.Join("dist", outFile),
		data,
	)
}
