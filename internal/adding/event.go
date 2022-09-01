package adding

import "time"

type Event struct {
	Title           string
	Description     string
	StartDate       time.Time
	OrganizerUserID int
	InvitedUserIDs  []int
	Duration        int
	Repeat          *EventRepeat
}

type EventRepeat struct {
	DaysOfWeek  []time.Weekday
	DayOfMonth  int
	MonthOfYear time.Month
	WeekOfMonth string
}
