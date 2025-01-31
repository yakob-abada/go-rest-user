package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// AMQPPublisher implements Publisher using RabbitMQ
type AMQPPublisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
}

// NewAMQPPublisher creates a new instance of AMQPPublisher
func NewAMQPPublisher(amqpURL, exchange string) (*AMQPPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange if it does not exist
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &AMQPPublisher{
		Connection: conn,
		Channel:    ch,
		Exchange:   exchange,
	}, nil
}

// Publish sends a message to the RabbitMQ exchange
func (p *AMQPPublisher) Publish(event EventType, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.Channel.Publish(
		p.Exchange,
		fmt.Sprintf("%s", event),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close closes the RabbitMQ connection
func (p *AMQPPublisher) Close() {
	p.Channel.Close()
	p.Connection.Close()
}
