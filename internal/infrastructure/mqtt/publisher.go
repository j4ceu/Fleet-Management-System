package mqtt

import (
	"FleetManagementSystem/internal/entity"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var vehicleIDs = []string{
	"B1234XYZ", "D5678DGV", "F9876CMP", "T1111CBC", "B2123KCM", "A5230KCM",
}

func loadGPSRoute(filePath string) ([][2]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var route [][2]float64

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV data: %v", err)
	}

	for _, record := range records {
		if len(record) < 1 {
			continue
		}

		latitude, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Println("Invalid latitude value:", record[0])
			continue
		}

		longitude, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Println("Invalid longitude value:", record[1])
			continue
		}

		route = append(route, [2]float64{latitude, longitude})
	}

	return route, nil
}

func addNoise(lat, lng float64) (float64, float64) {
	noise := rand.Float64() * 0.0005
	return lat + (rand.Float64()*2-1)*noise, lng + (rand.Float64()*2-1)*noise
}

func PublishToMQTT(client mqtt.Client, log *logrus.Logger, routeFile string) {
	const funcName = "PublishToMQTT"
	log.Infof("[INFO][%s] Starting publisher with GPS route from file: %s", funcName, routeFile)

	route, err := loadGPSRoute(routeFile)
	if err != nil {
		log.Fatalf("[ERROR][%s] Failed to load GPS route: %v", funcName, err)
		return
	}

	vehicleStartIndexes := make(map[string]int)
	vehicleDirections := make(map[string]int)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, vehicleID := range vehicleIDs {
			startIndex, exists := vehicleStartIndexes[vehicleID]
			if !exists {
				startIndex = rand.Intn(len(route))
				vehicleStartIndexes[vehicleID] = startIndex
				vehicleDirections[vehicleID] = 1
			}

			coords := route[startIndex]

			lat, lng := addNoise(coords[0], coords[1])

			payload := entity.Location{
				VehicleID: vehicleID,
				Latitude:  lat,
				Longitude: lng,
				Timestamp: time.Now().Unix(),
			}

			data, err := json.Marshal(payload)
			if err != nil {
				log.Errorf("[ERROR][%s] Failed to marshal payload for %s: %v", funcName, vehicleID, err)
				continue
			}

			topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)
			token := client.Publish(topic, 1, false, data)
			if token.Error() != nil {
				log.Errorf("[INFO][%s] Publish failed to %s: %v", funcName, topic, token.Error())
			} else {
				log.Infof("[INFO][%s] Published to %s: %s", funcName, topic, string(data))
			}

			if vehicleDirections[vehicleID] == 1 {
				startIndex = (startIndex + 1) % len(route) //Kendaraan Maju
			} else {
				startIndex = (startIndex - 1 + len(route)) % len(route) //kendaraan mundur
			}

			if startIndex == len(route)-1 {
				vehicleDirections[vehicleID] = -1
			} else if startIndex == 0 {
				vehicleDirections[vehicleID] = 1
			}

			vehicleStartIndexes[vehicleID] = startIndex
		}
	}
}
