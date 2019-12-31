package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseFiles("view.tmpl"))

type Dummy struct{}

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

func stdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t JSONInput
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		errorResponse(w, err)
		return
	}

	pretty_json, err := getPrettyJSONString([]byte(t.Data))
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(pretty_json))
}

func specHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	var d Dummy
	err = jsonapi.UnmarshalPayload(bytes.NewReader(b), &d)
	if err != nil {
		errorResponse(w, err)
		return
	}

	pretty_json, err := getPrettyJSONString(b)
	if err != nil {
		errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(pretty_json))
}

func getPrettyJSONString(src []byte) (string, error) {
	var json_pretty bytes.Buffer
	err := json.Indent(&json_pretty, src, "", "\t")
	if err != nil {
		return "", err
	}

	return string(json_pretty.Bytes()), nil
}

func errorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", viewHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/std", stdHandler).Methods(http.MethodPost)
	api.HandleFunc("/spec", specHandler).Methods(http.MethodPost)

	fmt.Println("Listening on Localhost:8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
