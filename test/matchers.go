package test

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/golang/mock/gomock"
)

type isRequestMethod struct{ m string }

// Matcher to check if http Request object has correct method
func IsRequestMethod(m string) gomock.Matcher {
	return &isRequestMethod{m}
}

func (i *isRequestMethod) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return req.Method == i.m
}

func (i *isRequestMethod) String() string {
	return fmt.Sprintf("HTTP method %s", i.m)
}

type isRequestURL struct{ u string }

// Matcher to check if http Request object has correct URL
func IsRequestURL(u string) gomock.Matcher {
	return &isRequestURL{u}
}

func (i *isRequestURL) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return req.URL.String() == i.u
}

func (i *isRequestURL) String() string {
	return fmt.Sprintf("URL %s", i.u)
}

type isRequestBody struct{ r *http.Request }

// Matcher to check if http Request object has correct body
func IsRequestBody(r *http.Request) gomock.Matcher {
	return &isRequestBody{r}
}

func (i *isRequestBody) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return reflect.DeepEqual(req.Body, i.r.Body)
}

func (i *isRequestBody) String() string {
	return fmt.Sprintf("Body: %s", i.r.Body)
}
