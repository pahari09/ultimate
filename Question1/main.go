/*
PS S:\Ultimate> go run main.go -country germany -city dortmund
URL string: https://api.openweathermap.org/data/2.
5/weather?appid=e3a6589e82b161dd83124dd1fd211243&q=dortmund,germany&units=metric
API Response as struct {Coord:{Lon:7.45 Lat:51.5167} Weather:[{ID:800 Main:Clear Description:clear sky Icon:01d}] Base:stations Main:{Temp:17.67 Pressure:1028 Humidity:36 TempMin:16.09
TempMax:19.94} Visibility:10000 Wind:{Speed:2.57 Deg:70} Clouds:{All:0} Dt:1648125629 Sys:{Type:2 ID:2007579 Message:0 Country:DE Sunrise:1648099447 Sunset:1648144128} ID:2935517 Name:Dortmund Cod:200}
{params: [{name: temperature, value: 17.67}]}
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

// template used for output as per Question

const desiredTemplate = `{params: [{name: temperature, value: {{.Main.Temp}}}]}`

// Pointers to hold the contents of the flag args.
var (
	apiKey      = "e3a6589e82b161dd83124dd1fd211243"
	countryFlag = flag.String("country", "", "Location to get weather.  "+
		"If country has a space, wrap the location in double quotes.")
	cityFlag = flag.String("city", "", "If country has a space, wrap the location in double quotes.")
)

// this is how Open weather api should look like
// const link = "https://api.openweathermap.org/data/2.5/weather/?q=London," +
//  "uk&appid=e3a6589e82b161dd83124dd1fd211243&units=metric"

// WeatherData to hold the result of the query

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func getWeatherData(country string, city string) {

	// concatenating for passing as query parameters
	query := city + "," + country
	// creating http client for handling request
	client := &http.Client{}

	// creating URL structure
	req, err := http.NewRequest(
		"GET", "https://api.openweathermap."+
			"org/data/2."+
			"5/weather?", nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// Query params
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("appid", apiKey)
	q.Add("units", "metric")

	req.URL.RawQuery = q.Encode()
	decoded, err := url.QueryUnescape(q.Encode())
	req.URL.RawQuery = decoded
	if err != nil {
		log.Fatalln(err)
	}
	// printing query string for reference
	fmt.Println("URL string:", req.URL.String())

	// Making HTTP GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to weatherData struct
	var weatherStruct WeatherData
	err = json.Unmarshal(bodyBytes, &weatherStruct)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("API Response as struct %+v\n", weatherStruct)

	// Desired output template
	tmpl, err := template.New("weather").Parse(desiredTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	// Render the template and display
	if err := tmpl.Execute(os.Stdout, weatherStruct); err != nil {
		log.Fatalln(err)
	}
}

func main() {

	flag.Parse()

	// Basic error handling of cli arguments
	if len(*countryFlag) <= 1 || len(*cityFlag) <= 1 {
		flag.Usage()
		os.Exit(1)
	}
	getWeatherData(*countryFlag, *cityFlag)
	os.Exit(0)
}
