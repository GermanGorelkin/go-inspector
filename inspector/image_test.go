package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImageService_UploadByURL(t *testing.T) {
	imgUrl := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		var got UploadByUrlRequest
		err = json.Unmarshal(b, &got)
		assert.NoError(t, err)
		want := UploadByUrlRequest{URL: imgUrl}
		assert.Equal(t, want, got)

		_, err = fmt.Fprintln(w, `{  
				"id": 156673,  
                "url": "https://test.inspector-cloud.com/media/12345678-1234-5678-1234567812345678.jpg",
				"width": 720,
				"height": 1280,
				"created_date": "2016-08-31T10:32:15.687287Z"  }`)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{
		Instance: ts.URL,
		APIKey:   "",
	})
	assert.NoError(t, err)

	img, err := client.Image.UploadByURL(context.Background(), imgUrl)
	assert.NoError(t, err)

	date, _ := time.Parse(time.RFC3339, "2016-08-31T10:32:15.687287Z")
	want := Image{
		ID:          156673,
		URL:         "https://test.inspector-cloud.com/media/12345678-1234-5678-1234567812345678.jpg",
		Width:       720,
		Height:      1280,
		CreatedDate: date,
	}
	assert.Equal(t, want, img)
}
