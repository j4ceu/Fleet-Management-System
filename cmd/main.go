package main

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/controller"
	"FleetManagementSystem/internal/infrastructure/config"
	"FleetManagementSystem/internal/infrastructure/db"
	"FleetManagementSystem/internal/infrastructure/logger"
	"FleetManagementSystem/internal/infrastructure/mqtt"
	"FleetManagementSystem/internal/infrastructure/rabbitmq"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Env)
	dbPool, err := db.NewPgxPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbPool.Close()

	mqttClient := mqtt.NewClient(cfg, log, "fleet-subscriber-mqtt")
	defer mqttClient.Disconnect(250)

	rabbitmq.Init(cfg.RabbitURL, log)

	locationRepo := location.NewRepository(dbPool)
	locationService := location.NewService(locationRepo, log)
	locationController := controller.NewLocationController(locationService)

	go mqtt.SubscribeLocation(mqttClient, log, locationService, *cfg)

	r := config.SetupRouter(locationController)
	r.Run()

}
