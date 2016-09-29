package xld

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux *http.ServeMux

	client *Client

	server *httptest.Server
)

var mockConfig = Config{
	User:     "admin",
	Password: "password",
	Host:     "localhost",
	Port:     "4516",
	Context:  "",
	Scheme:   "http",
}

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(&mockConfig)

	url, _ := url.Parse(server.URL)

	client.BaseURL = url

}

func teardown() {
	server.Close()
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Repository", "Meta",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		if cv.FieldByName(s).IsNil() {
			t.Errorf("c.%s shouldn't be nil", s)
		}
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if c.UserAgent != userAgent {
		t.Errorf("NewClick UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientServices(t, c)
	testClientDefaultUserAgent(t, c)
}

func TestNewClient(t *testing.T) {

	c := NewClient(&mockConfig)
	testClientDefaults(t, c)
}

func TestNew(t *testing.T) {
	c := New(&mockConfig)

	testClientDefaults(t, c)
}

func TestNewRequest(t *testing.T) {
	c := NewClient(&mockConfig)

	inURL, outURL := "/foo", c.BaseURL.String()+"/foo"

	req, _ := c.NewRequest(inURL, "GET", nil)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("/", "GET", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}
