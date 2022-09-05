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
						Frequency:   "daily",
						DayOfMonth:  1,
						MonthOfYear: 1,
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
						Frequency:   "daily",
						DaysOfWeek:  []int{1},
						DayOfMonth:  1,
						MonthOfYear: 1,
						WeekOfMonth: 5,
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
