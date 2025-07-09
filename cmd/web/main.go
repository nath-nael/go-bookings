package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/nath-nael/go-learn/pkg/config"
	"github.com/nath-nael/go-learn/pkg/handlers"
	"github.com/nath-nael/go-learn/pkg/render"
)
const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager



func main() {

	app.InProduction = false // Change to true in production

	session =scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // Set to true if using HTTPS

	app.Session = session
	
	tc, err:= render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache:", err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)


	render.NewTemplate(&app)

	fmt.Println(fmt.Sprintf("Starting the server on port http.localhost%s",portNumber))

	srv:= &http.Server{
		Addr: 	portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}