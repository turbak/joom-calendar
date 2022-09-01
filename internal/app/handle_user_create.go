package app

import (
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/creating"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}

func (a *App) handleCreateUser() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		var args CreateUserRequest

		if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
			return nil, CodableError{Err: err, StatusCode: http.StatusBadRequest}
		}

		user := creating.User{
			Name:  args.Name,
			Email: args.Email,
		}

		createdID, err := a.creator.CreateUser(req.Context(), user)
		if err != nil {
			if errors.Is(err, creating.ErrUserAlreadyExists) {
				return nil, CodableError{Err: err, StatusCode: http.StatusConflict}
			}

			return nil, err
		}

		return CreateUserResponse{ID: createdID}, nil
	}
}
