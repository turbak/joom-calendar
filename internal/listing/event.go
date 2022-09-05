package listing

import (
	"github.com/teambition/rrule-go"
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	Duration    int
	StartDate   time.Time
	Rrule       *rrule.RRule
	Attendees   []EventAttendee
	IsAllDay    bool
	IsRepeated  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type EventAttendee struct {
	EventID   int
	Status    EventAttendeeStatus
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventAttendeeStatus string

const (
	EventAttendeeStatusUnconfirmed EventAttendeeStatus = "unconfirmed"
	EventAttendeeStatusConfirmed   EventAttendeeStatus = "confirmed"
	EventAttendeeStatusOrganizer   EventAttendeeStatus = "organizer"
)
