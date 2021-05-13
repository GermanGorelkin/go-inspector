package inspector

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Sku struct {
	ID           int      `json:"id"`
	CID          string   `json:"cid"`
	EAN13        *string  `json:"ean13,omitempty"`
	Image        int      `json:"image"`
	Name         string   `json:"name"`
	Brand        *int     `json:"brand,omitempty"`
	Category     *int     `json:"category,omitempty"`
	Manufacturer *int     `json:"manufacturer,omitempty"`
	SizeXMM      *float64 `json:"size_x_mm,omitempty" mapstructure:"size_x_mm"`
	SizeYMM      *float64 `json:"size_y_mm,omitempty" mapstructure:"size_y_mm"`
	SizeZMM      *float64 `json:"size_z_mm,omitempty" mapstructure:"size_z_mm"`
}

type SkuService struct {
	client *Client
}

func (srv *SkuService) GetSKU(ctx context.Context, offset, limit int) (*Pagination, error) {
	path := "sku/"
	req, err := srv.client.httpClient.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf("limit=%d&offset=%d", limit, offset)
	req.URL.RawQuery += q

	var pag Pagination
	_, err = srv.client.httpClient.Do(ctx, req, &pag)
	if err != nil {
		return nil, err
	}

	return &pag, nil
}

func (srv *SkuService) ToSku(v interface{}) ([]Sku, error) {
	var r []Sku
	err := mapstructure.Decode(v, &r)
	if err != nil {
		err = errors.WithStack(err)
	}
	return r, err
}
