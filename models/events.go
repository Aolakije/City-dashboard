package models

type Event struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Date        string `json:"date"`
	Description string `json:"description"`
}
