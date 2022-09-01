package app

import (
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/creating"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
	"strconv"
	"time"
)

type CreateEventRequest struct {
	Title           string                    `json:"title"`
	Desc            string                    `json:"desc"`
	StartDate       time.Time                 `json:"start_date"`
	InvitedUserIDs  []int                     `json:"invited_user_ids"`
	OrganizerUserID int                       `json:"organizer_user_id"`
	Duration        int                       `json:"duration"`
	Repeat          *CreateEventRequestRepeat `json:"repeat"`
}

type CreateEventRequestRepeat struct {
	DayOfWeek   string `json:"day_of_week"`
	DayOfMonth  string `json:"day_of_month"`
	MonthOfYear string `json:"month_of_year"`
	WeekOfMonth string `json:"week_of_month"`
}

type CreateEventResponse struct {
	ID int `json:"id"`
}

func (a *App) handleCreateEvent() httputil.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) (interface{}, error) {
		var args CreateEventRequest

		if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
			return nil, CodableError{Err: err, StatusCode: http.StatusBadRequest}
		}

		if err := validateCreateEventRequest(args); err != nil {
			return nil, err
		}

		event := creating.Event{
			Title:           args.Title,
			Description:     args.Desc,
			StartDate:       args.StartDate,
			Duration:        args.Duration,
			OrganizerUserID: args.OrganizerUserID,
			InvitedUserIDs:  args.InvitedUserIDs,
			Repeat:          toAddingEventRepeat(args.Repeat),
		}

		createdID, err := a.creator.CreateEvent(req.Context(), event)
		if err != nil {
			return nil, err
		}

		return CreateEventResponse{ID: createdID}, nil
	}
}

func toAddingEventRepeat(repeat *CreateEventRequestRepeat) *creating.EventRepeat {
	if repeat == nil {
		return nil
	}

	return &creating.EventRepeat{
		DayOfWeek:   repeat.DayOfWeek,
		DayOfMonth:  repeat.DayOfMonth,
		MonthOfYear: repeat.MonthOfYear,
		WeekOfMonth: repeat.WeekOfMonth,
	}
}

func validateCreateEventRequest(args CreateEventRequest) error {
	if args.Title == "" {
		return CodableError{Err: errors.New("title is required"), StatusCode: http.StatusBadRequest}
	}
	if args.StartDate.IsZero() {
		return CodableError{Err: errors.New("start is required"), StatusCode: http.StatusBadRequest}
	}
	if args.Duration == 0 {
		return CodableError{Err: errors.New("duration is required"), StatusCode: http.StatusBadRequest}
	}
	if args.OrganizerUserID == 0 {
		return CodableError{Err: errors.New("organizer_user_id is required"), StatusCode: http.StatusBadRequest}
	}
	if err := validateCreateEventRequestRepeat(args.Repeat); err != nil {
		return err
	}
	return nil
}

func validateCreateEventRequestRepeat(repeat *CreateEventRequestRepeat) error {
	if repeat == nil {
		return nil
	}

	if repeat.MonthOfYear != "" && repeat.MonthOfYear != "*" {
		monthOfYear, err := strconv.ParseInt(repeat.MonthOfYear, 10, 64)
		if err != nil {
			return CodableError{Err: errors.New("month_of_year must be a number or *"), StatusCode: http.StatusBadRequest}
		}
		if monthOfYear < 1 || monthOfYear > 12 {
			return CodableError{Err: errors.New("month_of_year must be between 1 and 12"), StatusCode: http.StatusBadRequest}
		}
	}

	if repeat.WeekOfMonth != "" && repeat.WeekOfMonth != "first" && repeat.WeekOfMonth != "second" && repeat.WeekOfMonth != "third" && repeat.WeekOfMonth != "fourth" && repeat.WeekOfMonth != "last" {
		return CodableError{Err: errors.New("week_of_month must be one of first, second, third, fourth, last"), StatusCode: http.StatusBadRequest}
	}

	if repeat.DayOfMonth != "" && repeat.DayOfMonth != "*" {
		dayOfMonth, err := strconv.ParseInt(repeat.DayOfMonth, 10, 64)
		if err != nil {
			return CodableError{Err: errors.New("day_of_month must be a number or *"), StatusCode: http.StatusBadRequest}
		}
		if dayOfMonth < 1 || dayOfMonth > 31 {
			return CodableError{Err: errors.New("day_of_month must be between 1 and 31"), StatusCode: http.StatusBadRequest}
		}
	}

	if repeat.DayOfWeek != "" && repeat.DayOfWeek != "*" {
		dayOfWeek, err := strconv.ParseInt(repeat.DayOfWeek, 10, 64)
		if err != nil {
			return CodableError{Err: errors.New("day_of_week must be a number or *"), StatusCode: http.StatusBadRequest}
		}
		if dayOfWeek < 1 || dayOfWeek > 7 {
			return CodableError{Err: errors.New("day_of_week must be between 1 and 7"), StatusCode: http.StatusBadRequest}
		}
	}

	if repeat.DayOfWeek != "" && repeat.DayOfMonth != "" {
		return CodableError{Err: errors.New("day_of_week and day_of_month cannot be both set"), StatusCode: http.StatusBadRequest}
	}

	if repeat.DayOfMonth != "" && repeat.WeekOfMonth != "" {
		return CodableError{Err: errors.New("day_of_month and week_of_month cannot be both set"), StatusCode: http.StatusBadRequest}
	}

	return nil
}
