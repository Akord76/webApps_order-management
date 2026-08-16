package routes

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
)

// loadTemplates walks root and parses every .html file, naming each
// template after its path relative to root (with forward slashes), e.g.
// "category_template/create_update.html". gin's default LoadHTMLGlob
// names templates by base filename only, which would collide here since
// every module folder has its own create_update.html.
func loadTemplates(root string) render.HTMLRender {
	tmpl := template.New("")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = tmpl.New(rel).Parse(string(content))
		return err
	})
	if err != nil {
		log.Fatalf("routes: failed to load templates from %q: %v", root, err)
	}

	return render.HTMLProduction{Template: tmpl}
}
