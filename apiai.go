package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const apiEndpoint string = "https://api.api.ai/v1/%s?v=%s"
const apiVersion string = "20150910"

var apiAccessToken = os.Getenv("APIAI_ACCESS_TOKEN")

//APIAIRequest : Incoming request format from APIAI
type APIAIRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
		Source           string            `json:"source"`
		ResolvedQuery    string            `json:"resolvedQuery"`
		Action           string            `json:"action"`
		ActionIncomplete bool              `json:"actionIncomplete"`
		Parameters       map[string]string `json:"parameters"`
		Contexts         []struct {
			Name       string `json:"name"`
			Parameters struct {
				Name string `json:"name"`
			} `json:"parameters"`
			Lifespan int `json:"lifespan"`
		} `json:"contexts"`
		Metadata struct {
			IntentID   string `json:"intentId"`
			IntentName string `json:"intentName"`
		} `json:"metadata"`
		Fulfillment struct {
			Speech      string `json:"speech"`
			DisplayText string `json:"displayText"`
			Source      string `json:"source"`
		} `json:"fulfillment"`
	} `json:"result"`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
}

//APIAIMessage : Response Message Structure
type APIAIMessage struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Source      string `json:"source"`
}

// PerformRequestToAPIAi sends natural language text and information as query parameters to API.AI
func PerformRequestToAPIAi(text string) (APIAIRequest, error) {

	record := APIAIRequest{}
	myURL := fmt.Sprintf(apiEndpoint, "query", apiVersion)
	myURL = myURL + "&query=" + url.QueryEscape(text) + "&lang=en" + "&sessionId=1234567890"

	// Build the request
	req, err := http.NewRequest("GET", myURL, nil)
	if err != nil {
		log.Fatalf("The /GET request to %s failed", myURL)
		return record, err
	}

	// Replace authToken by your Client access token
	authValue := "Bearer " + apiAccessToken
	req.Header.Add("Authorization", authValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return record, err
	}

	// Callers should defer the close of resp.Body when done reading from it
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	return record, nil
}

// ResolveEmoji converts the weather code into an emoji
func ResolveEmoji(weatherCode string) (emoji string) {

	switch weatherCode {
	case "4": // thunderstorms
		return "⛈️⚡"
	case "11", "12":
		return "🌧️☔"
	case "16":
		return "🌨️❄️"
	case "20":
		return "🌫️"
	case "24":
		return "💨"
	case "25":
		return "🐧🐧"
	case "28": // mostly cloudy (day)
		return "☁️"
		case "30": // partly cloudy (day)
		return "⛅"
	case "32":
		return "☀️"
	case "36":
		return "🔥🔥"
	case "38", "39": // scattered thunderstorms
		return "⛈️"
	default:
		fmt.Printf("%s.", weatherCode)
		return ""
	}

}
