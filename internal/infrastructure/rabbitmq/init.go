package rabbitmq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var Ch *amqp.Channel
var Conn *amqp.Connection

func Init(url string, log *logrus.Logger) {
	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	Ch, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	err = Ch.ExchangeDeclare("fleet.events", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = Ch.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = Ch.QueueBind("geofence_alerts", "geofence_alerts", "fleet.events", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

}

func Close() {
	Ch.Close()
	Conn.Close()
}
