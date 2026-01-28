package inspector

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

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

// ReportService provides access to the SKU functions in the IC API.
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
func (srv *SkuService) ToSku(v any) ([]Sku, error) {
	var r []Sku
	if err := mapstructure.Decode(v, &r); err != nil {
		return r, fmt.Errorf("failed to Decode %v:%w", v, err)
	}
	return r, nil
}

// SKUIterator provides paginated iteration over SKUs.
type SKUIterator struct {
	client    *SkuService
	ctx       context.Context
	pageSize  int
	offset    int
	hasMore   bool
	seenPages map[int]bool
	maxPages  int
}

// IterateSKU returns an iterator for paginated SKU retrieval.
// pageSize controls how many items are fetched per page (default: 100).
// The iterator automatically handles pagination and includes safeguards
// against infinite loops.
func (srv *SkuService) IterateSKU(ctx context.Context, pageSize int) *SKUIterator {
	if pageSize <= 0 {
		pageSize = 100
	}
	return &SKUIterator{
		client:    srv,
		ctx:       ctx,
		pageSize:  pageSize,
		offset:    0,
		hasMore:   true,
		seenPages: make(map[int]bool),
		maxPages:  1000, // Safety limit to prevent infinite loops
	}
}

// Next returns the next page of SKUs.
// Returns nil, nil when no more pages are available.
func (it *SKUIterator) Next() ([]Sku, error) {
	if !it.hasMore {
		return nil, nil
	}

	// Check infinite loop safeguard
	if it.seenPages[it.offset] {
		return nil, fmt.Errorf("detected pagination loop at offset %d", it.offset)
	}
	if len(it.seenPages) >= it.maxPages {
		return nil, fmt.Errorf("exceeded maximum page limit of %d", it.maxPages)
	}

	// Mark this page as seen
	it.seenPages[it.offset] = true

	// Fetch the page
	pag, err := it.client.GetSKU(it.ctx, it.offset, it.pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SKU page at offset %d:%w", it.offset, err)
	}

	// Convert results to SKUs
	skus, err := it.client.ToSku(pag.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SKU page at offset %d:%w", it.offset, err)
	}

	// Check if we have more pages
	it.hasMore = pag.Next != nil

	// Update iterator state using next page offset when available
	if pag.Next != nil {
		nextOffset, ok := parseNextOffset(*pag.Next)
		if ok {
			it.offset = nextOffset
		} else {
			it.offset += len(skus)
		}
	} else {
		it.offset += len(skus)
	}

	return skus, nil
}

func parseNextOffset(nextURL string) (int, bool) {
	parsed, err := url.Parse(nextURL)
	if err != nil {
		return 0, false
	}
	offsetParam := parsed.Query().Get("offset")
	if offsetParam == "" {
		return 0, false
	}
	value, err := strconv.Atoi(offsetParam)
	if err != nil {
		return 0, false
	}
	return value, true
}

// GetAllSKU fetches all SKUs using automatic pagination.
// pageSize controls how many items are fetched per page (default: 100).
func (srv *SkuService) GetAllSKU(ctx context.Context, pageSize int) ([]Sku, error) {
	iterator := srv.IterateSKU(ctx, pageSize)
	var allSKUs []Sku

	for {
		page, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		if page == nil {
			break
		}
		allSKUs = append(allSKUs, page...)
	}

	return allSKUs, nil
}
