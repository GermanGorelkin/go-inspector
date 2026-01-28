package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const shortPageResults = `[
        {
            "id": 10,
            "cid": "SKU010",
            "ean13": "1234567890999",
            "image": 1010,
            "name": "Short Page Product",
            "brand": 3,
            "category": 3,
            "manufacturer": 3,
            "size_x_mm": 110.0,
            "size_y_mm": 210.0,
            "size_z_mm": 55.0
        }]`

const testPaginationResults = `[
        {
            "id": 26,
            "cid": "4601501027624",
            "ean13": null,
            "image": 3166335,
            "name": "Heineken банка 0,5 л",
            "brand": 25,
            "category": 7,
            "manufacturer": null,
            "size_x_mm": null,
            "size_y_mm": null,
            "size_z_mm": null
        },
        {
            "id": 12423,
            "cid": null,
            "ean13": null,
            "image": 3178440,
            "name": "AUTO:Mission_284537524/551756127",
            "brand": null,
            "category": 53,
            "manufacturer": null,
            "size_x_mm": null,
            "size_y_mm": null,
            "size_z_mm": null
        }]`

func TestSkuService_ToSku(t *testing.T) {
	b := []byte(testPaginationResults)
	var f interface{}
	err := json.Unmarshal(b, &f)
	assert.NoError(t, err)

	serv := SkuService{}

	t.Run("2 SKU", func(t *testing.T) {
		lsku, err := serv.ToSku(f)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(lsku))
	})

	t.Run("nil PaginationResults", func(t *testing.T) {
		lsku, err := serv.ToSku(nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(lsku))
	})

	t.Run("invalid structure", func(t *testing.T) {
		_, err := serv.ToSku(map[string]string{"bad": "data"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to Decode")
	})
}

func TestSkuService_IterateSKU(t *testing.T) {
	// Test data for multiple pages
	page1Results := `[
		{
			"id": 1,
			"cid": "SKU001",
			"ean13": "1234567890123",
			"image": 1001,
			"name": "Product 1",
			"brand": 1,
			"category": 1,
			"manufacturer": 1,
			"size_x_mm": 100.0,
			"size_y_mm": 200.0,
			"size_z_mm": 50.0
		},
		{
			"id": 2,
			"cid": "SKU002",
			"ean13": "1234567890124",
			"image": 1002,
			"name": "Product 2",
			"brand": 1,
			"category": 1,
			"manufacturer": 1,
			"size_x_mm": 150.0,
			"size_y_mm": 250.0,
			"size_z_mm": 60.0
		}
	]`

	page2Results := `[
		{
			"id": 3,
			"cid": "SKU003",
			"ean13": "1234567890125",
			"image": 1003,
			"name": "Product 3",
			"brand": 2,
			"category": 2,
			"manufacturer": 2,
			"size_x_mm": 120.0,
			"size_y_mm": 220.0,
			"size_z_mm": 55.0
		}
	]`

	t.Run("multiple pages", func(t *testing.T) {
		callCount := 0
		var serverURL string

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, methodGET, r.Method)
			assert.Equal(t, "/"+endpointSKU, r.URL.Path)

			query := r.URL.Query()
			limit := query.Get("limit")
			offset := query.Get("offset")

			callCount++

			// Return different pages based on offset
			var response string
			var nextURL *string

			switch offset {
			case "0":
				response = page1Results
				nextURLStr := fmt.Sprintf("%s/sku/?limit=%s&offset=2", serverURL, limit)
				nextURL = &nextURLStr
			case "2":
				response = page2Results
				nextURL = nil // This is the last page (partial page)
			default:
				t.Errorf("unexpected offset: %s", offset)
				response = `[]`
			}

			pagResponse := map[string]interface{}{
				"count":    3,
				"next":     nextURL,
				"previous": nil,
				"results":  json.RawMessage(response),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		})

		ts := httptest.NewServer(handler)
		defer ts.Close()
		serverURL = ts.URL

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		iterator := client.Sku.IterateSKU(context.Background(), 2)
		assert.NotNil(t, iterator)

		// Collect all SKUs
		var allSKUs []Sku
		for {
			page, err := iterator.Next()
			assert.NoError(t, err)
			if page == nil {
				break
			}
			allSKUs = append(allSKUs, page...)
		}

		assert.Equal(t, 3, len(allSKUs))
		assert.Equal(t, 2, callCount) // Should have made 2 calls (offsets 0 and 2)

		// Verify SKU data
		assert.Equal(t, 1, allSKUs[0].ID)
		assert.Equal(t, "SKU001", allSKUs[0].CID)
		assert.Equal(t, 3, allSKUs[2].ID)
		assert.Equal(t, "SKU003", allSKUs[2].CID)
	})

	t.Run("short page with next", func(t *testing.T) {
		callCount := 0
		var serverURL string

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, methodGET, r.Method)
			assert.Equal(t, "/"+endpointSKU, r.URL.Path)

			callCount++
			offset := r.URL.Query().Get("offset")

			var response string
			var nextURL *string

			switch offset {
			case "0":
				response = shortPageResults
				nextURLStr := fmt.Sprintf("%s/sku/?limit=2&offset=2", serverURL)
				nextURL = &nextURLStr
			case "2":
				response = page2Results
				nextURL = nil
			default:
				t.Errorf("unexpected offset: %s", offset)
				response = `[]`
			}

			pagResponse := map[string]interface{}{
				"count":    2,
				"next":     nextURL,
				"previous": nil,
				"results":  json.RawMessage(response),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		})

		ts := httptest.NewServer(handler)
		defer ts.Close()
		serverURL = ts.URL

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		iterator := client.Sku.IterateSKU(context.Background(), 2)
		assert.NotNil(t, iterator)

		var allSKUs []Sku
		for {
			page, err := iterator.Next()
			assert.NoError(t, err)
			if page == nil {
				break
			}
			allSKUs = append(allSKUs, page...)
		}

		assert.Equal(t, 2, len(allSKUs))
		assert.Equal(t, 2, callCount)
		assert.Equal(t, 10, allSKUs[0].ID)
		assert.Equal(t, 3, allSKUs[1].ID)
	})

	t.Run("single page", func(t *testing.T) {
		callCount := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			assert.Equal(t, "0", r.URL.Query().Get("offset"))

			pagResponse := map[string]interface{}{
				"count":    2,
				"next":     nil,
				"previous": nil,
				"results":  json.RawMessage(page1Results),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		}))
		defer ts.Close()

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		iterator := client.Sku.IterateSKU(context.Background(), 10)
		page1, err := iterator.Next()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(page1))

		page2, err := iterator.Next()
		assert.NoError(t, err)
		assert.Nil(t, page2) // No more pages

		assert.Equal(t, 1, callCount)
	})

	t.Run("empty result set", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pagResponse := map[string]interface{}{
				"count":    0,
				"next":     nil,
				"previous": nil,
				"results":  json.RawMessage(`[]`),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		}))
		defer ts.Close()

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		iterator := client.Sku.IterateSKU(context.Background(), 10)
		page, err := iterator.Next()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(page))

		nextPage, err := iterator.Next()
		assert.NoError(t, err)
		assert.Nil(t, nextPage)
	})

	t.Run("error handling", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, `{"error": "internal server error"}`)
		}))
		defer ts.Close()

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		iterator := client.Sku.IterateSKU(context.Background(), 10)
		_, err = iterator.Next()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch SKU page")
	})
}

