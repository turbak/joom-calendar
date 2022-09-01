package adding

import "errors"

// ErrUserAlreadyExists is returned when user already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrSomeUsersNotFound is returned when some users are not found
var ErrSomeUsersNotFound = errors.New("some users not found")
