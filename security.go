package xld

import (
	"errors"
	"fmt"
)

const (
	securityBasePath = "deployit/security"
)

//SecurityService represents the service for engaging the XL-Deploy security rest interface
type SecurityService interface {
	GetUser(n string) (User, error)
	UserExists(n string) bool
	CreateUser(n string, a bool) (User, error)
	SetPasswordForUser(n, p string) error
}

//SecurityServiceOp holds the communication service for the Security rest api
type SecurityServiceOp struct {
	client *Client
}

var _ SecurityService = &SecurityServiceOp{}

//User is a user in the xl-deploy repository
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

//UserExists check if a user exists in the XL-Deploy Repository
// returns true if it does
// returns false if anything goes wrong (!!!)
// only applies to the internal xldeploy repository
func (s SecurityServiceOp) UserExists(n string) bool {
	_, err := s.GetUser(n)
	if err != nil {
		return false
	}

	return true
}

//CreateUser creates a user in the XL-Deploy repository
// n is the name of the user
// a signified if the user should be admin
func (s SecurityServiceOp) CreateUser(n string, a bool) (User, error) {
	var u User

	// check if the user already exists. If so return an error saying just that
	if s.UserExists(n) {
		return u, errors.New("user already exists")
	}

	// if we made it this far it is time to set up the user
	u.Username = n
	u.Admin = a

	url := securityBasePath + "/user/" + n

	req, err := s.client.NewRequest(url, "POST", u)

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

//SetPasswordForUser updates an already existing user with a password
// this can be setting the password for the first time, or setting a new one
func (s SecurityServiceOp) SetPasswordForUser(n, p string) error {
	var u User

	if s.UserExists(n) == false {
		return errors.New("user does not exists, unable to set password")
	}

	u, _ = s.GetUser(n)

	u.Password = p
	url := securityBasePath + "/user/" + n

	fmt.Printf("%+v\n", u)
	fmt.Println(url)
	req, err := s.client.NewRequest(url, "PUT", u)

	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, &u)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil

}
