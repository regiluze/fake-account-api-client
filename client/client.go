package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/regiluze/form3-account-api-client/resources"
)

const DefaultMimeType = "application/vnd.api+json"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	httpClient HTTPClient
	urlBuilder URLBuilder
}

func NewForm3APIClient(baseURL string, httpClient HTTPClient) *Form3Client {
	urlBuilder := NewURLBuilder(baseURL)
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Form3Client{
		httpClient: httpClient,
		urlBuilder: urlBuilder,
	}
}

func NewForm3APIClientWithTimeout(URLBuilder URLBuilder, timeout time.Duration) *Form3Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	return &Form3Client{
		httpClient: httpClient,
		urlBuilder: URLBuilder,
	}
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
	req.Header.Set("Accept", DefaultMimeType)
	req.Header.Set("Content-Type", DefaultMimeType)
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
	if resp.StatusCode > http.StatusBadRequest {
		return NewErrFromServer(method, url, resp.StatusCode)
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
