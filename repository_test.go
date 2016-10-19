package xld

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCreateGeneric(t *testing.T) {
	setup()
	defer teardown()

	//setup the mock restinterfaces

	mux.HandleFunc("/deployit/repository/ci/Environments/testDictionary1", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, m.method)

		fmt.Fprint(w, mockTestDictionaryResponse)
	})

	mux.HandleFunc("/deployit/metadata/type/udm.Dictionary", func(w http.ResponseWriter, r *http.Request) {
		//testMethod(t, r, m.method)

		fmt.Fprint(w, mockTestDictionaryMetaResponse)
	})

	cases := []struct {
		name        string
		ciType      string
		properties  map[string]interface{}
		expectedCi  Ci
		expectedErr error
	}{
		{
			name:   "Environments/testDictionary1",
			ciType: "udm.Dictionary",
			properties: map[string]interface{}{
				"entries": map[string]string{"test": "test", "bank": "rabo"},
			},
			expectedCi:  getDictionaryCiStruct(),
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		_, err := client.Repository.CreateCi(c.name, c.ciType, c.properties)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

	}
}
func TestGetGeneric(t *testing.T) {
	setup()
	defer teardown()

	//setup mock rest interfaces
	mux.HandleFunc("/deployit/repository/ci/Environments/testDictionary1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, mockTestDictionaryResponse)
	})
	mux.HandleFunc("/deployit/metadata/type/udm.Dictionary", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, mockTestDictionaryMetaResponse)
	})

	//use the GetGeneric function
	acct, err := client.Repository.GetCi("Environments/testDictionary1")
	if err != nil {
		t.Errorf("repository.GetGeneric returned error: %v", err)
	}

	expected, _ := client.Repository.GetCi("Environments/testDictionary1")

	fields := []string{"entries", "encryptedEntries", "restrictToContainers", "restrictToApplications"}

	for _, f := range fields {
		if !reflect.DeepEqual(acct.Properties[f], expected.Properties[f]) {
			t.Errorf("Template.List returned %+v, expected %+v", acct.Properties[f], expected.Properties[f])
		}
	}
}

func TestListCis(t *testing.T) {
	setup()
	defer teardown()
	//setup mock rest interfaces
	mux.HandleFunc("/deployit/repository/query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, mockTestListResponse)
	})

	//use the GetGeneric function
	acct, err := client.Repository.ListCis("Environments")
	if err != nil {
		t.Errorf("repository.ListCis returned error: %v", err)
	}

	expected := getCiListStruct()
	for _, ci := range acct {
		if ciInList(ci, expected) == false {
			t.Errorf("repository.ListCis: expected: %v , found: %v", acct, expected)
		}
	}

}

func ciInList(c CiListEntry, l CiList) bool {
	for _, ci := range l {
		if c.ID == ci.ID && ci.Type == c.Type {
			return true
		}
	}
	return false
}

func TestCiExists(t *testing.T) {
	setup()
	defer teardown()

	// mock answers
	mockTestDictionary1ExistsResponse := `{ "boolean" : true}`
	mockTestDictionary2ExistsResponse := `{ "boolean" : false}`

	//setup mock rest interfaces
	mux.HandleFunc("/deployit/repository/exists/Environments/testDictionary1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, mockTestDictionary1ExistsResponse)
	})
	mux.HandleFunc("/deployit/repository/exists/Environments/testDictionary2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, mockTestDictionary2ExistsResponse)
	})

	cases := []struct {
		searchID     string
		expectedErr  error
		expectedBool bool
	}{
		{
			searchID:     "Environments/testDictionary1",
			expectedBool: true,
			expectedErr:  nil,
		}, {
			searchID:     "Environments/testDictionary2",
			expectedBool: false,
			expectedErr:  nil,
		},
	}

	for _, c := range cases {
		acct, err := client.Repository.CiExists(c.searchID)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}

		if c.expectedBool != acct {
			t.Errorf("Expected %v but got %v", c.expectedBool, acct)
		}
	}

}

// response variables
func getDictionaryCiStruct() Ci {

	return Ci{
		ID:             "Environments/testDictionary1",
		Type:           "udm.Dictionary",
		Token:          "7f5eeb79-73f9-4312-a4d3-0363402c109d",
		CreatedBy:      "admin",
		CreatedAt:      "2016-09-27T09:42:58.212+0200",
		LastModifiedBy: "admin",
		LastModifiedAt: "2016-09-27T09:42:58.212+0200",
		Properties: map[string]interface{}{"restrictToContainers": []string{"Infrastructure/testHost"},
			"restrictToApplications": []string{"Applications/testApp", "Applications/testApp2"},
			"entries":                map[string]string{"test": "test", "bank": "rabo"},
			"encryptedEntries":       map[string]string{"test": "test", "bank": "rabo"},
		},
	}

}

