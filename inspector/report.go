package inspector

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"time"
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
	ID          int                      `json:"id"`
	Status      string                   `json:"status"`
	ReportType  string                   `json:"report_type"`
	CreatedDate time.Time                `json:"created_date,omitempty"`
	UpdatedDate time.Time                `json:"updated_date,omitempty"`
	Visit       int                      `json:"visit,omitempty"`
	Json        []map[string]interface{} `json:"json,omitempty"`
}

type ReportPriceTagsJson struct {
	Brand        string  `json:"brand,omitempty"`
	Manufacturer string  `json:"manufacturer,omitempty"`
	Price        float64 `json:"price"`
	Name         string  `json:"name"`
	Category     string  `json:"category,omitempty"`
	SkuImageUrl  string  `json:"sku_image_url"`
	Promo        string  `json:"promo"`
	SkuId        int     `json:"sku_id"`
}

func (srv *ReportService) GetReport(id int) (*Report, error) {
	path := fmt.Sprintf("reports/%d/", id)
	req, err := srv.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var report Report
	_, err = srv.client.do(req, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (srv *ReportService) ToPriceTags(v []map[string]interface{}) ([]ReportPriceTagsJson, error) {
	var pt []ReportPriceTagsJson
	err := mapstructure.Decode(v, &pt)
	return pt, err
}
