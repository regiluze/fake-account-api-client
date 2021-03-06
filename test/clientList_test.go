// +build unit

package test

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

	. "github.com/onsi/gomega"
	. "github.com/regiluze/form3-account-api-client/client"
	"github.com/regiluze/form3-account-api-client/resources"
)

const (
	pageNumber      = 0
	pageSize        = 50
	id2             = "ad27e265-4402-3b3b-a0e3-6664ea9cc8dc"
	organisationID2 = "eb0bd6f9-c3f5-44b2-b644-acd23cdde73c"
)

var _ = Describe("Account api resource client LIST method", func() {
	var (
		client          *Form3Client
		mockCtrl        *gomock.Controller
		httpClientMock  *MockHTTPClient
		emptyFilter     map[string]interface{}
		queryParameters = fmt.Sprintf("?page[number]=%d&page[size]=%d", pageNumber, pageSize)
		expectedURL     = fmt.Sprintf(fmt.Sprintf("%s/organisation/accounts%s", baseURL, queryParameters))
		ctx             = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		httpClientMock = NewMockHTTPClient(mockCtrl)
		client = NewForm3APIClient(baseURL, httpClientMock)
	})

	Context("building request", func() {
		It("builds a request with GET method", func() {
			httpClientMock.EXPECT().Do(IsRequestMethod("GET")).Return(nil, errors.New("fake")).Times(1)

			client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)
		})
		It("builds a request with resource endpoint and page[number] and page[size] parameters", func() {
			httpClientMock.EXPECT().Do(IsRequestURL(expectedURL)).Return(nil, errors.New("fake")).Times(1)

			client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)
		})
	})
	Context("When getting succesful response", func() {
		It("returns ListDataContainer struct as response data", func() {
			account1 := BuildBasicAccountResource(id, organisationID)
			account2 := BuildBasicAccountResource(id2, organisationID2)
			data := resources.ListDataContainer{
				Data: []resources.Resource{
					account1,
					account2,
				},
			}
			dataBt, _ := json.Marshal(data)
			expectedResponseBody := ioutil.NopCloser(bytes.NewReader(dataBt))
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 200,
					Body:       expectedResponseBody,
				},
				nil,
			).Times(1)

			response, err := client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)

			Expect(err).To(BeNil())
			Expect(len(response.Data)).To(Equal(2))
		})
	})
	Context("When something goes wrong", func() {
		It("returns an error when http client return an error", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				nil,
				errors.New("error"),
			).Times(1)

			response, err := client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)

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

			response, err := client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)

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

			response, err := client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					NewErrResponseStatusCode("GET", expectedURL, 500)),
			)
		})
		It("returns ErrNotFound error when server responses an error 40X", func() {
			httpClientMock.EXPECT().Do(gomock.Any()).Return(
				&http.Response{
					StatusCode: 409,
				},
				nil,
			).Times(1)

			response, err := client.List(ctx, resources.Account, emptyFilter, pageNumber, pageSize)

			Expect(response).To(BeNil())
			Expect(err).Should(
				MatchError(
					NewErrResponseStatusCode("GET", expectedURL, 409)),
			)
		})
	})
})
