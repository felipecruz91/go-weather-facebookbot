package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID    string `json:"mid"`
	Seq    int64  `json:"seq"`
	Text   string `json:"text"`
	IsEcho bool   `json:"is_echo"`
}

type Payload struct {
	TemplateType string  `json:"template_type"`
	Text         string  `json:"text"`
	Buttons      Buttons `json:"buttons"`
}

type Buttons struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type ButtonMessageBody struct {
	Attachment Attachment `json:"attachment"`
}

type ButtonMessage struct {
	Recipient         Recipient         `json:"recipient"`
	ButtonMessageBody ButtonMessageBody `json:"message"`
}

type SendMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

const FacebookEndPoint = "https://graph.facebook.com/v2.6/me/messages"

// SendMessageToBot sends a message to the Facebook bot
func SendMessageToBot(botID string, response APIAIRequest) {

	recipient := new(Recipient)
	recipient.ID = botID
	sendMessage := new(SendMessage)
	sendMessage.Recipient = *recipient
	sendMessage.Message.Text = response.Result.Fulfillment.Speech
	sendMessageBody, err := json.Marshal(sendMessage)
	if err != nil {
		log.Print(err)
	}

	req, err := http.NewRequest("POST", FacebookEndPoint, bytes.NewBuffer(sendMessageBody))
	if err != nil {
		log.Print(err)
	}

	values := url.Values{}
	values.Add("access_token", accessToken)
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()
	var result map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print(err)
	}
	log.Print(result)
}
