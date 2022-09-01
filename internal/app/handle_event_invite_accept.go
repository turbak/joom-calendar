package app

import (
	"github.com/go-chi/chi"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
	"strconv"
)

func (a *App) handleAcceptInvite() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		inviteID, err := strconv.Atoi(chi.URLParam(req, "invite_id"))
		if err != nil {
			return nil, CodableError{Err: err, StatusCode: http.StatusBadRequest}
		}

		err = a.inviter.AcceptInvite(req.Context(), inviteID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
