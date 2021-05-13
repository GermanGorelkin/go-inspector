package inspector

import (
	"context"
	"fmt"
	"time"
)

type ImageService struct {
	client *Client
}

type Image struct {
	ID          int       `json:"id"`
	URL         string    `json:"url,omitempty"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	CreatedDate time.Time `json:"created_date"`
}

type UploadByUrlRequest struct {
	URL string `json:"url"`
}

// TODO
// func (srv *ImageService) Upload(r io.Reader, filename string) (*Image, error) {
// 	req, err := srv.client.newRequestFormFile("uploads/", r, filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var img Image
// 	_, err = srv.client.do(req, &img)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &img, nil
// }

// UploadByURL
func (srv *ImageService) UploadByURL(ctx context.Context, url string) (*Image, error) {
	body := UploadByUrlRequest{URL: url}

	req, err := srv.client.httpClient.NewRequest("POST", "uploads/upload_by_url/", body)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(POST, uploads/upload_by_url/, %v):%w", body, err)
	}

	img := new(Image)
	_, err = srv.client.httpClient.Do(ctx, req, img)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(POST, uploads/upload_by_url/, %v):%w", body, err)
	}

	return img, nil
}
