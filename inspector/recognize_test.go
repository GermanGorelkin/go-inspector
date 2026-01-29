package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecognizeService_Recognize(t *testing.T) {
	recReq := RecognizeRequest{
		Images:      []int{1, 2, 3},
		ReportTypes: []string{"FACING_COUN", "PLANOGRAM_COMPLIANCE"},
		Visit:       1,
		Webhook:     "webhook_test",
		CountryCode: "RU",
		RetailChain: "Magnit",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		var got RecognizeRequest
		err = json.Unmarshal(b, &got)
		assert.NoError(t, err)
		assert.Equal(t, recReq, got)

		_, err = fmt.Fprintln(w, `{
				"id": 11,
                "images": [1,2,3],
				"scene": "4d8b66992cd841f6922723afe9bd8cf8",
				"reports": {
							"FACING_COUN":22,
							"PLANOGRAM_COMPLIANCE":33
							}
					}`)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{
		Instance: ts.URL,
		APIKey:   "",
	})
	assert.NoError(t, err)

	recRes, err := client.Recognize.Recognize(context.Background(), recReq)
	assert.NoError(t, err)

	want := &RecognizeResponse{
		ID:     11,
		Images: []int{1, 2, 3},
		Scene:  "4d8b66992cd841f6922723afe9bd8cf8",
		Reports: map[string]int{
			"FACING_COUN":          22,
			"PLANOGRAM_COMPLIANCE": 33,
		},
	}

	assert.Equal(t, want, recRes)
}

func TestRecognizeService_RecognitionError(t *testing.T) {
	req := &RecognitionErrorRequest{
		Images:  []int{42},
		SkuId:   777,
		Scene:   "scene_test",
		Message: "wrong sku",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, methodPOST, r.Method)
		assert.Equal(t, "/"+endpointRecognitionError, r.URL.Path)
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		var got RecognitionErrorRequest
		err = json.Unmarshal(b, &got)
		assert.NoError(t, err)
		assert.Equal(t, *req, got)

		_, err = fmt.Fprintln(w, `{"recognition_error_id":99}`)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{Instance: ts.URL, APIKey: ""})
	assert.NoError(t, err)

	resp, err := client.Recognize.RecognitionError(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &RecognitionErrorResponse{RecognitionErrorID: 99}, resp)
}
