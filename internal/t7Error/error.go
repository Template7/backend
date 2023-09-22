package t7Error

import (
	"fmt"
	"net/http"
)

type Code int
type Type int

type Error struct {
	Code    Code   `json:"code" example:"1024"`
	Message string `json:"message" example:"token unauthorized"`
	Detail  string `json:"detail" example:"empty token"`
	Type    Type   `json:"type" example:"32"`
	status  int
}

func (e Error) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
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

func ToT7Error(err error) (*Error, bool) {
	if t7Err, ok := err.(*Error); ok {
		return t7Err, true
	}
	return nil, false
}
