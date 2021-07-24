package t7Error

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
	status  int
}

func (e Error) Error() string {
	return fmt.Sprintf("error code: %s, message: %s", e.Code, e.Message)
}

func (e Error) WithDetail(d string) *Error {
	e.Detail = d
	return &e
}

func (e Error) GetStatus() int {
	if e.status == 0 {
		return http.StatusBadRequest
	} else {
		return e.status
	}
}

func (e Error) WithStatus(s int) *Error {
	e.status = s
	return &e
}

func (e Error) WithDetailAndStatus(d string, s int) *Error {
	e.Detail = d
	e.status = s
	return &e
}
