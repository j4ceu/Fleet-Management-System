package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

func PublishRMQ(data any, routeKey, exchange string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Ch.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
