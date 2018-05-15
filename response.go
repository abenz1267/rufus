package rufus

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Response handles the rendering
type Response struct {
	Status       int
	TemplateFile string
	Data         interface{}
}

func (resp Response) renderHTML(w http.ResponseWriter, r *http.Request, template, strippedTemplate *template.Template) {
	if r.Header.Get("Accept") == "text/html-raw" {
		if err := strippedTemplate.Execute(w, resp.Data); err != nil {
			log.Println(err)
			http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := template.Execute(w, resp.Data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
	}
}

func (resp Response) renderJSON(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(resp.Data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
	}
}
