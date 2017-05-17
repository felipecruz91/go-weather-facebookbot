package main

import (
	"os"
	"testing"
)

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
