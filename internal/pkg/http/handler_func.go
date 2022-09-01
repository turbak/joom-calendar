package http

import (
	"encoding/json"
	"github.com/turbak/joom-calendar/internal/pkg/logger"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func Handler(f HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	type codableError interface {
		error
		Code() int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(w, r)
		if err != nil {
			code := http.StatusInternalServerError
			if err, ok := err.(codableError); ok {
				code = err.Code()
			}
			http.Error(w, err.Error(), code)
			logger.Errorf("error while performing %s %s: %s", r.Method, r.RequestURI, err.Error())
			return
		}

		if resp == nil {
			w.Write([]byte("{}"))
			return
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to encode response: %v", err)
		}
	}
}
