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

// FullEventAttendee is an event attendee with user data
type FullEventAttendee struct {
	EventID   int
	Status    EventAttendeeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
