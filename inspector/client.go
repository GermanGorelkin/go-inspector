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
	Instance    string
	APIKey      string
	httpClient  *httpclient.Client
	httpTimeout time.Duration

	Image     *ImageService
	Recognize *RecognizeService
	Report    *ReportService
	Sku       *SkuService
	Visit     *VisitService
}

// ClientConf holds all of the configuration options for Client.
type ClientConf struct {
	Instance   string
	APIKey     string
	Verbose    bool
	HTTPClient *http.Client
	Timeout    time.Duration
}

// ClintConf is kept for backward compatibility with the historical typo.
type ClintConf = ClientConf

// NewClient makes a new Client for IC API.
func NewClient(cfg ClientConf) (*Client, error) {
	var httpc *http.Client
	if cfg.HTTPClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		httpc = &http.Client{Timeout: timeout}
	} else {
		httpc = cfg.HTTPClient
	}

	cl, err := httpclient.New(
		httpc,
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
		APIKey:      cfg.APIKey,
		Instance:    cfg.Instance,
		httpClient:  cl,
		httpTimeout: httpc.Timeout,
	}
	c.Image = &ImageService{client: c}
	c.Recognize = &RecognizeService{client: c}
	c.Report = &ReportService{client: c}
	c.Sku = &SkuService{client: c}
	c.Visit = &VisitService{client: c}

	return c, nil
}
