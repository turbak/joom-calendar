package adding

import "errors"

// ErrUserAlreadyExists is returned when user already exists
var ErrUserAlreadyExists = errors.New("user already exists")
