package app

import (
	"errors"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
)

type AuthenticationResponse struct {
	Token string `json:"token"`
}

func (a *App) handleGithubCallback() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		code := req.URL.Query().Get("code")
		if code == "" {
			return nil, CodableError{Err: errors.New("code is empty"), StatusCode: http.StatusUnauthorized}
		}

		token, err := a.authenticator.AuthenticateGithub(req.Context(), code)
		if err != nil {
			return nil, CodableError{Err: err, StatusCode: http.StatusUnauthorized}
		}

		return AuthenticationResponse{Token: token}, nil
	}
}
