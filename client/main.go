package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"client/model"
	"client/services"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// App struct holds the gRPC client
type App struct {
	eventService services.EventService
}

// CORS middleware to handle CORS requests with specific origin and credentials support
func corsMiddleware(next http.Handler) http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in client/main/corsMiddleware")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Specify the allowed origin
		allowedOrigin := os.Getenv("FRONTEND_ROUTE")

		// Get the origin of the incoming request
		origin := r.Header.Get("Origin")

		// Check if the origin matches the allowed origin
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials (cookies, etc.)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow additional headers
		}

		// If it's a preflight request (OPTIONS), respond with a 200 status
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// HealthCheckHandler handles health check requests
func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) // Respond with a simple "OK" message
}

// GetAllEventsHandler handles fetching all events
func (app *App) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	var req model.GetAllEventsRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	res, err := app.eventService.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publicEvents := make([]*services.Event, 0)
	for _, event := range res.Events {
		if event.ClubId == "" {
			publicEvents = append(publicEvents, event)
		}
	}

	joinedClubEvents := make([]*services.Event, 0)
	for _, event := range res.Events {
		if slices.Contains(req.ClubIDs, event.ClubId) {
			joinedClubEvents = append(joinedClubEvents, event)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.GetAllEventsResponse{
		PublicEvents:     publicEvents,
		JoinedClubEvents: joinedClubEvents,
	}) // Return the response to the frontend
}

// CreateEventHandler handles the creation of a new event
func (app *App) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var req services.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// userID, err := util.GetUserIdFromRequestObject(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	res, err := app.eventService.CreateEvent(req.Title, req.Description, req.Datetime, req.Location, req.MaxParticipation, req.ClubId, req.CreatedById, req.CreatedByName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// GetEventHandler handles fetching an event by ID
func (app *App) getEventHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/event/")
	res, err := app.eventService.GetEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// GetAllEventsByUserHandler handles fetching all events by user ID
func (app *App) getAllEventsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/user/")
	userID = strings.TrimSuffix(userID, "/events")

	res, err := app.eventService.GetAllEventsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// GetAllEventsByClubHandler handles fetching all events by club ID
func (app *App) getAllEventsByClubHandler(w http.ResponseWriter, r *http.Request) {
	clubID := strings.TrimPrefix(r.URL.Path, "/club/")
	clubID = strings.TrimSuffix(clubID, "/events")
	res, err := app.eventService.GetAllEventsByClub(clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// UpdateEventHandler handles updating an existing event
func (app *App) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/event/")
	var req services.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.Id = id

	res, err := app.eventService.UpdateEvent(req.Id, req.Title, req.Description, req.Datetime, req.Location, req.MaxParticipation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// DeleteEventHandler handles deleting an event by ID
func (app *App) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/event/")
	res, err := app.eventService.DeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// GetAllParticipatedEventsHandler handles fetching all participated events for a user
func (app *App) getAllParticipatedEventsHandler(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimPrefix(r.URL.Path, "/user/")
	userID = strings.TrimSuffix(userID, "/participated-events")

	// userID, err := util.GetUserIdFromRequestObject(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	res, err := app.eventService.GetAllParticipatedEvents(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// JoinEventHandler handles joining an event
func (app *App) joinEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimPrefix(r.URL.Path, "/event/")
	eventID = strings.TrimSuffix(eventID, "/join")
	var req services.JoinEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.EventId = eventID

	// userID, err := util.GetUserIdFromRequestObject(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println("Join req: ", req.EventId, req.UserId)
	res, err := app.eventService.JoinEvent(req.EventId, req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// LeaveEventHandler handles leaving an event
func (app *App) leaveEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimPrefix(r.URL.Path, "/event/")
	eventID = strings.TrimSuffix(eventID, "/leave")
	var req services.LeaveEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.EventId = eventID

	// userID, err := util.GetUserIdFromRequestObject(r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	res, err := app.eventService.LeaveEvent(req.EventId, req.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// SearchEventsHandler handles searching for events
func (app *App) searchEventsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	clubID := r.URL.Query().Get("club_id")

	res, err := app.eventService.SearchEvents(query, clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Return the response to the frontend
}

// Main function to run the application
func main() {
	// Create a gRPC connection
	creds := insecure.NewCredentials()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in client/main/main")
	}

	cc, err := grpc.Dial(os.Getenv("GRPC_SERVER_PORT"), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer cc.Close()

	// Initialize the gRPC client
	eventClient := services.NewEventServiceClient(cc)
	eventService := services.NewEventService(eventClient)

	// _, err = eventService.CreateEvent("Test3", "this is the second test", "when", "where", 99, "club_id", "1234", "ping")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create an instance of App
	app := &App{
		eventService: eventService,
	}

	// Set up HTTP routes and wrap them with CORS middleware
	http.Handle("/health", corsMiddleware(http.HandlerFunc(app.healthCheckHandler)))         // Health check route
	http.Handle("/events", corsMiddleware(http.HandlerFunc(app.getAllEventsHandler)))        // Handler for get all events
	http.Handle("/event", corsMiddleware(http.HandlerFunc(app.createEventHandler)))          // Handler for create an event
	http.Handle("/event/", corsMiddleware(http.HandlerFunc(app.eventHandler)))               // Combine Handler for fetching/updating/deleting an event by ID and join/leave event
	http.Handle("/club/", corsMiddleware(http.HandlerFunc(app.getAllEventsByClubHandler)))   // Handler for getting all events by club ID
	http.Handle("/user/", corsMiddleware(http.HandlerFunc(app.usersHandler)))                // Combine Handler for user events and user participated-events
	http.Handle("/events/search", corsMiddleware(http.HandlerFunc(app.searchEventsHandler))) // Handler for searching events

	// Start the HTTP server
	port := ":" + os.Getenv("HTTP_PORT")
	log.Printf("HTTP server started on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// eventsHandler handles both GET, PUT, and DELETE requests for /event
func (app *App) eventHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getEventHandler(w, r) // Fetch event
	case http.MethodPut:
		app.updateEventHandler(w, r) // Update event
	case http.MethodDelete:
		app.deleteEventHandler(w, r) // Delete event
	case http.MethodPost:
		path := r.URL.Path

		if strings.HasSuffix(path, "/join") {
			app.joinEventHandler(w, r) // Join event
		} else if strings.HasSuffix(path, "/leave") {
			app.leaveEventHandler(w, r) // Leave event
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasSuffix(path, "/events") {
		app.getAllEventsByUserHandler(w, r)
	} else if strings.HasSuffix(path, "/participated-events") {
		app.getAllParticipatedEventsHandler(w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
