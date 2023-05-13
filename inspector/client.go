package inspector

import (
	"crypto/tls"
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
		GetHTTPClient(),
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

// GetHTTPClient
func GetHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
			},
		},
	}
}
