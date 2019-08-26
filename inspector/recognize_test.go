package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRecognizeService_Recognize(t *testing.T) {
	recReq := RecognizeRequest{
		Images:      []int{1, 2, 3},
		ReportTypes: []string{"FACING_COUN", "PLANOGRAM_COMPLIANCE"},
		Visit:       1,
		Webhook:     "webhook_test",
		CountryCode: "RU",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got RecognizeRequest
		err = json.Unmarshal(b, &got)
		if err != nil {
			t.Fatal(err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	inst, _ := url.Parse(ts.URL)
	client := NewClient(ClintConf{
		Instance: inst,
		APIKey:   "",
	})

	recRes, err := client.Recognize.Recognize(recReq)

	want := &RecognizeResponse{
		ID:     11,
		Images: []int{1, 2, 3},
		Scene:  "4d8b66992cd841f6922723afe9bd8cf8",
		Reports: map[string]int{
			"FACING_COUN":          22,
			"PLANOGRAM_COMPLIANCE": 33,
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, want, recRes)
}
