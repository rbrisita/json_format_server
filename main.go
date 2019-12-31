package main

import (
	"bytes"
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
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", json_format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view", &JSONInput{Data: ""})
}

func manHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t JSONInput
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	var json_pretty bytes.Buffer
	err = json.Indent(&json_pretty, []byte(t.Data), "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`'` + string(json_pretty.Bytes()) + `'`))
}

func libHandler(w http.ResponseWriter, r *http.Request) {
	var t JSONInput
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		t.Data = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`'` + t.Data + `'`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", viewHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/man", manHandler).Methods(http.MethodPost)
	api.HandleFunc("/lib", libHandler).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
