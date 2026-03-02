package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./statics"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/not_found.html"))
			tmpl.Execute(w, nil)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
		tmpl.Execute(w, nil)
	})

	// accounts
	// http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", signinHandler)
	http.HandleFunc("/logout", signout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

}
