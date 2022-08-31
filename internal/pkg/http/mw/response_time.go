package mw

import (
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
	"time"
)

func ResponseTimeLogging() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			logger.Debugf("%s request to %s took %s", r.Method, r.RequestURI, time.Since(start))
		})
	}
}
