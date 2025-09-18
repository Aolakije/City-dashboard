package models

type Crime struct {
	Type     string `json:"type"`
	Location string `json:"location"`
	Date     string `json:"date"`
	Severity int    `json:"severity"` // 1-5
}

