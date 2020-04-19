// +build contract

package accountclient

import (
	"fmt"
	"time"

	"./resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	id             = "ad69e865-4402-3b3b-a0e5-3004ea9cc8dc"
	organisationID = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	version        = 0
	baseURL        = "http://localhost:8080"
	apiVersion     = "v1"
	pageNumber     = 0
	pageSize       = 10
	headerAccept   = ""
	mimeType       = "application/vnd.api+json"
)

var _ = Describe("", func() {
	var (
		apiClient   *Form3Client
		emptyFilter = map[string]interface{}{}
		urlBuilder  URLBuilder
	)

	BeforeEach(func() {
		urlBuilder = NewURLBuilder(baseURL, apiVersion)
		apiClient = NewForm3APIClientWithTimeout(mimeType, urlBuilder, 5*time.Second)
	})

	Describe("Account resource operations", func() {
		Context("Create", func() {
			It("creates an account and return the new account data with links", func() {
				attributes := map[string]interface{}{
					"country":       "GB",
					"base_currency": "GBP",
					"bank_id":       "400300",
					"bank_id_code":  "GBDSC",
					"bic":           "NWBKGB22",
				}

				accountData := resources.NewAccount(id, organisationID, attributes)

				resp, err := apiClient.Create(resources.Account, accountData)
				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(id))
				fmt.Println(">>>> response data ", resp)
			})
			// create for another country
			// create without country
		})
		Context("Fetch", func() {
			It("fetch an account with provided 'id' parameter", func() {

				resp, err := apiClient.Fetch(resources.Account, id)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(id))
				fmt.Println(">>>> response data ", resp)
			})
			// bad request
			// not found
		})
		Context("List", func() {
			It("returns a collection of accounts", func() {

				resp, err := apiClient.List(
					resources.Account,
					emptyFilter,
					pageNumber,
					pageNumber,
				)

				Expect(err).To(BeNil())
				//Expect(resp.Data.ID).To(Equal(id))
				fmt.Println(">>>> response data ", len(resp.Data))
			})
			//
		})
		Context("Delete", func() {
			It("delete an account with provided 'id' parameter", func() {

				err := apiClient.Delete(resources.Account, id, version)

				Expect(err).To(BeNil())
			})
			// not found
			// bad request
		})
	})

})
