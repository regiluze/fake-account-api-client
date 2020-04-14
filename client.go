package accountclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"./resources"
)

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

func (fc Form3Client) Create(resource resources.Resource) (*resources.DataContainer, error) {
	endpoint := "organisation/accounts"
	data := resources.NewDataContainer(resource)
	dataBt, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", fc.url, endpoint),
		bytes.NewBuffer(dataBt),
	)

	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var responseData resources.DataContainer
	err = json.Unmarshal(body, &responseData)
	return &responseData, nil
}
