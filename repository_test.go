package xld

import (
	"fmt"
	"testing"
)

var mockConfig = Config{
	User:     "admin",
	Password: "admin01",
	Host:     "localhost",
	Port:     "4516",
	Context:  "",
	Scheme:   "http",
}

func TestGetCi(t *testing.T) {
	Client := NewClient(&mockConfig)
	// setup the api route we're going to need for this test

	ci, _ := Client.Repository.GetGeneric("Infrastructure/testhost")
	fmt.Printf("%+v", ci)

}

// func TestCreateCi(t *testing.T) {
// 	Client := NewClient(&mockConfig)
//
// 	props := map[string]interface{}{"entries": map[string]string{"bank": "rabo", "test": "test"}}
//
// 	props["restrictToContainers"] = []string{"Infrastructure/testhost", "Infrastructure/testhost2"}
// 	_, _ = Client.Repository.CreateCi("Environments/Wian2", "udm.Dictionary", props)
//
// }

//helper methods

// func getDictionary() EnvironmentCi {
// 	return EnvironmentCi{
// 		ID:             "Environments/test1",
// 		Type:           "udm.Dictionary",
// 		Token:          "a4061653-a502-4c8f-a3cc-8e663543a028",
// 		CreatedBy:      "admin",
// 		CreatedAt:      "2016-09-26T09:21:17.434+0200",
// 		LastModifiedBy: "admin",
// 		LastModifiedAt: "2016-09-26T09:40:16.429+0200",
// 		Entries: {
// 			"rabo": "bank",
// 		},
// 		EncryptedEntries: {},
// 		RestrictToContainers: {
// 			"Infrastructure/testhost",
// 		},
// 		RestrictToApplications: {
// 			"Applications/testapp1",
// 		},
// 	}
// }
