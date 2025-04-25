package mqtt

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/entity"
	"FleetManagementSystem/internal/infrastructure/config"
	"FleetManagementSystem/internal/infrastructure/rabbitmq"
	"FleetManagementSystem/pkg/utils"
	"context"
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func NewClient(cfg *config.Config, log *logrus.Logger, clientID string) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://127.0.0.1:1883").
		AddBroker("tcp://[::1]:1883")
	opts.SetClientID(clientID)
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

func SubscribeLocation(client mqtt.Client, log *logrus.Logger, service location.Service, cfg config.Config) {
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

		// Location 1
		loc1 := entity.Location{
			Latitude:     -6.22744,
			Longitude:    106.99696,
			LocationName: "La Terazza",
		}

		loc2 := entity.Location{
			Latitude:     -6.22801,
			Longitude:    107.00148,
			LocationName: "Summarecon Mal Bekasi",
		}

		for _, point := range []entity.Location{loc1, loc2} {
			if utils.IsInsideGeofence(loc.Latitude, loc.Longitude, point.Latitude, point.Longitude, cfg.GeofenceRadius) {
				alerts := map[string]interface{}{
					"vehicle_id":    loc.VehicleID,
					"message":       "Vehicle entered geofence",
					"latitude":      loc.Latitude,
					"longitude":     loc.Longitude,
					"timestamp":     loc.Timestamp,
					"location_name": point.LocationName,
				}
				log.Infof("[INFO][%s] Publishing RMQ to geofence_alerts", funcName)
				if err := rabbitmq.PublishRMQ(alerts, "geofence_alerts", "fleet.events"); err != nil {
					log.Errorf("Failed to publish to RabbitMQ: %v", err)
				}
			}
		}

	})

	log.Infof("[INFO][%s] Subscribed to topic: %s", funcName, topic)
}
