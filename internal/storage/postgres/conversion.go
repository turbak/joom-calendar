package postgres

import (
	"github.com/teambition/rrule-go"
	"github.com/turbak/joom-calendar/internal/listing"
)

func toListingEventAttendee(attendee FullEventAttendee) listing.EventAttendee {
	return listing.EventAttendee{
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
	}
}

func toListingEventAttendees(attendees []FullEventAttendee) []listing.EventAttendee {
	listingAttendees := make([]listing.EventAttendee, 0, len(attendees))
	for _, attendee := range attendees {
		listingAttendees = append(listingAttendees, toListingEventAttendee(attendee))
	}

	return listingAttendees
}

func toListingEvents(events []Event, attendees []FullEventAttendee) []listing.Event {
	eventAttendees := make(map[int][]listing.EventAttendee)
	for _, attendee := range attendees {
		eventAttendees[attendee.EventID] = append(eventAttendees[attendee.EventID], toListingEventAttendee(attendee))
	}

	listingEvents := make([]listing.Event, 0, len(events))
	for _, event := range events {
		listingEvents = append(listingEvents, toListingEvent(event, eventAttendees[event.ID]))
	}

	return listingEvents
}

func toListingEvent(event Event, eventAttendees []listing.EventAttendee) listing.Event {
	return listing.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Duration:    event.Duration,
		StartDate:   event.StartDate,
		IsAllDay:    event.IsAllDay,
		IsRepeated:  event.IsRepeated,
		Rrule:       toRrule(event.Rrule),
		Attendees:   eventAttendees,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func toRrule(rule string) *rrule.RRule {
	if rule == "" {
		return nil
	}

	ret, _ := rrule.StrToRRule(rule)
	return ret
}
