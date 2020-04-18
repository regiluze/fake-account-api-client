// +build unit

package accountclient

import (
	"errors"
	"fmt"
	"net/http"

	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"./test"
	. "github.com/onsi/gomega"
)

const (
	version = 1
)

var _ = Describe("Account api resource client DELETE method", func() {
	var (
		client         *Form3Client
		mockCtrl       *gomock.Controller
		httpClientMock *MockHTTPClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClientMock = NewMockHTTPClient(mockCtrl)
		client = NewForm3Client(httpClientMock, baseURL, fakeHeaders)
	})

	Context("building request", func() {
		It("builds a request with DELETE method", func() {
			httpClientMock.EXPECT().Do(test.IsRequestMethod("DELETE")).Return(nil, errors.New("fake")).Times(1)

			client.DeleteAccount(id, version)
		})
		It("builds a request with resource endpoint with resource id and 'version' query parameter", func() {
			httpClientMock.EXPECT().Do(test.IsRequestURL(
				fmt.Sprintf("%s/organisation/accounts/%s?version=%d", baseURL, id, version))).Return(nil, errors.New("fake")).Times(1)

			client.DeleteAccount(id, version)
		})
	})
	Context("When getting succesful response", func() {
		It("returns resourceContainer struct as response data", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 204,
				},
				nil,
			).Times(1)

			err := client.DeleteAccount(id, version)

			Expect(err).To(BeNil())
		})
	})
	Context("When something goes wrong", func() {
		It("returns an error when http client return an error", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				nil,
				errors.New("error"),
			).Times(1)

			err := client.DeleteAccount(id, version)

			Expect(err).NotTo(BeNil())
		})
	})
	Context("When getting error response from the server", func() {
		It("returns an error when server responses an error 50X", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 500,
				},
				nil,
			).Times(1)

			err := client.DeleteAccount(id, version)

			Expect(err).Should(
				MatchError(
					ErrFromServer{"DELETE",
						fmt.Sprintf("%s/organisation/accounts/%s?version=%d", baseURL, id, version),
						500}),
			)
		})
	})
})
