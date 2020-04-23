package client

import (
	"fmt"

	"github.com/regiluze/form3-account-api-client/resources"
)

var (
	basicErrorAPIStatusCodes = [...]int{401, 403, 405, 406, 409, 429, 500, 502, 503, 504}
)

// ErrNotFound is returned when getting a 404 status code from server.
type ErrNotFound struct {
	url string
}

func NewErrNotFound(url string) error {
	return ErrNotFound{url}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf(
		"Resource or endpoint not exists: %s",
		e.url,
	)
}

// ErrBadRequest is returned when getting a 400 status code from server.
type ErrBadRequest struct {
	method    string
	errorData resources.BadRequestData
}

func NewErrBadRequest(method string, errorData resources.BadRequestData) error {
	return ErrBadRequest{method, errorData}
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf(
		"Bad Request (%s): Error code %d, message: %s",
		e.method,
		e.errorData.ErrorCode,
		e.errorData.ErrorMessage,
	)
}

func (e ErrBadRequest) getErrorData() resources.BadRequestData {
	return e.errorData
}

// ErrFromServer is returned when getting a 500 status code from server.
type ErrFromServer struct {
	method     string
	url        string
	statusCode int
}

func NewErrFromServer(method, url string, statusCode int) error {
	return ErrFromServer{method, url, statusCode}
}

func (e ErrFromServer) Error() string {
	return fmt.Sprintf(
		"Error requesting (%s, %s): Status code %d",
		e.method,
		e.url,
		e.statusCode,
	)
}
