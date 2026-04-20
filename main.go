package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GoogleWeatherData struct {
	CurrentTime string `json:"currentTime"`
	TimeZone    struct {
		Id string `json:"id"`
	} `json:"timeZone"`
	IsDaytime        bool `json:"isDaytime"`
	WeatherCondition struct {
		IconBaseUri string `json:"iconBaseUri"`
		Description struct {
			Text         string `json:"text"`
			LanguageCode string `json:"languageCode"`
		} `json:"description"`
		Type string `json:"type"`
	}
	Temperature struct {
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	}
	FeelsLikeTemperature struct {
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	}
	DewPoint struct {
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	}
	HeatIndex struct {
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	}
	WindChill struct {
		Degrees float64 `json:"degrees"`
		Unit    string  `json:"unit"`
	}
	RelativeHumidity float64 `json:"relativeHumidity"`
	UvIndex          float64 `json:"uvIndex"`
	Precipitation    struct {
		Probability struct {
			Percent float64 `json:"percent"`
			Type    string  `json:"type"`
		}
		Qpf struct {
			Quantity float64 `json:"quantity"`
			Unit     string  `json:"unit"`
		}
	}
	ThunderstormProbability float64 `json:"thunderstormProbability"`
	AirPressure             struct {
		MeanSeaLevelMillibars float64 `json:"meanSeaLevelMillibars"`
	}
	Wind struct {
		Direction struct {
			Degrees  int    `json:"degrees"`
			Cardinal string `json:"cardinal"`
		}
		Speed struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		}
		Gust struct {
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		}
	}
	Visibility struct {
		Distance float64 `json:"distance"`
		Unit     string  `json:"unit"`
	}
	CloudCover               float64 `json:"cloudCover"`
	CurrentConditionsHistory struct {
		TemperatureChange struct {
			Degrees float64 `json:"degrees"`
			Unit    string  `json:"unit"`
		}
		MaxTemperature struct {
			Degrees float64 `json:"degrees"`
			Unit    string  `json:"unit"`
		}
		MinTemperature struct {
			Degrees float64 `json:"degrees"`
			Unit    string  `json:"unit"`
		}
		Qpf struct {
			Quantity float64 `json:"quantity"`
			Unit     string  `json:"unit"`
		}
	}
}

type GeoData struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

type WeartherApiData struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int64   `json:"dt"`
		Sunrise    int64   `json:"sunrise"`
		Sunset     int64   `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   float64 `json:"pressure"`
		Humidity   float64 `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Clouds     float64 `json:"clouds"`
		Uvi        float64 `json:"uvi"`
		Visibility float64 `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindGust   float64 `json:"wind_gust"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	}
}

type FormFields struct {
	City    string
	Country string
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Get city name and country name from command line argument
	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: go run main.go <city_name> <country_name>")
	// 	return
	// }
	// city := os.Args[1]
	// country := os.Args[2]

	r := gin.Default()
	r.GET("/weather", getWeatherApi)
	r.Run(":8080")
}

func getWeatherApi(c *gin.Context) {
	city := c.Query("city")
	country := c.Query("country")

	if city == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing city or country"})
		return
	}
	getWeather(city, country, c)
}

func getWeather(city string, country string, c *gin.Context) {
	// Get API key from environment variable
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	if weatherApiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WEATHER_API_KEY environment variable is not set"})
		return
	}

	//get lat and lon from city name using Google Geocoding API
	lat, lon, err := getLetLong(city, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error making request to Weather API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Weather API returned status code " + resp.Status})
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response from Weather API"})
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
	var data WeartherApiData
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing JSON from OpenWeather API"})
		return
	}
	parseResponseFromOpenWeather(data, c, city, country)
}

func getLetLong(city string, country string) (float64, float64, error) {
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
	var geoData GeoData
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

func parseResponseFromGoogle(data GoogleWeatherData, c *gin.Context, city string, country string) {
	// Print weather information
	c.JSON(http.StatusOK, gin.H{
		"city":        city,
		"country":     country,
		"temperature": data.Temperature.Degrees,
		"condition":   data.WeatherCondition.Description.Text,
		"wind":        data.Wind.Speed.Value,
		"humidity":    data.RelativeHumidity,
	})
}

func parseResponseFromOpenWeather(data WeartherApiData, c *gin.Context, city string, country string) {
	// Print weather information
	c.JSON(http.StatusOK, gin.H{
		"city":        city,
		"country":     country,
		"temperature": data.Current.Temp,
		"condition":   data.Current.Weather[0].Main,
		"wind":        data.Current.WindSpeed,
		"humidity":    data.Current.Humidity,
		"data":        data,
	})
}
