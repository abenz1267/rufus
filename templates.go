package rufus

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Templates holds cached template files
type Templates struct {
	TemplateFolder string `json:"template_folder,omitempty"`
	BaseTemplate   string `json:"base_template,omitempty"`
	live           map[string]*template.Template
	liveStripped   map[string]*template.Template
	test           map[string]*template.Template
	testStripped   map[string]*template.Template
}

// CacheFiles processes templates and saves them to the according map
func (t *Templates) CacheFiles(translation Translation) error {
	funcs := template.FuncMap{}
	funcs["translate"] = translation.Translate
	funcs["translateURL"] = translation.TranslateURL
	funcs["safeHTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}

	liveMap := make(map[string]*template.Template)
	strippedLiveMap := make(map[string]*template.Template)

	liveTemplateFiles, err := ioutil.ReadDir(t.TemplateFolder)
	if err != nil {
		fmt.Println(t.TemplateFolder)
		return err
	}

	baseTemplate := filepath.Join(t.TemplateFolder, t.BaseTemplate)

	for _, file := range liveTemplateFiles {
		filename := file.Name()
		templateFile := filepath.Join(t.TemplateFolder, filename)

		if filename != t.BaseTemplate {

			filenameNoHTML := strings.TrimSuffix(filename, ".html")

			switch strings.Contains(filenameNoHTML, "_raw") {
			case true:
				newTemplate, err := template.New(filename).Funcs(funcs).ParseFiles(templateFile)
				if err != nil {
					return err
				}
				strippedLiveMap[filenameNoHTML] = newTemplate
			default:
				newTemplate, err := template.New(t.BaseTemplate).Funcs(funcs).ParseFiles(baseTemplate, templateFile)
				if err != nil {
					return err
				}

				liveMap[filenameNoHTML] = newTemplate
			}
		}
	}

	t.live = liveMap
	t.liveStripped = strippedLiveMap

	return nil
}
