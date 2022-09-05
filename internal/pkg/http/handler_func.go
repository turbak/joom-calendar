package http

import (
	"context"
	"encoding/json"
	"errors"
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

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				code = http.StatusRequestTimeout
			}

			http.Error(w, err.Error(), code)
			logger.Errorf("error while performing %s %s: %s", r.Method, r.RequestURI, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")

		switch resp.(type) {
		case nil:
			w.WriteHeader(http.StatusNoContent)
		case []byte:
			_, err = w.Write(resp.([]byte))
		default:
			err = json.NewEncoder(w).Encode(resp)
		}

		if err != nil {
			logger.Error("failed to encode response: %v", err)
		}
	}
}
