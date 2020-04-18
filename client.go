package accountclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"./resources"
)

var basicErrorAPIStatusCodes = [...]int{401, 403, 405, 406, 409, 429, 500, 502, 503, 504}

// ErrNotFound is returned when getting a 404 status code from server.
type ErrNotFound struct {
	url string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf(
		"Resource or endpoint not exists: %s",
		e.url,
	)
}

// ErrBadRequest is returned when getting a 400 status code from server.
type ErrBadRequest struct {
	verb      string
	errorData resources.BadRequestData
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf(
		"Bad Request (%s): Error code %d, message: %s",
		e.verb,
		e.errorData.ErrorCode,
		e.errorData.ErrorMessage,
	)
}

// ErrFromServer is returned when getting a 500 status code from server.
type ErrFromServer struct {
	verb       string
	url        string
	statusCode int
}

func (e ErrFromServer) Error() string {
	return fmt.Sprintf(
		"Error requesting (%s, %s): Status code %d",
		e.verb,
		e.url,
		e.statusCode,
	)
}

func createError(verb, url string, expected, actual int) error {
	return fmt.Errorf(
		"Error requesting: %s %s (wrong http status code; expected:%d got:%d)",
		verb,
		url,
		expected,
		actual,
	)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	url        string
	httpClient HTTPClient
}

func NewForm3Client(apiBaseURL string, httpClient HTTPClient) *Form3Client {
	return &Form3Client{apiBaseURL, httpClient}
}

func (fc Form3Client) CreateAccount(resource resources.Resource) (*resources.DataContainer, error) {
	data := resources.NewDataContainer(resource)
	dataB, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fc.buildRequestURL(map[string]string{}, "accounts")
	req, _ := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(dataB),
	)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := fc.isResponseStatusCodeAnError(resp, http.MethodPost, url); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var responseData resources.DataContainer
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func (fc Form3Client) FetchAccount(id string) (*resources.DataContainer, error) {
	url := fc.buildRequestURL(map[string]string{}, "accounts", id)
	req, _ := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := fc.isResponseStatusCodeAnError(resp, http.MethodGet, url); err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var responseData resources.DataContainer
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func (fc Form3Client) ListAccount(filter map[string]interface{}, pageNumber, pageSize int) (*resources.ListDataContainer, error) {
	url := fc.buildRequestURL(map[string]string{
		"page[number]": strconv.Itoa(pageNumber),
		"page[size]":   strconv.Itoa(pageSize),
		// TODO add the filter query parameter
	}, "accounts")
	req, _ := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := fc.isResponseStatusCodeAnError(resp, http.MethodGet, url); err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var responseData resources.ListDataContainer
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func (fc Form3Client) DeleteAccount(id string, version int) error {
	url := fc.buildRequestURL(
		map[string]string{
			"version": strconv.Itoa(version),
		}, "accounts", id)
	req, _ := http.NewRequest(
		http.MethodDelete,
		url,
		nil,
	)
	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return err
	}
	if err := fc.isResponseStatusCodeAnError(resp, http.MethodDelete, url); err != nil {
		return err
	}
	return err
}

func (fc Form3Client) isResponseStatusCodeAnError(resp *http.Response, verb, url string) error {
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound{url}
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fc.buildBadRequestError(resp)
	}
	for _, errorStatusCode := range basicErrorAPIStatusCodes {
		if errorStatusCode == resp.StatusCode {
			return ErrFromServer{verb, url, resp.StatusCode}
		}
	}
	return nil
}

func (fc Form3Client) buildBadRequestError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var errorData resources.BadRequestData
	if err := json.Unmarshal(body, &errorData); err != nil {
		return err
	}
	return ErrBadRequest{http.MethodPost, errorData}
}

// TODO solve the resource endpoint problem
func (fc Form3Client) buildRequestURL(parameters map[string]string, paths ...string) string {
	endpoint := "organisation/accounts"
	idPath := ""
	queryParams := ""
	if len(paths) > 1 {
		idPath = fmt.Sprintf("/%s", paths[1])
	}
	if len(parameters) > 0 {
		flatParams := []string{}
		paramNames := []string{}
		for name := range parameters {
			paramNames = append(paramNames, name)
		}
		sort.Strings(paramNames)
		for _, name := range paramNames {
			flatParams = append(flatParams, fmt.Sprintf("%s=%s", name, parameters[name]))
		}
		queryParams = fmt.Sprintf("?%s", strings.Join(flatParams, "&"))
	}
	return fmt.Sprintf("%s/%s%s%s", fc.url, endpoint, idPath, queryParams)
}
