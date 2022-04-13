package main

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{message}
}
