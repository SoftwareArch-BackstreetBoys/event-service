package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event represents the event data stored in MongoDB
type MongoEvent struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`     // MongoDB ObjectID
	Title            string             `bson:"title"`             // Title of the event
	Description      string             `bson:"description"`       // Description of the event
	Datetime         string             `bson:"datetime"`          // Event date and time
	Location         string             `bson:"location"`          // Event location
	MaxParticipation int64              `bson:"max_participation"` // Maximum number of participants
	CurParticipation int64              `bson:"cur_participation"` // Current number of participants
	ClubId           string             `bson:"club_id"`           // Club ID associated with the event
	CreatedById      string             `bson:"created_by_id"`     // User ID of the event creator
	CreatedByName    string             `bson:"created_by_name"`   // User Name of the event creator
	CreatedAt        time.Time          `bson:"created_at"`        // Timestamp when the event was created
	UpdatedAt        time.Time          `bson:"updated_at"`        // Timestamp when the event was last updated
}
