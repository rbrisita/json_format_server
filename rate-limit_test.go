package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Possible Broken - WIP
func TestFiveHitsPerSec(t *testing.T) {
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(viewHandler)
		handler.ServeHTTP(res, req)
		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
		}

		time.Sleep(time.Second / 5)
	}
}

// Possible Broken - WIP
func TestTenHitsPerSec(t *testing.T) {
	for i := 0; i < 10; i++ {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(viewHandler)
		handler.ServeHTTP(res, req)
		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTooManyRequests)
		}

		time.Sleep(time.Second / 10)
	}
}

// Broken - WIP
func TestTwentyHitsPerSec(t *testing.T) {
	var status_codes [2000]int
	client := http.DefaultClient

	for i := 0; i < len(status_codes); i++ {
		resp, err := client.Get("http://localhost:8080")
		if err != nil {
			t.Fatal(err)
		}

		status_codes[i] = resp.StatusCode
	}

	too_many_requests := 0
	for _, v := range status_codes {
		if v == http.StatusTooManyRequests {
			too_many_requests++
		}
	}

	if too_many_requests == 0 {
		t.Errorf("test did not find any status codes of %v", http.StatusTooManyRequests)
	}
}
