package app

import (
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/adding"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"golang.org/x/exp/maps"
	"net/http"
	"time"
)

type CreateEventRequest struct {
	Title           string                    `json:"title"`
	Desc            string                    `json:"desc"`
	StartDate       time.Time                 `json:"start"`
	InvitedUserIDs  []int                     `json:"invited_user_ids"`
	OrganizerUserID int                       `json:"organizer_user_id"`
	Duration        int                       `json:"duration"`
	Repeat          *CreateEventRequestRepeat `json:"repeat"`
}

type CreateEventRequestRepeat struct {
	DaysOfWeek  []int  `json:"days_of_week"`
	DayOfMonth  int    `json:"day_of_month"`
	MonthOfYear int    `json:"month_of_year"`
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

		event := adding.Event{
			Title:           args.Title,
			Description:     args.Desc,
			StartDate:       args.StartDate,
			Duration:        args.Duration,
			OrganizerUserID: args.OrganizerUserID,
			InvitedUserIDs:  args.InvitedUserIDs,
			Repeat:          toAddingEventRepeat(args.Repeat),
		}

		createdID, err := a.addingService.CreateEvent(req.Context(), event)
		if err != nil {
			return nil, err
		}

		return CreateEventResponse{ID: createdID}, nil
	}
}

func toAddingEventRepeat(repeat *CreateEventRequestRepeat) *adding.EventRepeat {
	if repeat == nil {
		return nil
	}

	return &adding.EventRepeat{
		DaysOfWeek:  toAddingEventRepeatDaysOfWeek(repeat.DaysOfWeek),
		DayOfMonth:  repeat.DayOfMonth,
		MonthOfYear: time.Month(repeat.MonthOfYear),
		WeekOfMonth: repeat.WeekOfMonth,
	}
}

func toAddingEventRepeatDaysOfWeek(weekday []int) []time.Weekday {
	res := make(map[time.Weekday]struct{}, len(weekday))
	for _, day := range weekday {
		res[time.Weekday(day)] = struct{}{}
	}

	return maps.Keys(res)
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

	if repeat.MonthOfYear != 0 && repeat.MonthOfYear < 1 && repeat.MonthOfYear > 12 {
		return CodableError{Err: errors.New("month_of_year must be between 1 and 12"), StatusCode: http.StatusBadRequest}
	}

	if repeat.WeekOfMonth != "" && repeat.WeekOfMonth != "first" && repeat.WeekOfMonth != "second" && repeat.WeekOfMonth != "third" && repeat.WeekOfMonth != "fourth" && repeat.WeekOfMonth != "last" {
		return CodableError{Err: errors.New("week_of_month must be one of first, second, third, fourth, last"), StatusCode: http.StatusBadRequest}
	}

	for _, day := range repeat.DaysOfWeek {
		if day < 1 || day > 7 {
			return CodableError{Err: errors.New("days_of_week must be between 1 and 7"), StatusCode: http.StatusBadRequest}
		}
	}
	return nil
}
