package resources

type ResourceName string

const (
	Account ResourceName = "account"
)

type DataContainer struct {
	Data  Resource          `json:"data"`
	Links map[string]string `json:"links,omitempty"`
}

type ListDataContainer struct {
	Data []Resource `json:"data"`
}

type Resource struct {
	ResourceType   string                 `json:"type"`
	ID             string                 `json:"id"`
	Version        int                    `json:"version"`
	OrganisationID string                 `json:"organisation_id"`
	Attributes     map[string]interface{} `json:"attributes"`
	CreatedOn      string                 `json:"created_on:omitempty"`
	ModifiedOn     string                 `json:"modified_on,omitempty"`
	Relationships  map[string]interface{} `json:"relationships"`
}

type BadRequestData struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func NewDataContainer(resource Resource) DataContainer {
	return DataContainer{
		Data: resource,
	}
}

func NewAccount(id, organisationId string, attributes map[string]interface{}) Resource {
	return Resource{
		ResourceType:   "accounts",
		ID:             id,
		OrganisationID: organisationId,
		Attributes:     attributes,
	}
}
