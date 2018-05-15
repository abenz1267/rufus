package rufus

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unrolled/secure"
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
func (r *Router) RegisterRoutes(languages map[string]int, server server, csp string) {
	r.RoutesSender = make(chan string)
	r.RoutesReceiver = make(chan http.Handler)

	r.prependMiddleware(server, csp)

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
	r.RoutesSender <- language
	return <-r.RoutesReceiver
}

func (r *Router) prependMiddleware(server server, csp string) {
	host := server.ProductionHost

	r.Mux.Use(middleware.Compress(5, "application/octet-stream", "application/javascript", "application/json", "text/html", "text/css", "text/plain", "text/javascript", "image/svg+xml", "image/jpeg", "image/png", "image/x-icon"))

	if server.Dev {
		host = server.DevelopmentHost
	} else {
		r.Mux.Use(middleware.RequestID)
		r.Mux.Use(middleware.RealIP)
		r.Mux.Use(middleware.Recoverer)
	}

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{host, "www." + host},
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           true,
		SSLHost:               host,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ReferrerPolicy:        "no-referrer",
		ContentSecurityPolicy: csp,
	})

	r.Mux.Use(secureMiddleware.Handler)

	log.Println("Setting content type based on 'Accept' Header from request")
	r.Mux.Use(r.Middleware.setContentType())

	if r.Middleware.RedirectToNonWWW {
		log.Println("Redirecting to non-www")
		r.Mux.Use(r.Middleware.redirectWithoutWWW())
	}

	if r.Middleware.EnableResponseCache {
		log.Println("Using response cache")
		r.Mux.Use(r.Middleware.Cache.Check())
	}
}
