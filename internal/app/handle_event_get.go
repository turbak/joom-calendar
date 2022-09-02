package app

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/turbak/joom-calendar/internal/listing"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	ID          int             `json:"id,omitempty"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Duration    int             `json:"duration,omitempty"`
	Attendees   []EventAttendee `json:"attendees,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type EventAttendee struct {
	EventID   int       `json:"event_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	User      User      `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *App) handleGetEvent() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		eventID, err := strconv.Atoi(chi.URLParam(req, "event_id"))
		if err != nil {
			return nil, CodableError{Err: errors.New("invalid event id"), StatusCode: http.StatusBadRequest}
		}

		event, err := a.lister.GetEventByID(req.Context(), eventID)
		if err != nil {
			if errors.Is(err, listing.ErrEventNotFound) {
				return nil, CodableError{Err: err, StatusCode: http.StatusNotFound}
			}
			return nil, err
		}

		return toEvent(event), nil
	}
}

func toEvent(event *listing.Event) Event {
	return Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Duration:    event.Duration,
		Attendees:   toEventAttendees(event.Attendees),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func toEventAttendees(attendees []listing.EventAttendee) []EventAttendee {
	var result []EventAttendee
	for _, attendee := range attendees {
		result = append(result, toEventAttendee(attendee))
	}
	return result
}

func toEventAttendee(attendee listing.EventAttendee) EventAttendee {
	return EventAttendee{
		EventID:   attendee.EventID,
		Status:    string(attendee.Status),
		User:      toUser(attendee.User),
		CreatedAt: attendee.CreatedAt,
		UpdatedAt: attendee.UpdatedAt,
	}
}

func toUser(user listing.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
