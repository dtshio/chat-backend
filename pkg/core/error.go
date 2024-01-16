package core

import (
	"fmt"
)

const (
	InvalidData = "Invalid data format"
	CreatingError = "Error creating"
	GettingError = "Error getting"
	DuplicateError = "Aleady exists"
	UpdateError = "Error updating"
	DeleteError = "Error deleting"
)

func GenError(err string, model interface{}) error {
	newError := fmt.Errorf(fmt.Sprintf("%s: %v", err, model))
	return newError
}
