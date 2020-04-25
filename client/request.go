package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/regiluze/form3-account-api-client/resources"
)

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
