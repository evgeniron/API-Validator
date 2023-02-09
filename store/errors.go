package store

import "fmt"

const recordNotFoundErrorMessage = "record not found"

type RecordNotFoundError struct {
	Key string
}

func (e *RecordNotFoundError) Error() string {
	return fmt.Sprintf("%s: %s", recordNotFoundErrorMessage, e.Key)
}
