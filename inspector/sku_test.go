package inspector

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
