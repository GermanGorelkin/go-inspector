package inspector

import (
	"context"
	"fmt"
	"time"
)

type RecognizeService struct {
	client *Client
}

type RecognizeRequest struct {
	Images      []int    `json:"images"`
	ReportTypes []string `json:"report_types"`

	Display     int        `json:"display,omitempty"`
	Visit       int        `json:"visit,omitempty"`
	Datetime    *time.Time `json:"datetime,omitempty"`
	Webhook     string     `json:"webhook,omitempty"`
	CountryCode string     `json:"country_code,omitempty"`
}
type RecognizeResponse struct {
	ID      int            `json:"id"`
	Images  []int          `json:"images"`
	Display int            `json:"display,omitempty"`
	Scene   string         `json:"scene"`
	Reports map[string]int `json:"reports"`
}

func (srv *RecognizeService) Recognize(ctx context.Context, rr RecognizeRequest) (*RecognizeResponse, error) {
	req, err := srv.client.httpClient.NewRequest("POST", "recognize/", rr)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(POST, recognize/, %v):%w", rr, err)
	}

	var rec RecognizeResponse
	_, err = srv.client.httpClient.Do(ctx, req, &rec)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(POST, recognize/, %v):%w", rr, err)
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

func (srv *RecognizeService) RecognitionError(ctx context.Context, rr *RecognitionErrorRequest) (*RecognitionErrorResponse, error) {
	req, err := srv.client.httpClient.NewRequest("POST", "recognition_error/", rr)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(POST, recognition_error/, %v):%w", rr, err)
	}

	var rec RecognitionErrorResponse
	_, err = srv.client.httpClient.Do(ctx, req, &rec)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(POST, recognition_error/, %v):%w", rr, err)
	}

	return &rec, nil
}
