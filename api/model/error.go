package model

import "fmt"

type ApiError struct {
	Message string `json:"message"`
}

func (apiErr ApiError) Error() string {
	return fmt.Sprintf("%s", apiErr.Message)
}
