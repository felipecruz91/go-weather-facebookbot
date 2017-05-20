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

func TestBuildForecastURL(t *testing.T) {
	// Arrange
	var city = "Manchester"
	var duration = "5"
	var expectedURL = "https://query.yahooapis.com/v1/public/yql?format=json&q=select+item.forecast+from+weather.forecast%280%2C5%29+where+woeid+in+%28select+woeid+from+geo.places%281%29+where+text%3D%22Manchester%22%29++and+u%3D%27c%27"
	// Act
	var url = BuildForecastURL(city, duration)

	// Assert
	if url != expectedURL {
		t.Errorf("Url was incorrect, got: %s, expected: %s", url, expectedURL)
	}
}
