package inspector

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient_DefaultsAndServices(t *testing.T) {
	c, err := NewClient(ClientConf{Instance: "https://example.com", APIKey: "abc"})
	assert.NoError(t, err)

	assert.Equal(t, DefaultHTTPTimeout, c.httpTimeout)
	assert.NotNil(t, c.Image)
	assert.NotNil(t, c.Recognize)
	assert.NotNil(t, c.Report)
	assert.NotNil(t, c.Sku)
	assert.NotNil(t, c.Visit)
}

func TestNewClient_UsesCustomHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c, err := NewClient(ClientConf{Instance: "https://example.com", APIKey: "abc", HTTPClient: custom})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, custom.Timeout, c.httpTimeout)
}

func TestNewClient_VerboseAddsInterceptor(t *testing.T) {
	_, err := NewClient(ClientConf{Instance: "https://example.com", APIKey: "abc", Verbose: true})
	assert.NoError(t, err)
}
