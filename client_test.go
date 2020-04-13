package accountclient

import (
	"fmt"
	"net/http"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	//. "github.com/onsi/gomega"
	"./resources"
)

const (
	id             = "ad27e265-4402-3b3b-a0e5-3004ea9cc8dc"
	organisationId = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
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
		client = NewForm3Client("http://localhost:8080/v1/organisation/accounts", http.DefaultClient)
		//client = NewForm3Client("http://localhost:8080/v1/organisation/accounts", httpClientMock)
	})

	Describe("Resource operations", func() {
		Context("Create", func() {
			It("blablabla", func() {
				httpClientMock.EXPECT().Do(gomock.Any()).Return(nil, nil).Times(1)

				attributes := map[string]interface{}{
					"country":       "GB",
					"base_currency": "GBP",
					"bank_id":       "400300",
					"bank_id_code":  "GBDSC",
					"bic":           "NWBKGB22",
				}
				resource := resources.NewAccount(id, organisationId, attributes)

				response, _ := client.Create(resource)
				fmt.Printf(">>>>> response %+v", response)
			})
		})
	})
})
