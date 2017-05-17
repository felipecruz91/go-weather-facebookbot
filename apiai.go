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

// Send request to API.AI
func SendTextToApiAi(text string) (APIAIRequest, error) {

	record := APIAIRequest{}
	myUrl := fmt.Sprintf(apiEndpoint, "query", apiVersion)
	myUrl = myUrl + "&query=" + url.QueryEscape(text) + "&lang=en" + "&sessionId=1234567890"

	fmt.Println(myUrl)

	// Build the request
	req, err := http.NewRequest("GET", myUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
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

	fmt.Println("Status = ", record.Status.Code)
	fmt.Println("ErrorType = ", record.Status.ErrorType)
	fmt.Println("Response = ", record.Result.Fulfillment.Speech)

	return record, nil
}

// API.AI -> Weather
func HandleRequestFromApiAi(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	var t APIAIRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error in decoding the Request data", http.StatusInternalServerError)
	}

	fmt.Println(t.Result.Action)
	fmt.Println(t.Result.Parameters["location"])

	if t.Result.Action == "weather" {

		city := t.Result.Parameters["location"]
		queryURL := BuildUrl(city)
		z := MakeQuery(queryURL)
		if w == nil {
			fmt.Printf("Program Error")
			log.Printf("Program Error")
		} else {
			apiResponseText := "The weather in " + city + " is " + z.Temp + "º" + z.Tp + " and " + z.Humidity + "% humidity"
			msg := APIAIMessage{Source: "Weather Agent System", Speech: apiResponseText, DisplayText: apiResponseText}
			json.NewEncoder(w).Encode(msg)
		}

	}
}
