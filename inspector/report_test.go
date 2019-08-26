package inspector

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestReportService_GetReport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, `{
							"id":2831638,
							"status":"READY",
							"report_type":"FACING_COUNT_1_5",
							"created_date":"2019-08-26T16:33:30.563548Z",
							"updated_date":"2019-08-26T16:34:12.241644Z",
							"visit":115604,
							"json":[{
									"count": 2,
									"sku_id": 2176
								}]
						}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	inst, _ := url.Parse(ts.URL)
	client := NewClient(ClintConf{
		Instance: inst,
		APIKey:   "",
	})

	report, err := client.Report.GetReport(1)
	assert.NoError(t, err)

	cdate, err := time.Parse(time.RFC3339, "2019-08-26T16:33:30.563548Z")
	assert.NoError(t, err)
	udate, err := time.Parse(time.RFC3339, "2019-08-26T16:34:12.241644Z")
	assert.NoError(t, err)
	want := &Report{
		ID:          2831638,
		Status:      "READY",
		ReportType:  "FACING_COUNT_1_5",
		CreatedDate: cdate,
		UpdatedDate: udate,
		Visit:       115604,
		Json: []map[string]interface{}{
			{
				"count":  float64(2),
				"sku_id": float64(2176),
			},
		},
	}

	assert.Equal(t, want, report)
}

func TestReportService_ToFacingCount(t *testing.T) {
	b := `[
        {
            "count": 2,
            "sku_id": 2176
        },
        {
            "count": 83,
            "sku_id": 2
        },
        {
            "count": 1,
            "sku_id": 1989
        }]`
	var m []map[string]interface{}
	err := json.Unmarshal([]byte(b), &m)
	if err != nil {
		t.Fatal(err)
	}

	var srv ReportService
	got, err := srv.ToFacingCount(m)
	assert.NoError(t, err)

	want := []ReportFacingCountJson{
		{
			Count: 2,
			SkuId: 2176,
		},
		{
			Count: 83,
			SkuId: 2,
		},
		{
			Count: 1,
			SkuId: 1989,
		},
	}

	assert.Equal(t, want, got)
}
