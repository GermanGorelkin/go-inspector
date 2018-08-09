package inspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Pagination struct {
	Count    int                      `json:"count"`
	Next     *string                  `json:"next,omitempty"`
	Previous *string                  `json:"previous,omitempty"`
	Results  []map[string]interface{} `json:"results"`
}

type Client struct {
	Instance   *url.URL
	APIKey     string
	httpClient *http.Client

	Image     *ImageService
	Recognize *RecognizeService
	Report    *ReportService
	Sku       *SkuService
}

type ClintConf struct {
	Instance *url.URL
	APIKey   string
}

func NewClient(cfg ClintConf) *Client {
	c := &Client{
		APIKey:     cfg.APIKey,
		Instance:   cfg.Instance,
		httpClient: http.DefaultClient,
	}
	c.Image = &ImageService{client: c}
	c.Recognize = &RecognizeService{client: c}
	c.Report = &ReportService{client: c}
	c.Sku = &SkuService{client: c}

	return c
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

	//dump, err := httputil.DumpRequest(req, true)
	//fmt.Println(string(dump))

	return req, nil
}

func (c *Client) newRequestFormFile(path string, r io.Reader, filename string) (*http.Request, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("datafile", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fileWriter, r)
	if err != nil {
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	rel := &url.URL{Path: path}
	u := c.Instance.ResolveReference(rel)

	req, err := http.NewRequest("POST", u.String(), bodyBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIKey))
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//dump, err := httputil.DumpResponse(resp, true)
	//fmt.Println(string(dump))

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
