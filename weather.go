package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	simplejson "github.com/bitly/go-simplejson"
)

type WeatherInfo struct {
	Temp     string
	Humidity string
	Weth     string
	Units
}

type Units struct {
	Tp string
}

func BuildUrl(loc string) (urlParsed string) {
	Url, _ := url.Parse("https://query.yahooapis.com/v1/public/yql")
	parameters := url.Values{}
	parameters.Add("q", "select * from weather.forecast where woeid in (select woeid from geo.places(1) where text=\""+loc+"\")  and u='c'")
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
