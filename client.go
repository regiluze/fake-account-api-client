package accountclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"./resources"
)

var basicErrorAPIStatusCodes = [...]int{401, 403, 404, 405, 406, 409, 429, 500, 502, 503, 504}

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
	url := fc.buildRequestURL("accounts")
	req, _ := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(dataB),
	)

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := fc.isResponseStatusCodeAnError(resp, url); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	// TODO test for this
	//if err != nil {
	//	return nil, err
	//}
	var responseData resources.DataContainer
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func (fc Form3Client) isResponseStatusCodeAnError(resp *http.Response, url string) error {
	for _, errorStatusCode := range basicErrorAPIStatusCodes {
		if errorStatusCode == resp.StatusCode {
			return ErrFromServer{http.MethodPost, url, resp.StatusCode}
		}
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fc.buildBadRequestError(resp)
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

func (fc Form3Client) buildRequestURL(resource string) string {
	endpoint := "organisation/accounts"
	return fmt.Sprintf("%s/%s", fc.url, endpoint)
}
