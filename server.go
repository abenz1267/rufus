package rufus

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

type server struct {
	Dev             bool   `json:"dev,omitempty"`
	ProductionHost  string `json:"production_host,omitempty"`
	ProductionPort  string `json:"production_port,omitempty"`
	DevelopmentHost string `json:"development_host,omitempty"`
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
