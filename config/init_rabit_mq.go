package config

import (
	"github.com/yakob-abada/go-rest-user/pkg/publisher"
	"os"
)

func InitRabbitMQ() (*publisher.AMQPPublisher, error) {
	url := os.Getenv("RABBITMQ_URL")
	return publisher.NewAMQPPublisher(url, "user.events")
}
