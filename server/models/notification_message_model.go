package models

// Define a struct for the notification message
type NotificationMessage struct {
	NotificationType string `json:"notification_type"`
	Sender           string `json:"sender"`
	Receiver         string `json:"receiver"`
	Subject          string `json:"subject"`
	BodyMessage      string `json:"body_message"`
	Status					 string `json:"status"`
}
