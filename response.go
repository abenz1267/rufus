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
			log.Printf("Error in response.renderHTML: %e", err)
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

func (resp Response) renderHTMLDev(w http.ResponseWriter, r *http.Request, funcs template.FuncMap, templateFolder, baseTemplate string) {

	if r.Header.Get("Accept") == "text/html-raw" {
		tmpl, err := template.New(baseTemplate).Funcs(funcs).ParseFiles(templateFolder+"/"+baseTemplate, templateFolder+"/"+resp.TemplateFile+"_raw.html")
		if err != nil {
			resp.printTemplateError("renderHTMLDev", err, w)
		}

		err = tmpl.ExecuteTemplate(w, baseTemplate, resp.Data)
		if err != nil {
			resp.printTemplateError("renderHTMLDev", err, w)
		}
		return
	}

	tmpl, err := template.New(baseTemplate).Funcs(funcs).ParseFiles(templateFolder+"/"+baseTemplate, templateFolder+"/"+resp.TemplateFile+".html")
	if err != nil {
		resp.printTemplateError("renderHTMLDev", err, w)
	}

	err = tmpl.ExecuteTemplate(w, baseTemplate, resp.Data)
	if err != nil {
		resp.printTemplateError("renderHTMLDev", err, w)
	}
}

func (resp Response) printTemplateError(where string, err error, w http.ResponseWriter) {
	log.Printf("%s: %g", where, err)
	http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
	return
}

func (resp Response) renderJSON(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(resp.Data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
	}
}
