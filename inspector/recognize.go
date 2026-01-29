package inspector

import (
	"context"
	"fmt"
	"time"
)

// RecognizeService provides access to the Image recognition functions in the IC API.
type RecognizeService struct {
	client *Client
}

// RecognizeRequest represents a payload of request recognize
type RecognizeRequest struct {
	Images      []int    `json:"images"`       // list of IC image IDs
	ReportTypes []string `json:"report_types"` // list of reports to be generated

	Display     int        `json:"display,omitempty"`      // IC Display ID
	Visit       int        `json:"visit,omitempty"`        // IC Visit ID
	Datetime    *time.Time `json:"datetime,omitempty"`     // date and time of the visit
	Webhook     string     `json:"webhook,omitempty"`      // reports will be sent by POST request to this URL immediately when ready (one by one).
	CountryCode string     `json:"country_code,omitempty"` // two-character сountry code to specify from which country the sku on the photo eg 'RU', 'KZ'.
	RetailChain string     `json:"retail_chain,omitempty"` // retail chain identifier for the store.
}

// RecognizeResponse represents a payload of response recognize
type RecognizeResponse struct {
	ID      int            `json:"id"`                // IC Request ID
	Images  []int          `json:"images"`            // copy from the request
	Display int            `json:"display,omitempty"` // copy from the request or default
	Scene   string         `json:"scene"`             // uuid of the created scene.
	Reports map[string]int `json:"reports"`           // IDs of reports to be generated for the scene
}

// Recognize starts the asynchronous process of recognizing a group of images and returns IDs of reports
func (srv *RecognizeService) Recognize(ctx context.Context, rr RecognizeRequest) (*RecognizeResponse, error) {
	req, err := srv.client.httpClient.NewRequest(methodPOST, endpointRecognize, rr)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(%s, %s, %v):%w", methodPOST, endpointRecognize, rr, err)
	}

	var rec RecognizeResponse
	_, err = srv.client.httpClient.Do(ctx, req, &rec)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(%s, %s, %v):%w", methodPOST, endpointRecognize, rr, err)
	}

	return &rec, nil
}

// RecognitionErrorRequest represents a payload of request Recognition error message
type RecognitionErrorRequest struct {
	Images  []int  `json:"images"`  // list of image IDs
	SkuId   int    `json:"sku_gid"` // global id that identifies sku that was mistakenly recognized.
	Scene   string `json:"scene"`   // scene uuid that identifies sku that was mistakenly recognized
	Message string `json:"message"` // error description.
}

// RecognitionErrorRequest represents a payload of response Recognition error message
type RecognitionErrorResponse struct {
	RecognitionErrorID int `json:"recognition_error_id"` // IC ID of the saved recognition error request
}

// RecognitionError creates a recognition error message
func (srv *RecognizeService) RecognitionError(ctx context.Context, rr *RecognitionErrorRequest) (*RecognitionErrorResponse, error) {
	req, err := srv.client.httpClient.NewRequest(methodPOST, endpointRecognitionError, rr)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(%s, %s, %v):%w", methodPOST, endpointRecognitionError, rr, err)
	}

	var rec RecognitionErrorResponse
	_, err = srv.client.httpClient.Do(ctx, req, &rec)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(%s, %s, %v):%w", methodPOST, endpointRecognitionError, rr, err)
	}

	return &rec, nil
}
