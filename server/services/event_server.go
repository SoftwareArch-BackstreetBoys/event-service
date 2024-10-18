package services

import (
	context "context"
	"log"
	"server/configs"
	"server/models"
	"server/queue"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetAllEvents retrieves all events from MongoDB, sorted by created_at in descending order
func (eventServiceServer) GetAllEvents(ctx context.Context, req *GetAllEventsRequest) (*GetAllEventsResponse, error) {
	// Define sorting options to sort by created_at in descending order
	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	// Find all events in the collection with sorting
	cur, err := eventCollection.Find(ctx, bson.M{}, findOptions)
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
			CreatedAt:        timestamppb.New(event.CreatedAt), // Convert time.Time to google.protobuf.Timestamp
			UpdatedAt:        timestamppb.New(event.UpdatedAt), // Convert time.Time to google.protobuf.Timestamp
		})
	}

	// Check for errors that occurred during iteration
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &GetAllEventsResponse{Events: events}, nil
}

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
			CreatedAt:        timestamppb.New(event.CreatedAt),
			UpdatedAt:        timestamppb.New(event.UpdatedAt),
		},
	}, nil
}

// GetAllEventsByUser retrieves all events created by a specific user
func (eventServiceServer) GetAllEventsByUser(ctx context.Context, req *GetAllEventsByUserRequest) (*GetAllEventsByUserResponse, error) {
	// Define filter to find events created by the specified user
	filter := bson.M{"created_by": req.UserId}

	// Define sorting options to sort by created_at in descending order
	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	// Find all events in the collection for the specified user with sorting
	cur, err := eventCollection.Find(ctx, filter, findOptions)
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
			CreatedAt:        timestamppb.New(event.CreatedAt),
			UpdatedAt:        timestamppb.New(event.UpdatedAt),
		})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &GetAllEventsByUserResponse{Events: events}, nil
}

// GetAllEventsByClub retrieves all events associated with a specific club
func (eventServiceServer) GetAllEventsByClub(ctx context.Context, req *GetAllEventsByClubRequest) (*GetAllEventsByClubResponse, error) {
	// Define filter to find events associated with the specified club
	filter := bson.M{"club_id": req.ClubId}

	// Define sorting options to sort by created_at in descending order
	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	// Find all events in the collection for the specified club with sorting
	cur, err := eventCollection.Find(ctx, filter, findOptions)
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
			CreatedAt:        timestamppb.New(event.CreatedAt),
			UpdatedAt:        timestamppb.New(event.UpdatedAt),
		})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &GetAllEventsByClubResponse{Events: events}, nil
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

	// After successfully updating the event, send a notification
	notification := models.NotificationMessage{
		NotificationType: "event_update",
    Sender:           "soeisoftarch@gmail.com",
    Receiver:         "jcypher2121@gmail.com",
    Subject:          "Event Update",
    BodyMessage:      "The event details have been updated.",
    Status:           "pending",
	}

	err = queue.SendMessage(&notification)
    if err != nil {
        log.Println("Failed to publish message to RabbitMQ:", err)
        // Handle error if needed
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

		// After successfully deleting the event, send a notification
	notification := models.NotificationMessage{
		NotificationType: "event_delete",
		Sender:           "soeisoftarch@gmail.com",
		Receiver:         "jcypher2121@gmail.com",
		Subject:          "Event Delete",
		BodyMessage:      "The event details have been updated.",
		Status:           "pending",
	}

	err = queue.SendMessage(&notification)
	if err != nil {
		log.Println("Failed to publish message to RabbitMQ:", err)
		// Handle error if needed
	}

	return &DeleteEventResponse{Success: true}, nil
}

// GetAllParticipatedEvents retrieves events where the user is a participant
func (eventServiceServer) GetAllParticipatedEvents(ctx context.Context, req *GetAllParticipatedEventsRequest) (*GetAllParticipatedEventsResponse, error) {
	// Convert user ID from string to ObjectID
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	// Find all events where the user is a participant
	cur, err := eventParticipationCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Extract event IDs from the participation documents
	var eventIDs []primitive.ObjectID
	for cur.Next(ctx) {
		var participation models.MongoEventParticipation
		if err := cur.Decode(&participation); err != nil {
			return nil, err
		}
		eventIDs = append(eventIDs, participation.EventId)
	}

	// Find the details of the events
	eventCur, err := eventCollection.Find(ctx, bson.M{"_id": bson.M{"$in": eventIDs}})
	if err != nil {
		return nil, err
	}
	defer eventCur.Close(ctx)

	// Parse the results into a slice of Event
	var events []*Event
	for eventCur.Next(ctx) {
		var event models.MongoEvent
		if err := eventCur.Decode(&event); err != nil {
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
			CreatedAt:        timestamppb.New(event.CreatedAt), // Convert time.Time to google.protobuf.Timestamp
			UpdatedAt:        timestamppb.New(event.UpdatedAt), // Convert time.Time to google.protobuf.Timestamp
		})
	}

	return &GetAllParticipatedEventsResponse{Events: events}, nil
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

		// After successfully joining the event, send a notification
	notification := models.NotificationMessage{
		NotificationType: "event_update",
		Sender:           "soeisoftarch@gmail.com",
		Receiver:         "6430386821@student.chula.ac.th",
		Subject:          "Event Join",
		BodyMessage:      "The event details have been updated.",
		Status:           "pending",
	}

	err = queue.SendMessage(&notification)
		if err != nil {
			log.Println("Failed to publish message to RabbitMQ:", err)
			// Handle error if needed
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

		// After successfully leaving the event, send a notification
	notification := models.NotificationMessage{
		NotificationType: "event_update",
		Sender:           "soeisoftarch@gmail.com",
		Receiver:         "jcypher2121@gmail.com",
		Subject:          "Event Leave",
		BodyMessage:      "The event details have been updated.",
		Status:           "pending",
	}

	err = queue.SendMessage(&notification)
	if err != nil {
		log.Println("Failed to publish message to RabbitMQ:", err)
		// Handle error if needed
	}

	return &LeaveEventResponse{Success: true}, nil
}

// SearchEvents searches for events by title and description, with an optional filter by club ID
func (eventServiceServer) SearchEvents(ctx context.Context, req *SearchEventsRequest) (*SearchEventsResponse, error) {
	// Define the search query
	searchQuery := req.SearchQuery

	// Build the search filter
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"description": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	// Apply the optional club_id filter if provided
	if req.ClubId != "" {
		filter["club_id"] = req.ClubId
	}

	// Define sorting options to sort by created_at in descending order
	findOptions := options.Find().SetSort(bson.M{"created_at": -1})

	// Find all events matching the search criteria with sorting
	cur, err := eventCollection.Find(ctx, filter, findOptions)
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
			CreatedAt:        timestamppb.New(event.CreatedAt), // Convert time.Time to google.protobuf.Timestamp
			UpdatedAt:        timestamppb.New(event.UpdatedAt), // Convert time.Time to google.protobuf.Timestamp
		})
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &SearchEventsResponse{Events: events}, nil
}
