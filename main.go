package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"goWeatherApi/internal"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.GET("/weather", getWeatherApi)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func getWeatherApi(c *echo.Context) error {
	form := new(internal.FormFields)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid form data")
	}
	city := form.City
	country := form.Country

	if city == "" || country == "" {
		return c.JSON(http.StatusBadRequest, "Missing city or country")
	}
	getWeather(city, country, c)
	return nil
}

func getWeather(city string, country string, c *echo.Context) {
	// Get API key from environment variable
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	if weatherApiKey == "" {
		c.JSON(http.StatusInternalServerError, "WEATHER_API_KEY environment variable is not set")
		return
	}

	//get lat and lon from city name using Google Geocoding API
	lat, lon, err := getLatLong(city, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// Construct API URL for Google Maps Weather API
	// if useGoogleWeather {
	// url := fmt.Sprintf("https://weather.googleapis.com/v1/currentConditions:lookup?key=%s&location.latitude=%f&location.longitude=%f", weatherApiKey, lat, lon)
	// } else {
	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=%s&exclude=hourly,daily,minutely,alerts&units=metric", lat, lon, weatherApiKey)
	// }

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Error making request to Weather API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Weather API returned status code " + resp.Status})
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Error reading response from Weather API"})
		return
	}

	// Parse JSON
	// if useGoogleWeather {
	// 	var data GoogleWeatherData
	// 	if err := json.Unmarshal(body, &data); err != nil {
	// 		fmt.Printf("Error parsing JSON from Google Maps Weather API: %v\n", err)
	// 		return
	// 	}
	// 	parseResponseFromGoogle(data, w)
	// } else {
	var data internal.WeartherApiData
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{"error": "Error parsing JSON from OpenWeather API"})
		return
	}
	parseResponseFromOpenWeather(data, c, city, country)
}

func getLatLong(city string, country string) (float64, float64, error) {
	//get lat and lon from city name using Google Geocoding API
	googleMapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsApiKey == "" {
		return 0, 0, fmt.Errorf("GOOGLE_MAPS_API_KEY environment variable is not set")
	}
	geoUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,%s&key=%s", city, country, googleMapsApiKey)
	geoResp, err := http.Get(geoUrl)
	if err != nil {
		return 0, 0, fmt.Errorf("error making request to Google Maps API: %v", err)
	}
	defer geoResp.Body.Close()
	if geoResp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Google Maps API returned status code %s", geoResp.Status)
	}
	geoBody, err := io.ReadAll(geoResp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response from Google Maps API: %v", err)
	}
	var geoData internal.GeoData
	if err := json.Unmarshal(geoBody, &geoData); err != nil {
		return 0, 0, fmt.Errorf("error parsing JSON from Google Maps API: %v", err)
	}

	if len(geoData.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for the city")
	}

	lat := geoData.Results[0].Geometry.Location.Lat
	lon := geoData.Results[0].Geometry.Location.Lng

	return lat, lon, nil
}

func parseResponseFromGoogle(data internal.GoogleWeatherData, c *echo.Context, city string, country string) {
	// Print weather information
	c.JSON(http.StatusOK, map[string]any{
		"city":        city,
		"country":     country,
		"temperature": data.Temperature.Degrees,
		"condition":   data.WeatherCondition.Description.Text,
		"wind":        data.Wind.Speed.Value,
		"humidity":    data.RelativeHumidity,
	})
}

func parseResponseFromOpenWeather(data internal.WeartherApiData, c *echo.Context, city string, country string) {
	// Print weather information
	c.JSON(http.StatusOK, map[string]any{
		"city":        city,
		"country":     country,
		"temperature": data.Current.Temp,
		"condition":   data.Current.Weather[0].Main,
		"wind":        data.Current.WindSpeed,
		"humidity":    data.Current.Humidity,
		"data":        data,
	})
}
