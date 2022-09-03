package creating

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
	DaysOfWeek  []int
	DayOfMonth  int
	MonthOfYear int
	WeekOfMonth int
}
