package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/regiluze/form3-account-api-client/resources"
)

const DefaultMimeType = "application/vnd.api+json"

type Client interface {
	Fetch(ctx context.Context, resourceName resources.ResourceName, id string) (*resources.DataContainer, error)
	Create(ctx context.Context, resourceName resources.ResourceName, resource resources.Resource) (*resources.DataContainer, error)
	List(ctx context.Context, resourceName resources.ResourceName, filter map[string]interface{}, pageNumber, pageSize int) (*resources.ListDataContainer, error)
	Delete(ctx context.Context, resourceName resources.ResourceName, id string, version int) error
}

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

func (fc Form3Client) Create(ctx context.Context, resourceName resources.ResourceName, resource resources.Resource) (*resources.DataContainer, error) {
	data := resources.NewDataContainer(resource)
	dataB, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fc.urlBuilder.DoForResource(resourceName)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataB))
	if err != nil {
		return nil, err
	}

	responseData := &resources.DataContainer{}
	if err := fc.makeRequest(ctx, req, responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (fc Form3Client) Fetch(ctx context.Context, resourceName resources.ResourceName, id string) (*resources.DataContainer, error) {
	url := fc.urlBuilder.DoForResourceWithID(resourceName, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

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
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return fc.makeRequest(ctx, req, nil)
}
