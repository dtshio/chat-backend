package core

import (
	"fmt"
)

const (
	InvalidData = "Invalid data format"
	CreateError = "Error creating entry"
	NotFoundError = "Entry not found"
	DuplicateError = "Entry already exists"
	UpdateError = "Error updating entry"
	DeleteError = "Error deleting entry"
)

func GenError(err string, model interface{}) error {
	newError := fmt.Errorf(fmt.Sprintf("%s: %v", err, model))
	return newError
}
