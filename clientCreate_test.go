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

var _ = Describe("Account api resource client CREATE method", func() {
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

	Context("Building the request", func() {
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
	})
	Context("When getting succesful response", func() {
		// TODO status code
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
	Context("When something goes wrong", func() {
		It("returns an error when http client return an error", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				nil,
				errors.New("error"),
			).Times(1)

			response, err := client.CreateAccount(account)

			Expect(response).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
		It("returns an error when response body unmarshal fails", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			json := `{}}`
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader([]byte(json)))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 201,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.CreateAccount(account)

			Expect(response).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})
	Context("When getting error response from the server", func() {
		It("returns an error when server responses an error 50X", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 500,
				},
				nil,
			).Times(1)

			response, err := client.CreateAccount(account)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrFromServer{"POST",
						fmt.Sprintf("%s/organisation/accounts", baseURL),
						500}),
			)
		})
		It("returns an error when server responses an error 40X", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 403,
				},
				nil,
			).Times(1)

			response, err := client.CreateAccount(account)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrFromServer{"POST",
						fmt.Sprintf("%s/organisation/accounts", baseURL),
						403}),
			)
		})
		It("returns an error with information to indentify the problem when server responses an error 400", func() {
			account := test.BuildBasicAccountResource(id, organisationID)
			errorData := resources.BadRequestData{
				ErrorCode:    400,
				ErrorMessage: "mandatory",
			}
			json := `{"error_code": 400, "error_message": "mandatory"}`
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader([]byte(json)))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 400,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.CreateAccount(account)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrBadRequest{"POST",
						errorData}),
			)
		})
	})
})
