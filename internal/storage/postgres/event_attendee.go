package postgres

import "time"

type EventAttendee struct {
	EventID   int
	UserID    int
	Status    EventAttendeeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventAttendeeStatus string

const (
	EventAttendeeStatusUnconfirmed EventAttendeeStatus = "unconfirmed"
	EventAttendeeStatusConfirmed   EventAttendeeStatus = "confirmed"
	EventAttendeeStatusOrganizer   EventAttendeeStatus = "organizer"
)
