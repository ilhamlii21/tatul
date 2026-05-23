package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"method/models"
)

type Config struct {
	SupabaseURL     string `json:"SUPABASE_URL"`
	SupabaseAnonKey string `json:"SUPABASE_ANON_KEY"`
}

// LoadConfig loads credentials from config.json at root or local folder
func LoadConfig() (Config, error) {
	var cfg Config
	file, err := os.Open("../config.json")
	if err != nil {
		file, err = os.Open("config.json")
		if err != nil {
			return cfg, err
		}
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
}

// FetchAndParseData fetches the nested data from Supabase and parses it, measuring both network latency and parsing latency.
func FetchAndParseData(client *http.Client, baseURL string, anonKey string) ([]models.Fleet, time.Duration, time.Duration, error) {
	// Querying resource embedding: fleet with its nested tripsand trip_telemetry
	url := baseURL + "fleet?select=*,trips(*,trip_telemetry(*))"

	reqStart := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, 0, err
	}

	req.Header.Set("apikey", anonKey)
	req.Header.Set("Authorization", "Bearer "+anonKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, 0, err
	}
	networkDuration := time.Since(reqStart)

	// Measure parsing latency specifically
	parseStart := time.Now()
	var fleets []models.Fleet
	err = json.Unmarshal(body, &fleets)
	if err != nil {
		return nil, 0, 0, err
	}
	parseDuration := time.Since(parseStart)

	return fleets, networkDuration, parseDuration, nil
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Printf("Gagal memuat config.json: %v\n", err)
		return
	}

	client := &http.Client{}
	fleets, netTime, parseTime, err := FetchAndParseData(client, cfg.SupabaseURL, cfg.SupabaseAnonKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Koneksi sukses!\n")
	fmt.Printf("Jumlah Fleet Terbaca: %d\n", len(fleets))
	fmt.Printf("Network Latency: %v\n", netTime)
	fmt.Printf("Parsing Latency: %v\n", parseTime)

	if len(fleets) > 0 {
		f := fleets[0]
		fmt.Printf("\nContoh Data Level 1 (Fleet): %s (Status: %s)\n", f.Name, f.Status)
		if len(f.Trips) > 0 {
			t := f.Trips[0]
			fmt.Printf("  └─ Level 2 (Trip): Route %s -> %s\n", t.Origin, t.Destination)
			if len(t.TripTelemetry) > 0 {
				tel := t.TripTelemetry[0]
				fmt.Printf("      └─ Level 3 (Telemetry): Lat %f, Lng %f, Speed %f km/h\n", 
					tel.SensorData.Location.Lat, 
					tel.SensorData.Location.Lng, 
					tel.SensorData.Location.Speed,
				)
			}
		}
	}
}
