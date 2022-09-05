package app

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/turbak/joom-calendar/internal/listing"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
	"strconv"
)

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