func getCiListStruct() CiList {
	return CiList{CiListEntry{
		ID: "Environments/Wian", Type: "udm.Dictionary",
	}, {
		ID: "Environments/Wian2", Type: "udm.Dictionary",
	}, {
		ID: "Environments/merged_test1_test2", Type: "udm.Dictionary",
	}, {
		ID: "Environments/test1", Type: "udm.Dictionary",
	}, {
		ID: "Environments/test2", Type: "udm.Dictionary",
	}, {
		ID: "Environments/test2env", Type: "udm.Environment",
	}}

}

var mockTestDictionaryResponse = `{
  "id": "Environments/testDictionary1",
  "type": "udm.Dictionary",
  "$token": "7f5eeb79-73f9-4312-a4d3-0363402c109d",
  "$createdBy": "admin",
  "$createdAt": "2016-09-27T09:42:58.212+0200",
  "$lastModifiedBy": "admin",
  "$lastModifiedAt": "2016-09-27T09:42:58.212+0200",
  "entries": {
    "test": "test",
    "bank": "rabo"
  },
  "encryptedEntries": {
	  "test": "test",
	  "bank": "rabo"
	},
  "restrictToContainers": ["Infrastructure/testHost"],
  "restrictToApplications": ["Applications/testApp", "Applications/testApp2"]
}`

var mockTestDictionaryMetaResponse = `{
  "type": "udm.Dictionary",
  "virtual": false,
  "icon": "icons/types/udm.Dictionary.svg",
  "root": "Environments",
  "description": "A Dictionary contains key-value pairs that can be replaced",
  "properties": [
    {
      "name": "entries",
      "fqn": "udm.Dictionary.entries",
      "label": "Entries",
      "kind": "MAP_STRING_STRING",
      "description": "The dictionary entries",
      "category": "Common",
      "asContainment": false,
      "inspection": false,
      "required": false,
      "requiredInspection": false,
      "password": false,
      "transient": false,
      "size": "DEFAULT",
      "referencedType": null,
      "default": null
    },
    {
      "name": "encryptedEntries",
      "fqn": "udm.Dictionary.encryptedEntries",
      "label": "Encrypted Entries",
      "kind": "MAP_STRING_STRING",
      "description": "The encrypted dictionary entries",
      "category": "Common",
      "asContainment": false,
      "inspection": false,
      "required": false,
      "requiredInspection": false,
      "password": true,
      "transient": false,
      "size": "DEFAULT",
      "referencedType": null,
      "default": null
    },
    {
      "name": "restrictToContainers",
      "fqn": "udm.Dictionary.restrictToContainers",
      "label": "Restrict to containers",
      "kind": "SET_OF_CI",
      "description": "Only apply this dictionary to the containers mentioned",
      "category": "Restrictions",
      "asContainment": false,
      "inspection": false,
      "required": false,
      "requiredInspection": false,
      "password": false,
      "transient": false,
      "size": "DEFAULT",
      "referencedType": "udm.Container",
      "default": null
    },
    {
      "name": "restrictToApplications",
      "fqn": "udm.Dictionary.restrictToApplications",
      "label": "Restrict to applications",
      "kind": "SET_OF_CI",
      "description": "Only apply this dictionary to the applications mentioned",
      "category": "Restrictions",
      "asContainment": false,
      "inspection": false,
      "required": false,
      "requiredInspection": false,
      "password": false,
      "transient": false,
      "size": "DEFAULT",
      "referencedType": "udm.Application",
      "default": null
    }
  ],
  "interfaces": [
    "udm.ConfigurationItem"
  ],
  "superTypes": [
    "xld.BaseDictionary",
    "udm.BaseConfigurationItem"
  ]
}`

var mockTestListResponse = `
[{"ref":"Environments/Wian","type":"udm.Dictionary"},
{"ref":"Environments/Wian2","type":"udm.Dictionary"},
{"ref":"Environments/merged_test1_test2","type":"udm.Dictionary"},
{"ref":"Environments/test1","type":"udm.Dictionary"},
{"ref":"Environments/test2","type":"udm.Dictionary"},
{"ref":"Environments/test2env","type":"udm.Environment"}]`

