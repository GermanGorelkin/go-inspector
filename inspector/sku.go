package inspector

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// ReportPriceTagsJson represents a IC SKU
type Sku struct {
	ID           int      `json:"id"`                                           // Unique SKU ID
	CID          string   `json:"cid"`                                          // Client-specific SKU ID, e.g. internal code
	EAN13        *string  `json:"ean13,omitempty"`                              // European Article Number
	Image        int      `json:"image"`                                        // SKU image ID,
	Name         string   `json:"name"`                                         // Human readable SKU name
	Brand        *int     `json:"brand,omitempty"`                              // Brand ID
	Category     *int     `json:"category,omitempty"`                           // Category ID
	Manufacturer *int     `json:"manufacturer,omitempty"`                       // Manufacturer ID
	SizeXMM      *float64 `json:"size_x_mm,omitempty" mapstructure:"size_x_mm"` // Product width in mm
	SizeYMM      *float64 `json:"size_y_mm,omitempty" mapstructure:"size_y_mm"` // Product height in mm
	SizeZMM      *float64 `json:"size_z_mm,omitempty" mapstructure:"size_z_mm"` // Product depth in mm
}

//ReportService provides access to the SKU functions in the IC API.
type SkuService struct {
	client *Client
}

// GetSKU requests list of SKU.
// Return Pagination for the given offset and limit
func (srv *SkuService) GetSKU(ctx context.Context, offset, limit int) (*Pagination, error) {
	path := "sku/"
	req, err := srv.client.httpClient.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(GET, %s):%w", path, err)
	}

	q := fmt.Sprintf("limit=%d&offset=%d", limit, offset)
	req.URL.RawQuery += q

	var pag Pagination
	_, err = srv.client.httpClient.Do(ctx, req, &pag)
	if err != nil {
		return nil, fmt.Errorf("failed to Do with Request(GET, %s):%w", req.URL.RawQuery, err)
	}

	return &pag, nil
}

// ToSku parses json to []Sku
func (srv *SkuService) ToSku(v interface{}) ([]Sku, error) {
	var r []Sku
	if err := mapstructure.Decode(v, &r); err != nil {
		return r, fmt.Errorf("failed to Decode %v:%w", v, err)
	}
	return r, nil
}