func TestSkuService_GetAllSKU(t *testing.T) {
	page1Results := `[
		{"id": 1, "cid": "SKU001", "name": "Product 1", "image": 1001},
		{"id": 2, "cid": "SKU002", "name": "Product 2", "image": 1002}
	]`

	page2Results := `[
		{"id": 3, "cid": "SKU003", "name": "Product 3", "image": 1003}
	]`

	t.Run("fetches all pages", func(t *testing.T) {
		callCount := 0
		var serverURL string

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			offset := r.URL.Query().Get("offset")

			var response string
			var nextURL *string

			if offset == "0" {
				response = page1Results
				nextURLStr := fmt.Sprintf("%s/sku/?limit=2&offset=2", serverURL)
				nextURL = &nextURLStr
			} else {
				response = page2Results
				nextURL = nil
			}

			pagResponse := map[string]interface{}{
				"count":    3,
				"next":     nextURL,
				"previous": nil,
				"results":  json.RawMessage(response),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		})

		ts := httptest.NewServer(handler)
		defer ts.Close()
		serverURL = ts.URL

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		allSKUs, err := client.Sku.GetAllSKU(context.Background(), 2)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(allSKUs))
		assert.Equal(t, 2, callCount)

		assert.Equal(t, 1, allSKUs[0].ID)
		assert.Equal(t, "SKU001", allSKUs[0].CID)
		assert.Equal(t, 3, allSKUs[2].ID)
		assert.Equal(t, "SKU003", allSKUs[2].CID)
	})

	t.Run("handles empty result", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pagResponse := map[string]interface{}{
				"count":    0,
				"next":     nil,
				"previous": nil,
				"results":  json.RawMessage(`[]`),
			}

			w.Header().Set(headerContentType, contentTypeJSON)
			json.NewEncoder(w).Encode(pagResponse)
		}))
		defer ts.Close()

		client, err := NewClient(ClintConf{
			Instance: ts.URL,
			APIKey:   "test-key",
		})
		assert.NoError(t, err)

		allSKUs, err := client.Sku.GetAllSKU(context.Background(), 10)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(allSKUs))
	})
}
