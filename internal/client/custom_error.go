package client

import "fmt"

type CustomError struct {
	code    int
	message string
}

func NewCustomError(code int, format string, params ...interface{}) *CustomError {
	return &CustomError{
		code:    code,
		message: fmt.Sprintf(format, params...),
	}
}

// CustomError implements the error interface.
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error code %d: %s", e.code, e.message)
}
