package entity

import (
	"errors"
	"fmt"
)

var ErrUserConflict = errors.New("user with same login already exists")

type (
	ModuleNotFoundError struct {
		UUID string
	}
)

func (err *ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module with uuid=\"%s\" does not exist", err.UUID)
}
