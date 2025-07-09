package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/nath-nael/go-learn/pkg/config"
	"github.com/nath-nael/go-learn/pkg/models"
)

var functions = template.FuncMap{

}
var app *config.AppConfig
func NewTemplate(a *config.AppConfig){
	app = a
}
func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}
//Complex template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache{
		tc = app.TemplateCache
	} else{
		tc, _ = CreateTemplateCache() // Error handling can be improved
	}

	//get requested template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(w, "Template not found", http.StatusNotFound)
		return
	}
	buf :=new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to response:", err)
		return
	}


}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
		log.Println("Adding template to cache:", name)
	}
	return myCache, nil
}




//Simpler template cache
var tc = make(map[string]*template.Template)

func RenderTemplate2(w http.ResponseWriter, t string){
	var tmpl *template.Template
	var err error 
	//check if template is already parsed
	_,inMap := tc[t]
	if !inMap {
		log.Println("Creating template cache for:", t)
		err = createTemplateCache2(t)
		if err != nil {
			log.Println("Error creating template cache:", err)
		}
	} else {
		log.Println("Using cached template:", t)
	}
	tmpl =tc[t]
	err = tmpl.Execute(w, nil)
}

func createTemplateCache2(t string) error{
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",	
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err 
	}
	tc[t] = tmpl
	return nil
}