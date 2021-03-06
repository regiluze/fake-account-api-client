package client

import (
	"fmt"
	"sort"
	"strings"

	"github.com/regiluze/form3-account-api-client/resources"
)

var (
	resourcesEndpointsMap = map[resources.ResourceName]string{
		resources.Account: "organisation/accounts",
	}
)

type URLBuilder struct {
	baseURL string
}

func NewURLBuilder(baseURL string) URLBuilder {
	return URLBuilder{
		baseURL: baseURL,
	}
}

func (u URLBuilder) DoForResource(resourceName resources.ResourceName) string {
	endpoint := resourcesEndpointsMap[resourceName]
	return fmt.Sprintf("%s/%s", u.baseURL, endpoint)
}

func (u URLBuilder) DoForResourceWithID(resourceName resources.ResourceName, id string) string {
	resourceEndpoint := u.DoForResource(resourceName)
	return fmt.Sprintf("%s/%s", resourceEndpoint, id)
}

func (u URLBuilder) DoForResourceWithParameters(resourceName resources.ResourceName, parameters map[string]string) string {
	resourceEndpoint := u.DoForResource(resourceName)
	return fmt.Sprintf("%s%s", resourceEndpoint, u.buildQueryParameters(parameters))
}

func (u URLBuilder) DoForResourceWithIDAndParameters(resourceName resources.ResourceName, id string, parameters map[string]string) string {
	resourceEndpoint := u.DoForResourceWithID(resourceName, id)
	return fmt.Sprintf("%s%s", resourceEndpoint, u.buildQueryParameters(parameters))
}

func (u URLBuilder) buildQueryParameters(parameters map[string]string) string {
	flatParams := []string{}
	paramNames := []string{}
	for name := range parameters {
		paramNames = append(paramNames, name)
	}
	sort.Strings(paramNames)
	for _, name := range paramNames {
		flatParams = append(flatParams, fmt.Sprintf("%s=%s", name, parameters[name]))
	}
	return fmt.Sprintf("?%s", strings.Join(flatParams, "&"))
}
