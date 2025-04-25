package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	DbURL          string
	MQTTBroker     string
	RabbitURL      string
	GeofenceRadius float64
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)

	}
	return &Config{
		Env:            getEnv("APP_ENV", "development"),
		DbURL:          getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/fleetdb?sslmode=disable"),
		MQTTBroker:     getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		RabbitURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		GeofenceRadius: getEnvFloat("GEOFENCE_RADIUS", 50),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvFloat(key string, fallback float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fallback
	}

	return f
}
