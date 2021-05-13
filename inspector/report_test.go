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
		if err != nil {
			t.Fatal(err)
		}
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
	_ = json.Unmarshal([]byte(jr), &v)

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

//func TestReportService_ToPriceTags(t *testing.T) {
//	b := `[
//        {
//            "sku_image_url": "http://henkel.inspector-cloud.ru/media/f8837b16-3987-4248-b244-dda1960a7b38.jpg",
//            "price_tag_colors": [
//                "white"
//            ],
//            "category": "LAUNDRY",
//            "sku_id": 12176,
//            "colors": [
//                {
//                    "color": "white",
//                    "score": 0.97222222
//                }
//            ],
//            "min_price": 40.0,
//            "name": "Vernel 910/1000 Детский",
//            "promo": "No",
//            "price": 40.0,
//            "brand": "vernel",
//            "manufacturer": "Henkel",
//            "max_price": 40.0
//        },
//        {
//            "sku_image_url": "http://henkel.inspector-cloud.ru/media/40327d62-b4e5-4d52-bf37-63152a1b4e66.jpg",
//            "price_tag_colors": [
//                "white"
//            ],
//            "category": "LAUNDRY",
//            "sku_id": 12177,
//            "colors": [
//                {
//                    "color": "white",
//                    "score": 0.87152778
//                }
//            ],
//            "min_price": 177.0,
//            "name": "Vernel 910/1000 Гибискус и Роза",
//            "promo": "No",
//            "price": 177.0,
//            "brand": "vernel",
//            "manufacturer": "Henkel",
//            "max_price": 177.0
//        }]`
//	var m []map[string]interface{}
//	err := json.Unmarshal([]byte(b), &m)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var srv ReportService
//	got, err := srv.ToPriceTags(m)
//	assert.NoError(t, err)
//
//	want := []ReportPriceTagsJson{
//		{
//			Brand:        "vernel",
//			Manufacturer: "Henkel",
//			Price:        40.0,
//			Name:         "Vernel 910/1000 Детский",
//			Category:     "LAUNDRY",
//			SkuImageUrl:  "http://henkel.inspector-cloud.ru/media/f8837b16-3987-4248-b244-dda1960a7b38.jpg",
//			Promo:        "0",
//			SkuId:        12176,
//		},
//		{
//			Brand:        "vernel",
//			Manufacturer: "Henkel",
//			Price:        177.0,
//			Name:         "Vernel 910/1000 Гибискус и Роза",
//			Category:     "LAUNDRY",
//			SkuImageUrl:  "http://henkel.inspector-cloud.ru/media/40327d62-b4e5-4d52-bf37-63152a1b4e66.jpg",
//			Promo:        "0",
//			SkuId:        12177,
//		},
//	}
//
//	assert.Equal(t, want, got)
//}

func TestReportService_ToPriceTags_newdata(t *testing.T) {
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
	if err != nil {
		t.Fatal(err)
	}

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
			}]
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
		},
	}

	assert.Equal(t, want, got)
}
