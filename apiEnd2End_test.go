// +build contract

package accountclient

import (
	"fmt"
	"net/http"

	"./resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	id             = "ad69e865-4402-3b3b-a0e5-3004ea9cc8dc"
	organisationID = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"
	version        = 0
	baseURL        = "http://localhost:8080/v1"
)

var _ = Describe("", func() {
	var (
		apiClient *Form3Client
	)

	BeforeEach(func() {
		httpClient := http.DefaultClient
		apiClient = NewForm3Client(baseURL, httpClient)
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

				account := resources.NewAccount(id, organisationID, attributes)

				resp, err := apiClient.CreateAccount(account)
				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(id))
				fmt.Println(">>>> response data ", resp)
			})
		})
		Context("Fetch", func() {
			It("fetch an account with provided 'id' parameter", func() {

				resp, err := apiClient.FetchAccount(id)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(id))
				fmt.Println(">>>> response data ", resp)
			})
		})
		Context("Delete", func() {
			It("delete an account with provided 'id' parameter", func() {

				err := apiClient.DeleteAccount(id, version)

				Expect(err).To(BeNil())
			})
		})
	})

})
