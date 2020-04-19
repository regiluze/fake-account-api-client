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

type isRequestHeaderValues struct{ r *http.Request }

// Matcher to check if http Request header has correct accept and content-type
func IsRequestHeaderValues(r *http.Request) gomock.Matcher {
	return &isRequestHeaderValues{r}
}

func (i *isRequestHeaderValues) Matches(x interface{}) bool {
	req := x.(*http.Request)
	authorization := true
	if len(i.r.Header["Authorization"]) > 0 {
		authorization = req.Header["Authorization"][0] == i.r.Header["Authorization"][0]
	}
	return req.Header["Accept"][0] == i.r.Header["Accept"][0] &&
		req.Header["Content-Type"][0] == i.r.Header["Content-Type"][0] &&
		authorization
}

func (i *isRequestHeaderValues) String() string {
	return fmt.Sprintf("Headers : %s", i.r.Header)
}
