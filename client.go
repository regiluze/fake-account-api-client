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

type User struct {
	ResourceType   string                 `json:"type"`
	Id             string                 `json:"id"`
	OrganisationId string                 `json:"organizarion_id"`
	Attributes     map[string]interface{} `json:"attributes"`
}

func (fc Form3Client) Create(resource resources.Resource) (*resources.DataContainer, error) {
	data := resources.NewDataContainer(resource)
	dataBt, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(
		"POST",
		fc.url,
		bytes.NewBuffer(dataBt),
	)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	resp, err := fc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(">>>> status code ", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(">>>> resp body ", string(body))
	var responseData resources.DataContainer
	err = json.Unmarshal(body, &responseData)
	fmt.Println(">>>> unmarshal error ", err)
	return &responseData, nil
}
