package creating

import (
	"github.com/teambition/rrule-go"
	"time"
)

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
	StartDate   time.Time
	Frequency   EventRepeatFrequency
	DaysOfWeek  []int
	DayOfMonth  int
	MonthOfYear int
	WeekOfMonth int
}

type EventRepeatFrequency string

const (
	EventRepeatFrequencyDaily   EventRepeatFrequency = "daily"
	EventRepeatFrequencyWeekly  EventRepeatFrequency = "weekly"
	EventRepeatFrequencyMonthly EventRepeatFrequency = "monthly"
	EventRepeatFrequencyYearly  EventRepeatFrequency = "yearly"
)

func (e *EventRepeat) ToRrule(startDate time.Time) string {
	if e == nil {
		return ""
	}

	opts := rrule.ROption{
		Freq:    e.Frequency.ToRrule(),
		Dtstart: startDate,
	}

	for _, day := range e.DaysOfWeek {
		opts.Byweekday = append(opts.Byweekday, toRruleWeekDay(day))
	}

	if e.DayOfMonth > 0 {
		opts.Bymonthday = []int{e.DayOfMonth}
	}

	if e.MonthOfYear > 0 {
		opts.Bymonth = []int{e.MonthOfYear}
	}

	if e.WeekOfMonth > 0 {
		opts.Byweekno = []int{e.WeekOfMonth}
	}

	rule, _ := rrule.NewRRule(opts)

	return rule.String()
}

func toRruleWeekDay(day int) rrule.Weekday {
	switch day {
	case 1:
		return rrule.MO
	case 2:
		return rrule.TU
	case 3:
		return rrule.WE
	case 4:
		return rrule.TH
	case 5:
		return rrule.FR
	case 6:
		return rrule.SA
	case 7:
		return rrule.SU
	}
	return rrule.SU
}

func (e EventRepeatFrequency) ToRrule() rrule.Frequency {
	switch e {
	case EventRepeatFrequencyDaily:
		return rrule.DAILY
	case EventRepeatFrequencyWeekly:
		return rrule.WEEKLY
	case EventRepeatFrequencyMonthly:
		return rrule.MONTHLY
	case EventRepeatFrequencyYearly:
		return rrule.YEARLY
	}
	return rrule.YEARLY
}
