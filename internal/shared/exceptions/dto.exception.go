package exceptions

import (
	"errors"
	"strings"
)

func ErrInvalidDTO(dto ...string) error {
	invalid := strings.Join(dto, ", ")
	return errors.New("invalid dto" + invalid)
}

func ErrDuplicated(dto ...string) error {
	invalid := strings.Join(dto, ", ")
	return errors.New("duplicated data" + invalid)
}
