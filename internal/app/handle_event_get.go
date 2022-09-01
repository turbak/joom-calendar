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
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

		return event, nil
	}
}
