package accountclient

import (
	"fmt"
	"sort"
	"strings"

	"./resources"
)

var (
	resourcesEndpointsMap = map[resources.ResourceName]string{
		resources.Account: "organisation/accounts",
	}
)

type URLBuilder struct {
	baseURL string
	version string
}

func NewURLBuilder(baseURL, version string) URLBuilder {
	return URLBuilder{
		baseURL: baseURL,
		version: version,
	}
}

func (u URLBuilder) DoForResource(resourceName resources.ResourceName) string {
	endpoint := resourcesEndpointsMap[resourceName]
	return fmt.Sprintf("%s/%s/%s", u.baseURL, u.version, endpoint)
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
