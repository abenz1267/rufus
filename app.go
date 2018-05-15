package rufus

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"golang.org/x/text/language"
)

// App is the main Rufus instance
type App struct {
	Server      server `json:"server,omitempty"`
	Router      `json:"router,omitempty"`
	Translation `json:"-"`
	Templates   `json:"templates,omitempty"`
	Response    `json:"-"`
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
	a.Router.RegisterRoutes(a.Translation.Languages, a.Server, a.CSPPolicy)

	if err := a.Templates.CacheFiles(a.Translation); err != nil {
		return err
	}

	go a.Server.startFileServer(a.Router.Mux, a.Server.StaticURL, http.Dir(a.Server.StaticFolder))

	if a.Server.Dev {
		return a.Server.startDevelopment(a.Router.Mux)
	}

	return a.Server.startProduction(a.Router.Mux)
}

// Render handles and generates the reponse
func (a *App) Render(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" {
		a.Response.renderJSON(w, r)
		return
	}

	if a.Response.Status != http.StatusOK {
		a.Response.TemplateFile = "error"
	}

	var altName strings.Builder
	altName.WriteString(a.Response.TemplateFile)
	altName.WriteString("_raw")

	a.Response.renderHTML(w, r, a.Templates.live[a.Response.TemplateFile], a.Templates.liveStripped[altName.String()])
}
