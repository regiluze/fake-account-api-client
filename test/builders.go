package test

import (
	"github.com/regiluze/form3-account-api-client/resources"
)

func BuildBasicAccountResource(id, organisationID string) resources.Resource {
	return resources.NewAccount(id, organisationID, map[string]interface{}{})
}

func BuildUKAccountWithoutCoP(id, organisationID string) resources.Resource {
	return resources.NewAccount(id, organisationID, buildUKAccountWithoutCoP())
}

func BuildUKAccountWithCoP(id, organisationID string) resources.Resource {
	return resources.NewAccount(id, organisationID, buildUKAccountWithCoP())
}

func BuildUKSampleAccountWithoutCountry(id, organisationID string) resources.Resource {
	attributes := buildUKAccountWithoutCoP()
	delete(attributes, "country")
	return resources.NewAccount(id, organisationID, attributes)
}

func buildUKAccountWithoutCoP() map[string]interface{} {
	return map[string]interface{}{
		"country":       "GB",
		"base_currency": "GBP",
		"bank_id":       "400300",
		"bank_id_code":  "GBDSC",
		"bic":           "NWBKGB22",
	}
}

func buildUKAccountWithCoP() map[string]interface{} {
	return map[string]interface{}{
		"country":       "GB",
		"base_currency": "GBP",
		"bank_id":       "400300",
		"bank_id_code":  "GBDSC",
		"bic":           "NWBKGB22",
		"name": []string{
			"Samantha Holder",
		},
		"alternative_names": []string{
			"Sam Holder",
		},
		"account_classification":   "Personal",
		"joint_account":            false,
		"account_matching_opt_out": false,
		"secondary_identification": "A1B2C3D4",
	}
}
