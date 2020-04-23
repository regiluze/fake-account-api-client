// +build e2e

package test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/regiluze/form3-account-api-client/client"
	"github.com/regiluze/form3-account-api-client/resources"
)

const (
	invalidUUID       = "bd6f8-c1f2-11b2-b677-acd23cdde73c"
	defaultVersion    = 0
	defaultBaseURL    = "http://localhost:8080"
	defaultApiVersion = "v1"
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
		apiClient *Form3Client
		//apiClientWithTimeout *Form3Client
		emptyFilter = map[string]interface{}{}
		ctx         = context.Background()
	)

	BeforeEach(func() {
		apiClient = NewForm3APIClient(baseURL, apiVersion, http.DefaultClient)
	})

	AfterSuite(func() {
	})

	Describe("Account resource operations", func() {
		//Context("Timeout", func() {
		//	It("returns an error when request return status code is Timeout", func() {
		//		apiClientWithTimeout = NewForm3APIClientWithTimeout(mimeType, authToken, urlBuilder, 1*time.Millisecond)

		//		_, err := apiClientWithTimeout.Fetch(ctx, resources.Account, ukAccountID3)

		//		Expect(err).NotTo(BeNil())
		//	})
		//})
		Context("Create", func() {
			It("creates a UK account without CoP (non SEPA Indirect) and return the new account data with links", func() {
				ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())

				resp, err := apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID, ukOrganisationID),
				)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID))
				Expect(resp.Links).NotTo(BeEmpty())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)
				Expect(err).To(BeNil())
			})
			It("creates a UK account with CoP (non SEPA Indirect) and return the new account data with links", func() {
				ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())

				resp, err := apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithCoP(ukAccountID, ukOrganisationID),
				)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID))
				Expect(resp.Links).NotTo(BeEmpty())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)
				Expect(err).To(BeNil())
			})
			Context("unhappy path", func() {
				It("returns ErrBadRequest error when missing country account attributes", func() {
					ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
					Expect(err).To(BeNil())

					resp, err := apiClient.Create(
						ctx,
						resources.Account,
						BuildUKSampleAccountWithoutCountry(ukAccountID, ukOrganisationID),
					)

					Expect(resp).To(BeNil())
					Expect(err).Should(
						MatchError(
							NewErrBadRequest("POST",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "validation failure list:\nvalidation failure list:\nvalidation failure list:\ncountry in body is required",
								},
							)),
					)
				})
			})
			It("returns ErrFromServer error when account id already exists", func() {
				ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				accountData := BuildUKAccountWithCoP(ukAccountID, ukOrganisationID)
				resp, err := apiClient.Create(
					ctx,
					resources.Account,
					accountData,
				)
				Expect(err).To(BeNil())

				resp, err = apiClient.Create(ctx, resources.Account, accountData)

				expectedURL := fmt.Sprintf("%s/%s/organisation/accounts", baseURL, apiVersion)
				Expect(resp).To(BeNil())
				Expect(err).Should(
					MatchError(
						NewErrFromServer(
							"POST",
							expectedURL,
							http.StatusConflict,
						)),
				)
				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)
				Expect(err).To(BeNil())
			})
		})
		Context("Fetch", func() {
			It("fetch an account with provided 'id' parameter", func() {
				ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				resp, err := apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID, ukOrganisationID),
				)

				resp, err = apiClient.Fetch(ctx, resources.Account, ukAccountID)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID))
				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)
				Expect(err).To(BeNil())
			})
			Context("unhappy path", func() {
				It("returns ErrNotFound error when account id not found", func() {
					accountID, _, err := BuildRandomUUIDs()
					Expect(err).To(BeNil())

					resp, err := apiClient.Fetch(ctx, resources.Account, accountID)

					Expect(resp).To(BeNil())
					expectedURL := fmt.Sprintf("%s/%s/organisation/accounts/%s", baseURL, apiVersion, accountID)
					Expect(err).Should(
						MatchError(
							NewErrNotFound(
								expectedURL,
							)),
					)
				})
				It("returns an error when account id is invalid", func() {
					resp, err := apiClient.Fetch(ctx, resources.Account, invalidUUID)

					Expect(resp).To(BeNil())
					Expect(err).Should(
						MatchError(
							NewErrBadRequest("GET",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "id is not a valid uuid",
								},
							)),
					)
				})
			})
		})
		Context("List", func() {
			It("returns a collection of 3 accounts when page size is 3", func() {
				ukAccountID1, ukOrganisationID1, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID1, ukOrganisationID1),
				)
				Expect(err).To(BeNil())
				ukAccountID2, ukOrganisationID2, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID2, ukOrganisationID2),
				)
				Expect(err).To(BeNil())
				ukAccountID3, ukOrganisationID3, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID3, ukOrganisationID3),
				)
				Expect(err).To(BeNil())

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
				Expect(resp.Data[0].ID).To(Equal(ukAccountID1))
				Expect(resp.Data[1].ID).To(Equal(ukAccountID2))
				Expect(resp.Data[2].ID).To(Equal(ukAccountID3))
				err = apiClient.Delete(ctx, resources.Account, ukAccountID1, defaultVersion)
				Expect(err).To(BeNil())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID2, defaultVersion)
				Expect(err).To(BeNil())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID3, defaultVersion)
				Expect(err).To(BeNil())
			})
			It("returns second account when page size is 1 and number 0 when there are 3 accounts", func() {
				ukAccountID1, ukOrganisationID1, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID1, ukOrganisationID1),
				)
				Expect(err).To(BeNil())
				ukAccountID2, ukOrganisationID2, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID2, ukOrganisationID2),
				)
				Expect(err).To(BeNil())
				ukAccountID3, ukOrganisationID3, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithoutCoP(ukAccountID3, ukOrganisationID3),
				)
				Expect(err).To(BeNil())

				pageNumber := 1
				pageSize := 1
				resp, err := apiClient.List(
					ctx,
					resources.Account,
					emptyFilter,
					pageNumber,
					pageSize,
				)

				Expect(err).To(BeNil())
				Expect(resp.Data[0].ID).To(Equal(ukAccountID2))
				err = apiClient.Delete(ctx, resources.Account, ukAccountID1, defaultVersion)
				Expect(err).To(BeNil())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID2, defaultVersion)
				Expect(err).To(BeNil())
				err = apiClient.Delete(ctx, resources.Account, ukAccountID3, defaultVersion)
				Expect(err).To(BeNil())
			})
			Context("unhappy path", func() {
				It("returns ErrFromServer error when page number and size are negative numbers", func() {
					pageNumber := -1
					pageSize := -1
					resp, err := apiClient.List(
						ctx,
						resources.Account,
						emptyFilter,
						pageNumber,
						pageSize,
					)

					queryParameters := fmt.Sprintf("?page[number]=%d&page[size]=%d", pageNumber, pageSize)
					expectedURL := fmt.Sprintf(fmt.Sprintf("%s/%s/organisation/accounts%s", baseURL, apiVersion, queryParameters))
					Expect(resp).To(BeNil())
					Expect(err).Should(
						MatchError(
							NewErrFromServer(
								"GET",
								expectedURL,
								http.StatusInternalServerError,
							)),
					)
				})

			})
		})
		Context("Delete", func() {
			It("delete an account with provided 'id' parameter", func() {
				ukAccountID, ukOrganisationID, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())
				_, err = apiClient.Create(
					ctx,
					resources.Account,
					BuildUKAccountWithCoP(ukAccountID, ukOrganisationID),
				)

				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)

				Expect(err).To(BeNil())
			})
			It("returns nil error when account id not exists", func() {
				ukAccountID, _, err := BuildRandomUUIDs()
				Expect(err).To(BeNil())

				err = apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)

				Expect(err).To(BeNil())
			})
			Context("unhappy path", func() {
				It("returns ErrBadRequest error when account id is invalid", func() {

					err := apiClient.Delete(ctx, resources.Account, invalidUUID, defaultVersion)

					Expect(err).Should(
						MatchError(
							NewErrBadRequest("DELETE",
								resources.BadRequestData{
									ErrorCode:    0,
									ErrorMessage: "id is not a valid uuid",
								},
							)),
					)
				})
			})
		})
	})
})
