package rabbitmq

import (
	"FleetManagementSystem/internal/api/location"
	"FleetManagementSystem/internal/entity"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

func ConsumeToSaveDB(log *logrus.Logger, service location.Service) {
	msgs, err := Ch.Consume(
		"save_db",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[ERROR][SaveDBWorker] Failed to register a consumer: %v", err)
	}

	fmt.Println("[INFO][SaveDBWorker] Waiting for saving data location to database...")
	for msg := range msgs {
		var loc entity.Location
		if err := json.Unmarshal(msg.Body, &loc); err != nil {
			log.Errorln("[ERROR][SaveDBWorker] Failed to Unmarshal body: ", err)
			continue
		}

		if err := service.StoreLocation(context.Background(), loc); err != nil {
			log.Errorf("[ERROR][SaveDBWorker] Failed to store location: %v", err)
			return
		}
	}
}
