package inspector

import (
	"encoding/json"
	"testing"
)

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
	if err != nil {
		t.Fatal(err)
	}

	serv := SkuService{}

	t.Run("2 SKU", func(t *testing.T) {
		lsku, err := serv.ToSku(f)
		if err != nil {
			t.Fatal(err)
		}

		if len(lsku) != 2 {
			t.Errorf("len(lsku) got %d want %d", len(lsku), 2)
		}
	})

	t.Run("nil PaginationResults", func(t *testing.T) {
		lsku, err := serv.ToSku(nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(lsku) != 0 {
			t.Errorf("len(lsku) got %d want %d", len(lsku), 0)
		}
	})
}
