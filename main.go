package main

import (
	"html/template"
	"log"
	"net/http"

	"city-dashboard/handlers"
	"city-dashboard/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env automatically
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Load stop names into memory
	utils.InitStops("./stops.txt") // <-- fix path to your file

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/api/weather", handlers.WeatherHandler)
	mux.HandleFunc("/api/weather/comprehensive", handlers.ComprehensiveWeatherHandler)
	mux.HandleFunc("/api/events", handlers.EventHandler)
	mux.HandleFunc("/api/crime", handlers.CrimeHandler)
	mux.HandleFunc("/api/transport", handlers.TransportHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/dashboard.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})
	log.Println("Starting new server on :7070")
	err := http.ListenAndServe(":7070", mux)
	log.Fatal(err)
}