// func TestTemplatesList(t *testing.T) {
// 	setup()
// 	defer teardown()
//
// 	mux.HandleFunc("/api/v1/templates", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		fmt.Fprint(w, mockTemplateListResponse)
// 	})
//
// 	acct, err := client.Templates.List()
// 	if err != nil {
// 		t.Errorf("Templates.List returned error: %v", err)
// 	}
//
// 	expected := getCis()
//
// 	if !reflect.DeepEqual(acct, expected) {
// 		t.Errorf("Template.List returned %+v, expected %+v", acct, expected)
// 	}
// }
//
// func TestTemplatesShow(t *testing.T) {
// 	setup()
// 	defer teardown()
//
// 	// setup the api route we're going to need for this test
// 	mux.HandleFunc("/api/v1/templates", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		fmt.Fprint(w, mockTemplateListResponse)
// 	})
//
// 	cases := []struct {
// 		searchID         string
// 		byTitle          bool
// 		expectedTemplate template.Template
// 		expectedErr      error
// 	}{
// 		{
//
// 			searchID:         "Release6999264",
// 			byTitle:          false,
// 			expectedTemplate: getFullTemplate(),
// 			expectedErr:      nil,
// 		}, {
// 			searchID:         "Applications/Release6999264",
// 			byTitle:          false,
// 			expectedTemplate: getFullTemplate(),
// 			expectedErr:      nil,
// 		}, {
// 			searchID:         "test_template",
// 			byTitle:          true,
// 			expectedTemplate: getFullTemplate(),
// 			expectedErr:      nil,
// 		}, {
// 			searchID:         "bogus",
// 			byTitle:          false,
// 			expectedTemplate: template.Template{},
// 			expectedErr:      errors.New("unable to find template with id: Applications/bogus"),
// 		},
// 	}
//
// 	for _, c := range cases {
// 		acct, err := client.Templates.Show(c.searchID, c.byTitle)
// 		if !reflect.DeepEqual(err, c.expectedErr) {
// 			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
// 		}
//
// 		if c.expectedTemplate.ID != acct.ID {
// 			t.Errorf("Expected %v but got %v", c.expectedTemplate, acct)
// 		}
// 	}
//
// }
//
// func TestTemplateCreate(t *testing.T) {
// 	setup()
// 	defer teardown()
//
// 	mux.HandleFunc("/api/v1/templates", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		fmt.Fprint(w, mockTemplateListResponse)
// 	})
//
// 	mux.HandleFunc("/api/v1/templates/import", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")
//
// 	})
//
// 	cases := []struct {
// 		template         template.Template
// 		overWrite        bool
// 		expectedTemplate template.Template
// 		expectedErr      error
// 	}{
// 		{
// 			template:         getFullTemplate(),
// 			overWrite:        false,
// 			expectedTemplate: template.Template{},
// 			expectedErr:      errors.New("Template with the same title already exists:test_template"),
// 		}, {
// 			template:         getFullTemplate(),
// 			overWrite:        true,
// 			expectedTemplate: getFullTemplate(),
// 			expectedErr:      nil,
// 		},
// 	}
//
// 	for _, c := range cases {
// 		acct, err := client.Templates.CreateTemplate(c.template, c.overWrite)
// 		if !reflect.DeepEqual(err, c.expectedErr) {
// 			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
// 		}
//
// 		if c.expectedTemplate.ID != acct.ID {
// 			t.Errorf("Expected %v but got %v", c.expectedTemplate, acct)
// 		}
// 	}
//
// }
//
// // private helper functions ..
// // usually to keep the code nice, clean and readable
//
// func getFullTemplate() template.Template {
// 	return template.Template{
// 		Ci: ci.Ci{
// 			ID:             "Applications/Release6999264",
// 			Type:           "xlrelease.Release",
// 			Token:          "8198a254-ce39-4075-9581-a65ec2ab72f1",
// 			CreatedBy:      "admin",
// 			CreatedAt:      "2016-08-01T16:26:29.298+0000",
// 			LastModifiedBy: "admin",
// 			LastModifiedAt: "2016-08-05T14:51:55.858+0000",
// 			Title:          "test_template",
// 		},
// 		ScheduledStartDate:                 "2016-08-01T09:00:00Z",
// 		FlagStatus:                         "OK",
// 		MaxConcurrentReleases:              100,
// 		QueryableStartDate:                 "2016-08-01T09:00:00Z",
// 		RealFlagStatus:                     "OK",
// 		Status:                             "TEMPLATE",
// 		CalendarPublished:                  false,
// 		Tutorial:                           false,
// 		AbortOnFailure:                     false,
// 		AllowConcurrentReleasesFromTrigger: false,
// 		RunningTriggeredReleasesCount:      0,
// 		CreatedFromTrigger:                 false,
// 	}
// }
//
// func getCis() []ci.Ci {
// 	return []ci.Ci{
// 		ci.Ci{ID: "Applications/Release6999264",
// 			CreatedAt:      "2016-08-01T16:26:29.298+0000",
// 			CreatedBy:      "admin",
// 			LastModifiedAt: "2016-08-05T14:51:55.858+0000",
// 			LastModifiedBy: "admin",
// 			Token:          "8198a254-ce39-4075-9581-a65ec2ab72f1",
// 			Title:          "test_template",
// 			Type:           "xlrelease.Release"},
// 		ci.Ci{
// 			ID:             "Applications/Release6999266",
// 			CreatedAt:      "2016-08-01T16:26:29.298+0000",
// 			CreatedBy:      "admin",
// 			LastModifiedAt: "2016-08-05T14:51:55.858+0000",
// 			LastModifiedBy: "admin",
// 			Token:          "8198a254-ce39-4075-9581-a65ec2ab72f1",
// 			Title:          "test_template2",
// 			Type:           "xlrelease.Release"},
// 	}
// }
