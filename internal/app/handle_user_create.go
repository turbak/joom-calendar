package app

import (
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/adding"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
)

func (a *App) handleCreateUser() httputil.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type response struct {
		ID int `json:"id"`
	}

	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		var args request

		if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
			return nil, CodableError{Err: err, StatusCode: http.StatusBadRequest}
		}

		user := adding.User{
			Name:  args.Name,
			Email: args.Email,
		}

		createdID, err := a.addingService.CreateUser(req.Context(), user)
		if err != nil {
			if errors.Is(err, adding.ErrUserAlreadyExists) {
				return nil, CodableError{Err: err, StatusCode: http.StatusConflict}
			}

			return nil, err
		}

		return response{ID: createdID}, nil
	}
}
