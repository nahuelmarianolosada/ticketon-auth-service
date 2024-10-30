package model

import "fmt"

type ApiError struct {
	Message string `json:"message"`
	Err     error  `json:"error"`
}

func (apiErr ApiError) Error() string {
	return fmt.Sprintf("%s", apiErr.Message)
}
