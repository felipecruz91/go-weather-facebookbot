package main

import "testing"
import "net/http"
import "net/http/httptest"

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
