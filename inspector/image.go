package inspector

import (
	"io"
	"time"
)

type ImageService struct {
	client *Client
}

type Image struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	CreatedDate time.Time `json:"created_date"`
}

func (srv *ImageService) Upload(r io.Reader, filename string) (*Image, error) {
	req, err := srv.client.newRequestFormFile("uploads/", r, filename)
	if err != nil {
		return nil, err
	}

	var img Image
	_, err = srv.client.do(req, &img)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

//upload_by_url

func (srv *ImageService) UploadByURL(url string) (*Image, error) {
	rd := struct {
		URL string `json:"url"`
	}{
		URL: url,
	}

	req, err := srv.client.newRequest("POST", "uploads/upload_by_url/", rd)
	if err != nil {
		return nil, err
	}

	var img Image
	_, err = srv.client.do(req, &img)
	if err != nil {
		return nil, err
	}

	return &img, nil
}
