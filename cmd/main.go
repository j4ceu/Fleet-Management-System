package main

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/controller"
	"FleetManagementSystem/internal/infrastructure/config"
	"FleetManagementSystem/internal/infrastructure/db"
	"FleetManagementSystem/internal/infrastructure/logger"
	"FleetManagementSystem/internal/infrastructure/mqtt"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Env)
	dbPool, err := db.NewPgxPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer dbPool.Close()

	mqttClient := mqtt.NewClient(cfg, log)
	defer mqttClient.Disconnect(250)

	locationRepo := location.NewRepository(dbPool)
	locationService := location.NewService(locationRepo, log)
	locationController := controller.NewLocationController(locationService)

	go mqtt.SubscribeLocation(mqttClient, log, locationService)
	// go scripts.PublishToMQTT(mqttClient, log)

	r := config.SetupRouter(locationController)
	r.Run()

}
