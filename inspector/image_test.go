package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestImageService_Upload(t *testing.T) {
	const (
		filename = "shelf.jpg"
		payload  = "test-image-data"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, methodPOST, r.Method)
		assert.Equal(t, "/"+endpointUploads, r.URL.Path)
		assert.Equal(t, authSchemeToken+" test-key", r.Header.Get(headerAuthorization))

		mediaType, params, err := mime.ParseMediaType(r.Header.Get(headerContentType))
		assert.NoError(t, err)
		assert.Equal(t, contentTypeMultipartFormData, mediaType)

		reader := multipart.NewReader(r.Body, params["boundary"])
		part, err := reader.NextPart()
		assert.NoError(t, err)
		assert.Equal(t, formFieldFile, part.FormName())
		assert.Equal(t, filename, part.FileName())

		data, err := io.ReadAll(part)
		assert.NoError(t, err)
		assert.Equal(t, payload, string(data))

		_, err = reader.NextPart()
		assert.Equal(t, io.EOF, err)

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
		APIKey:   "test-key",
	})
	assert.NoError(t, err)

	img, err := client.Image.Upload(context.Background(), strings.NewReader(payload), filename)
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
