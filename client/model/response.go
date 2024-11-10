package model

import "client/services"

type GetAllEventsResponse struct {
	PublicEvents     []*services.Event `json:"public_events"`
	JoinedClubEvents []*services.Event `json:"joined_club_events"`
}
