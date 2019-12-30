package main

import (
	"log"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("view.tmpl"))

type JSONFormat struct {
	Data string
}

func renderTemplate(w http.ResponseWriter, tmpl string, json_format *JSONFormat) {
	err := templates.ExecuteTemplate(w, tmpl + ".tmpl", json_format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view", &JSONFormat{Data: ""})
}

func main() {
	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
