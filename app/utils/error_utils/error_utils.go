package error_utils

import (
	"encoding/json"
	"net/http"
)

type MessageErr interface {
	Message() string
	Status() int
	Error() string
}

type messageErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

func (err *messageErr) Error() string {
	return err.ErrError
}

func (err *messageErr) Message() string {
	return err.ErrMessage
}

func (err *messageErr) Status() int {
	return err.ErrStatus
}

func NewNotFoundError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusNotFound,
		ErrError: "not_found",
	}
}

func NewBadRequestError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusBadRequest,
		ErrError: "bad_request",
	}
}

func NewUnprocessibleEntityError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusUnprocessableEntity,
		ErrError: "invalid_request",
	}
}

func NewApiErrFromBytes(body []byte) (MessageErr, error) {
	var result messageErr
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func NewInternalServerError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusInternalServerError,
		ErrError: "server_error",
	}
}

func NewUnauthorizedError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusUnauthorized,
		ErrError: "unauthorized",
	} 
}

func NewForbiddenError(msg string) MessageErr {
	return &messageErr{
		ErrMessage: msg,
		ErrStatus: http.StatusForbidden,
		ErrError: "forbidden",
	} 
}
