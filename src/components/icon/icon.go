package icon

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

var IconFontHref string

//go:embed *.svg
var svgFiles embed.FS

// SVG holds every .svg in this package keyed by filename without extension, so
// templates can inline a logo or mark without a second HTTP request.
var SVG = map[string]template.HTML{}

func init() {
	entries, err := fs.ReadDir(svgFiles, ".")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		content, err := svgFiles.ReadFile(entry.Name())
		if err != nil {
			panic(err)
		}
		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		SVG[name] = template.HTML(content)
	}
}

// Names lists every icon used anywhere in the app's templates. An icon must be
// listed here or it is not included in the subsetted font and will not render.
var Names = []string{
	"chevron_left",
	"chevron_right",
	"close",
	"first_page",
	"home",
	"keyboard_arrow_right",
	"last_page",
	"menu",
	"search",
	"widgets",
}

func init() {
	// The fonts API requires icon_names to be alphabetically sorted.
	slices.Sort(Names)

	IconFontHref = fmt.Sprintf("https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200&icon_names=%s&display=block", strings.Join(Names, ","))
}
