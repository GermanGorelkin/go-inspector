package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	// Report types
	ReportTypeFACING_COUNT         = "FACING_COUNT"
	ReportTypeSHARE_OF_SPAC        = "SHARE_OF_SPAC"
	ReportTypeREALOGRAM            = "REALOGRAM"
	ReportTypePRICE_TAGS           = "PRICE_TAGS"
	ReportTypeMHL_COMPLIANCE       = "MHL_COMPLIANCE"
	ReportTypePLANOGRAM_COMPLIANCE = "PLANOGRAM_COMPLIANCE"

	// Report statuses
	ReportStatusNOT_READY = "NOT_READY" // the report in the process of preparation. The client must repeat the request later;
	ReportStatusREADY     = "READY"     // the report was successfully prepared. The client can use the “data” field, see below;
	ReportStatusERROR     = "ERROR"     // error in the process of preparing the report. The error message is available in the 'error' field.

	// Report polling defaults
	ReportWaitDefaultInterval = 2 * time.Second
	ReportWaitDefaultTimeout  = 60 * time.Second
)

// ReportService provides access to the Reports functions in the IC API.
type ReportService struct {
	client *Client
}

// ReportWaitProgressFunc receives the latest report during polling.
type ReportWaitProgressFunc func(report *Report)

// ReportWaitBackoffFunc returns the next interval based on attempts.
type ReportWaitBackoffFunc func(attempt int, prevInterval time.Duration) time.Duration

// ReportWaitOptions configures polling behavior for WaitForReport.
type ReportWaitOptions struct {
	Interval   time.Duration          // base polling interval (default: ReportWaitDefaultInterval)
	Timeout    time.Duration          // overall timeout (default: ReportWaitDefaultTimeout)
	Backoff    ReportWaitBackoffFunc  // optional interval backoff
	OnProgress ReportWaitProgressFunc // optional progress callback
}

// Report represents a payload of report
type Report struct {
	ID          int       `json:"id"`                     // unique report ID
	Status      string    `json:"status"`                 // report status
	ReportType  string    `json:"report_type"`            // report type
	CreatedDate time.Time `json:"created_date,omitempty"` // date and time of report generation
	UpdatedDate time.Time `json:"updated_date,omitempty"` // date and time of report update
	Visit       int       `json:"visit,omitempty"`        // IC Visit ID
	Json        any       `json:"json,omitempty"`         // Report data
}

// WebhookReports represents a payload of report from webhook
// Once the reports are generated, JSON will be sent to it with a POST request
type WebhookReports struct {
	ID      int `json:"id"`
	Display int `json:"display"`
	Reports struct {
		FacingCount []ReportFacingCountJson `json:"FACING_COUNT_1_5"`
		PriceTags   []ReportPriceTagsJson   `json:"PRICE_TAGS"`
		Realogram   []ReportRealogramJson   `json:"REALOGRAM_1_5"`
	}
}

// ReportPriceTagsJson represents a unit of data of PRICE_TAGS report
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

// ReportFacingCountJson represents a unit of data of FACING_COUNT report
type ReportFacingCountJson struct {
	Count int `json:"count"`
	SkuId int `json:"sku_id" mapstructure:"sku_id"`
}

// ReportRealogramJson represents a data of PLANOGRAM_COMPLIANCE report
type ReportRealogramJson struct {
	Image            int                               `json:"image"`
	Annotations      []ReportRealogramAnnotations      `json:"annotations"`
	ShelfAnnotations []ReportRealogramShelfAnnotations `json:"shelf_annotations" mapstructure:"shelf_annotations"`
}

// ReportRealogramAnnotations represents a unit of data of RealogramAnnotations report
type ReportRealogramAnnotations struct {
	H         int    `json:"h"`
	W         int    `json:"w"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Name      string `json:"name"`
	SkuId     int    `json:"sku_id" mapstructure:"sku_id"`
	Duplicate bool   `json:"duplicate"`
}

// ReportRealogramShelfAnnotations represents a unit of data of RealogramShelfAnnotations report
type ReportRealogramShelfAnnotations struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

// GetReport requests data of report for the given reportID
func (srv *ReportService) GetReport(ctx context.Context, id int) (*Report, error) {
	path := fmt.Sprintf("reports/%d/", id)
	req, err := srv.client.httpClient.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest(GET, %s):%w", path, err)
	}

	var report Report
	if _, err = srv.client.httpClient.Do(ctx, req, &report); err != nil {
		return nil, fmt.Errorf("failed to Do with Request(GET, %s):%w", path, err)
	}

	return &report, nil
}

// ToPriceTags parses json from Report.Json to []ReportPriceTagsJson
func (srv *ReportService) ToPriceTags(v any) ([]ReportPriceTagsJson, error) {
	var r []ReportPriceTagsJson
	if err := mapstructure.WeakDecode(v, &r); err != nil {
		return r, fmt.Errorf("failed to WeakDecode %v:%w", v, err)
	}
	return r, nil
}

// ToFacingCount parses json from Report.Json to []ReportFacingCountJson
func (srv *ReportService) ToFacingCount(v any) ([]ReportFacingCountJson, error) {
	var r []ReportFacingCountJson
	if err := mapstructure.WeakDecode(v, &r); err != nil {
		return r, fmt.Errorf("failed to WeakDecode %v:%w", v, err)
	}
	return r, nil
}

// ToRealogram parses json from Report.Json to []ReportRealogramJson
func (srv *ReportService) ToRealogram(v any) ([]ReportRealogramJson, error) {
	var r []ReportRealogramJson
	if err := mapstructure.WeakDecode(v, &r); err != nil {
		return r, fmt.Errorf("failed to WeakDecode %v:%w", v, err)
	}
	return r, nil
}

// ParseWebhookReports parses json from Webhook request to WebhookReports
func (srv *ReportService) ParseWebhookReports(b []byte) (*WebhookReports, error) {
	var reports WebhookReports
	if err := json.Unmarshal(b, &reports); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal %q:%w", b, err)
	}
	return &reports, nil
}

// WaitForReport polls until the report is READY or ERROR.
// Context ctx is used for cancellation and timeout.
func (srv *ReportService) WaitForReport(ctx context.Context, id int, opts *ReportWaitOptions) (*Report, error) {
	options := applyReportWaitDefaults(opts)
	ctx, cancel := withReportWaitTimeout(ctx, options.Timeout)
	if cancel != nil {
		defer cancel()
	}

	interval := options.Interval
	for attempt := 1; ; attempt++ {
		report, err := srv.GetReport(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to GetReport(%d):%w", id, err)
		}
		if options.OnProgress != nil {
			options.OnProgress(report)
		}
		switch report.Status {
		case ReportStatusREADY:
			return report, nil
		case ReportStatusERROR:
			return nil, fmt.Errorf("failed to WaitForReport(%d) with status ERROR", id)
		}

		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("failed to WaitForReport(%d):%w", id, ctx.Err())
		case <-timer.C:
			timer.Stop()
		}

		if options.Backoff != nil {
			interval = options.Backoff(attempt, interval)
		}
	}
}

func applyReportWaitDefaults(opts *ReportWaitOptions) ReportWaitOptions {
	if opts == nil {
		return ReportWaitOptions{
			Interval: ReportWaitDefaultInterval,
			Timeout:  ReportWaitDefaultTimeout,
		}
	}

	options := *opts
	if options.Interval == 0 {
		options.Interval = ReportWaitDefaultInterval
	}
	if options.Timeout == 0 {
		options.Timeout = ReportWaitDefaultTimeout
	}
	return options
}

func withReportWaitTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, timeout)
}
