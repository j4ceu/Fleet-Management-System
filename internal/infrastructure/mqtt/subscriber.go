package mqtt

import (
	"FleetManagementSystem/internal/entity"
	"FleetManagementSystem/internal/infrastructure/api/location"
	"FleetManagementSystem/internal/infrastructure/config"
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type LocationPayload struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func NewClient(cfg *config.Config, log *logrus.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(cfg.MQTTBroker)
	opts.SetClientID("fleet-subscriber")
	opts.SetConnectTimeout(5 * time.Second)
	opts.OnConnect = func(c mqtt.Client) {
		log.Info("Connected to MQTT broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Errorf("MQTT connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection failed: %v", token.Error())
	}
}

func SubscribeLocation(client mqtt.Client, log *logrus.Logger, service location.Service) {
	topic := "/fleet/vehicle/+/location"

	client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		log.Infof("Received message on topic: %s", m.Topic())
		var payload LocationPayload
		if err := json.Unmarshal(m.Payload(), &payload); err != nil {
			log.Errorf("Failed to unmarshal MQTT message: %v", err)
			return
		}
		log.Infof("Location Data: %+v", payload)

		loc := entity.Location{
			VehicleID: payload.VehicleID,
			Latitude:  payload.Latitude,
			Longitude: payload.Longitude,
			Timestamp: payload.Timestamp,
		}

		if err := service.StoreLocation(context.Background(), loc); err != nil {
			log.Errorf("Failed to store location: %v", err)
			return
		}

		log.Infof("Location stored for vehicle: %s", loc.VehicleID)

		// TODO: Trigger geofence check and publish to RabbitMQ if necessary
	})

	log.Infof("Subscribed to topic: %s", topic)
}
