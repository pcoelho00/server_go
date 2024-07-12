package jsondecoders

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithError(t *testing.T) {
	// Create a new HTTP request with a response recorder
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the RespondWithError function with a custom error message
	RespondWithError(rr, http.StatusInternalServerError, "Internal Server Error")

	// Check the response status code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, rr.Code)
	}

	// Check the response body
	expected := `{"error":"Internal Server Error"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}
}
func TestRespondWithJson(t *testing.T) {
	// Create a new HTTP request with a response recorder
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Define the payload for the response
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	// Call the RespondWithJson function with the payload
	RespondWithJson(rr, http.StatusOK, payload)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expected := `{"message":"Hello, World!"}`
	if rr.Body.String() != expected {
		t.Errorf("Expected response body %q, but got %q", expected, rr.Body.String())
	}
}
func TestRespondWithNoBody(t *testing.T) {
	// Create a new HTTP request with a response recorder
	_, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the RespondWithNoBody function
	RespondWithNoBody(rr, http.StatusNoContent)

	// Check the response status code
	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, but got %d", http.StatusNoContent, rr.Code)
	}

	// Check the response body
	if rr.Body.String() != "" {
		t.Errorf("Expected empty response body, but got %q", rr.Body.String())
	}
}
