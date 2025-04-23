package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env            string
	DbURL          string
	MQTTBroker     string
	RabbitURL      string
	GeofenceLat    float64
	GeofenceLong   float64
	GeofenceRadius float64
}

func Load() *Config {
	return &Config{
		Env:            getEnv("APP_ENV", "development"),
		DbURL:          getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/fleetdb?sslmode=disable"),
		MQTTBroker:     getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		RabbitURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		GeofenceLat:    getEnvFloat("GEOFENCE_LAT", -6.2088),
		GeofenceLong:   getEnvFloat("GEOFENCE_LONG", 106.8456),
		GeofenceRadius: getEnvFloat("GEOFENCE_RADIUS", 50),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if val, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return fallback
}
