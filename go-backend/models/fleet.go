package models

type Engine struct {
	RPM       int `json:"rpm"`
	Temp      int `json:"temp"`
	FuelLevel int `json:"fuel_level"`
}

type Location struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Speed float64 `json:"speed"`
}

type Diagnostics struct {
	Alerts       []string `json:"alerts"`
	SystemHealth string   `json:"system_health"`
}

type SensorData struct {
	Engine      Engine      `json:"engine"`
	Location    Location    `json:"location"`
	Diagnostics Diagnostics `json:"diagnostics"`
}

// Level 3: Trip Telemetry
type TripTelemetry struct {
	ID         string     `json:"id"`
	TripID     string     `json:"trip_id"`
	RecordedAt string     `json:"recorded_at"`
	SensorData SensorData `json:"sensor_data"`
}

// Level 2: Trips
type Trip struct {
	ID             string          `json:"id"`
	FleetID        string          `json:"fleet_id"`
	Origin         string          `json:"origin"`
	Destination    string          `json:"destination"`
	DepartureTime  string          `json:"departure_time"`
	ArrivalTime    string          `json:"arrival_time"`
	PassengerCount int             `json:"passenger_count"`
	CreatedAt      string          `json:"created_at"`
	TripTelemetry  []TripTelemetry `json:"trip_telemetry"`
}

// Level 1: Fleet
type Fleet struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	LicensePlate string `json:"license_plate"`
	VehicleType  string `json:"vehicle_type"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	Trips        []Trip `json:"trips"`
}
