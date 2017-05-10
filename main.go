package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"io"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

//APIAIRequest : Incoming request format from APIAI
type APIAIRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
		Parameters map[string]string `json:"parameters"`
		Contexts   []interface{}     `json:"contexts"`
		Metadata   struct {
			IntentID                  string `json:"intentId"`
			WebhookUsed               string `json:"webhookUsed"`
			WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
			IntentName                string `json:"intentName"`
		} `json:"metadata"`
		Score float32 `json:"score"`
	} `json:"result"`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
	SessionID       string      `json:"sessionId"`
	OriginalRequest interface{} `json:"originalRequest"`
}

//APIAIMessage : Response Message Structure
type APIAIMessage struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Data        Data   `json:"data"`
	Source      string `json:"source"`
}

type Data struct {
	Facebook Facebook `json:"facebook"`
}

type Facebook struct {
	Text string `"json:text"`
}

type WeatherInfo struct {
	Temp     string
	Humidity string
	Weth     string
	Units
}

type Units struct {
	Tp string
}

type Location struct {
	City  string
	State string
}

func BuildLocation(city string, state string) (loc *Location) {
	return &Location{
		city,
		state,
	}
}

func BuildUrl(loc *Location) (urlParsed string) {
	Url, _ := url.Parse("https://query.yahooapis.com/v1/public/yql")
	parameters := url.Values{}
	parameters.Add("q", "select * from weather.forecast where woeid in (select woeid from geo.places(1) where text=\""+loc.City+", "+loc.State+"\")  and u='c'")
	parameters.Add("format", "json")
	Url.RawQuery = parameters.Encode()
	urlParsed = Url.String()
	return
}

func MakeQuery(weatherUrl string) (w *WeatherInfo) {
	resp, err := http.Get(weatherUrl)
	if err != nil {
		fmt.Println("Connected Error")
		return nil
	}

	defer resp.Body.Close()
	body, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		fmt.Println("Cannot Read Information")
		return nil
	}

	js, e := simplejson.NewJson(body)
	if e != nil {
		fmt.Println("Parsing Json Error")
		return nil
	}

	//parse json
	w = new(WeatherInfo)
	w.Tp, _ = js.Get("query").Get("results").Get("channel").Get("units").Get("temperature").String()
	w.Temp, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("temp").String()
	w.Weth, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("text").String()
	w.Humidity, _ = js.Get("query").Get("results").Get("channel").Get("atmosphere").Get("humidity").String()
	return
}

func HealthCheckEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

//WebhookEndpoint - HTTP Request Handler for /webhook
func WebhookEndpoint(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		// Not sure but I think it is needed to validate webhook from Facebook
		hubmode := req.URL.Query().Get("hub.mode")
		verifyToken := req.URL.Query().Get("hub.verify_token")
		hubchallenge := req.URL.Query().Get("hub.challenge")

		fmt.Printf("Validating webhook. hubmode %s verifyToken: %s hubchallenge: %s", hubmode, verifyToken, hubchallenge)

		if hubmode == "subscribe" && verifyToken == "123456" {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, hubchallenge)
		} else {
			http.Error(w, "Failed validation. Make sure the validation tokens match.", http.StatusForbidden)
		}
	} else if req.Method == "POST" {

		decoder := json.NewDecoder(req.Body)

		var t APIAIRequest
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error in decoding the Request data", http.StatusInternalServerError)
		}

		loc := BuildLocation("Manchester", "Greater Manchester")
		queryURL := BuildUrl(loc)
		z := MakeQuery(queryURL)
		if w == nil {
			fmt.Printf("Program Error")
			log.Printf("Program Error")
		} else {
			fmt.Printf("Temperature: %s %s, %s, Humidity: %s", z.Temp, z.Tp, z.Weth, z.Humidity)
			log.Printf("Temperature: %s %s, %s, Humidity: %s", z.Temp, z.Tp, z.Weth, z.Humidity)

			data := Data{
				Facebook: Facebook{
					Text: "abc",
				},
			}

			log.Printf("data %s", data)

			msg := APIAIMessage{Source: "Weather Agent System", Speech: "Temperature: " + z.Temp + z.Tp, DisplayText: "Temperature: " + z.Temp + z.Tp, Data: data}
			log.Printf("msg %s", msg)

			json.NewEncoder(w).Encode(msg)
		}
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
	router.HandleFunc("/healthcheck", HealthCheckEndpoint).Methods("GET")
	router.HandleFunc("/webhook", WebhookEndpoint).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":"+GetPortOrDefault("4747"), router))
}
