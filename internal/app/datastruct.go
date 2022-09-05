package app

import "time"

type Event struct {
	ID          int             `json:"id,omitempty"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Duration    int             `json:"duration,omitempty"`
	Attendees   []EventAttendee `json:"attendees,omitempty"`
	Rrule       string          `json:"rrule,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type EventAttendee struct {
	EventID   int       `json:"event_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	User      User      `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	Frequency   string `json:"frequency"`
	DaysOfWeek  []int  `json:"days_of_week"`
	DayOfMonth  int    `json:"day_of_month"`
	MonthOfYear int    `json:"month_of_year"`
	WeekOfMonth int    `json:"week_of_month"`
}

type CreateEventResponse struct {
	ID int `json:"id"`
}
