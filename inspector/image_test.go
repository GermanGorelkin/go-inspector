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
	"time"
)

func TestImageService_UploadByURL(t *testing.T) {
	imgUrl := "test"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got UploadByUrlRequest
		err = json.Unmarshal(b, &got)
		if err != nil {
			t.Fatal(err)
		}
		want := UploadByUrlRequest{URL: imgUrl}
		assert.Equal(t, want, got)

		_, err = fmt.Fprintln(w, `{  
				"id": 156673,  
                "url": "https://test.inspector-cloud.com/media/12345678-1234-5678-1234567812345678.jpg",
				"width": 720,
				"height": 1280,
				"created_date": "2016-08-31T10:32:15.687287Z"  }`)
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

	img, err := client.Image.UploadByURL(imgUrl)
	assert.NoError(t, err)

	date, _ := time.Parse(time.RFC3339, "2016-08-31T10:32:15.687287Z")
	want := &Image{
		ID:          156673,
		URL:         "https://test.inspector-cloud.com/media/12345678-1234-5678-1234567812345678.jpg",
		Width:       720,
		Height:      1280,
		CreatedDate: date,
	}
	assert.Equal(t, want, img)
}
