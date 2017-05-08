package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthCheckEndpoint(t *testing.T) {
	// Create a request to pass to our endpoint. We don't have any query parameters for now,
	// so we'll pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckEndpoint)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetPortWhenPortIsDefined(t *testing.T) {

	// Arrange
	var defaultPort = "4747"
	var expectedPort = "4747"
	os.Setenv("PORT", expectedPort)

	// Act
	var port = GetPortOrDefault(defaultPort)

	// Assert
	if port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, expected: %s", port, expectedPort)
	}
}
func TestGetPortWhenPortIsNotDefined(t *testing.T) {

	// Arrange
	var defaultPort = "4747"
	os.Setenv("PORT", "")

	// Act
	var port = GetPortOrDefault(defaultPort)

	// Assert
	if port != defaultPort {
		t.Errorf("Port was incorrect, got: %s, expected: %s", port, defaultPort)
	}
}
