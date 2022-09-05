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

func (a *App) handleGetUserEvents() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		userID, err := strconv.Atoi(chi.URLParam(req, "user_id"))
		if err != nil {
			return nil, CodableError{Err: errors.New("invalid user id"), StatusCode: http.StatusBadRequest}
		}

		from, err := time.Parse(time.RFC3339, req.URL.Query().Get("from"))
		if err != nil {
			return nil, CodableError{Err: errors.New("invalid from date"), StatusCode: http.StatusBadRequest}
		}

		to, err := time.Parse(time.RFC3339, req.URL.Query().Get("to"))
		if err != nil {
			return nil, CodableError{Err: errors.New("invalid to date"), StatusCode: http.StatusBadRequest}
		}

		events, err := a.lister.ListUsersEvents(req.Context(), userID, from, to)
		if err != nil {
			if errors.Is(err, listing.ErrEventNotFound) {
				return nil, CodableError{Err: errors.New("no events found"), StatusCode: http.StatusNotFound}
			}
			return nil, err
		}

		return toEvents(events), nil
	}
}
