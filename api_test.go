package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

const (
	API_VER           = "/api/v1"
	CONTENT_TYPE_JSON = "application/json"
)

/**
* curl http://localhost:8080
**/
func TestViewPage(t *testing.T) {
	requestEndpoint(
		t,
		http.MethodGet,
		"/",
		"",
		"",
		viewHandler,
		http.StatusOK,
		"JSON Formatter",
	)
}

/**
* curl -d '{"data":"{}"}' http://localhost:8080/api/v1/std
**/
func TestStdEndpoint(t *testing.T) {
	requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/std",
		`{"data":"{}"}`,
		CONTENT_TYPE_JSON,
		stdHandler,
		http.StatusCreated,
		"{}",
	)
}

/**
* curl -d "{}" http://localhost:8080/api/v1/spec
**/
func TestSpecEndpoint(t *testing.T) {
	res := requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/spec",
		"{}",
		jsonapi.MediaType,
		specHandler,
		http.StatusCreated,
		"{}",
	)

	// Check content type is up to JSON:API spec
	result := res.Result()
	if content_type := result.Header.Get("Content-type"); content_type != jsonapi.MediaType {
		t.Errorf("handler returned unexpected body: got %v want %v", content_type, jsonapi.MediaType)
	}
}

// curl -d '{"data":"{\"test\":\"test\"}"}' http://localhost:8080/api/v1/std
func TestStdFormatsCorrectly(t *testing.T) {
	pretty_json := `{
	"test": "test"
}`

	requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/std",
		`{"data":"{\"test\":\"test\"}"}`,
		CONTENT_TYPE_JSON,
		stdHandler,
		http.StatusCreated,
		pretty_json,
	)
}

// curl -d '{"data":{"test":"test"}}' http://localhost:8080/api/v1/spec
func TestSpecFormatsCorrectly(t *testing.T) {
	pretty_json := `{
	"data": {
		"test": "test"
	}
}`

	res := requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/spec",
		`{"data":{"test":"test"}}`,
		jsonapi.MediaType,
		specHandler,
		http.StatusCreated,
		pretty_json,
	)

	// Check content type is up to JSON:API spec
	result := res.Result()
	if content_type := result.Header.Get("Content-type"); content_type != jsonapi.MediaType {
		t.Errorf("handler returned unexpected body: got %v want %v", content_type, jsonapi.MediaType)
	}
}

// curl -X POST http://localhost:8080/api/v1/std
func TestStdFailure(t *testing.T) {
	requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/std",
		"",
		CONTENT_TYPE_JSON,
		stdHandler,
		http.StatusUnprocessableEntity,
		`{"error": "EOF"}`,
	)
}

// curl -X POST http://localhost:8080/api/v1/spec
func TestSpecFailure(t *testing.T) {
	res := requestEndpoint(
		t,
		http.MethodPost,
		API_VER+"/spec",
		"",
		jsonapi.MediaType,
		specHandler,
		http.StatusUnprocessableEntity,
		`{"error": "EOF"}`,
	)

	// Check content type is up to JSON:API spec
	result := res.Result()
	if content_type := result.Header.Get("Content-type"); content_type != jsonapi.MediaType {
		t.Errorf("handler returned unexpected body: got %v want %v", content_type, jsonapi.MediaType)
	}
}

/**
* requestEndpoint is a boilerplate function to make testing endpoints a one liner.
**/
func requestEndpoint(
	t *testing.T,
	method string,
	url string,
	body string,
	content_type string,
	req_handler http.HandlerFunc,
	expected_status int,
	expected_data string) *httptest.ResponseRecorder {

	var r io.Reader
	if body != "" {
		var json_str = []byte(body)
		r = bytes.NewBuffer(json_str)
	} else {
		r = strings.NewReader(body)
	}

	// Request valid?
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		t.Fatal(err)
	}

	if content_type != "" {
		req.Header.Set("Content-Type", content_type)
	}

	// Check status code
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(req_handler)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != expected_status {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected_status)
	}

	// Check the response body is what we expect.
	if body := res.Body.String(); !strings.Contains(body, expected_data) {
		t.Errorf("handler returned unexpected body: got %v want %v", body, expected_data)
	}

	return res
}
