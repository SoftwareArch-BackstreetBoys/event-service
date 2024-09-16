package services

import (
	context "context"
	"log"
	"server/configs"
	"server/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type eventServiceServer struct {
}

func NewEventServiceServer() EventServiceServer {
	return eventServiceServer{}
}

// Initialize MongoDB collections
var eventCollection *mongo.Collection = configs.GetCollection(configs.DB, "events")
var eventParticipationCollection *mongo.Collection = configs.GetCollection(configs.DB, "event_participation")

func (eventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// CreateEvent inserts a new event into MongoDB
func (eventServiceServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
	// Set the current time for created_at and updated_at fields
	currentTime := time.Now()

	event := models.MongoEvent{
		Title:            req.Title,
		Description:      req.Description,
		Datetime:         req.Datetime,
		Location:         req.Location,
		MaxParticipation: req.MaxParticipation,
		CurParticipation: 0,
		ClubId:           req.ClubId,
		CreatedBy:        req.CreatedBy,
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
	}

	// Insert the event into the MongoDB collection
	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		log.Println("Failed to create event:", err)
		return nil, err
	}

	// Convert the inserted ID to string
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	return &CreateEventResponse{Id: insertedID}, nil
}

// GetEvent retrieves an event from MongoDB by its _id
func (s eventServiceServer) GetEvent(ctx context.Context, req *GetEventRequest) (*GetEventResponse, error) {
	// Convert event ID from string to ObjectID
	eventID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Find the event in the MongoDB collection
	var event models.MongoEvent
	err = eventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	// Convert created_at and updated_at to protobuf Timestamp
	createdAtProto := timestamppb.New(event.CreatedAt)
	updatedAtProto := timestamppb.New(event.UpdatedAt)

	// Return the event in the GetEventResponse
	return &GetEventResponse{
		Event: &Event{
			Id:               eventID.Hex(),
			Title:            event.Title,
			Description:      event.Description,
			Datetime:         event.Datetime,
			Location:         event.Location,
			MaxParticipation: event.MaxParticipation,
			CurParticipation: event.CurParticipation,
			ClubId:           event.ClubId,
			CreatedBy:        event.CreatedBy,
			CreatedAt:        createdAtProto,
			UpdatedAt:        updatedAtProto,
		},
	}, nil
}

// UpdateEvent updates an event's information in MongoDB and returns the updated event
func (eventServiceServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {
	// Convert event ID from string to ObjectID
	eventID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Define the update
	update := bson.M{
		"$set": bson.M{
			"title":             req.Title,
			"description":       req.Description,
			"datetime":          req.Datetime,
			"location":          req.Location,
			"max_participation": req.MaxParticipation,
			"updated_at":        time.Now(), // Update the timestamp
		},
	}

	// Update the event in the MongoDB collection
	_, err = eventCollection.UpdateOne(ctx, bson.M{"_id": eventID}, update)
	if err != nil {
		log.Println("Failed to update event:", err)
		return nil, err
	}

	// Retrieve the updated event
	var updatedEvent models.MongoEvent
	err = eventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&updatedEvent)
	if err != nil {
		return nil, err
	}

	// Return the updated event in the UpdateEventResponse
	return &UpdateEventResponse{
		Event: &Event{
			Id:               updatedEvent.Id.Hex(),
			Title:            updatedEvent.Title,
			Description:      updatedEvent.Description,
			Datetime:         updatedEvent.Datetime,
			Location:         updatedEvent.Location,
			MaxParticipation: updatedEvent.MaxParticipation,
			CurParticipation: updatedEvent.CurParticipation,
			ClubId:           updatedEvent.ClubId,
			CreatedBy:        updatedEvent.CreatedBy,
			CreatedAt:        timestamppb.New(updatedEvent.CreatedAt), // Convert time.Time to google.protobuf.Timestamp
			UpdatedAt:        timestamppb.New(updatedEvent.UpdatedAt), // Convert time.Time to google.protobuf.Timestamp
		},
	}, nil
}

