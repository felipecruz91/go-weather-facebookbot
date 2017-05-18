package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	simplejson "github.com/bitly/go-simplejson"
)

// WeatherInfo struct
type WeatherInfo struct {
	Temp     string
	Scale    string
	Humidity string
	Text     string
	Code     string
}

const yahooWeatherAPIURL = "https://query.yahooapis.com/v1/public/yql"

// BuildWeatherURL builds the Yahoo API weather URL
func BuildWeatherURL(city string) (urlParsed string) {
	URL, _ := url.Parse(yahooWeatherAPIURL)
	parameters := url.Values{}
	parameters.Add("q", "select * from weather.forecast where woeid in (select woeid from geo.places(1) where text=\""+city+"\")  and u='c'")
	parameters.Add("format", "json")
	URL.RawQuery = parameters.Encode()
	urlParsed = URL.String()
	return
}

// RequestWeather performs a call to the Yahoo Weather API and returns the weather information.
func RequestWeather(city string) (w *WeatherInfo) {

	weatherURL := BuildWeatherURL(city)

	resp, err := http.Get(weatherURL)
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

	weatherInfo := new(WeatherInfo)

	weatherInfo.Temp, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("temp").String()
	weatherInfo.Scale, _ = js.Get("query").Get("results").Get("channel").Get("units").Get("temperature").String()
	weatherInfo.Humidity, _ = js.Get("query").Get("results").Get("channel").Get("atmosphere").Get("humidity").String()
	weatherInfo.Text, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("text").String()
	weatherInfo.Code, _ = js.Get("query").Get("results").Get("channel").Get("item").Get("condition").Get("code").String()

	return weatherInfo
}
