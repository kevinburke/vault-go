// Package vault is a client around the Vault v1 API.
//
// Operations center around a Client. Create a client by specifying a base URL,
// a token, and (optionally) a http.Client instance:
//
//     client := vault.NewClient("https://localhost:8200", "my-token", nil)
//
// Once you have a client, you can perform various operations.
//
//     client.Transit.Create(context.TODO(), "key", map[string]interface{}{
//        "algorithm": "ed25519",
//     })
//
// POST operations take all required parameters as required and optional
// parameters as map[string]interface{} inputs. Consult the API documentation to
// determine what optional parameters are available for each endpoint.
package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	types "github.com/kevinburke/go-types"
	"github.com/kevinburke/rest"
)

// The client Version.
const Version = "0.1"

// The server Version.
const APIVersion = "v1"
const userAgent = "vault-go/" + Version

const defaultTimeout = 7 * time.Second

var defaultHttpClient *http.Client

func init() {
	defaultHttpClient = &http.Client{
		Timeout:   defaultTimeout,
		Transport: rest.DefaultTransport,
	}
}

// Client makes requests to Vault.
type Client struct {
	*rest.Client
	token string

	Transit *TransitService
}

// NewClient creates a new Client.
func NewClient(baseURL, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = defaultHttpClient
	}
	restClient := rest.NewClient("", "", baseURL)
	restClient.Client = httpClient
	restClient.UploadType = rest.JSON
	c := &Client{Client: restClient, token: token}
	c.Transit = &TransitService{
		client: c,
	}
	return c
}

// GetResource retrieves an instance resource with the given path part (e.g.
// "/Messages") and sid (e.g. "MM123").
func (c *Client) GetResource(ctx context.Context, pathPart string, sid string, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "GET", sidPart, nil, v)
}

// CreateResource makes a POST request to the given resource.
func (c *Client) CreateResource(ctx context.Context, pathPart string, data map[string]interface{}, v interface{}) error {
	return c.MakeRequest(ctx, "POST", pathPart, data, v)
}

func (c *Client) UpdateResource(ctx context.Context, pathPart string, sid string, data map[string]interface{}, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "POST", sidPart, data, v)
}

func (c *Client) DeleteResource(ctx context.Context, pathPart string, sid string) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	err := c.MakeRequest(ctx, "DELETE", sidPart, nil, nil)
	if err == nil {
		return nil
	}
	rerr, ok := err.(*rest.Error)
	if ok && rerr.StatusCode == http.StatusNotFound {
		return nil
	}
	return err
}

func (c *Client) ListResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "GET", pathPart, data, v)
}

// Make a request to the Vault API. pathPart is presumed to begin with a slash.
func (c *Client) MakeRequest(ctx context.Context, method string, pathPart string, data interface{}, v interface{}) error {
	var r io.Reader
	if data != nil && (method == "POST" || method == "PUT") {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(data)
		if err != nil {
			return err
		}
		r = b
	}
	if method == "GET" && data != nil {
		udata, ok := data.(url.Values)
		if ok {
			pathPart = pathPart + "?" + udata.Encode()
		}
	}
	req, err := c.NewRequest(method, "/"+APIVersion+pathPart, r)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", c.token)
	req = req.WithContext(ctx)
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", userAgent+" "+ua)
	}
	return c.Do(req, &v)
}

type Response struct {
	Data      interface{}      `json:"data"`
	RequestID types.PrefixUUID `json:"request_id"`
	Renewable bool             `json:"renewable"`
	LeaseID   string           `json:"lease_id"`
	WrapInfo  json.RawMessage  `json:"wrap_info"`
	Warnings  json.RawMessage  `json:"warnings"`
	Auth      json.RawMessage  `json:"auth"`
}
