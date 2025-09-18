// utils/load_stops.go
package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadStopsCSV reads a stops.txt-like file and returns a map[stopID]stopName
func LoadStopsCSV(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stopsMap := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, ",")
		if len(cols) < 3 {
			continue // skip invalid line
		}
		stopID := strings.TrimSpace(cols[0])
		stopName := strings.Trim(cols[2], `"`) // remove quotes
		stopsMap[stopID] = stopName
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopsMap, nil
}

// StopNames holds the global mapping
var StopNames map[string]string

// InitStops loads the stop names into the global map
func InitStops(filePath string) {
	m, err := LoadStopsCSV(filePath)
	if err != nil {
		log.Println("Warning: could not load stops file:", err)
		StopNames = make(map[string]string)
		return
	}
	StopNames = m
	log.Printf("Loaded %d stops\n", len(StopNames))
}
