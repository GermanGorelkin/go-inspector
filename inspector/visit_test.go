package inspector

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitService_AddVisit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, methodPOST, r.Method)
		assert.Equal(t, "/"+endpointVisits, r.URL.Path)
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{}`, string(body))

		_, err = fmt.Fprint(w, `{"id":999}`)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{Instance: ts.URL, APIKey: ""})
	assert.NoError(t, err)

	resp, err := client.Visit.AddVisit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 999, resp.ID)
}
