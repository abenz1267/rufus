package rufus

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/acme/autocert"
)

type server struct {
	Dev             bool   `json:"dev,omitempty"`
	ProductionHost  string `json:"production_host,omitempty"`
	ProductionPort  string `json:"production_port,omitempty"`
	DevelopmentHost string `json:"development_host,omitempty"`
	StaticURL       string `json:"static_url,omitempty"`
	StaticFolder    string `json:"static_folder,omitempty"`
}

func (s server) startDevelopment(r http.Handler) error {
	if _, err := os.Stat("./certs/"); os.IsNotExist(err) {
		generateCert()
	}

	log.Println("Starting development server...")

	server := &http.Server{Addr: s.DevelopmentHost, Handler: r}

	if err := server.ListenAndServeTLS(filepath.Join("certs", "cert.pem"), filepath.Join("certs", "key.pem")); err != nil {
		return err
	}

	return nil
}

func (s server) startProduction(r http.Handler) error {
	var host string

	if strings.Contains(s.ProductionHost, "www.") {
		host = strings.TrimPrefix(host, "www.")
	} else {
		host = "www." + host
	}

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(s.ProductionHost, host),
		Cache:      autocert.DirCache("certs"),
	}

	// server config
	server := &http.Server{
		Addr:    ":" + s.ProductionPort,
		Handler: r,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	// TLS redirect + LetsEncrypt
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))

	log.Println("Starting production server...")

	if err := server.ListenAndServeTLS("", ""); err != nil {
		return err
	}

	return nil
}

func (s server) startFileServer(r chi.Router, path string, root http.FileSystem) {
	var cacheHeaders = map[string]string{
		"Vary":          "Accept-Encoding",
		"date":          time.Now().Format(http.TimeFormat),
		"Expires":       time.Now().AddDate(1, 0, 0).Format(http.TimeFormat),
		"Cache-Control": "public, max-age=31536000",
	}

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range cacheHeaders {
			w.Header().Set(k, v)
		}

		fs.ServeHTTP(w, r)
	}))
}
