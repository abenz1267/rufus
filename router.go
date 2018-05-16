package rufus

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
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

	r.Mux.Use(middleware.Compress(5, "application/octet-stream", "application/javascript", "application/json", "text/html", "text/css", "text/plain", "text/javascript", "image/svg+xml", "image/jpeg", "image/png", "image/x-icon"))

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

	http.Redirect(w, req, b.String(), http.StatusMovedPermanently)
}

func (r *Router) getRoutes(language string) http.Handler {
	r.RoutesSender <- language
	return <-r.RoutesReceiver
}

// PrependMiddleware is used to add basic middleware to routers
func (r *Router) PrependMiddleware(router *chi.Mux, server server, csp string) {
	host := server.ProductionHost

	//zerolog config
	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.TimeFieldFormat = time.RFC822

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	router.Use(hlog.NewHandler(logger))
	router.Use(hlog.RemoteAddrHandler("ip"))
	router.Use(r.Middleware.logRequests())

	if server.Dev {
		host = server.DevelopmentHost
	} else {
		router.Use(middleware.RequestID)
		router.Use(middleware.RealIP)
		router.Use(middleware.Recoverer)
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

	router.Use(secureMiddleware.Handler)

	if r.Middleware.RedirectToNonWWW {
		router.Use(r.Middleware.redirectWithoutWWW())
	}

	if r.Middleware.EnableResponseCache {
		router.Use(r.Middleware.Cache.Check())
	}
}
