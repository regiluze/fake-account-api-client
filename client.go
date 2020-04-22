package form3apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	httpClient HTTPClient
	urlBuilder URLBuilder
	mimeType   string
	authToken  string
}

func NewForm3APIClient(mimeType, authToken string, urlBuilder URLBuilder, httpClient HTTPClient) *Form3Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Form3Client{
		httpClient: httpClient,
		mimeType:   mimeType,
		urlBuilder: urlBuilder,
		authToken:  authToken,
	}
}

func NewForm3APIClientWithTimeout(mimeType, authToken string, URLBuilder URLBuilder, timeout time.Duration) *Form3Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &Form3Client{
		httpClient: httpClient,
		mimeType:   mimeType,
		urlBuilder: URLBuilder,
		authToken:  authToken,
	}
}

func (fc *Form3Client) SetMimeType(mimeType string) {
	fc.mimeType = mimeType
}

func (fc *Form3Client) SetAuthToken(token string) {
	fc.authToken = token
}

func (fc Form3Client) Create(ctx context.Context, resourceName resources.ResourceName, resource resources.Resource) (*resources.DataContainer, error) {
	data := resources.NewDataContainer(resource)
	dataB, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fc.urlBuilder.DoForResource(resourceName)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataB))

	responseData := &resources.DataContainer{}
	if err := fc.makeRequest(ctx, req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) Fetch(ctx context.Context, resourceName resources.ResourceName, id string) (*resources.DataContainer, error) {
	url := fc.urlBuilder.DoForResourceWithID(resourceName, id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	responseData := &resources.DataContainer{}
	if err := fc.makeRequest(ctx, req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) List(ctx context.Context, resourceName resources.ResourceName, filter map[string]interface{}, pageNumber, pageSize int) (*resources.ListDataContainer, error) {
	url := fc.urlBuilder.DoForResourceWithParameters(
		resourceName,
		map[string]string{
			"page[number]": strconv.Itoa(pageNumber),
			"page[size]":   strconv.Itoa(pageSize),
			// TODO add the filter query parameter
		},
	)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	responseData := &resources.ListDataContainer{}
	if err := fc.makeRequest(ctx, req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) Delete(ctx context.Context, resourceName resources.ResourceName, id string, version int) error {
	url := fc.urlBuilder.DoForResourceWithIDAndParameters(
		resourceName,
		id,
		map[string]string{
			"version": strconv.Itoa(version),
		},
	)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	return fc.makeRequest(ctx, req, nil)
}

func (fc Form3Client) makeRequest(ctx context.Context, req *http.Request, responseData interface{}) error {
	req.Header.Set("Accept", fc.mimeType)
	req.Header.Set("Content-Type", fc.mimeType)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", fc.authToken))
	cReq := req.WithContext(ctx)

	resp, err := fc.httpClient.Do(cReq)
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

func (fc Form3Client) isResponseStatusCodeAnError(resp *http.Response, method, url string) error {
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound{url}
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fc.buildBadRequestError(method, resp)
	}
	for _, errorStatusCode := range basicErrorAPIStatusCodes {
		if errorStatusCode == resp.StatusCode {
			return NewErrFromServer(method, url, resp.StatusCode)
		}
	}
	return nil
}

func (fc Form3Client) buildBadRequestError(method string, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var errorData resources.BadRequestData
	if err := json.Unmarshal(body, &errorData); err != nil {
		return err
	}
	return NewErrBadRequest(method, errorData)
}
