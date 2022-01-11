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

// Client provides IC API Client.
// Contains services for access functions in the IC API
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

// ClintConf holds all of the configuration options for Client
type ClintConf struct {
	Instance string
	APIKey   string
	Verbose  bool
}

// NewClient makes a new Client for IC API.
func NewClient(cfg ClintConf) (*Client, error) {
	cl, err := httpclient.New(
		&http.Client{Timeout: 30 * time.Second},
		httpclient.WithBaseURL(cfg.Instance),
		httpclient.WithAuthorization(fmt.Sprintf("%s %s", "Token", cfg.APIKey)),
		httpclient.WithInterceptor(httpclient.ResponseInterceptor))
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
