package handlers

import (
	"city-dashboard/models"
	"encoding/json"
	"net/http"
)

func EventHandler(w http.ResponseWriter, r *http.Request) {
	mockData := []models.Event{
		{Name: "Concert", Location: "City Hall", Date: "2025-09-05", Description: "Jazz evening"},
		{Name: "Food Festival", Location: "Market Square", Date: "2025-09-10", Description: "Local delicacies"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
}
