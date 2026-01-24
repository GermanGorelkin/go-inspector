package inspector

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVisitService_AddVisit(t *testing.T) {
	started := time.Date(2026, 1, 24, 10, 0, 0, 0, time.UTC)
	request := &Visit{Shop: 321, Agent: "John", StartedDate: started, Latitude: 1.23, Longitude: 4.56}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/visits/", r.URL.Path)
		_, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		_, err = fmt.Fprintf(w, `{"id":999,"shop":321,"agent":"John","started_date":"%s"}`, started.Format(time.RFC3339Nano))
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{Instance: ts.URL, APIKey: ""})
	assert.NoError(t, err)

	resp, err := client.Visit.AddVisit(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 999, resp.ID)
	assert.Equal(t, request.Shop, resp.Shop)
	assert.Equal(t, request.Agent, resp.Agent)
	assert.Equal(t, started, resp.StartedDate)
}
