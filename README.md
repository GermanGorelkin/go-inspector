# go-inspector

# go-inspector

Go SDK for [Inspector Cloud](https://inspector-cloud.ru/) retail image recognition API. Provides a service-oriented client with typed helpers for uploads, recognition, reports, SKU data, and visits. Tested on Go 1.24.

- **Current version:** v1.1.0
- **License:** MIT
- **Status:** Production-ready client

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

### 3. Upload an image from a reader

```go
file, err := os.Open("/path/to/shelf.jpg")
if err != nil {
	log.Fatal(err)
}
defer file.Close()

img, err = cli.Image.Upload(ctx, file, "shelf.jpg")
if err != nil {
	log.Fatalf("upload failed: %v", err)
}
```

> Direct uploads use multipart/form-data and buffer the file in memory via the HTTP client helper.

### 4. Trigger recognition

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

### 5. Poll for report readiness

```go
// Use defaults (2s interval, 60s timeout)
report, err := cli.Report.WaitForReport(ctx, reportID, nil)
if err != nil {
	log.Fatalf("wait for report failed: %v", err)
}

// Or customize with options
report, err = cli.Report.WaitForReport(ctx, reportID, &inspector.ReportWaitOptions{
	Interval: 3 * time.Second,
	OnProgress: func(r *inspector.Report) {
		log.Printf("status: %s", r.Status)
	},
})
if err != nil {
	log.Fatalf("wait for report failed: %v", err)
}

facing, err := cli.Report.ToFacingCount(report.Json)
if err != nil {
	log.Fatalf("parse facing: %v", err)
}
log.Printf("facing count: %+v", facing)
```

Webhook users can parse payloads with `inspector.ParseWebhookReports(body)`.

### Additional helpers

- **Visits:** `cli.Visit.AddVisit(ctx)` (request body `{}`; server assigns defaults)
- **SKU pagination:**

  ```go
  // Manual pagination
  page, _ := cli.Sku.GetSKU(ctx, offset, limit)
  skus, _ := cli.Sku.ToSku(page.Results)
  
  // Automatic pagination with iterator
  iterator := cli.Sku.IterateSKU(ctx, 100) // 100 items per page
  for {
      pageSkus, err := iterator.Next()
      if err != nil {
          // handle error
      }
      if pageSkus == nil {
          break // no more pages
      }
      // process pageSkus...
  }
  
  // Or fetch all SKUs at once (uses automatic pagination)
  allSkus, err := cli.Sku.GetAllSKU(ctx, 100)
  ```

  Pagination responses follow the standard `count/next/previous/results` format documented here:
  `https://help.inspector-cloud.com/docs/api/backend/methods/pagination`.
  SKU endpoint reference: `https://help.inspector-cloud.com/docs/api/backend/methods/v1.5/catalog/sku`.

- **Report converters:** `ToPriceTags`, `ToFacingCount`, `ToRealogram`, `ToSku`

## Architecture & Services

| Service | Purpose | Key methods |
| --- | --- | --- |
| `ImageService` | Upload shelf photos | `UploadByURL`, `Upload` |
| `RecognizeService` | Trigger recognition jobs | `Recognize` |
| `ReportService` | Retrieve/parse reports | `GetReport`, `ToFacingCount`, `ToPriceTags`, `ToRealogram`, `ParseWebhookReports` |
| `SkuService` | Work with SKU catalogs | `GetSKU`, `ToSku`, `IterateSKU`, `GetAllSKU` |
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

Standard workflow (see `AGENTS.md` for full guidelines). You can use the included `Makefile` for common tasks:

```bash
# Format, lint, test, and build
make all

# Individual targets
make fmt        # Format code
make lint       # Run go vet and gofmt
make test       # Run tests
make coverage   # Run tests with coverage report
make build      # Build library and CLI
make clean      # Clean up artifacts
```

- Run tests on Go 1.24 or newer
- Use `go test ./inspector -run TestImageService_UploadByURL -v` for targeted checks
- Keep imports organized (std lib → blank line → third-party)

## Resources

- `specs/spec.md` – complete specification and workflows
- `AGENTS.md` – contributor guidelines
- `cmd/cli/examples.http` – raw API examples
- `inspector/testdata/` – sample JSON fixtures
- SKU API reference: `https://help.inspector-cloud.com/docs/api/backend/methods/v1.5/catalog/sku`

## Support

- File issues or feature requests via GitHub Issues
- Contributions welcome (open a PR following the planning workflow)
