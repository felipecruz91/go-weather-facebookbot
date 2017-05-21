package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var accessToken = os.Getenv("ACCESS_TOKEN")
var verifyToken = os.Getenv("VERIFY_TOKEN")
var port = os.Getenv("PORT")

func verifyTokenAction(w http.ResponseWriter, req *http.Request) {

	hubmode := req.URL.Query().Get("hub.mode")
	token := req.URL.Query().Get("hub.verify_token")
	hubchallenge := req.URL.Query().Get("hub.challenge")

	if hubmode == "subscribe" && token == verifyToken {
		log.Print("verify token success.")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, hubchallenge)
	} else {
		log.Print("Error: verify token failed.")
		http.Error(w, "Failed validation. Make sure the validation tokens match.", http.StatusForbidden)
	}
}

// APIAiHandler handles the requests from API.AI
func APIAiHandler(w http.ResponseWriter, req *http.Request) {

	var t APIAIRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error in decoding the Request data", http.StatusInternalServerError)
	}

	if t.Result.Action == "weather" {

		city := t.Result.Parameters["location"]

		weather := RequestWeather(city)
		if weather == nil {
			log.Printf("Error requesting the weather information for location %s", city)
		} else {
			emoji := ResolveEmoji(weather.Code)
			apiResponseText := fmt.Sprintf(
				"The weather in %s is %s %s! The temperature is  %sº%s and %s humidity.",
				city, weather.Text, emoji, weather.Temp, weather.Scale, weather.Humidity)

			msg := APIAIMessage{
				Source:      "Weather Agent System",
				Speech:      apiResponseText,
				DisplayText: apiResponseText,
			}

			json.NewEncoder(w).Encode(msg)
		}
	}

	if t.Result.Action == "forecast" {

		city := t.Result.Parameters["location"]
		duration := t.Result.Parameters["duration"]

		forecast := RequestForecast(city, duration)

		var buffer bytes.Buffer
		for _, element := range forecast {
			item := element.Item.Forecast
			emoji := ResolveEmoji(item.Code)
			buffer.WriteString(item.Day + " " + item.Date + " " + item.Text + " " + emoji + " (" + item.Low + "º/" + item.High + "º)" + "\n")
		}

		apiResponseText := fmt.Sprintf(
			"The forecast in %s for the next %s days is:\n %s",
			city, duration, buffer.String())

		msg := APIAIMessage{
			Source:      "Weather Agent System",
			Speech:      apiResponseText,
			DisplayText: apiResponseText,
		}

		json.NewEncoder(w).Encode(msg)
	}
}

//WebhookEndpoint - HTTP Request Handler for /webhook
func WebhookEndpoint(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		verifyTokenAction(w, req)
	} else if req.Method == "POST" {

		var receivedMessage ReceivedMessage

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&receivedMessage)
		if err != nil {
			log.Print(err)
			http.Error(w, "Error in decoding the Request data", http.StatusInternalServerError)
		}

		messagingEvents := receivedMessage.Entry[0].Messaging
		for _, event := range messagingEvents {
			if &event.Message != nil && event.Message.Text != "" && !event.Message.IsEcho {

				// Send request to API.AI
				response, err := PerformRequestToAPIAi(event.Message.Text)
				if err != nil {
					log.Print(err)
				}

				// Send response back to Facebook bot
				botID := event.Sender.ID
				SendMessageToBot(botID, response)
			}
		}
	} else {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
}

// GetPortOrDefault gets the Port from the environment or sets it to a default value.
func GetPortOrDefault(defaultPort string) string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return port
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", WebhookEndpoint).Methods("GET", "POST")
	router.HandleFunc("/ai", APIAiHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+GetPortOrDefault("4747"), router))
}
