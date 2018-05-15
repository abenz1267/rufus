package rufus

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-chi/chi"
	"golang.org/x/text/language"
)

// App is the main Rufus instance
type App struct {
	Server      server `json:"server,omitempty"`
	Router      `json:"router,omitempty"`
	Translation `json:"-"`
	Templates   `json:"templates,omitempty"`
	Language    string `json:"language,omitempty"`
	CSPPolicy   string `json:"csp_policy,omitempty"`
}

// LoadConfigAndRouter loads data from 'config.json' and sets chi router
func (a *App) LoadConfigAndRouter() error {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, a); err != nil {
		return err
	}

	a.Router.Mux = chi.NewRouter()

	if err := a.Translation.loadData(); err != nil {
		return err
	}

	a.addLanguageMatcher()

	return nil
}

func (a *App) addLanguageMatcher() {
	var languageTags []language.Tag

	for k := range a.Translation.Languages {
		switch k {
		case "en":
			languageTags = append(languageTags, language.English)
		case "de":
			languageTags = append(languageTags, language.German)
		case "es":
			languageTags = append(languageTags, language.Spanish)
		case "fr":
			languageTags = append(languageTags, language.French)
		case "it":
			languageTags = append(languageTags, language.Italian)
		}
	}

	a.Router.LanguageMatcher = language.NewMatcher(languageTags)
}

// Start checks the environment set (development or production) and starts an according server
func (a *App) Start() error {
	a.Router.RegisterRoutes(a.Translation.Languages)

	if err := a.Templates.CacheFiles(a.Translation); err != nil {
		return err
	}

	if a.Server.Dev {
		return a.Server.startDevelopment(a.Router.Mux)
	}

	return a.Server.startProduction(a.Router.Mux)
}
