package accountclient

import (
	"fmt"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Account api resource client", func() {
	var (
		client         *Form3Client
		mockCtrl       *gomock.Controller
		httpClientMock = NewMockHTTPClient(mockCtrl)
	)

	BeforeEach(func() {
		client = NewForm3Client(httpClientMock)
	})

	Describe("Resource operations", func() {
		Context("Create", func() {
			It("blablabla", func() {
				fmt.Println(">>>> client ", client)
				//httpClient.Mock.Expecs
			})
		})
	})
})
