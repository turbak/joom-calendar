package listing

import "errors"

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ErrEventNotFound is returned when an event is not found.
var ErrEventNotFound = errors.New("event not found")
