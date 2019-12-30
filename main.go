package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseFiles("view.tmpl"))

type JSONInput struct {
	Data string
}

func renderTemplate(w http.ResponseWriter, tmpl string, json_format *JSONInput) {
	err := templates.ExecuteTemplate(w, tmpl + ".tmpl", json_format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view", &JSONInput{Data: ""})
}

func manHandler(w http.ResponseWriter, r *http.Request) {
	var t JSONInput
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		t.Data = err.Error()
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"message": "post man called", "output": "` + t.Data + `"}`))
}

func libHandler(w http.ResponseWriter, r *http.Request) {
	var t JSONInput
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		t.Data = err.Error()
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "post lib called", "output": "` + t.Data + `"}`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", viewHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/man", manHandler).Methods(http.MethodPost)
	api.HandleFunc("/lib", libHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
