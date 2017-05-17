package main

import (
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

	//TODO: Update verify token
	if hubmode == "subscribe" && token == verifyToken {
		log.Print("verify token success.")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, hubchallenge)
	} else {
		log.Print("Error: verify token failed.")

		http.Error(w, "Failed validation. Make sure the validation tokens match.", http.StatusForbidden)
	}
}

func webhookPostAction(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)

	var receivedMessage ReceivedMessage
	err := decoder.Decode(&receivedMessage)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error in decoding the Request data", http.StatusInternalServerError)
	}

	messagingEvents := receivedMessage.Entry[0].Messaging
	for _, event := range messagingEvents {
		if &event.Message != nil && event.Message.Text != "" && !event.Message.IsEcho {
			sendTextMessage(event)
		}
	}

}

func AiEndpoint(w http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		HandleRequestFromApiAi(w, req)
	} else {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
}

//WebhookEndpoint - HTTP Request Handler for /webhook
func WebhookEndpoint(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		verifyTokenAction(w, req)
	} else if req.Method == "POST" {
		webhookPostAction(w, req)
	} else {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
}

// Get the Port from the environment so we can run on Heroku
func GetPortOrDefault(defaultPort string) string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = defaultPort
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return port
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", WebhookEndpoint).Methods("GET", "POST")
	router.HandleFunc("/ai", AiEndpoint).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":"+GetPortOrDefault("4747"), router))
}
