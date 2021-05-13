package inspector

import (
	"fmt"
	"net/http"
	"time"

	httpclient "github.com/germangorelkin/http-client"
)

type Pagination struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next,omitempty"`
	Previous *string     `json:"previous,omitempty"`
	Results  interface{} `json:"results"`
}

type Client struct {
	Instance   string
	APIKey     string
	httpClient *httpclient.Client

	Image     *ImageService
	Recognize *RecognizeService
	Report    *ReportService
	Sku       *SkuService
	Visit     *VisitService
}

type ClintConf struct {
	Instance string
	APIKey   string
	Verbose  bool
}

func NewClient(cfg ClintConf) (*Client, error) {
	cl, err := httpclient.New(
		&http.Client{Timeout: 30 * time.Second},
		httpclient.SetBaseURL(cfg.Instance),
		httpclient.SetAuthorization(cfg.APIKey, "Token"),
		httpclient.SetInterceptor(httpclient.ResponseInterceptor))
	if err != nil {
		return nil, fmt.Errorf("failed to build http-client:%w", err)
	}
	if cfg.Verbose {
		if err := cl.AddInterceptor(httpclient.DumpInterceptor); err != nil {
			return nil, err
		}
	}

	c := &Client{
		APIKey:     cfg.APIKey,
		Instance:   cfg.Instance,
		httpClient: cl,
	}
	c.Image = &ImageService{client: c}
	c.Recognize = &RecognizeService{client: c}
	c.Report = &ReportService{client: c}
	c.Sku = &SkuService{client: c}
	c.Visit = &VisitService{client: c}

	return c, nil
}
