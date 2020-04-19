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

func (u URLBuilder) Do(resourceName resources.ResourceName, id string, parameters map[string]string) string {
	endpoint := resourcesEndpointsMap[resourceName]
	idPath := ""
	queryParams := ""
	if len(id) > 0 {
		idPath = fmt.Sprintf("/%s", id)
	}
	if len(parameters) > 0 {
		flatParams := []string{}
		paramNames := []string{}
		for name := range parameters {
			paramNames = append(paramNames, name)
		}
		sort.Strings(paramNames)
		for _, name := range paramNames {
			flatParams = append(flatParams, fmt.Sprintf("%s=%s", name, parameters[name]))
		}
		queryParams = fmt.Sprintf("?%s", strings.Join(flatParams, "&"))
	}
	return fmt.Sprintf("%s/%s/%s%s%s", u.baseURL, u.version, endpoint, idPath, queryParams)
}
