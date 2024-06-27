package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestSubmitHandler tests the SubmitHandler function
func TestSubmitHandler(t *testing.T) {
	// Ensure the working directory is the project root
	// This is necessary to correctly locate the template files
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	// Create a URLRequest payload to be sent in the POST request
	urlRequestPayload := URLRequest{URL: "http://example.com"}
	// Marshal the payload into JSON format
	urlRequestPayloadBytes, err := json.Marshal(urlRequestPayload)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new POST request to the /submit endpoint with the JSON payload
	postRequest, err := http.NewRequest("POST", "/submit", bytes.NewBuffer(urlRequestPayloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	// Set the Content-Type header to application/json
	postRequest.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	responseRecorder := httptest.NewRecorder()
	// Wrap the SubmitHandler function with http.HandlerFunc
	submitHandler := http.HandlerFunc(SubmitHandler)

	// Serve the HTTP request using the handler
	submitHandler.ServeHTTP(responseRecorder, postRequest)

	// Check if the status code returned by the handler is 200 OK
	if statusCode := responseRecorder.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, http.StatusOK)
		return
	}

	// Get the response body as a string
	responseBody := responseRecorder.Body.String()
	// Check if the response body contains the expected HTML elements and content
	if !strings.Contains(responseBody, `<h2>Response</h2>`) ||
		!strings.Contains(responseBody, `<pre>{success http://example.com}</pre>`) {
		t.Errorf("handler returned unexpected body: got %v", responseBody)
	}
}
