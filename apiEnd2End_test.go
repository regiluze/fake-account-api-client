// +build e2e

package accountclient

import (
	"context"
	"os"
	"time"

	"./resources"
	"./test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ukAccountID          = "ad27e266-9605-4b4b-a0e5-4443ea8cc4dc"
	ukAccountID2         = "ad69e555-2202-3b3b-a0e5-3004ea9cc8dc"
	ukAccountID3         = "ad96e666-1102-3b3b-a0e5-3004ea9cc8dc"
	ukAccountID4         = "ad95e777-3402-3b3b-a0e5-3004ea9cc8dc"
	ukAccountID5         = "ad95e888-3502-3b3b-a0e5-3004ea9cc8dc"
	ukAccountID6         = "ad95e999-3602-3b3b-a0e5-3004ea9cc8dc"
	ukAccountID7         = "ad93e777-3702-3b3b-a0e5-3004ea9cc8dc"
	ukOrganisationID     = "eb0bd6f5-c2f5-44b2-b677-acd23cdde73c"
	wrongUkAccountID     = "eb0bd6f9-c1f5-11b2-b677-acd23cdde73c"
	notExistsUkAccountID = "eb0bd6f8-c4f1-11b2-b677-acd23cdde73c"
	invalidUUID          = "bd6f8-c1f2-11b2-b677-acd23cdde73c"

	australiaOrganisationID = "eb0bd2f4-c3f5-44b2-b677-acd23cdde73c"
	version                 = 0
	mimeType                = "application/vnd.api+json"
	authToken               = ""
	defaultBaseURL          = "http://localhost:8080"
	defaultApiVersion       = "v1"
)

var (
	baseURL    string
	apiVersion string
)

func init() {
	baseURL = os.Getenv("FORM3_API_BASE_URL")
	if len(baseURL) == 0 {
		baseURL = defaultBaseURL
	}
	apiVersion = os.Getenv("FORM3_API_VERSION")
	if len(apiVersion) == 0 {
		apiVersion = defaultApiVersion
	}
}

