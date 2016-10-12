package xld

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	repositoryBasePath  = "deployit/repository"
	environmentCiPrefix = "Environments"
)

//RepositoryService is an interface representing the repository service
type RepositoryService interface {
	//GetDictionary(n string) (DictionaryCI, error)
	//GetGeneric(n string)
	CreateCi(n string, t string, p map[string]interface{}) (Ci, error)
	GetCi(n string) (Ci, error)
	CiExists(n string) (bool, error)
}

//RepositoryServiceOp holds the communication service for Repositorys
type RepositoryServiceOp struct {
	client *Client
}

var _ RepositoryService = &RepositoryServiceOp{}

//Ci representation of a xldeploy configuration item
type Ci struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Token          string `json:"$token,omitempty"`
	CreatedBy      string `json:"$createdBy,omitempty"`
	CreatedAt      string `json:"$createdAt,omitempty"`
	LastModifiedBy string `json:"$lastModifiedBy,omitempty"`
	LastModifiedAt string `json:"$lastModifiedAt,omitempty"`
	Properties     map[string]interface{}
}

type ciTrue struct {
	Exists bool `json:"boolean"`
}

//GetGeneric retrieves a CI form xld
func (r RepositoryServiceOp) GetCi(n string) (Ci, error) {

	var e map[string]interface{}
	ri := make(map[string]interface{})

	var err error
	var c Ci

	url := repositoryBasePath + "/" + "ci" + "/" + n

	req, err := r.client.NewRequest(url, "GET", nil)

	resp, err := r.client.Do(req, &e)
	defer resp.Body.Close()

	//Pull out the ci parts
	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(e)
	if err != nil {
		return c, err
	}

	err = json.NewDecoder(data).Decode(&c)
	if err != nil {
		return c, err
	}

	// handle properties
	//get property metadata for intended type
	properties, _ := r.client.Meta.GetProperties(c.Type)

	// loop over the properties and check if they where in the requested ci
	for k := range properties {
		if val, ok := e[k]; ok {
			ri[k] = val
		}
	}

	c.Properties = ri

	return c, nil
}

//CreateCi  creates/updates a CI
// n: name
// t: type
// p: properties
func (r RepositoryServiceOp) CreateCi(n string, t string, p map[string]interface{}) (Ci, error) {

	ci := make(map[string]interface{})
	var dc Ci
	var verb string

	// validate the id: it needs to contain either Environments, Infrastructure, Applications
	_, err := validateID(n)
	if err != nil {
		return dc, err
	}

	ci["id"] = n
	ci["type"] = t

	//get metadata for intended type
	metaData, _ := r.client.Meta.GetProperties(t)

	//validate Properties
	//loop over the metadata and see if the properties we got handed are actually the right type
	// it they are the right type put them in the final map
	for k, v := range p {
		propType := metaData[k]
		switch v := v.(type) {
		default:
			fmt.Printf("unexpected type %T\n", v) // %T prints whatever type t has
		case string:
			if propType == "STRING" || propType == "CI" {
				ci[k] = v
			}
		case bool:
			if propType == "BOOLEAN" {
				ci[k] = v
			}
		case int:
			if propType == "INTEGER" {
				ci[k] = int(v)
			}
		case map[string]interface{}, map[string]string:
			if propType == "MAP_STRING_STRING" {
				ci[k] = v
			}
		case []string:
			if propType == "SET_OF_STRING" || propType == "SET_OF_CI" {
				ci[k] = v
			}
		}

	}

	//marshall the json and send it
	url := repositoryBasePath + "/ci/" + n

	exists, _ := r.CiExists(n)

	if exists == true {
		verb = "PUT"
	} else {
		verb = "POST"
	}

	req, err := r.client.NewRequest(url, verb, ci)
	if err != nil {
		return dc, err
	}

	_, err = r.client.Do(req, &dc)
	if err != nil {
		return dc, err
	}

	return dc, nil

}

//CiExists checks if a CI exists
func (r RepositoryServiceOp) CiExists(n string) (bool, error) {

	var e ciTrue

	_, err := validateID(n)

	if err != nil {
		return false, err
	}

	url := repositoryBasePath + "/" + "exists" + "/" + n

	req, err := r.client.NewRequest(url, "GET", nil)

	resp, err := r.client.Do(req, &e)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return e.Exists, nil

}

//private functions

func validateID(i string) (bool, error) {
	validPrefix := [3]string{"Environments", "Infrastructure", "Applications"}

	for _, p := range validPrefix {
		if strings.Contains(i, p) {
			return true, nil
		}
	}

	return false, errors.New("invalid ci id")
}
