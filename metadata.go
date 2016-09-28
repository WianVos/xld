package xld

const (
	MetaDataBasePath = "deployit/metadata"
)

type MetaDataService interface {
	GetProperties(t string) (map[string]string, error)
	GetType(t string) (MetaData, error)
}

//RepositoryServiceOp holds the communication service for Repositorys
type MetaDataServiceOp struct {
	client *Client
}

var _ MetaDataService = &MetaDataServiceOp{}

type MetaDataList []MetaData

type MetaData struct {
	Type           string        `json:"type, omitempty"`
	Virtual        bool          `json:"virtual, omitempty"`
	Icon           string        `json:"icon, omitempty"`
	Root           string        `json:"root,omitempty"`
	Description    string        `json:"description, omitempty"`
	Properties     []Property    `json:"properties, omitempty"`
	Interfaces     []string      `json:"interfaces, omitempty"`
	SuperTypes     []string      `json:"superTypes, omitempty"`
	DeployableType string        `json:"deployableType,omitempty"`
	ContainerType  string        `json:"containerType,omitempty"`
	ControlTasks   []ControlTask `json:"control-tasks,omitempty"`
}

type Property struct {
	Name               string      `json:"name, omitempty"`
	Fqn                string      `json:"fqn, omitempty"`
	Label              string      `json:"label, omitempty"`
	Kind               string      `json:"kind, omitempty"`
	Description        string      `json:"description, omitempty"`
	Category           string      `json:"category, omitempty"`
	AsContainment      bool        `json:"asContainment, omitempty"`
	Inspection         bool        `json:"inspection, omitempty"`
	Required           bool        `json:"required, omitempty"`
	RequiredInspection bool        `json:"requiredInspection, omitempty"`
	Password           bool        `json:"password, omitempty"`
	Transient          bool        `json:"transient, omitempty"`
	Size               string      `json:"size, omitempty"`
	ReferencedType     string      `json:"referencedType, omitempty"`
	Default            interface{} `json:"default, omitempty"`
}

type ControlTask struct {
	Name        string `json:"name, omitempty"`
	Fqn         string `json:"fqn, omitempty"`
	Description string `json:"description, omitempty"`
	Label       string `json:"label, omitempty"`
}

//GetType retrieve MetaData
func (m MetaDataServiceOp) GetType(t string) (MetaData, error) {

	var meta MetaData

	url := MetaDataBasePath + "/" + "type" + "/" + t

	req, err := m.client.NewRequest(url, "GET", nil)

	_, err = m.client.Do(req, &meta)

	return meta, err

}

func (m MetaDataServiceOp) GetProperties(t string) (map[string]string, error) {

	p := make(map[string]string)

	d, err := m.GetType(t)
	if err != nil {
		return p, err
	}

	for _, v := range d.Properties {

		p[v.Name] = v.Kind

	}

	return p, nil

}
