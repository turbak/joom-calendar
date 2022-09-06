package app

import (
	"fmt"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
	"strconv"
	"time"
)

type FindNearestIntervalResponse struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (a *App) handleFindNearestTimeInterval() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		duration, err := time.ParseDuration(req.URL.Query().Get("duration"))
		if err != nil {
			return nil, CodableError{Err: fmt.Errorf("failed to parse duration: %v", err), StatusCode: http.StatusBadRequest}
		}

		if duration < time.Minute {
			return nil, CodableError{Err: fmt.Errorf("duration is too small: %v", duration), StatusCode: http.StatusBadRequest}
		}

		userIDs, err := parseUserIDs(req.URL.Query()["user_ids"])
		if err != nil || len(userIDs) == 0 {
			return nil, CodableError{Err: fmt.Errorf("failed to parse user ids: %v", err), StatusCode: http.StatusBadRequest}
		}

		min, max, err := a.lister.GetNearestEmptyTimeInterval(req.Context(), userIDs, duration)
		if err != nil {
			return nil, err
		}

		return FindNearestIntervalResponse{
			Start: min,
			End:   max,
		}, nil
	}
}

func parseUserIDs(userIDs []string) ([]int, error) {
	ret := make([]int, 0, len(userIDs))
	for _, id := range userIDs {
		userID, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, userID)
	}

	return ret, nil
}
