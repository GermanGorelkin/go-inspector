package inspector

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	ReportTypeFACING_COUNT         = "FACING_COUNT"
	ReportTypeSHARE_OF_SPAC        = "SHARE_OF_SPAC"
	ReportTypeREALOGRAM            = "REALOGRAM"
	ReportTypePRICE_TAGS           = "PRICE_TAGS"
	ReportTypeMHL_COMPLIANCE       = "MHL_COMPLIANCE"
	ReportTypePLANOGRAM_COMPLIANCE = "PLANOGRAM_COMPLIANCE"

	ReportStatusNOT_READY = "NOT_READY"
	ReportStatusREADY     = "READY"
	ReportStatusERROR     = "ERROR"
)

type ReportService struct {
	client *Client
}

type Report struct {
	ID          int         `json:"id"`
	Status      string      `json:"status"`
	ReportType  string      `json:"report_type"`
	CreatedDate time.Time   `json:"created_date,omitempty"`
	UpdatedDate time.Time   `json:"updated_date,omitempty"`
	Visit       int         `json:"visit,omitempty"`
	Json        interface{} `json:"json,omitempty"`
}

type WebhookReports struct {
	ID      int `json:"id"`
	Display int `json:"display"`
	Reports struct {
		FacingCount []ReportFacingCountJson `json:"FACING_COUNT_1_5"`
		PriceTags   []ReportPriceTagsJson   `json:"PRICE_TAGS"`
	}
}

type ReportPriceTagsJson struct {
	Brand        string  `json:"brand,omitempty"`
	Manufacturer string  `json:"manufacturer,omitempty"`
	Price        float64 `json:"price"`
	Name         string  `json:"name"`
	Category     string  `json:"category,omitempty"`
	SkuImageUrl  string  `json:"sku_image_url" mapstructure:"sku_image_url"`
	Promo        string  `json:"-"`
	SkuId        int     `json:"sku_id" mapstructure:"sku_id"`
}

type ReportFacingCountJson struct {
	Count int `json:"count"`
	SkuId int `json:"sku_id" mapstructure:"sku_id"`
}

func (srv *ReportService) GetReport(id int) (*Report, error) {
	path := fmt.Sprintf("reports/%d/", id)
	req, err := srv.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request %q:%w", path, err)
	}

	var report Report
	if _, err = srv.client.do(req, &report); err != nil {
		return nil, fmt.Errorf("failed to do request %q:%w", path, err)
	}

	return &report, nil
}

func (srv *ReportService) ToPriceTags(v interface{}) ([]ReportPriceTagsJson, error) {
	var r []ReportPriceTagsJson
	if err := mapstructure.WeakDecode(v, &r); err != nil {
		return r, fmt.Errorf("failed to WeakDecode %v:%w", v, err)
	}
	return r, nil
}

func (srv *ReportService) ToFacingCount(v interface{}) ([]ReportFacingCountJson, error) {
	var r []ReportFacingCountJson
	if err := mapstructure.WeakDecode(v, &r); err != nil {
		return r, fmt.Errorf("failed to WeakDecode %v:%w", v, err)
	}
	return r, nil
}

func (srv *ReportService) ParseWebhookReports(b []byte) (*WebhookReports, error) {
	var reports WebhookReports
	if err := json.Unmarshal(b, &reports); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal %q:%w", b, err)
	}
	return &reports, nil
}
