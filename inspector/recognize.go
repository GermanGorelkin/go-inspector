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

	Display  *int       `json:"display,omitempty"`
	Visit    *int       `json:"visit,omitempty"`
	Datetime *time.Time `json:"datetime,omitempty"`
	Webhook  *string    `json:"webhook,omitempty"`
}

type RecognizeResponse struct {
	ID      int            `json:"id"`
	Images  []int          `json:"images"`
	Display int            `json:"display,omitempty"`
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
