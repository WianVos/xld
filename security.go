package xld

const (
	securityBasePath = "deployit/security"
)

type SecurityService interface {
}

//RepositoryServiceOp holds the communication service for Repositorys
type SecurityServiceOp struct {
	client *Client
}

var _ SecurityService = &SecurityServiceOp{}

type User struct {
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	Password string `json:"password,omitempty"`
}

//GetUser returns a user from xld
func (s SecurityServiceOp) GetUser(n string) (User, error) {

	var u User

	url := securityBasePath + "/user/" + n

	req, err := s.client.NewRequest(url, "GET", nil)

	if err != nil {
		return u, err
	}

	resp, err := s.client.Do(req, &u)

	if err != nil {
		return u, err
	}
	defer resp.Body.Close()

	return u, nil
}
