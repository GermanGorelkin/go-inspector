package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewClient(ClintConf{
		Instance: ts.URL,
		APIKey:   "",
	})
	assert.NoError(t, err)

	report, err := client.Report.GetReport(context.Background(), 1)
	assert.NoError(t, err)

	// prepare dates
	cdate, err := time.Parse(time.RFC3339, "2019-08-26T16:33:30.563548Z")
	assert.NoError(t, err)
	udate, err := time.Parse(time.RFC3339, "2019-08-26T16:34:12.241644Z")
	assert.NoError(t, err)
	// prepare json report
	jr := `[{"count": 2,"sku_id": 2176}]`
	var v []interface{}
	err = json.Unmarshal([]byte(jr), &v)
	assert.NoError(t, err)

	want := &Report{
		ID:          2831638,
		Status:      "READY",
		ReportType:  "FACING_COUNT_1_5",
		CreatedDate: cdate,
		UpdatedDate: udate,
		Visit:       115604,
		Json:        v,
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
	assert.NoError(t, err)

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

func TestReportService_ToPriceTags(t *testing.T) {
	b := `[
              {
        "promo": true,
        "category": "HOME & HYGIENE",
        "sku_image_url": "http://henkel.inspector-cloud.ru/media/2019/08/28/0ffbe108-2ab2-4b1d-a296-dca361ce1cd4.jpg",
        "manufacturer": "Henkel",
        "price": 360.0,
        "result_pricetag": 245411316,
        "colors": [
          {
            "score": 0.26388889,
            "color": "white"
          },
          {
            "score": 0.51041667,
            "color": "red"
          }
        ],
        "min_price": 360.0,
        "max_price": 360.0,
        "name": "Bref",
        "result_object": 245411277,
        "sku_id": 9859,
        "brand": "Bref",
        "price_tag_colors": [
          "white",
          "red"
        ]
      }]`
	var m []map[string]interface{}
	err := json.Unmarshal([]byte(b), &m)
	assert.NoError(t, err)

	var srv ReportService
	got, err := srv.ToPriceTags(m)
	assert.NoError(t, err)

	want := []ReportPriceTagsJson{
		{
			Brand:        "Bref",
			Manufacturer: "Henkel",
			Price:        360.0,
			Name:         "Bref",
			Category:     "HOME & HYGIENE",
			SkuImageUrl:  "http://henkel.inspector-cloud.ru/media/2019/08/28/0ffbe108-2ab2-4b1d-a296-dca361ce1cd4.jpg",
			Promo:        "1",
			SkuId:        9859,
		},
	}

	assert.Equal(t, want, got)
}

func TestReportService_ToRealogram(t *testing.T) {
	b := `[
		{
			"image": 55801587,
			"annotations": [
				{
					"h": 250,
					"w": 131,
					"x": 948,
					"y": 1214,
					"name": "Losk 2190 Gel Indian Jasmine 30WL",
					"sku_id": 53733,
					"duplicate": false
				}
			],
			"shelf_annotations": [
				{
					"x1": -14,
					"x2": 1118,
					"y1": 1362,
					"y2": 1349
				}
			]
		}
	]`
	var m []map[string]interface{}
	err := json.Unmarshal([]byte(b), &m)
	assert.NoError(t, err)

	var srv ReportService
	got, err := srv.ToRealogram(m)
	assert.NoError(t, err)

	want := []ReportRealogramJson{
		{
			Image: 55801587,
			Annotations: []ReportRealogramAnnotations{
				{
					H:         250,
					W:         131,
					X:         948,
					Y:         1214,
					Name:      "Losk 2190 Gel Indian Jasmine 30WL",
					SkuId:     53733,
					Duplicate: false,
				},
			},
			ShelfAnnotations: []ReportRealogramShelfAnnotations{
				{
					X1: -14,
					Y1: 1362,
					X2: 1118,
					Y2: 1349,
				},
			},
		},
	}

	assert.Equal(t, want, got)
}

func TestReportService_ParseWebhookReports(t *testing.T) {
	b := []byte(`{
		"id": 406907,
		"display": 1,
		"reports": {
		  "FACING_COUNT_1_5": [
			{
			  "sku_id": 9857,
			  "count": 4
			}
		  ],
		  "PRICE_TAGS": [
              {
				"promo": true,
				"category": "HOME & HYGIENE",
				"sku_image_url": "http://henkel.inspector-cloud.ru/media/2019/08/28/0ffbe108-2ab2-4b1d-a296-dca361ce1cd4.jpg",
				"manufacturer": "Henkel",
				"price": 360.0,
				"result_pricetag": 245411316,
				"colors": [
				{
					"score": 0.26388889,
					"color": "white"
				},
				{
					"score": 0.51041667,
					"color": "red"
				}
				],
				"min_price": 360.0,
				"max_price": 360.0,
				"name": "Bref",
				"result_object": 245411277,
				"sku_id": 9859,
				"brand": "Bref",
				"price_tag_colors": [
				"white",
				"red"
				]
			}],
		"REALOGRAM_1_5":[
			{
				"image": 55801587,
				"annotations": [
					{
						"h": 250,
						"w": 131,
						"x": 948,
						"y": 1214,
						"name": "Losk 2190 Gel Indian Jasmine 30WL",
						"sku_id": 53733,
						"duplicate": false
					}
				],
				"shelf_annotations": [
					{
						"x1": -14,
						"x2": 1118,
						"y1": 1362,
						"y2": 1349
					}
				]
			}
		]
		}
	  }`)

	var srv ReportService

	got, err := srv.ParseWebhookReports(b)
	assert.NoError(t, err)

	want := &WebhookReports{
		ID:      406907,
		Display: 1,
		Reports: struct {
			FacingCount []ReportFacingCountJson `json:"FACING_COUNT_1_5"`
			PriceTags   []ReportPriceTagsJson   `json:"PRICE_TAGS"`
			Realogram   []ReportRealogramJson   `json:"REALOGRAM_1_5"`
		}{
			FacingCount: []ReportFacingCountJson{
				{
					Count: 4,
					SkuId: 9857,
				},
			},
			PriceTags: []ReportPriceTagsJson{
				{
					Brand:        "Bref",
					Manufacturer: "Henkel",
					Price:        360.0,
					Name:         "Bref",
					Category:     "HOME & HYGIENE",
					SkuImageUrl:  "http://henkel.inspector-cloud.ru/media/2019/08/28/0ffbe108-2ab2-4b1d-a296-dca361ce1cd4.jpg",
					Promo:        "",
					SkuId:        9859,
				},
			},
			Realogram: []ReportRealogramJson{
				{
					Image: 55801587,
					Annotations: []ReportRealogramAnnotations{
						{
							H:         250,
							W:         131,
							X:         948,
							Y:         1214,
							Name:      "Losk 2190 Gel Indian Jasmine 30WL",
							SkuId:     53733,
							Duplicate: false,
						},
					},
					ShelfAnnotations: []ReportRealogramShelfAnnotations{
						{
							X1: -14,
							Y1: 1362,
							X2: 1118,
							Y2: 1349,
						},
					},
				},
			},
		},
	}
	assert.Equal(t, want, got)
}

func TestReportService_ParseWebhookReports_Error(t *testing.T) {
	var srv ReportService
	_, err := srv.ParseWebhookReports([]byte(`{`))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to Unmarshal")
}
