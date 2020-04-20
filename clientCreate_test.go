// +build unit

package accountclient

import (
	"bytes"
	"context"
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
		urlBuilder     URLBuilder
		expectedURL    = fmt.Sprintf("%s/%v/organisation/accounts", baseURL, apiVersion)
		ctx            = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClientMock = NewMockHTTPClient(mockCtrl)
		urlBuilder = NewURLBuilder(baseURL, apiVersion)
		client = NewForm3APIClient(fakeMimeType, fakeToken, urlBuilder, httpClientMock)
	})

	Context("Building the request", func() {
		It("builds a request with POST method", func() {
			httpClientMock.EXPECT().Do(test.IsRequestMethod("POST")).Return(nil, errors.New("fake")).Times(1)

			accountData := resources.NewAccount(id, organisationID, map[string]interface{}{})

			client.Create(ctx, resources.Account, accountData)
		})
		It("builds a request with resource endpoint", func() {
			httpClientMock.EXPECT().Do(test.IsRequestURL(expectedURL)).Return(nil, errors.New("fake")).Times(1)

			accountData := resources.NewAccount(id, organisationID, map[string]interface{}{})

			client.Create(ctx, resources.Account, accountData)
		})
		It("builds a request with dataContainer struct data", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataB, _ := json.Marshal(data)
			req, _ := http.NewRequest(
				"POST",
				expectedURL,
				bytes.NewBuffer(dataB),
			)
			httpClientMock.EXPECT().Do(test.IsRequestBody(req)).Return(nil, errors.New("fake")).Times(1)

			client.Create(ctx, resources.Account, accountData)
		})
		It("builds a request with accept and content type values in header", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataB, _ := json.Marshal(data)
			req, _ := http.NewRequest(
				"POST",
				expectedURL,
				bytes.NewBuffer(dataB),
			)
			req.Header.Set("Accept", fakeMimeType)
			req.Header.Set("Content-Type", fakeMimeType)
			httpClientMock.EXPECT().Do(test.IsRequestHeaderValues(req)).Return(nil, errors.New("fake")).Times(1)

			client.Create(ctx, resources.Account, accountData)
		})
		It("changes client headers setting different ones", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataB, _ := json.Marshal(data)
			req, _ := http.NewRequest(
				"POST",
				expectedURL,
				bytes.NewBuffer(dataB),
			)
			req.Header.Set("Accept", anotherFakeMimeType)
			req.Header.Set("Content-Type", anotherFakeMimeType)
			client.SetMimeType(anotherFakeMimeType)

			httpClientMock.EXPECT().Do(test.IsRequestHeaderValues(req)).Return(nil, errors.New("fake")).Times(1)

			client.Create(ctx, resources.Account, accountData)
		})
		It("builds a request with bearen token authorization", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataB, _ := json.Marshal(data)
			req, _ := http.NewRequest(
				"POST",
				expectedURL,
				bytes.NewBuffer(dataB),
			)
			req.Header.Set("Accept", fakeMimeType)
			req.Header.Set("Content-Type", fakeMimeType)
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", fakeToken))
			httpClientMock.EXPECT().Do(test.IsRequestHeaderValues(req)).Return(nil, errors.New("fake")).Times(1)

			client.Create(ctx, resources.Account, accountData)
		})
		It("builds a request with different bearen token authorization when setting another one", func() {
			client.SetAuthToken(anotherFakeToken)
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataB, _ := json.Marshal(data)
			req, _ := http.NewRequest(
				"POST",
				expectedURL,
				bytes.NewBuffer(dataB),
			)
			req.Header.Set("Accept", fakeMimeType)
			req.Header.Set("Content-Type", fakeMimeType)
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", anotherFakeToken))
			httpClientMock.EXPECT().Do(test.IsRequestHeaderValues(req)).Return(nil, errors.New("fake")).Times(1)

			client.Create(ctx, resources.Account, accountData)
		})
	})
	Context("When getting succesful response", func() {
		It("returns resourceContainer struct as response data", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			data := resources.NewDataContainer(accountData)
			dataBt, _ := json.Marshal(data)
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader(dataBt))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 201,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(err).To(BeNil())
			Expect(response.Data.ID).To(Equal(id))
		})
	})
	Context("When something goes wrong", func() {
		It("returns an error when http client return an error", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				nil,
				errors.New("error"),
			).Times(1)

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(response).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
		It("returns an error when response body unmarshal fails", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			json := `{}}`
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader([]byte(json)))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 201,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(response).To(BeNil())
			Expect(err).NotTo(BeNil())
		})
	})
	Context("When getting error response from the server", func() {
		It("returns an error when server responses an error 50X", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 500,
				},
				nil,
			).Times(1)

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrFromServer{"POST",
						expectedURL,
						500}),
			)
		})
		It("returns an error when server responses an error 40X", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 403,
				},
				nil,
			).Times(1)

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrFromServer{"POST",
						expectedURL,
						403}),
			)
		})
		It("returns an error with information to indentify the problem when server responses an error 400", func() {
			accountData := test.BuildBasicAccountResource(id, organisationID)
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

			response, err := client.Create(ctx, resources.Account, accountData)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					ErrBadRequest{"POST",
						errorData}),
			)
		})
	})
})
