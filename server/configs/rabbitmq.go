package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Function to connect to RabbitMQ and return the connection and channel
func ConnectRabbitMQ() (*amqp.Channel, error) {
	// Connect to RabbitMQ server
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in server/rabbitmq")
	}
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_CONNECTION"))
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %s", err)
		conn.Close() // Close the connection if channel creation fails
		return nil, err
	}

	return ch, nil
}
