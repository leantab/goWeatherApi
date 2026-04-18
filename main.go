package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var formFields FormFields
		formFields.City = r.URL.Query().Get("city")
		formFields.Country = r.URL.Query().Get("country")

		if formFields.City == "" || formFields.Country == "" {
			http.Error(w, "Missing city or country", http.StatusBadRequest)
			return
		}
		getWeather(formFields.City, formFields.Country, w)
	})

	http.ListenAndServe(":8080", nil)
}

func getWeather(city string, country string, w http.ResponseWriter) {
	// Get API key from environment variable
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	if weatherApiKey == "" {
		http.Error(w, "Error: WEATHER_API_KEY environment variable is not set", http.StatusInternalServerError)
		return
	}

	//get lat and lon from city name using Google Geocoding API
	lat, lon := getLetLong(city, country, w)
	if lat == 0 || lon == 0 {
		http.Error(w, "Error: No results found for the city", http.StatusInternalServerError)
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
		http.Error(w, "Error making request to Google Maps Weather API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		http.Error(w, "Error: Google Maps Weather API returned status code "+resp.Status, http.StatusInternalServerError)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response from Google Maps Weather API", http.StatusInternalServerError)
		return
	}

	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Weather in %s, %s:\n", city, country)))

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
		http.Error(w, "Error parsing JSON from OpenWeather API", http.StatusInternalServerError)
		return
	}
	parseResponseFromOpenWeather(data, w)
}

func getLetLong(city string, country string, w http.ResponseWriter) (float64, float64) {
	//get lat and lon from city name using Google Geocoding API
	googleMapsApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsApiKey == "" {
		http.Error(w, "Error: GOOGLE_MAPS_API_KEY environment variable is not set", http.StatusInternalServerError)
		return 0, 0
	}
	geoUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,%s&key=%s", city, country, googleMapsApiKey)
	geoResp, err := http.Get(geoUrl)
	if err != nil {
		http.Error(w, "Error making request to Google Maps API", http.StatusInternalServerError)
		return 0, 0
	}
	defer geoResp.Body.Close()
	if geoResp.StatusCode != 200 {
		http.Error(w, "Error: Google Maps API returned status code "+geoResp.Status, http.StatusInternalServerError)
		return 0, 0
	}
	geoBody, err := io.ReadAll(geoResp.Body)
	if err != nil {
		http.Error(w, "Error reading response from Google Maps API", http.StatusInternalServerError)
		return 0, 0
	}
	var geoData GeoData
	if err := json.Unmarshal(geoBody, &geoData); err != nil {
		http.Error(w, "Error parsing JSON from Google Maps API", http.StatusInternalServerError)
		return 0, 0
	}

	if len(geoData.Results) == 0 {
		http.Error(w, "Error: No results found for the city", http.StatusInternalServerError)
		return 0, 0
	}

	lat := geoData.Results[0].Geometry.Location.Lat
	lon := geoData.Results[0].Geometry.Location.Lng

	return lat, lon
}

func parseResponseFromGoogle(data GoogleWeatherData, w http.ResponseWriter) {
	// Print weather information
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Temperature: %.1f°C\n", data.Temperature.Degrees)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Condition: %s\n", data.WeatherCondition.Description.Text)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Wind: %.1f km/h\n", data.Wind.Speed.Value)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Humidity: %.1f%%\n", data.RelativeHumidity)))
}

func parseResponseFromOpenWeather(data WeartherApiData, w http.ResponseWriter) {
	// Print weather information
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Temperature: %.1f°C\n", data.Current.Temp)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Condition: %s\n", data.Current.Weather[0].Main)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Wind: %.1f km/h\n", data.Current.WindSpeed)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Humidity: %.1f%%\n", data.Current.Humidity)))
	http.ResponseWriter.Write(w, []byte(fmt.Sprintf("Data: %v\n", data)))
}
