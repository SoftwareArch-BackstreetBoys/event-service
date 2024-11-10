package model

type GetAllEventsRequestBody struct {
	ClubIDs []string `json:"joined_club_ids"`
}
