package xld

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "v1"
	basePath       = ""
	userAgent      = "goxld"
	mediaType      = "application/json"
	format         = "json"
)

//Config holds the configuration for the xlrelease server configuration
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Context  string
	Scheme   string
}

//Client holds all the settings needed to communicate with xl-release
type Client struct {
	// HTTP client used to communicate with the Veracross API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Client Config
	Config *Config

	Repository RepositoryService
	Meta       MetaDataService
}

//NewClient returns a new functional client struct
func NewClient(config *Config) *Client {
	// create the base url out of the stuff given
	var baseURL url.URL
	finalHost := config.Host + ":" + config.Port
	baseURL.Host = finalHost
	baseURL.Path = basePath
	baseURL.Scheme = config.Scheme

	c := &Client{client: http.DefaultClient, BaseURL: &baseURL, UserAgent: userAgent, Config: config}

	c.Repository = &RepositoryServiceOp{client: c}
	c.Meta = &MetaDataServiceOp{client: c}

	return c
}

//New returns a new XLR API client instance.
// This function is here for api completeness and just passes through to NewClient for now
func New(config *Config) *Client {
	c := NewClient(config)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(urlStr string, method string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)

	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Config.User, c.Config.Password)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	fmt.Printf("%+v\n", v)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}

	return resp, err
}
