package rufus

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"golang.org/x/text/language"
)

// Router contains the router and all needed methods for automated translation
type Router struct {
	Mux             *chi.Mux          `json:"-"`
	LanguageMatcher language.Matcher  `json:"-"`
	RoutesSender    chan string       `json:"-"`
	RoutesReceiver  chan http.Handler `json:"-"`
	Middleware      `json:"middleware,omitempty"`
}

// RegisterRoutes handles routes and different languages
func (r *Router) RegisterRoutes(languages map[string]int) {
	r.RoutesSender = make(chan string)
	r.prependMiddleware()

	if languages != nil {
		r.Mux.Get("/", r.getBrowserLanguagePreferenceAndRedirect)

		for k := range languages {
			r.Mux.Mount("/"+k, r.getRoutes(k))
		}
	} else {
		r.Mux.Mount("/", r.getRoutes("default"))
	}
}

func (r *Router) getBrowserLanguagePreferenceAndRedirect(w http.ResponseWriter, req *http.Request) {
	accept := req.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(r.LanguageMatcher, accept)
	lang, _ := tag.Base()

	var b strings.Builder
	b.WriteString("https://")
	b.WriteString(req.Host)
	b.WriteString("/")
	b.WriteString(lang.String())
	b.WriteString(req.RequestURI)

	http.Redirect(w, req, b.String(), http.StatusMovedPermanently)
}

func (r *Router) getRoutes(language string) http.Handler {
	r.RoutesReceiver = make(chan http.Handler)

	r.RoutesSender <- language
	return <-r.RoutesReceiver
}

func (r *Router) prependMiddleware() {
	if r.Middleware.RedirectToNonWWW {
		r.Mux.Use(r.Middleware.redirectWithoutWWW())
	}
}
