package configs

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Function to connect to RabbitMQ and return the connection and channel
func ConnectRabbitMQ() (*amqp.Channel, error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@shared-rabbitmq:5672/")
	if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
			return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
			conn.Close() // Close the connection if channel creation fails
			return nil, err
	}

	return ch, nil
}
