package inspector

import (
	"context"
	"fmt"
	"io"
	"time"

	httpclient "github.com/germangorelkin/http-client"
)

// ImageService provides access to the Image Uploads functions in the IC API.
type ImageService struct {
	client *Client
}

// Image represents a JPEG image
type Image struct {
	ID          int       `json:"id"`            // unique Image ID
	URL         string    `json:"url,omitempty"` // image URL
	Width       int       `json:"width"`         // image width in pixels
	Height      int       `json:"height"`        // image height in pixels
	CreatedDate time.Time `json:"created_date"`  // date and time of uploading the image
}

// UploadByUrlRequest represents a payload of upload_by_url
type UploadByUrlRequest struct {
	URL string `json:"url"` // Image URL
}

// Upload uploads Image to IC API via multipart/form-data.
func (srv *ImageService) Upload(ctx context.Context, r io.Reader, filename string) (Image, error) {
	var img Image
	form := httpclient.NewMultipartForm()
	form.AddFile("file", filename, r)

	req, err := srv.client.httpClient.NewMultipartRequest("POST", "uploads/", form)
	if err != nil {
		return img, fmt.Errorf("failed to NewMultipartRequest(POST, uploads/, %v):%w", filename, err)
	}

	_, err = srv.client.httpClient.Do(ctx, req, &img)
	if err != nil {
		return img, fmt.Errorf("failed to Do with Request(POST, uploads/, %v):%w", filename, err)
	}

	return img, nil
}

// UploadByURL uploads Image to IC API by photos url
func (srv *ImageService) UploadByURL(ctx context.Context, url string) (Image, error) {
	var img Image
	body := UploadByUrlRequest{URL: url}

	req, err := srv.client.httpClient.NewRequest("POST", "uploads/upload_by_url/", body)
	if err != nil {
		return img, fmt.Errorf("failed to NewRequest(POST, uploads/upload_by_url/, %v):%w", body, err)
	}

	_, err = srv.client.httpClient.Do(ctx, req, &img)
	if err != nil {
		return img, fmt.Errorf("failed to Do with Request(POST, uploads/upload_by_url/, %v):%w", body, err)
	}

	return img, nil
}
