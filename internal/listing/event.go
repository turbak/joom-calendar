package listing

import "time"

type Event struct {
	ID          int
	Title       string
	Description string
	Duration    int
	Attendees   []EventAttendee
	Repeat      EventRepeat
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

type EventRepeat struct {
	StartDate   time.Time
	DayOfWeek   string
	DayOfMonth  string
	MonthOfYear string
	WeekOfMonth string
}

type EventAttendeeStatus string

const (
	EventAttendeeStatusUnconfirmed EventAttendeeStatus = "unconfirmed"
	EventAttendeeStatusConfirmed   EventAttendeeStatus = "confirmed"
	EventAttendeeStatusOrganizer   EventAttendeeStatus = "organizer"
)
