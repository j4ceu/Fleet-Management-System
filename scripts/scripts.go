package scripts

import (
	"FleetManagementSystem/internal/entity"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var vehicleIDs = []string{
	"B1234XYZ", "D5678DGV", "F9876CMP", "T1111CBC", "B2123KCM", "A5230KCM",
}

func PublishToMQTT(client mqtt.Client, log *logrus.Logger) {
	const funcName = "PublishToMQTT"
	log.Infof("[INFO][%s] Starting publisher...", funcName)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, vehicleID := range vehicleIDs {
			// 1. Prepare payload dengan epoch timestamp
			payload := entity.Location{
				VehicleID: vehicleID,
				Latitude:  -6.2088 + rand.Float64()/100,
				Longitude: 106.8456 + rand.Float64()/100,
				Timestamp: time.Now().Unix(), // Format epoch (int64)
			}

			// 2. Marshal payload ke JSON
			data, err := json.Marshal(payload)
			if err != nil {
				log.Errorf("[ERROR][%s] Failed to marshal payload for %s: %v", funcName, vehicleID, err)
				continue
			}

			// 3. Publish ke MQTT
			topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
			token := client.Publish(topic, 1, false, data)
			if token.Error() != nil {
				log.Errorf("[INFO][%s] Publish failed to %s: %v", funcName, topic, token.Error())
			} else {
				log.Infof("[INFO][%s] Published to %s: %s", funcName, topic, string(data))
			}
		}
	}
}
