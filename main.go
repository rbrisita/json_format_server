package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

/**
* Flags for command line configuration.
**/
var (
	port = flag.Uint64("port", 8080, "Port to listen on.")
	host = flag.String("addr", "", "The server's host.")
)

// templates holds pre-parsed template files.
var templates = template.Must(template.ParseFiles("view.tmpl"))

// Dummy is an empty struct for JSON:API unmarshalling.
type Dummy struct{}

// JSONInput holds the JSON string to be formatted.
type JSONInput struct {
	Data string
}

/**
* Main exposes 3 API endpoints on a host and port.
**/
func main() {
	flag.Parse()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", viewHandler)

	api := rtr.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/std", stdHandler).Methods(http.MethodPost)
	api.HandleFunc("/spec", specHandler).Methods(http.MethodPost)

	fmt.Printf("Listening on %s:%d...\n", *host, *port)

	log.Fatal(http.ListenAndServe(*host+":"+strconv.FormatUint(*port, 10), rateLimit(rtr)))
}

/**
* viewHandler presents a view to the client.
**/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "view")
}

/**
* renderTemplate renders template that was parsed at program boot.
**/
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".tmpl", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
* stdHandler verifies given JSON and formats it for client.
**/
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

/**
* specHandler verifies given JSON by the JSON:API spec
* and formats it for client.
**/
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

/**
* errorResponse creates an error message for the client to process.
**/
func errorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

/**
* getPrettyJSONString indents given JSON byte array into a returned string.
**/
func getPrettyJSONString(src []byte) (string, error) {
	var json_pretty bytes.Buffer
	err := json.Indent(&json_pretty, src, "", "\t")
	if err != nil {
		return "", err
	}

	return string(json_pretty.Bytes()), nil
}
