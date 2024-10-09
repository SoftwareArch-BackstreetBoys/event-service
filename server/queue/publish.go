package queue

import (
	"encoding/json"
	"server/configs"
	"server/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Helper function to send a message to RabbitMQ
func SendMessage(notification *models.NotificationMessage) error {
	// Use connectRabbitMQ to get the connection and channel
	ch, err := configs.ConnectRabbitMQ()
	if err != nil {
			return err
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
			"event_notifications", // Queue name
			true,                  // Durable
			false,                 // Delete when unused
			false,                 // Exclusive
			false,                 // No-wait
			nil,                   // Arguments
	)
	if err != nil {
			return err
	}

	// Marshal the notification message to JSON
	messageBody, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	// Publish the message to the queue
	err = ch.Publish(
			"",     // Exchange
			q.Name, // Routing key
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
					ContentType: "application/json",
					Body:        messageBody,
			})
	return err
}
