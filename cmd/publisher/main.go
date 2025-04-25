package main

import (
	"FleetManagementSystem/internal/infrastructure/config"
	"FleetManagementSystem/internal/infrastructure/logger"

	"FleetManagementSystem/internal/infrastructure/mqtt"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Env)

	mqttClient := mqtt.NewClient(cfg, log, "fleet-publisher-mqtt")
	defer mqttClient.Disconnect(250)

	mqtt.PublishToMQTT(mqttClient, log, "route.csv")
}
