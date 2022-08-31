package mw

import (
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
)

func RequestBodyLogging() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debugf("%s %s %s", r.Method, r.RequestURI, r.Body)
			h.ServeHTTP(w, r)
		})
	}
}
