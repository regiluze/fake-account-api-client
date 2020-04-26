package client

import (
	"fmt"

	"github.com/regiluze/form3-account-api-client/resources"
)

// ErrNotFound is returned when getting a 404 status code.
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

// ErrBadRequest is returned when getting a 400 status code.
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

// ErrResponseStatusCode is returned when getting a 50X and 40X status codes,
// less for 400 and 404 status codes.
type ErrResponseStatusCode struct {
	method     string
	url        string
	StatusCode int
}

func NewErrResponseStatusCode(method, url string, statusCode int) error {
	return ErrResponseStatusCode{method, url, statusCode}
}

func (e ErrResponseStatusCode) Error() string {
	return fmt.Sprintf(
		"Error requesting (%s, %s): Status code %d",
		e.method,
		e.url,
		e.StatusCode,
	)
}
