package accountclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"./resources"
	. "github.com/onsi/gomega"
)

const (
	id             = "ad27e265-4402-3b3b-a0e5-3004ea9cc8dc"
	organisationID = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	baseURL        = "api_base_url"
)

var _ = Describe("Account api resource client", func() {
	var (
		client         *Form3Client
		mockCtrl       *gomock.Controller
		httpClientMock *MockHTTPClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClientMock = NewMockHTTPClient(mockCtrl)
		client = NewForm3Client(baseURL, httpClientMock)
	})

	Describe("Account resource operations", func() {
		Context("Create", func() {
			It("builds a request with POST method", func() {
				httpClientMock.EXPECT().Do(IsRequestMethod("POST")).Return(nil, errors.New("fake")).Times(1)

				resource := resources.NewAccount(id, organisationID, map[string]interface{}{})

				client.Create(resource)
			})
			It("builds a request with resource endpoint", func() {
				httpClientMock.EXPECT().Do(IsRequestURL(fmt.Sprintf("%s/organisation/accounts", baseURL))).Return(nil, errors.New("fake")).Times(1)

				resource := resources.NewAccount(id, organisationID, map[string]interface{}{})

				client.Create(resource)
			})
			It("builds a request with dataContainer struct data", func() {
				resource := resources.NewAccount(id, organisationID, map[string]interface{}{})
				data := resources.NewDataContainer(resource)
				dataBt, _ := json.Marshal(data)
				req, _ := http.NewRequest(
					"POST",
					fmt.Sprintf("%s/organisation/accounts", baseURL),
					bytes.NewBuffer(dataBt),
				)
				httpClientMock.EXPECT().Do(IsRequestBody(req)).Return(nil, errors.New("fake")).Times(1)

				client.Create(resource)
			})
			Context("When getting succesful response", func() {
				It("returns resourceContainer struct as response data", func() {
					resource := resources.NewAccount(id, organisationID, map[string]interface{}{})
					data := resources.NewDataContainer(resource)
					dataBt, _ := json.Marshal(data)
					expectedResponseBody := ioutil.NopCloser(bytes.NewReader(dataBt))
					httpClientMock.EXPECT().Do(gomock.Any()).Return(
						&http.Response{
							StatusCode: 201,
							Body:       expectedResponseBody,
						},
						nil,
					).Times(1)

					response, err := client.Create(resource)

					Expect(err).To(BeNil())
					Expect(response.Data.ID).To(Equal(id))
				})

			})
		})
	})
})

type isRequestMethod struct{ m string }

func IsRequestMethod(m string) gomock.Matcher {
	return &isRequestMethod{m}
}

func (i *isRequestMethod) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return req.Method == i.m
}

func (i *isRequestMethod) String() string {
	return fmt.Sprintf("HTTP method %s", i.m)
}

type isRequestURL struct{ u string }

func IsRequestURL(u string) gomock.Matcher {
	return &isRequestURL{u}
}

func (i *isRequestURL) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return req.URL.String() == i.u
}

func (i *isRequestURL) String() string {
	return fmt.Sprintf("URL %s", i.u)
}

type isRequestBody struct{ r *http.Request }

func IsRequestBody(r *http.Request) gomock.Matcher {
	return &isRequestBody{r}
}

func (i *isRequestBody) Matches(x interface{}) bool {
	req := x.(*http.Request)
	return reflect.DeepEqual(req.Body, i.r.Body)
}

func (i *isRequestBody) String() string {
	return fmt.Sprintf("Body: %s", i.r.Body)
}
