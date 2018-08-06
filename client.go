package go_inspector

import (
	"net/url"
	"net/http"
	"io"
	"bytes"
	"encoding/json"
	"fmt"
)

type Client struct {
	Instance   *url.URL
	APIKey string

	httpClient *http.Client
}

type ClintConf struct{
	Instance   *url.URL
	APIKey string
}

func NewClient(cfg ClintConf) *Client{
	return &Client{
		APIKey:     cfg.APIKey,
		Instance:   cfg.Instance,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.Instance.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIKey))
	return req, nil
}
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}