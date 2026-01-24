package inspector

import (
	"context"
	"fmt"
	"time"
)

// ReportPriceTagsJson represents a IC Visit
type Visit struct {
	ID          int       `json:"id"`                  // Unique visit ID
	Shop        int       `json:"shop"`                // Client-specific customer/outlet/shop
	Agent       string    `json:"agent"`               // Client-specific agent name/id/route
	StartedDate time.Time `json:"started_date"`        // Date and time of the visit start (UTC time).
	Latitude    float64   `json:"latitude,omitempty"`  // Location of the merchandiser at the time of the visit.
	Longitude   float64   `json:"longitude,omitempty"` // Location of the merchandiser at the time of the visit.
}

// VisitService provides access to the Visit functions in the IC API.
type VisitService struct {
	client *Client
}

// AddVisit adds new IC visit
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
