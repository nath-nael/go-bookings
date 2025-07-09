package main

import (
	"fmt"
	"net/http"
	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the endpoint")
		next.ServeHTTP(w, r)
})
}
//NoSurf is a CSRF protection middleware
func NoSurf(next http.Handler) http.Handler{
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path: "/",
			Secure: app.InProduction, // Set to true if using HTTPS
			SameSite: http.SameSiteLaxMode,})
	return csrfHandler
}
// SessionLoad loads and saves the session on each request
func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}