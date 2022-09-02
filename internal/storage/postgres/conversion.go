package postgres

import (
	"github.com/turbak/joom-calendar/internal/listing"
)

func toListingEventAttendees(attendees []FullEventAttendee) []listing.EventAttendee {
	res := make([]listing.EventAttendee, 0, len(attendees))
	for _, attendee := range attendees {
		res = append(res, listing.EventAttendee{
			Status: listing.EventAttendeeStatus(attendee.Status),
			User: listing.User{
				ID:        attendee.User.ID,
				Name:      attendee.User.Name,
				Email:     attendee.User.Email,
				CreatedAt: attendee.User.CreatedAt,
				UpdatedAt: attendee.User.UpdatedAt,
			},
			CreatedAt: attendee.CreatedAt,
			UpdatedAt: attendee.UpdatedAt,
		})
	}
	return res
}

func toListingEventRepeat(repeat EventRepeat) listing.EventRepeat {
	return listing.EventRepeat{
		StartDate:   repeat.StartDate,
		DayOfWeek:   repeat.DayOfWeek,
		DayOfMonth:  repeat.DayOfMonth,
		MonthOfYear: repeat.MonthOfYear,
		WeekOfMonth: repeat.WeekOfMonth,
	}
}

func toListingEvents(events []FullEvent, attendees []FullEventAttendee) []listing.Event {
	listingEvents := make([]listing.Event, 0, len(events))
	for _, event := range events {
		listingEvents = append(listingEvents, listing.Event{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Duration:    event.Duration,
			Attendees:   toListingEventAttendees(attendees),
			Repeat:      toListingEventRepeat(event.Repeat),
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		})
	}

	return listingEvents
}
