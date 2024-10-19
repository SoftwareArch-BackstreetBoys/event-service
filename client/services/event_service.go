package services

import (
	"context"
)

type EventService interface {
	CreateEvent(title string, description string, datetime string, location string, max_participation int64, club_id string, created_by string) (*CreateEventResponse, error)
	GetEvent(id string) (*GetEventResponse, error)
	UpdateEvent(id string, title string, description string, datetime string, location string, max_participation int64, club_id string) (*UpdateEventResponse, error)
	DeleteEvent(id string) (*DeleteEventResponse, error)
	GetAllEvents() (*GetAllEventsResponse, error)
	JoinEvent(event_id string, user_id string) (*JoinEventResponse, error)
	LeaveEvent(event_id string, user_id string) (*LeaveEventResponse, error)
	GetAllEventsByUser(user_id string) (*GetAllEventsByUserResponse, error)
	GetAllEventsByClub(club_id string) (*GetAllEventsByClubResponse, error)
	GetAllParticipatedEvents(user_id string) (*GetAllParticipatedEventsResponse, error)
	SearchEvents(search_query string, club_id string) (*SearchEventsResponse, error)
}

type eventService struct {
	eventServiceClient EventServiceClient
}

func NewEventService(eventServiceClient EventServiceClient) EventService {
	return eventService{eventServiceClient}
}

func (base eventService) CreateEvent(title string, description string, datetime string, location string, max_participation int64, club_id string, created_by string) (*CreateEventResponse, error) {
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
		return nil, err
	}

	return res, nil
}

func (base eventService) GetEvent(id string) (*GetEventResponse, error) {
	req := GetEventRequest{
		Id: id,
	}

	res, err := base.eventServiceClient.GetEvent(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) UpdateEvent(id string, title string, description string, datetime string, location string, max_participation int64, club_id string) (*UpdateEventResponse, error) {
	req := UpdateEventRequest{
		Id:               id,
		Title:            title,
		Description:      description,
		Datetime:         datetime,
		Location:         location,
		MaxParticipation: max_participation,
		ClubId:           club_id,
	}

	res, err := base.eventServiceClient.UpdateEvent(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) DeleteEvent(id string) (*DeleteEventResponse, error) {
	req := DeleteEventRequest{
		Id: id,
	}

	res, err := base.eventServiceClient.DeleteEvent(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) GetAllEvents() (*GetAllEventsResponse, error) {
	req := GetAllEventsRequest{}

	res, err := base.eventServiceClient.GetAllEvents(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) JoinEvent(event_id string, user_id string) (*JoinEventResponse, error) {
	req := JoinEventRequest{
		EventId: event_id,
		UserId:  user_id,
	}

	res, err := base.eventServiceClient.JoinEvent(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) LeaveEvent(event_id string, user_id string) (*LeaveEventResponse, error) {
	req := LeaveEventRequest{
		EventId: event_id,
		UserId:  user_id,
	}

	res, err := base.eventServiceClient.LeaveEvent(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) GetAllEventsByUser(user_id string) (*GetAllEventsByUserResponse, error) {
	req := GetAllEventsByUserRequest{
		UserId: user_id,
	}

	res, err := base.eventServiceClient.GetAllEventsByUser(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) GetAllEventsByClub(club_id string) (*GetAllEventsByClubResponse, error) {
	req := GetAllEventsByClubRequest{
		ClubId: club_id,
	}

	res, err := base.eventServiceClient.GetAllEventsByClub(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) GetAllParticipatedEvents(user_id string) (*GetAllParticipatedEventsResponse, error) {
	req := GetAllParticipatedEventsRequest{
		UserId: user_id,
	}

	res, err := base.eventServiceClient.GetAllParticipatedEvents(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (base eventService) SearchEvents(search_query string, club_id string) (*SearchEventsResponse, error) {
	req := SearchEventsRequest{
		SearchQuery: search_query,
		ClubId:      club_id,
	}

	res, err := base.eventServiceClient.SearchEvents(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
