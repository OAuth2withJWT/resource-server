package app

import (
	"fmt"
)

type InvalidUserIdError struct {
	UserId int
}

func (e *InvalidUserIdError) Error() string {
	return fmt.Sprintf("Invalid user id: %d", e.UserId)
}
