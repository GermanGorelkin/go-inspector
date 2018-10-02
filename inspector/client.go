package inspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Pagination struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next,omitempty"`
	Previous *string     `json:"previous,omitempty"`
	Results  interface{} `json:"results"`
}

type Client struct {
	Instance   *url.URL
	APIKey     string
	httpClient *http.Client

	Image     *ImageService
	Recognize *RecognizeService
	Report    *ReportService
	Sku       *SkuService
	Visit     *VisitService
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
	c.Visit = &VisitService{client: c}

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
			err = errors.WithStack(err)
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIKey))

	// dump request for debug
	dump, err := httputil.DumpRequest(req, true)
	if err == nil {
		logrus.Debugln(string(dump))
	}
	// ----------------

	return req, nil
}

func (c *Client) newRequestFormFile(path string, r io.Reader, filename string) (*http.Request, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("datafile", filename)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	_, err = io.Copy(fileWriter, r)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	rel := &url.URL{Path: path}
	u := c.Instance.ResolveReference(rel)

	req, err := http.NewRequest("POST", u.String(), bodyBuf)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIKey))

	// dump request for debug
	dump, err := httputil.DumpRequest(req, true)
	if err == nil {
		logrus.Debugln(string(dump))
	}
	// --------------------

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	defer resp.Body.Close()

	// dump response for debug
	dump, err := httputil.DumpResponse(resp, true)
	if err == nil {
		logrus.Debugln(string(dump))
	}
	// -------------------------

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			err = errors.WithStack(err)
			return resp, err
		}
	} else {
		return resp, errors.New(fmt.Sprintf("Response Status %s\n", resp.Status))
	}

	return resp, nil
}
