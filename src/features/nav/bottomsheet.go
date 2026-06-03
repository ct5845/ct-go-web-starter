package nav

import (
	"ct-go-web-starter/src/components/bottomsheet"
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed bottomsheet.html
	bottomsheetHTML string
	bottomsheetTpl  = component.New("bottomsheet.html", bottomsheetHTML)
)

func createBottomSheet(activeTab string) (template.HTML, error) {
	content, err := bottomsheetTpl.Render(activeTab)
	if err != nil {
		return "", err
	}

	return bottomsheet.Render(bottomsheet.Options{
		Id:      moreSheetID,
		Content: content,
	})
}
