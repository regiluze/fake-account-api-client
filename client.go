package accountclient

func NewForm3Client(httpClient HTTPClient) *Form3Client {
	return &Form3Client{}
}

type HTTPClient interface {
}

type Form3Client struct {
	httpClient HTTPClient
}
