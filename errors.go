package activitypub

import (
	"errors"
)

// ErrNotImplemented is an error indicating that a feature or specific functionality is not implemented.
var ErrNotImplemented = errors.New("Not implemented")

// ErrNotFound is an error indicating that an item is not present or was not found.
var ErrNotFound = errors.New("Not found")
