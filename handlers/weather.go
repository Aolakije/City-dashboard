package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"city-dashboard/models"
	"city-dashboard/utils"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get city from query params, default to Rouen
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Rouen"
	}

	// API key from environment
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API key not configured",
		})
		return
	}

	// Build API URL
	apiURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		url.QueryEscape(city),
		apiKey,
	)

	var raw models.Weather
	if err := utils.FetchJSON(apiURL, &raw); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch weather data: " + err.Error(),
		})
		return
	}

	// Build response for frontend
	response := struct {
		City        string   `json:"city"`
		Lat         float64  `json:"lat"`
		Lon         float64  `json:"lon"`
		Temperature float64  `json:"temperature"`
		FeelsLike   float64  `json:"feels_like"`
		TempMin     float64  `json:"temp_min"`
		TempMax     float64  `json:"temp_max"`
		Pressure    int      `json:"pressure"`
		Humidity    int      `json:"humidity"`
		Conditions  []string `json:"conditions"`
		WindSpeed   float64  `json:"wind_speed"`
		WindDeg     int      `json:"wind_deg"`
		Visibility  int      `json:"visibility"`
		Rain1h      float64  `json:"rain_1h"`
		Clouds      int      `json:"clouds"`
		Sunrise     int64    `json:"sunrise"`
		Sunset      int64    `json:"sunset"`
	}{
		City:        raw.Name,
		Lat:         raw.Coord.Lat,
		Lon:         raw.Coord.Lon,
		Temperature: raw.Main.Temp,
		FeelsLike:   raw.Main.FeelsLike,
		TempMin:     raw.Main.TempMin,
		TempMax:     raw.Main.TempMax,
		Pressure:    raw.Main.Pressure,
		Humidity:    raw.Main.Humidity,
		Conditions: func() []string {
			arr := []string{}
			for _, w := range raw.Weather {
				arr = append(arr, w.Description)
			}
			return arr
		}(),
		WindSpeed:  raw.Wind.Speed,
		WindDeg:    raw.Wind.Deg,
		Visibility: raw.Visibility,
		Rain1h:     raw.Rain.OneH,
		Clouds:     raw.Clouds.All,
		Sunrise:    raw.Sys.Sunrise,
		Sunset:     raw.Sys.Sunset,
	}

	json.NewEncoder(w).Encode(response)
}

// models/weather.go (you probably already have this)
type AirPollution struct {
	List []struct {
		Main struct {
			AQI int `json:"aqi"`
		} `json:"main"`
		Components struct {
			CO   float64 `json:"co"`
			NO   float64 `json:"no"`
			NO2  float64 `json:"no2"`
			O3   float64 `json:"o3"`
			SO2  float64 `json:"so2"`
			PM25 float64 `json:"pm2_5"`
			PM10 float64 `json:"pm10"`
			NH3  float64 `json:"nh3"`
		} `json:"components"`
	} `json:"list"`
}

type UVIndex struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Value float64 `json:"value"`
}

func ComprehensiveWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	city := r.URL.Query().Get("city")
	if city == "" {
		city = "Rouen"
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API key not configured",
		})
		return
	}

	// First get basic weather for coordinates
	weatherURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		url.QueryEscape(city), apiKey,
	)

	var weather models.Weather
	if err := utils.FetchJSON(weatherURL, &weather); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch weather data",
		})
		return
	}

	// Get Air Quality data
	aqiURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/air_pollution?lat=%f&lon=%f&appid=%s",
		weather.Coord.Lat, weather.Coord.Lon, apiKey,
	)

	var airPollution models.AirPollution
	utils.FetchJSON(aqiURL, &airPollution) // Don't fail if AQI unavailable

	// Get UV Index data
	uvURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/uvi?lat=%f&lon=%f&appid=%s",
		weather.Coord.Lat, weather.Coord.Lon, apiKey,
	)

	var uvData models.UVIndex
	utils.FetchJSON(uvURL, &uvData) // Don't fail if UV unavailable

	// Build comprehensive response
	response := struct {
		Weather struct {
			WindSpeed float64 `json:"wind_speed"` // km/h
		} `json:"weather"`
		AirQuality struct {
			AQI int `json:"aqi"`
		} `json:"air_quality"`
		UVIndex float64 `json:"uv_index"`
	}{
		Weather: struct {
			WindSpeed float64 `json:"wind_speed"`
		}{
			WindSpeed: weather.Wind.Speed * 3.6, // Convert m/s to km/h
		},
		UVIndex: uvData.Value,
	}

	// Convert OpenWeatherMap AQI (1-5) to standard AQI scale
	if len(airPollution.List) > 0 {
		owmAQI := airPollution.List[0].Main.AQI
		response.AirQuality.AQI = convertAQIToStandard(owmAQI)
	}

	json.NewEncoder(w).Encode(response)
}

// Helper function to convert OpenWeatherMap AQI to standard scale
func convertAQIToStandard(owmAQI int) int {
	conversion := map[int]int{
		1: 25,  // Good (1-50)
		2: 75,  // Fair (51-100)
		3: 125, // Moderate (101-150)
		4: 175, // Poor (151-200)
		5: 250, // Very Poor (201-300)
	}
	if val, ok := conversion[owmAQI]; ok {
		return val
	}
	return 42 // Default fallback
}
