package accountclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"./resources"
	"./test"
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
				httpClientMock.EXPECT().Do(test.IsRequestMethod("POST")).Return(nil, errors.New("fake")).Times(1)

				resource := resources.NewAccount(id, organisationID, map[string]interface{}{})

				client.CreateAccount(resource)
			})
			It("builds a request with resource endpoint", func() {
				httpClientMock.EXPECT().Do(test.IsRequestURL(fmt.Sprintf("%s/organisation/accounts", baseURL))).Return(nil, errors.New("fake")).Times(1)

				resource := resources.NewAccount(id, organisationID, map[string]interface{}{})

				client.CreateAccount(resource)
			})
			It("builds a request with dataContainer struct data", func() {
				account := test.BuildBasicAccountResource(id, organisationID)
				data := resources.NewDataContainer(account)
				dataB, _ := json.Marshal(data)
				req, _ := http.NewRequest(
					"POST",
					fmt.Sprintf("%s/organisation/accounts", baseURL),
					bytes.NewBuffer(dataB),
				)
				httpClientMock.EXPECT().Do(test.IsRequestBody(req)).Return(nil, errors.New("fake")).Times(1)

				client.CreateAccount(account)
			})
			Context("When getting succesful response", func() {
				It("returns resourceContainer struct as response data", func() {
					account := test.BuildBasicAccountResource(id, organisationID)
					data := resources.NewDataContainer(account)
					dataBt, _ := json.Marshal(data)
					expectedResponseBody := ioutil.NopCloser(bytes.NewReader(dataBt))
					httpClientMock.EXPECT().Do(gomock.Any()).Return(
						&http.Response{
							StatusCode: 201,
							Body:       expectedResponseBody,
						},
						nil,
					).Times(1)

					response, err := client.CreateAccount(account)

					Expect(err).To(BeNil())
					Expect(response.Data.ID).To(Equal(id))
				})
			})
		})
	})
})
