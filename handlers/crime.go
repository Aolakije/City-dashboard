package handlers

import (
	"city-dashboard/models"
	"encoding/json"
	"net/http"
)

func CrimeHandler(w http.ResponseWriter, r *http.Request) {
	mockData := []models.Crime{
		{Type: "Theft", Location: "Main Street", Date: "2025-09-02", Severity: 3},
		{Type: "Assault", Location: "Central Park", Date: "2025-09-01", Severity: 4},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
}
