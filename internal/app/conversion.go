package app

import (
	"github.com/turbak/joom-calendar/internal/creating"
	"github.com/turbak/joom-calendar/internal/listing"
	"time"
)

func toEvent(event *listing.Event) Event {
	if event == nil {
		return Event{}
	}

	rrule := ""
	if event.Rrule != nil {
		rrule = event.Rrule.String()
	}

	return Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Duration:    event.Duration,
		Rrule:       rrule,
		Attendees:   toEventAttendees(event.Attendees),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func toEventAttendees(attendees []listing.EventAttendee) []EventAttendee {
	var result []EventAttendee
	for _, attendee := range attendees {
		result = append(result, toEventAttendee(attendee))
	}
	return result
}

func toEventAttendee(attendee listing.EventAttendee) EventAttendee {
	return EventAttendee{
		EventID:   attendee.EventID,
		Status:    string(attendee.Status),
		User:      toUser(attendee.User),
		CreatedAt: attendee.CreatedAt,
		UpdatedAt: attendee.UpdatedAt,
	}
}

func toUser(user listing.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toEvents(events []listing.Event) []Event {
	var es []Event
	for _, e := range events {
		es = append(es, toEvent(&e))
	}
	return es
}

func toCreatingEventRepeat(repeat *CreateEventRequestRepeat, startDate time.Time) *creating.EventRepeat {
	if repeat == nil {
		return nil
	}

	return &creating.EventRepeat{
		Frequency:   creating.EventRepeatFrequency(repeat.Frequency),
		StartDate:   startDate,
		DaysOfWeek:  repeat.DaysOfWeek,
		DayOfMonth:  repeat.DayOfMonth,
		MonthOfYear: repeat.MonthOfYear,
		WeekOfMonth: repeat.WeekOfMonth,
	}
}
