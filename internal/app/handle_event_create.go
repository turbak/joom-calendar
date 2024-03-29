package app

import (
	"encoding/json"
	"errors"
	"github.com/turbak/joom-calendar/internal/creating"
	httputil "github.com/turbak/joom-calendar/internal/pkg/http"
	"net/http"
)

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
			Repeat:          toCreatingEventRepeat(args.Repeat, args.StartDate),
		}

		createdID, err := a.creator.CreateEvent(req.Context(), event)
		if err != nil {
			return nil, err
		}

		return CreateEventResponse{ID: createdID}, nil
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

	if repeat.MonthOfYear != 0 {
		if repeat.MonthOfYear < 1 || repeat.MonthOfYear > 12 {
			return CodableError{Err: errors.New("month_of_year must be between 1 and 12"), StatusCode: http.StatusBadRequest}
		}
	}

	if repeat.WeekOfMonth != 0 &&
		repeat.WeekOfMonth != 1 &&
		repeat.WeekOfMonth != 2 &&
		repeat.WeekOfMonth != 3 &&
		repeat.WeekOfMonth != 4 &&
		repeat.WeekOfMonth != -1 {
		return CodableError{Err: errors.New("week_of_month must be in [1, 2, 3, 4, -1]"), StatusCode: http.StatusBadRequest}
	}

	if repeat.DayOfMonth != 0 {
		if repeat.DayOfMonth < 1 || repeat.DayOfMonth > 31 {
			return CodableError{Err: errors.New("day_of_month must be between 1 and 31"), StatusCode: http.StatusBadRequest}
		}
	}

	for _, day := range repeat.DaysOfWeek {
		if day < 0 || day > 6 {
			return CodableError{Err: errors.New("days_of_week must be between 0 and 6"), StatusCode: http.StatusBadRequest}
		}
	}

	if len(repeat.DaysOfWeek) != 0 && repeat.DayOfMonth != 0 {
		return CodableError{Err: errors.New("day_of_week and day_of_month cannot be both set"), StatusCode: http.StatusBadRequest}
	}

	if repeat.DayOfMonth != 0 && repeat.WeekOfMonth != 0 {
		return CodableError{Err: errors.New("day_of_month and week_of_month cannot be both set"), StatusCode: http.StatusBadRequest}
	}

	if repeat.Frequency == "" {
		return CodableError{Err: errors.New("frequency is required"), StatusCode: http.StatusBadRequest}
	}

	if repeat.Frequency != "daily" &&
		repeat.Frequency != "weekly" &&
		repeat.Frequency != "monthly" {
		return CodableError{Err: errors.New("frequency must be in [daily, weekly, monthly]"), StatusCode: http.StatusBadRequest}
	}

	return nil
}
