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

var (
	emptyID                  = ""
	emptyParameters          = map[string]string{}
	basicErrorAPIStatusCodes = [...]int{401, 403, 405, 406, 409, 429, 500, 502, 503, 504}
	resourcesEndpointsMap    = map[resources.ResourceName]string{
		resources.Account: "organisation/accounts",
	}
)

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
	httpClient HTTPClient
	baseURL    string
	headers    map[string]string
}

func NewClientHeaders(accept, contentType string) map[string]string {
	return map[string]string{
		"Accept":       accept,
		"Content-Type": contentType,
	}
}

func NewForm3Client(apiBaseURL string, headers map[string]string, httpClient HTTPClient) *Form3Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Form3Client{
		httpClient: httpClient,
		baseURL:    apiBaseURL,
		headers:    headers,
	}
}

func (fc *Form3Client) SetHeaders(headers map[string]string) {
	fc.headers = headers
}

func (fc Form3Client) Create(resourceName resources.ResourceName, resource resources.Resource) (*resources.DataContainer, error) {
	data := resources.NewDataContainer(resource)
	dataB, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fc.buildRequestURL(resourceName, emptyID, emptyParameters)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataB))

	responseData := &resources.DataContainer{}
	if err := fc.makeRequest(req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) Fetch(resourceName resources.ResourceName, id string) (*resources.DataContainer, error) {
	url := fc.buildRequestURL(resourceName, id, emptyParameters)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	responseData := &resources.DataContainer{}
	if err := fc.makeRequest(req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) List(resourceName resources.ResourceName, filter map[string]interface{}, pageNumber, pageSize int) (*resources.ListDataContainer, error) {
	url := fc.buildRequestURL(
		resourceName,
		emptyID,
		map[string]string{
			"page[number]": strconv.Itoa(pageNumber),
			"page[size]":   strconv.Itoa(pageSize),
			// TODO add the filter query parameter
		})
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	responseData := &resources.ListDataContainer{}
	if err := fc.makeRequest(req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) Delete(resourceName resources.ResourceName, id string, version int) error {
	url := fc.buildRequestURL(
		resourceName,
		id,
		map[string]string{
			"version": strconv.Itoa(version),
		},
	)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	return fc.makeRequest(req, nil)
}

func (fc Form3Client) makeRequest(req *http.Request, responseData interface{}) error {
	req.Header.Set("Accept", fc.headers["Accept"])
	req.Header.Set("Content-Type", fc.headers["Content-Type"])

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return err
	}
	if err := fc.isResponseStatusCodeAnError(resp, req.Method, req.URL.String()); err != nil {
		return err
	}
	if responseData != nil {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &responseData); err != nil {
			return err
		}
	}
	return nil
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

func (fc Form3Client) buildRequestURL(resourceName resources.ResourceName, id string, parameters map[string]string) string {
	endpoint := resourcesEndpointsMap[resourceName]
	idPath := ""
	queryParams := ""
	if len(id) > 0 {
		idPath = fmt.Sprintf("/%s", id)
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
	return fmt.Sprintf("%s/%s%s%s", fc.baseURL, endpoint, idPath, queryParams)
}
