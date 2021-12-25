package exceptions

import "errors"

var ErrFullRoom error = errors.New("room is full")

var ErrNoRoom error = errors.New("no room exist")
