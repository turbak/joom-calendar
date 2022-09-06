package listing

import (
	"context"
	"github.com/teambition/rrule-go"
	"time"
)

type Storage interface {
	GetEventByID(ctx context.Context, ID int) (*Event, error)
	ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]Event, error)
	BatchGetEventsByUserIDs(ctx context.Context, userIDs []int) ([]Event, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) GetEventByID(ctx context.Context, ID int) (*Event, error) {
	event, err := s.storage.GetEventByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *Service) ListUsersEvents(ctx context.Context, userID int, from, to time.Time) ([]Event, error) {
	events, err := s.storage.ListUsersEvents(ctx, userID, from, to)
	if err != nil {
		return nil, err
	}

	j := 0
	for i := range events {
		if events[i].IsRepeated {
			isEventInPeriod := len(events[i].Rrule.Between(from, to, true)) > 0
			if isEventInPeriod {
				events[j] = events[i]
				j++
			}
		} else {
			events[j] = events[i]
			j++
		}
	}

	res := events[:j]
	if len(res) == 0 {
		return nil, ErrEventNotFound
	}

	return res, nil
}

func (s *Service) GetNearestEmptyTimeInterval(ctx context.Context, userIDs []int, minDuration time.Duration) (time.Time, time.Time, error) {
	const minIntervalTimeout = 5 * time.Second

	events, err := s.storage.BatchGetEventsByUserIDs(ctx, userIDs)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	rrules := pluckRRules(events)

	ctx, cancel := context.WithTimeout(ctx, minIntervalTimeout)
	defer cancel()

	min, err := findMinInterval(ctx, rrules, minDuration)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	max := time.Time{}
	for _, rule := range rrules {
		if after := rule.After(min, true); !after.IsZero() && (max.IsZero() || after.Before(max)) {
			max = after
		}
	}

	if max.IsZero() {
		max = min.Add(minDuration)
	}

	return min, max, nil
}

func findMinInterval(ctx context.Context, rrules []*rrule.RRule, minDuration time.Duration) (time.Time, error) {
	min := time.Now()
	minPlusDuration := min.Add(minDuration)
	for i := 0; i < len(rrules); i++ {
		select {
		case <-ctx.Done():
			return time.Time{}, ctx.Err()
		default:
		}

		if between := rrules[i].Between(min, minPlusDuration, true); len(between) > 0 {
			min = between[len(between)-1]
			minPlusDuration = min.Add(minDuration)
			i = -1
		}
	}

	return min, nil
}

func pluckRRules(events []Event) []*rrule.RRule {
	rrules := make([]*rrule.RRule, 0, len(events))
	for _, event := range events {
		if event.IsRepeated {
			rrules = append(rrules, event.Rrule)
		}
	}
	return rrules
}
