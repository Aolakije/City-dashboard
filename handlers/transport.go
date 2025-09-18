package handlers

import (
	"city-dashboard/models"
	"city-dashboard/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func TransportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	labelFilter := strings.ToLower(r.URL.Query().Get("label"))
	// Fetch vehicle positions
	vehicleURL := "https://gtfs.bus-tracker.fr/gtfs-rt/tcar/vehicle-positions.json"
	var raw models.GTFSRaw
	if err := utils.FetchJSON(vehicleURL, &raw); err != nil {
		http.Error(w, "Failed to fetch vehicle positions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch trip updates
	updateURL := "https://gtfs.bus-tracker.fr/gtfs-rt/tcar/trip-updates.json"
	var updates models.GTFSUpdates
	_ = utils.FetchJSON(updateURL, &updates) // ignore error if fails

	vehicles := []models.VehicleGTFS{}
	for _, e := range raw.Entity {
		v := e.Vehicle

		if labelFilter == "" || !strings.Contains(strings.ToLower(v.Vehicle.Label), labelFilter) {
			continue
		}

		vehicle := models.VehicleGTFS{
			ID:            v.Vehicle.ID,
			Label:         v.Vehicle.Label,
			RouteID:       v.Trip.RouteID,
			DirectionID:   v.Trip.DirectionID,
			Lat:           v.Position.Latitude,
			Lon:           v.Position.Longitude,
			Bearing:       v.Position.Bearing,
			CurrentStatus: v.CurrentStatus,
			Occupancy:     v.Occupancy,
			NextStops:     []models.StopETA{},
		}

		// Attach next stops with names
		for _, u := range updates.Entity {
			if u.TripUpdate.Trip.TripID == v.Trip.TripID {
				for _, s := range u.TripUpdate.StopTimeUpdate {
					eta := time.Unix(s.Arrival.Time, 0).Format("15:04")

					// Replace StopID with human-readable name if available
					stopName := utils.StopNames[s.StopID]
					if stopName == "" {
						stopName = s.StopID
					}

					vehicle.NextStops = append(vehicle.NextStops, models.StopETA{
						StopID: stopName,
						ETA:    eta,
					})
				}
			}
		}

		vehicles = append(vehicles, vehicle)
	}

	if err := json.NewEncoder(w).Encode(vehicles); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
