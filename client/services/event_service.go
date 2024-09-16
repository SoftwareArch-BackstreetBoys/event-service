package services

import (
	"context"
	"fmt"
)

type EventService interface {
	CreateEvent(title string, description string, datetime string, location string, max_participation int64, club_id string, created_by string) error
	GetEvent(id string) error
	UpdateEvent(id string, title string, description string, datetime string, location string, max_participation int64, club_id string) error
	DeleteEvent(id string) error
	GetAllEvents() error
	JoinEvent(event_id string, user_id string) error
	LeaveEvent(event_id string, user_id string) error
}

type eventService struct {
	eventServiceClient EventServiceClient
}

func NewEventService(eventServiceClient EventServiceClient) EventService {
	return eventService{eventServiceClient}
}

func (base eventService) CreateEvent(title string, description string, datetime string, location string, max_participation int64, club_id string, created_by string) error {
	req := CreateEventRequest{
		Title:            title,
		Description:      description,
		Datetime:         datetime,
		Location:         location,
		MaxParticipation: max_participation,
		ClubId:           club_id,
		CreatedBy:        created_by,
	}

	res, err := base.eventServiceClient.CreateEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: CreateEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) GetEvent(id string) error {
	req := GetEventRequest{
		Id: id,
	}

	res, err := base.eventServiceClient.GetEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: GetEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) UpdateEvent(id string, title string, description string, datetime string, location string, max_participation int64, club_id string) error {
	req := UpdateEventRequest{
		Id:               id,
		Description:      description,
		Datetime:         datetime,
		Location:         location,
		MaxParticipation: max_participation,
		ClubId:           club_id,
	}

	res, err := base.eventServiceClient.UpdateEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: UpdateEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) DeleteEvent(id string) error {
	req := DeleteEventRequest{
		Id: id,
	}

	res, err := base.eventServiceClient.DeleteEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: DeleteEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) GetAllEvents() error {
	req := GetAllEventsRequest{}

	res, err := base.eventServiceClient.GetAllEvents(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: GetAllEvents\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) JoinEvent(event_id string, user_id string) error {
	req := JoinEventRequest{
		EventId: event_id,
		UserId:  user_id,
	}

	res, err := base.eventServiceClient.JoinEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: JoinEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}

func (base eventService) LeaveEvent(event_id string, user_id string) error {
	req := LeaveEventRequest{
		EventId: event_id,
		UserId:  user_id,
	}

	res, err := base.eventServiceClient.LeaveEvent(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: LeaveEvent\n")
	fmt.Printf("Response: %v\n", res)
	return nil
}
