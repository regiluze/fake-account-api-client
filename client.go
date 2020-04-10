package accountclient

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	httpClient HTTPClient
}

func NewForm3Client(httpClient HTTPClient) *Form3Client {
	return &Form3Client{httpClient}
}
