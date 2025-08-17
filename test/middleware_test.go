package test

import (
	"devices-api/internal/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a simple handler that returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap the handler with logging middleware
	loggedHandler := middleware.LoggingMiddleware(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Execute the request
	loggedHandler.ServeHTTP(rr, req)

	// Check that the response is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "test response"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoggingMiddleware_StatusCode(t *testing.T) {
	// Create a handler that returns 404 Not Found
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	// Wrap the handler with logging middleware
	loggedHandler := middleware.LoggingMiddleware(testHandler)

	// Create a test request
	req := httptest.NewRequest("POST", "/api/v1/devices", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-client/1.0")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Execute the request
	loggedHandler.ServeHTTP(rr, req)

	// Check that the response status is correct
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := "not found"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoggingMiddleware_DifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			// Create a simple handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			})

			// Wrap with logging middleware
			loggedHandler := middleware.LoggingMiddleware(testHandler)

			// Create test request with the specific method
			var req *http.Request
			if method == "POST" || method == "PUT" || method == "PATCH" {
				req = httptest.NewRequest(method, "/test", strings.NewReader(`{"test":"data"}`))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(method, "/test", nil)
			}

			rr := httptest.NewRecorder()
			loggedHandler.ServeHTTP(rr, req)

			// Verify the response
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Handler returned wrong status code for %s: got %v want %v", method, status, http.StatusOK)
			}
		})
	}
}
