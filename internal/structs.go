package internal

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
