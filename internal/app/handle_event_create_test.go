package app

import (
	"testing"
	"time"
)

func Test_validateCreateEventRequest(t *testing.T) {
	type args struct {
		args CreateEventRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				args: CreateEventRequest{
					Title:           "title",
					Desc:            "desc",
					StartDate:       time.Now(),
					InvitedUserIDs:  []int{1, 2, 3},
					OrganizerUserID: 1,
					Duration:        1,
					Repeat: &CreateEventRequestRepeat{
						DaysOfWeek:  []int{1, 2, 3},
						DayOfMonth:  1,
						MonthOfYear: 1,
						WeekOfMonth: "first",
					},
				},
			},
		},
		{
			name: "invalid title",
			args: args{
				args: CreateEventRequest{
					Title:           "",
					Desc:            "desc",
					StartDate:       time.Now(),
					InvitedUserIDs:  []int{1, 2, 3},
					OrganizerUserID: 1,
					Duration:        1,
					Repeat: &CreateEventRequestRepeat{
						DaysOfWeek:  []int{1, 2, 3},
						DayOfMonth:  1,
						MonthOfYear: 1,
						WeekOfMonth: "first",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid repeat",
			args: args{
				args: CreateEventRequest{
					Title:           "title",
					Desc:            "desc",
					StartDate:       time.Now(),
					InvitedUserIDs:  []int{1, 2, 3},
					OrganizerUserID: 1,
					Duration:        1,
					Repeat: &CreateEventRequestRepeat{
						DaysOfWeek:  []int{1, 2, 3},
						DayOfMonth:  1,
						MonthOfYear: 1,
						WeekOfMonth: "invalid",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCreateEventRequest(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validateCreateEventRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
