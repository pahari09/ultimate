package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
)

const (
	apiKey = "e3a6589e82b161dd83124dd1fd211243"
)

var (
	router = gin.Default()
)

// MODEL
const desiredTemplate = `{params: [{name: temperature, value: {{.Main.Temp}}}]}`

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"mains"`
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
	} `json:"mains"`
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

type RequestBody struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

// Main Function (Setting the route)
func main() {
	println("about to start the application...")

	router.GET("/temperature", GetTemperature)

	err := router.Run(":8082")
	if err != nil {
		return
	}
}

// Controller
func GetTemperature(c *gin.Context) {
	var input RequestBody
	var data WeatherData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	data, statusCode := GetWeatherData(c)

	tmpl, err := template.New("weather").Parse(desiredTemplate)
	if err != nil {
		log.Fatalln(err)
	}
	// Render the template and display
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		log.Fatalln(err)
	}
	print("HTTP status: %s", statusCode)
}

// Intermediate function called from handler
func GetWeatherData(c *gin.Context) (data WeatherData, status string) {
	city := c.Param("city")
	country := c.Param("country")
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
	fmt.Println("API Response as struct %+v\n", weatherStruct)

	return weatherStruct, resp.Status
}
