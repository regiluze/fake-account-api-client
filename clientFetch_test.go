// +build unit

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

var _ = Describe("Account api resource client FETCH method", func() {
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
		It("builds a request with GET method", func() {
			httpClientMock.EXPECT().Do(test.IsRequestMethod("GET")).Return(nil, errors.New("fake")).Times(1)

			client.Fetch(resources.Account, id)
		})
		It("builds a request with resource endpoint and resource id", func() {
			httpClientMock.EXPECT().Do(test.IsRequestURL(fmt.Sprintf("%s/organisation/accounts/%s", baseURL, id))).Return(nil, errors.New("fake")).Times(1)

			client.Fetch(resources.Account, id)
		})
	})
	Context("When getting succesful response", func() {
		It("returns DataContainer struct as response data", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(account)
			dataBt, _ := json.Marshal(data)
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader(dataBt))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 200,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.Fetch(resources.Account, id)

			Expect(err).To(BeNil())
			Expect(response.Data.ID).To(Equal(id))
		})
	})
	Context("When something goes wrong", func() {
		It("returns an error when http client return an error", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				nil,
				errors.New("error"),
			).Times(1)

			response, err := client.Fetch(resources.Account, id)

			Expect(response).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
		It("returns an error when response body unmarshal fails", func() {
			json := `{}}`
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader([]byte(json)))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 200,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.Fetch(resources.Account, id)

			Expect(response).To(BeNil())
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

			response, err := client.Fetch(resources.Account, id)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrFromServer{"GET",
						fmt.Sprintf("%s/organisation/accounts/%s", baseURL, id),
						500}),
			)
		})
		It("returns ErrNotFound error when server responses an error 404", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 404,
				},
				nil,
			).Times(1)

			response, err := client.Fetch(resources.Account, id)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrNotFound{fmt.Sprintf("%s/organisation/accounts/%s",
						baseURL,
						id)}),
			)
		})
	})
})
