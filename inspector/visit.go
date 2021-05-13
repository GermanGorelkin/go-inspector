package inspector

import (
	"context"
	"fmt"
	"time"
)

type Visit struct {
	ID          int       `json:"id"`
	Shop        int       `json:"shop"`
	Agent       string    `json:"agent"`
	StartedDate time.Time `json:"started_date"`
	Latitude    float64   `json:"latitude,omitempty"`
	Longitude   float64   `json:"longitude,omitempty"`
}

type VisitService struct {
	client *Client
}

func (srv *VisitService) AddVisit(ctx context.Context, v *Visit) (*Visit, error) {
	rd := struct{}{}

	req, err := srv.client.httpClient.NewRequest("POST", "visits/", rd)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(POST, visits/):%w", err)
	}

	_, err = srv.client.httpClient.Do(ctx, req, v)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(POST, visits/):%w", err)
	}

	return v, nil
}
