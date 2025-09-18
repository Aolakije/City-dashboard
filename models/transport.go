package models

// Raw vehicle positions from GTFS API
type GTFSRaw struct {
	Entity []struct {
		ID      string `json:"id"`
		Vehicle struct {
			CurrentStatus string `json:"currentStatus"`
			Occupancy     string `json:"occupancyStatus"`
			Position      struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Bearing   int     `json:"bearing"`
			} `json:"position"`
			Timestamp int64 `json:"timestamp"`
			Trip      struct {
				TripID      string `json:"tripId"`
				RouteID     string `json:"routeId"`
				DirectionID int    `json:"directionId"`
			} `json:"trip"`
			Vehicle struct {
				ID    string `json:"id"`
				Label string `json:"label"`
			} `json:"vehicle"`
		} `json:"vehicle"`
	} `json:"entity"`
}

// Trip updates for ETA information
type GTFSUpdates struct {
	Entity []struct {
		TripUpdate TripUpdate `json:"tripUpdate"`
	} `json:"entity"`
}

// Separate TripUpdate struct for better organization
type TripUpdate struct {
	Trip struct {
		TripID  string `json:"tripId"`
		RouteID string `json:"routeId"`
	} `json:"trip"`
	StopTimeUpdate []struct {
		StopID  string `json:"stopId"`
		Arrival struct {
			Time int64 `json:"time"`
		} `json:"arrival"`
		Departure struct {
			Time int64 `json:"time"`
		} `json:"departure"`
	} `json:"stopTimeUpdate"`
}

// Next stop ETA structure for frontend
type StopETA struct {
	StopID   string `json:"stop_id"`   // Original GTFS stop ID
	StopName string `json:"stop_name"` // Human-readable stop name
	ETA      string `json:"eta"`       // Estimated time of arrival
}

// Simplified vehicle for frontend display
type VehicleGTFS struct {
	ID            string    `json:"id"`
	Label         string    `json:"label"`
	RouteID       string    `json:"route_id"`
	DirectionID   int       `json:"direction_id"`
	Lat           float64   `json:"lat"`
	Lon           float64   `json:"lon"`
	Bearing       int       `json:"bearing"`
	CurrentStatus string    `json:"current_status"`
	Occupancy     string    `json:"occupancy"`
	NextStops     []StopETA `json:"next_stops"`
}
