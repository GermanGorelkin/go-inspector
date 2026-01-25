# go-inspector

# go-inspector

Go SDK for [Inspector Cloud](https://inspectorcloud.com) retail image recognition API. Provides a service-oriented client with typed helpers for uploads, recognition, reports, SKU data, and visits. Tested on Go 1.24.

- **Current version:** v1.1.0
- **License:** MIT
- **Status:** Production-ready client, still missing direct file upload (URL uploads supported)

## Requirements

- Go 1.24+
- Inspector Cloud instance URL and API key
- Network access to `https://{instance}.inspector-cloud.ru/api/v1.5/`

## Installation

```bash
go get github.com/germangorelkin/go-inspector/inspector
```

Add the module to your project and run `go mod tidy`. The package namespace is `inspector`.

## Getting Started

### 1. Initialize the client

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	instance := os.Getenv("INSTANCE") // e.g. https://demo.inspector-cloud.ru/api/v1.5/

	cli, err := inspector.NewClient(inspector.ClientConf{
		APIKey:   apiKey,
		Instance: instance,
		Verbose:  true,              // optional HTTP dump interceptor
		Timeout:  45 * time.Second,  // optional custom timeout (defaults to 30s)
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	_ = ctx // reuse across calls
}
# `ClintConf` alias remains available for backward compatibility.

### 2. Upload an image by URL

```go
img, err := cli.Image.UploadByURL(ctx, "https://example.com/shelf.jpg")
if err != nil {
	log.Fatalf("upload failed: %v", err)
}
```

> Direct file uploads (`Upload(r io.Reader, filename string)`) are not implemented yet.

### 3. Trigger recognition

```go
resp, err := cli.Recognize.Recognize(ctx, inspector.RecognizeRequest{
	Images:      []int{img.ID},
	ReportTypes: []string{inspector.ReportTypeFACING_COUNT, inspector.ReportTypePRICE_TAGS},
	Webhook:     "https://myapp.example/webhooks/inspector",
})
if err != nil {
	log.Fatalf("recognize failed: %v", err)
}

reportID := resp.Reports[inspector.ReportTypeFACING_COUNT]
```

### 4. Poll for report readiness

```go
for {
	report, err := cli.Report.GetReport(ctx, reportID)
	if err != nil {
		log.Fatalf("get report failed: %v", err)
	}
	if report.Status == inspector.ReportStatusREADY {
		facing, err := cli.Report.ToFacingCount(report.Json)
		if err != nil {
			log.Fatalf("parse facing: %v", err)
		}
		log.Printf("facing count: %+v", facing)
		break
	}
	time.Sleep(2 * time.Second)
}
```

Webhook users can parse payloads with `inspector.ParseWebhookReports(body)`.

### Additional helpers

- **Visits:** `cli.Visit.AddVisit(ctx)` (request body `{}`; server assigns defaults)
- **SKU pagination:**

  ```go
  page, _ := cli.Sku.GetSKU(ctx, offset, limit)
  skus, _ := cli.Sku.ToSku(page.Results)
  ```

- **Report converters:** `ToPriceTags`, `ToFacingCount`, `ToRealogram`, `ToSku`

## Architecture & Services

| Service | Purpose | Key methods |
| --- | --- | --- |
| `ImageService` | Upload shelf photos | `UploadByURL`, `Upload` (TODO) |
| `RecognizeService` | Trigger recognition jobs | `Recognize` |
| `ReportService` | Retrieve/parse reports | `GetReport`, `ToFacingCount`, `ToPriceTags`, `ToRealogram`, `ParseWebhookReports` |
| `SkuService` | Work with SKU catalogs | `GetSKU`, `ToSku` |
| `VisitService` | Create visits for merchandisers | `AddVisit` |

All services share the same authenticated HTTP client created by `NewClient`. Every API call accepts `context.Context` as the first parameter for cancellation and deadlines.

### Report types & constants

- `ReportTypeFACING_COUNT`
- `ReportTypePRICE_TAGS`
- `ReportTypeREALOGRAM`
- `ReportTypeSHARE_OF_SPAC`
- `ReportTypeMHL_COMPLIANCE`
- `ReportTypePLANOGRAM_COMPLIANCE`

Statuses: `ReportStatusNOT_READY`, `ReportStatusREADY`, `ReportStatusERROR`.

## Development & Testing

Standard workflow (see `AGENTS.md` for full guidelines):

```bash
go fmt ./...
go vet ./...
go test ./... -v
go build ./...
```

- Run tests on Go 1.24 or newer
- Use `go test ./inspector -run TestImageService_UploadByURL -v` for targeted checks
- Keep imports organized (std lib → blank line → third-party)

## Resources

- `specs/spec.md` – complete specification and workflows
- `AGENTS.md` – contributor guidelines
- `cmd/cli/examples.http` – raw API examples
- `inspector/testdata/` – sample JSON fixtures

## Support

- File issues or feature requests via GitHub Issues
- Contributions welcome (open a PR following the planning workflow)
