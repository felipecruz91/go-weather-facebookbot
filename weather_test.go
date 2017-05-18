package main

import (
	"testing"
)

func TestBuildWeatherURL(t *testing.T) {

	// Arrange
	var city = "Manchester"
	var expectedURL = "https://query.yahooapis.com/v1/public/yql?format=json&q=select+%2A+from+weather.forecast+where+woeid+in+%28select+woeid+from+geo.places%281%29+where+text%3D%22" + city + "%22%29++and+u%3D%27c%27"

	// Act
	var url = BuildWeatherURL(city)

	// Assert
	if url != expectedURL {
		t.Errorf("Url was incorrect, got: %s, expected: %s", url, expectedURL)
	}
}
