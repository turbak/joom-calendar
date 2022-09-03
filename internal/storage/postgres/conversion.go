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

func toListingEvents(events []Event, attendees []FullEventAttendee) []listing.Event {
	listingEvents := make([]listing.Event, 0, len(events))
	for _, event := range events {
		listingEvents = append(listingEvents, listing.Event{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Duration:    event.Duration,
			StartDate:   event.StartDate,
			DaysOfWeek:  event.DaysOfWeek,
			DayOfMonth:  event.DayOfMonth,
			MonthOfYear: event.MonthOfYear,
			WeekOfMonth: event.WeekOfMonth,
			Attendees:   toListingEventAttendees(attendees),
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		})
	}

	return listingEvents
}
