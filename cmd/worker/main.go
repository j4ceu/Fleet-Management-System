package main

import (
	"FleetManagementSystem/internal/infrastructure/config"
	"FleetManagementSystem/internal/infrastructure/logger"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Env)
	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		log.Fatalf("[ERROR][GeofenceAlertWorker] Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[ERROR][GeofenceAlertWorker] Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"geofence_alerts",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[ERROR][GeofenceAlertWorker] Failed to register a consumer: %v", err)
	}

	fmt.Println("[INFO][GeofenceAlertWorker] Waiting for geofence alerts...")
	for msg := range msgs {
		var alerts map[string]interface{}
		if err := json.Unmarshal(msg.Body, &alerts); err != nil {
			log.Errorln("[ERROR][GeofenceAlertWorker] Failed to Unmarshal body: ", err)
			continue
		}

		fmt.Printf("[INFO][GeofenceAlertWorker] 🚨 Fleet with police number %s is entering the %s area\n", alerts["vehicle_id"], alerts["location_name"])
	}
}