// DeleteEvent removes an event from MongoDB by its ID
func (eventServiceServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {
	// Convert event ID from string to ObjectID
	eventID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Delete the event from the MongoDB collection
	_, err = eventCollection.DeleteOne(ctx, bson.M{"_id": eventID})
	if err != nil {
		return &DeleteEventResponse{Success: false}, err
	}

	// Delete all related event participations from the eventParticipationCollection
	_, err = eventParticipationCollection.DeleteMany(ctx, bson.M{"event_id": eventID})
	if err != nil {
		return &DeleteEventResponse{Success: false}, err
	}

	return &DeleteEventResponse{Success: true}, nil
}

// GetAllEvents retrieves all events from MongoDB
func (eventServiceServer) GetAllEvents(ctx context.Context, req *GetAllEventsRequest) (*GetAllEventsResponse, error) {
	// Find all events in the collection
	cur, err := eventCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Parse the results into a slice of Event
	var events []*Event
	for cur.Next(ctx) {
		var event models.MongoEvent
		if err := cur.Decode(&event); err != nil {
			return nil, err
		}

		events = append(events, &Event{
			Id:               event.Id.Hex(),
			Title:            event.Title,
			Description:      event.Description,
			Datetime:         event.Datetime,
			Location:         event.Location,
			MaxParticipation: event.MaxParticipation,
			CurParticipation: event.CurParticipation,
			ClubId:           event.ClubId,
			CreatedBy:        event.CreatedBy,
		})
	}

	return &GetAllEventsResponse{Events: events}, nil
}

// JoinEvent adds a user to an event and increments participation count
func (eventServiceServer) JoinEvent(ctx context.Context, req *JoinEventRequest) (*JoinEventResponse, error) {
	// Convert string IDs to ObjectIDs
	eventID, err := primitive.ObjectIDFromHex(req.EventId)
	if err != nil {
		return &JoinEventResponse{Success: false}, err
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &JoinEventResponse{Success: false}, err
	}

	// Check if the event exists and can accommodate more participants
	var event models.MongoEvent
	err = eventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
	if err != nil {
		return &JoinEventResponse{Success: false}, err
	}

	if event.CurParticipation >= event.MaxParticipation {
		return &JoinEventResponse{Success: false}, nil // Event is full
	}

	// Check if the user is already participating
	var participation models.MongoEventParticipation
	err = eventParticipationCollection.FindOne(ctx, bson.M{"event_id": eventID, "user_id": userID}).Decode(&participation)
	if err == nil {
		return &JoinEventResponse{Success: false}, nil // User already joined
	}

	// Add user to event participation
	participation = models.MongoEventParticipation{
		EventId: eventID,
		UserId:  userID,
	}
	_, err = eventParticipationCollection.InsertOne(ctx, participation)
	if err != nil {
		return &JoinEventResponse{Success: false}, err
	}

	// Update the event's current participation count
	filter := bson.M{"_id": eventID}
	update := bson.M{"$inc": bson.M{"cur_participation": 1}}
	_, err = eventCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &JoinEventResponse{Success: false}, err
	}

	return &JoinEventResponse{Success: true}, nil
}

func (eventServiceServer) LeaveEvent(ctx context.Context, req *LeaveEventRequest) (*LeaveEventResponse, error) {
	// Convert string IDs to ObjectIDs
	eventID, err := primitive.ObjectIDFromHex(req.EventId)
	if err != nil {
		return &LeaveEventResponse{Success: false}, err
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &LeaveEventResponse{Success: false}, err
	}

	// Check if the user is participating
	var participation models.MongoEventParticipation
	err = eventParticipationCollection.FindOne(ctx, bson.M{"event_id": eventID, "user_id": userID}).Decode(&participation)
	if err != nil {
		return &LeaveEventResponse{Success: false}, nil // User not found
	}

	// Remove user from event participation
	_, err = eventParticipationCollection.DeleteOne(ctx, bson.M{"event_id": eventID, "user_id": userID})
	if err != nil {
		return &LeaveEventResponse{Success: false}, err
	}

	// Decrease the event's current participation count
	filter := bson.M{"_id": eventID}
	update := bson.M{"$inc": bson.M{"cur_participation": -1}}
	_, err = eventCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &LeaveEventResponse{Success: false}, err
	}

	return &LeaveEventResponse{Success: true}, nil
}
