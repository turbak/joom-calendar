package postgres

import "time"

type EventInvite struct {
	ID        int
	EventID   int
	UserID    int
	Status    EventInviteStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventInviteStatus string

const (
	EventInviteStatusPending  EventInviteStatus = "pending"
	EventInviteStatusAccepted EventInviteStatus = "accepted"
	EventInviteStatusDeclined EventInviteStatus = "declined"
)
