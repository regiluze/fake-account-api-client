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
	invalidUUID    = "bd6f8-c1f2-11b2-b677-acd23cdde73c"
	defaultVersion = 0
	defaultBaseURL = "http://localhost:8080/v1"
)

var (
	baseURL string
)

func init() {
	baseURL = os.Getenv("FORM3_API_BASE_URL")
	if len(baseURL) == 0 {
		baseURL = defaultBaseURL
	}
}

var _ = Describe("Account API e2e test suite", func() {
	var (
		apiClient   *Form3Client
		emptyFilter = map[string]interface{}{}
		ctx         = context.Background()
	)

	BeforeEach(func() {
		apiClient = NewForm3APIClient(baseURL, http.DefaultClient)
	})

	AfterSuite(func() {
	})

	Describe("Account resource operations", func() {
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
				defer removeResources(ctx, apiClient, ukAccountID)
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
				defer removeResources(ctx, apiClient, ukAccountID)
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
			It("returns ErrResponseStatusCode error when account id already exists", func() {
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

				expectedURL := fmt.Sprintf("%s/organisation/accounts", baseURL)
				Expect(resp).To(BeNil())
				Expect(err).Should(
					MatchError(
						NewErrFromServer(
							"POST",
							expectedURL,
							http.StatusConflict,
						)),
				)
				defer removeResources(ctx, apiClient, ukAccountID)
			})
		})
		Context("Fetch", func() {
			It("fetch an account with provided 'id' parameter", func() {
				ukAccountID := addResource(ctx, apiClient)

				resp, err := apiClient.Fetch(ctx, resources.Account, ukAccountID)

				Expect(err).To(BeNil())
				Expect(resp.Data.ID).To(Equal(ukAccountID))
				defer removeResources(ctx, apiClient, ukAccountID)
			})
			Context("unhappy path", func() {
				It("returns ErrNotFound error when account id not found", func() {
					accountID, _, err := BuildRandomUUIDs()
					Expect(err).To(BeNil())

					resp, err := apiClient.Fetch(ctx, resources.Account, accountID)

					Expect(resp).To(BeNil())
					expectedURL := fmt.Sprintf("%s/organisation/accounts/%s", baseURL, accountID)
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
				ukAccountID1 := addResource(ctx, apiClient)
				ukAccountID2 := addResource(ctx, apiClient)
				ukAccountID3 := addResource(ctx, apiClient)

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
				defer removeResources(ctx, apiClient, ukAccountID1, ukAccountID2, ukAccountID3)
			})
			It("returns second account when page size is 1 and number 0 when there are 3 accounts", func() {
				ukAccountID1 := addResource(ctx, apiClient)
				ukAccountID2 := addResource(ctx, apiClient)
				ukAccountID3 := addResource(ctx, apiClient)

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
				defer removeResources(ctx, apiClient, ukAccountID1, ukAccountID2, ukAccountID3)
			})
			Context("unhappy path", func() {
				It("returns ErrResponseStatusCode error when page number and size are negative numbers", func() {
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
					expectedURL := fmt.Sprintf(fmt.Sprintf("%s/organisation/accounts%s", baseURL, queryParameters))
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
				ukAccountID := addResource(ctx, apiClient)

				err := apiClient.Delete(ctx, resources.Account, ukAccountID, defaultVersion)

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

func addResource(ctx context.Context, apiClient *Form3Client) string {
	accountID, organisationID, err := BuildRandomUUIDs()
	Expect(err).To(BeNil())
	_, err = apiClient.Create(
		ctx,
		resources.Account,
		BuildUKAccountWithoutCoP(accountID, organisationID),
	)
	Expect(err).To(BeNil())
	return accountID
}

func removeResources(ctx context.Context, apiClient *Form3Client, ids ...string) {
	for _, id := range ids {
		err := apiClient.Delete(ctx, resources.Account, id, defaultVersion)
		Expect(err).To(BeNil())
	}
}
