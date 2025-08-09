package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	apiKey     string
}

type Config struct {
	HTTPClient *http.Client
	BaseURL    string
	Username   string
	Password   string
	APIKey     string
}

func NewClient(config *Config) *Client {
	return &Client{
		httpClient: config.HTTPClient,
		baseURL:    config.BaseURL,
		username:   config.Username,
		password:   config.Password,
		apiKey:     config.APIKey,
	}
}

type Request struct {
	Method string
	Path   string
	Body   interface{}
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (c *Client) Do(req *Request) (*Response, error) {
	url := c.baseURL + req.Path

	var reqBody io.Reader
	var hasBody bool
	if req.Body != nil {
		// Check if it's a nil map (which would marshal to "null")
		if m, ok := req.Body.(map[string]interface{}); ok && m == nil {
			// Don't send anything for nil maps
			hasBody = false
		} else {
			bodyBytes, err := json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(bodyBytes)
			hasBody = true
		}
	}

	httpReq, err := http.NewRequest(req.Method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Only set Content-Type when there's actually a body to send
	if hasBody {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "ApiKey "+c.apiKey)
	} else if c.username != "" && c.password != "" {
		httpReq.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) Get(path string) (*Response, error) {
	return c.Do(&Request{
		Method: "GET",
		Path:   path,
	})
}

func (c *Client) Post(path string, body interface{}) (*Response, error) {
	return c.Do(&Request{
		Method: "POST",
		Path:   path,
		Body:   body,
	})
}

func (c *Client) Put(path string, body interface{}) (*Response, error) {
	return c.Do(&Request{
		Method: "PUT",
		Path:   path,
		Body:   body,
	})
}

func (c *Client) Delete(path string) (*Response, error) {
	return c.Do(&Request{
		Method: "DELETE",
		Path:   path,
	})
}