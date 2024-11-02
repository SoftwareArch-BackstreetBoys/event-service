package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// EventParticipation represents a user's participation in an event
type MongoEventParticipation struct {
	EventId primitive.ObjectID `bson:"event_id,omitempty"` // MongoDB ObjectID
	UserId  string             `bson:"user_id,omitempty"`  // MongoDB ObjectID
}
