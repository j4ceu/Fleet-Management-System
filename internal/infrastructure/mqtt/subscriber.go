package mqtt

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/entity"
	"FleetManagementSystem/internal/infrastructure/config"
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func NewClient(cfg *config.Config, log *logrus.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://127.0.0.1:1883").
		AddBroker("tcp://[::1]:1883")
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

	return client
}

func SubscribeLocation(client mqtt.Client, log *logrus.Logger, service location.Service) {
	topic := "/fleet/vehicle/+/location"
	funcName := "SubscribeLocation"

	client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		log.Infof("[INFO][%s] Received message on topic: %s", funcName, m.Topic())
		var payload entity.Location
		if err := json.Unmarshal(m.Payload(), &payload); err != nil {
			log.Errorf("[ERROR][%s] Failed to unmarshal MQTT message: %v", funcName, err)
			return
		}
		log.Infof("[INFO][%s] Location Data: %+v", funcName, payload)

		loc := entity.Location{
			VehicleID: payload.VehicleID,
			Latitude:  payload.Latitude,
			Longitude: payload.Longitude,
			Timestamp: payload.Timestamp,
		}

		if err := service.StoreLocation(context.Background(), loc); err != nil {
			log.Errorf("[ERROR][%s] Failed to store location: %v", funcName, err)
			return
		}

		log.Infof("[INFO][%s] Location stored for vehicle: %s", funcName, loc.VehicleID)

		// TODO: Trigger geofence check and publish to RabbitMQ if necessary
	})

	log.Infof("[INFO][%s] Subscribed to topic: %s", funcName, topic)
}
