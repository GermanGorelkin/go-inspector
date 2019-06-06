package inspector

import (
	"time"
)

type RecognizeService struct {
	client *Client
}

type RecognizeRequest struct {
	Images      []int    `json:"images"`
	ReportTypes []string `json:"report_types"`

	Display     *int       `json:"display,omitempty"`
	Visit       *int       `json:"visit,omitempty"`
	Datetime    *time.Time `json:"datetime,omitempty"`
	Webhook     *string    `json:"webhook,omitempty"`
	CountryCode *string    `json:"country_code"`
}
type RecognizeResponse struct {
	ID      int            `json:"id"`
	Images  []int          `json:"images"`
	Display int            `json:"display,omitempty"`
	Scene   string         `json:"scene,omitempty"`
	Reports map[string]int `json:"reports"`
}

func (srv *RecognizeService) Recognize(rr *RecognizeRequest) (*RecognizeResponse, error) {
	req, err := srv.client.newRequest("POST", "recognize/", rr)
	if err != nil {
		return nil, err
	}

	var rec RecognizeResponse
	_, err = srv.client.do(req, &rec)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

type RecognitionErrorRequest struct {
	Images  []int  `json:"images"`
	SkuId   int    `json:"sku_gid"`
	Scene   string `json:"scene"`
	Message string `json:"message"`
}
type RecognitionErrorResponse struct {
	RecognitionErrorID int `json:"recognition_error_id"`
}

func (srv *RecognizeService) RecognitionError(rr *RecognitionErrorRequest) (*RecognitionErrorResponse, error) {
	req, err := srv.client.newRequest("POST", "recognition_error/", rr)
	if err != nil {
		return nil, err
	}

	var rec RecognitionErrorResponse
	_, err = srv.client.do(req, &rec)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}
