package app

import (
	"fmt"
	"net/http"
	"os"
)

func (a *App) handleGithubLogin() http.HandlerFunc {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	callbackURL := fmt.Sprintf("http://localhost:%s/callbacks/github", os.Getenv("PORT"))
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		clientID,
		callbackURL,
	)

	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
	}
}
