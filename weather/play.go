package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const API_KEY = "80582e456f6cb4d3cfb9c2e9570b35b6"

type Main struct {
	Temp                                    float64 `json:"temp"`
	TempMin                                 float64 `json:"temp_min"`
	TempMax                                 float64 `json:"temp_max"`
	FeelsLike                               float64 `json:"feels_like"`
	Pressure, Humidity, SeaLevel, GrndLevel int
}

type WeatherData struct {
	Main  Main
	Coord Location
}

type Location struct {
	Lat, Lon float64
}

func FetchWeather(location Location, ch chan<- WeatherData, wg *sync.WaitGroup) (WeatherData, error) {
	defer wg.Done()

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&appid=%v&units=metric", location.Lat, location.Lon, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v", err)
	}

	ch <- weatherData

	return weatherData, nil
}

func GetGeoCoordinates(city string, ch chan<- Location, wg *sync.WaitGroup) (Location, error) {
	defer wg.Done()
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&&appid=%s", city, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		return Location{}, fmt.Errorf("API call failed!")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
	}

	// Parse the JSON response
	var locations []Location
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %v", err)
	}

	ch <- locations[0]
	return locations[0], nil
}

func DisplayWeatherData(weatherdata *WeatherData) {
	fmt.Printf("Location:  %v\n", weatherdata.Coord)
	fmt.Printf("The temperature is %vC but feels like %vC\n", weatherdata.Main.Temp, weatherdata.Main.FeelsLike)
	fmt.Printf("Temperature ranges from %vC min to %vC maximum\n", weatherdata.Main.TempMin, weatherdata.Main.TempMax)
	fmt.Println()
}

func main() {
	cities := []string{"Berlin", "Lahore", "London"}

	location_channel := make(chan Location, len(cities))
	var wg1 sync.WaitGroup

	for _, city := range cities {
		wg1.Add(1)
		go GetGeoCoordinates(city, location_channel, &wg1)
	}

	go func() {
		wg1.Wait()
		close(location_channel)
	}()

	weather_channel := make(chan WeatherData)
	var wg2 sync.WaitGroup
	for loc := range location_channel {
		wg2.Add(1)
		go FetchWeather(loc, weather_channel, &wg2)

	}

	go func() {
		wg2.Wait()
		close(weather_channel)
	}()

	for result := range weather_channel {
		DisplayWeatherData(&result)
	}

}
