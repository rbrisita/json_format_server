package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFiveHitsPerSec(t *testing.T) {
	status_codes, err := runRequests(5)
	if err != nil {
		t.Fatal(err)
	}

	if containsStatusCode429(status_codes) {
		t.Errorf("test found status codes of %v", http.StatusTooManyRequests)
	}
}

func TestTenHitsPerSec(t *testing.T) {
	// Wait for prior test connections to close
	time.Sleep(time.Millisecond * 500)

	status_codes, err := runRequests(10)
	if err != nil {
		t.Fatal(err)
	}

	if containsStatusCode429(status_codes) {
		t.Errorf("test found status codes of %v", http.StatusTooManyRequests)
	}
}

// TestTooManyRequests check for status code 429
func TestTooManyRequests(t *testing.T) {
	// Wait for prior test connections to close
	time.Sleep(time.Millisecond * 500)

	status_codes, err := runRequests(11)
	if err != nil {
		t.Fatal(err)
	}

	if !containsStatusCode429(status_codes) {
		t.Errorf("test did not find any status codes of %v", http.StatusTooManyRequests)
	}
}

// runRequests loops requesting total number given.
func runRequests(total int) ([]int, error) {
	status_codes := make([]int, total)

	for i := 0; i < total; i++ {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		if err != nil {
			return status_codes, err
		}

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(viewHandler)
		rateLimit(handler).ServeHTTP(res, req)
		status_codes[i] = res.Code
	}

	return status_codes, nil
}

// containsStatusCode429 loops through given array for integer value of 429.
func containsStatusCode429(status_codes []int) bool {
	too_many_requests := 0
	for _, v := range status_codes {
		if v == http.StatusTooManyRequests {
			too_many_requests++
		}
	}

	return too_many_requests > 0
}
