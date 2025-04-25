package entity

type Location struct {
	VehicleID    string  `json:"vehicle_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	LocationName string  `json:"location_name,omitempty"`
	Timestamp    int64   `json:"timestamp"`
}

type LocationHistoryPayload struct {
	VehicleID string `json:"vehicle_id"`
	StartTime int64
	EndTime   int64
}
