package errors

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %d Message: %s", e.Code, e.Message)
}

func New(Code int, Message string) Error {
	return Error{Code, Message}
}
