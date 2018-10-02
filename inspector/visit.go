package inspector

import (
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

func (srv *VisitService) AddVisit(v *Visit) (*Visit, error) {
	rd := struct{}{}

	req, err := srv.client.newRequest("POST", "visits/", rd)
	if err != nil {
		return nil, err
	}

	_, err = srv.client.do(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
