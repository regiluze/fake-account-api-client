package test

import (
	"../resources"
)

func BuildBasicAccountResource(id, organisationID string) resources.Resource {
	return resources.NewAccount(id, organisationID, map[string]interface{}{})
}