var _ = Describe("Account API e2e test suite", func() {
	var (
		apiClient            *Form3Client
		apiClientWithTimeout *Form3Client
		emptyFilter          = map[string]interface{}{}
		urlBuilder           URLBuilder
		ctx                  = context.Background()
	)

	BeforeSuite(func() {
		urlBuilder = NewURLBuilder(baseURL, apiVersion)
		apiClient = NewForm3APIClientWithTimeout(mimeType, authToken, urlBuilder, 5*time.Second)

		// set the repository
		accountData := test.BuildUKAccountWithCoP(ukAccountID3, ukOrganisationID)
		_, err := apiClient.Create(ctx, resources.Account, accountData)
		Expect(err).To(BeNil())
		accountData1 := test.BuildUKAccountWithCoP(ukAccountID4, ukOrganisationID)
		_, err = apiClient.Create(ctx, resources.Account, accountData1)
		Expect(err).To(BeNil())
		accountData2 := test.BuildUKAccountWithCoP(ukAccountID5, ukOrganisationID)
		_, err = apiClient.Create(ctx, resources.Account, accountData2)
		Expect(err).To(BeNil())
		accountData3 := test.BuildUKAccountWithCoP(ukAccountID6, ukOrganisationID)
		_, err = apiClient.Create(ctx, resources.Account, accountData3)
		Expect(err).To(BeNil())
		accountData4 := test.BuildUKAccountWithCoP(ukAccountID7, ukOrganisationID)
		_, err = apiClient.Create(ctx, resources.Account, accountData4)
		Expect(err).To(BeNil())
	})

	AfterSuite(func() {
		// Clean repository
		err := apiClient.Delete(ctx, resources.Account, ukAccountID, version)
		Expect(err).To(BeNil())
		err = apiClient.Delete(ctx, resources.Account, ukAccountID2, version)
		Expect(err).To(BeNil())
		err = apiClient.Delete(ctx, resources.Account, ukAccountID3, version)
		Expect(err).To(BeNil())
		err = apiClient.Delete(ctx, resources.Account, ukAccountID4, version)
		Expect(err).To(BeNil())
		err = apiClient.Delete(ctx, resources.Account, ukAccountID5, version)
		Expect(err).To(BeNil())
		err = apiClient.Delete(ctx, resources.Account, ukAccountID6, version)
		Expect(err).To(BeNil())
	})

	Describe("Account resource operations", func() {
		Context("Timeout", func() {
			It("returns an error when request return status code is Timeout", func() {
				apiClientWithTimeout = NewForm3APIClientWithTimeout(mimeType, authToken, urlBuilder, 1*time.Millisecond)

				_, err := apiClientWithTimeout.Fetch(ctx, resources.Account, ukAccountID3)

				Expect(err).NotTo(BeNil())
			})
		})
		Context("Create", func() {
			It("creates a UK account without CoP (non SEPA Indirect) and return the new account data with links", func() {
				accountData := test.BuildUKAccountWithoutCoP(ukAccountID, ukOrganisationID)

				resp, err := apiClient.Create(ctx, resources.Account, accountData)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID))
				Expect(resp.Links).NotTo(BeEmpty())

			})
			It("creates a UK account with CoP (non SEPA Indirect) and return the new account data with links", func() {
				accountData := test.BuildUKAccountWithCoP(ukAccountID2, ukOrganisationID)

				resp, err := apiClient.Create(ctx, resources.Account, accountData)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID2))
				Expect(resp.Links).NotTo(BeEmpty())

			})
			Context("Bad Request", func() {
				It("returns Bad Request status code when missing country", func() {
					accountData := test.BuildUKSampleAccountWithoutCountry(wrongUkAccountID, ukOrganisationID)

					resp, err := apiClient.Create(ctx, resources.Account, accountData)

					Expect(resp).To(BeNil())
					Expect(err).Should(
						MatchError(
							ErrBadRequest{"POST",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "validation failure list:\nvalidation failure list:\nvalidation failure list:\ncountry in body is required",
								},
							}),
					)
				})
			})
		})
		Context("Fetch", func() {
			It("fetch an account with provided 'id' parameter", func() {

				resp, err := apiClient.Fetch(ctx, resources.Account, ukAccountID3)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID3))

			})
			Context("Error from server", func() {
				It("returns an error when account not found", func() {
					resp, err := apiClient.Fetch(ctx, resources.Account, notExistsUkAccountID)

					Expect(resp).To(BeNil())
					expectedURL := urlBuilder.DoForResourceWithID(resources.Account, notExistsUkAccountID)
					Expect(err).Should(
						MatchError(
							ErrNotFound{
								expectedURL,
							}),
					)
				})
				It("returns an error when account id is invalid", func() {
					resp, err := apiClient.Fetch(ctx, resources.Account, invalidUUID)

					Expect(resp).To(BeNil())
					Expect(err).Should(
						MatchError(
							ErrBadRequest{"GET",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "id is not a valid uuid",
								},
							}),
					)
				})
			})
		})
		Context("List", func() {
			It("returns a collection of accounts", func() {

				pageNumber := 0
				pageSize := 3
				resp, err := apiClient.List(
					ctx,
					resources.Account,
					emptyFilter,
					pageNumber,
					pageSize,
				)

				Expect(err).To(BeNil())
				Expect(len(resp.Data)).To(Equal(pageSize))

			})
		})
		Context("Delete", func() {
			It("delete an account with provided 'id' parameter", func() {

				err := apiClient.Delete(ctx, resources.Account, ukAccountID7, version)

				Expect(err).To(BeNil())
			})
			Context("Error from server", func() {
				It("returns nil error when account id not exists", func() {

					err := apiClient.Delete(ctx, resources.Account, notExistsUkAccountID, version)

					Expect(err).To(BeNil())
				})
				It("returns an error when getting Bad Request status code", func() {

					err := apiClient.Delete(ctx, resources.Account, invalidUUID, version)

					Expect(err).Should(
						MatchError(
							ErrBadRequest{"DELETE",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "id is not a valid uuid",
								},
							}),
					)
				})
			})
		})
	})

})
