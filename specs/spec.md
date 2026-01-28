# go-inspector Project Specification

## Project Overview

**go-inspector** is a Go SDK/API client library for Inspector Cloud, a retail image recognition and analysis service. It provides a clean, service-oriented architecture for interacting with the Inspector Cloud API.

- **Repository:** github.com/germangorelkin/go-inspector
- **License:** MIT License (Copyright 2021 German Gorelkin)
- **Go Version:** 1.24
- **Current Version:** v1.1.0
- **Package:** `inspector`

## Architecture Overview

### Service-Oriented Design

The SDK follows a service-oriented architecture pattern where a central `Client` orchestrates specialized service modules:

```
Client
├── ImageService      → Upload images
├── RecognizeService  → Trigger recognition
├── ReportService     → Retrieve and parse reports
├── SkuService        → Manage SKU data
└── VisitService      → Create visits
```

Each service:
- Is instantiated by the `Client` constructor
- Shares the same HTTP client instance
- Provides domain-specific methods
- Accepts `context.Context` as first parameter for cancellation/timeout support

### HTTP Communication Layer

- Uses custom HTTP client wrapper: `github.com/germangorelkin/http-client`
- All API calls go through the shared `httpClient` instance
- Supports verbose logging when `Verbose: true` in configuration
- API base URL format: `https://{instance}.inspector-cloud.ru/api/v1.5/`

## Data Models

### Core Types

#### Client Configuration
```go
type ClientConf struct {
    Instance   string        // API base URL
    APIKey     string        // Authentication key
    Verbose    bool          // Enable HTTP logging
    HTTPClient *http.Client  // Optional custom HTTP client
    Timeout    time.Duration // Optional HTTP timeout (default 30s)
}

// Historical typo retained via alias for backward compatibility.
type ClintConf = ClientConf
```

#### Pagination
```go
type Pagination struct {
    Count    int          // Total count of items
    Next     *string      // URL to next page (nil if last page)
    Previous *string      // URL to previous page (nil if first page)
    Results  interface{}  // Actual results (type varies by endpoint)
}
```

### Domain Models

#### Image
```go
type Image struct {
    ID          int
    URL         string
    Width       int
    Height      int
    CreatedDate time.Time
}
```

**Upload Methods:**
- ✅ `UploadByURL(ctx, url)` - Upload from URL
- ✅ `Upload(ctx, reader, filename)` - Direct file upload via multipart/form-data

#### Recognition

**Request:**
```go
type RecognizeRequest struct {
    Images      []int      // List of IC image IDs
    ReportTypes []string   // Types of reports to generate
    Display     int        // Display ID (optional)
    Visit       int        // Visit ID (optional)
    Datetime    *time.Time // Recognition timestamp (optional)
    Webhook     string     // Webhook URL for async notification
    CountryCode string     // Country code for recognition
}
```

**Response:**
```go
type RecognizeResponse struct {
    ID      int            // Recognition ID
    Images  []int          // Processed image IDs
    Display int            // Display ID
    Scene   string         // Scene type detected
    Reports map[string]int // Report type → Report ID mapping
}
```

#### Report

**Report Types (Constants):**
- `FACING_COUNT` - Count product facings
- `SHARE_OF_SPAC` - Share of space analysis
- `REALOGRAM` - Visual shelf layout with annotations
- `PRICE_TAGS` - Price tag recognition
- `MHL_COMPLIANCE` - Minimum handling level compliance
- `PLANOGRAM_COMPLIANCE` - Planogram compliance check

**Report Status (Constants):**
- `NOT_READY` - Processing in progress
- `READY` - Report available
- `ERROR` - Processing failed

**Base Report:**
```go
type Report struct {
    ID          int         // Report ID
    Status      string      // Status constant
    Recognition int         // Recognition ID
    ReportType  string      // Report type constant
    Json        interface{} // Report-specific data
}
```

**Specialized Report Types:**

1. **Price Tags:**
```go
type ReportPriceTagsJson struct {
    Version     int
    CategoryID  int
    PriceTag    []PriceTag
    Annotations []Annotation
}

type PriceTag struct {
    Price float64
    Sku   *int      // SKU ID (nullable)
    Box   []float64 // Bounding box [x1, y1, x2, y2]
}
```

2. **Facing Count:**
```go
type ReportFacingCountJson struct {
    Version int
    Sku     *int
    Count   int
}
```

3. **Realogram:**
```go
type ReportRealogramJson struct {
    Version     int
    CategoryID  int
    Products    []Product
    Annotations []Annotation
}

type Product struct {
    Sku         *int
    Box         []float64 // Bounding box
    Facing      int
    Probability float64
    Layer       int
}
```

4. **Webhook Reports:**
```go
type WebhookReports struct {
    RecognitionID int
    Status        string
    Reports       []WebhookReport
}

type WebhookReport struct {
    ID         int
    ReportType string
    Status     string
}
```

#### SKU (Stock Keeping Unit)
```go
type Sku struct {
    ID           int
    CID          string     // Custom ID
    EAN13        *string    // Barcode (nullable)
    Image        int        // Image ID
    Name         string     // Product name
    Brand        *int       // Brand ID (nullable)
    Category     *int       // Category ID (nullable)
    Manufacturer *int       // Manufacturer ID (nullable)
    SizeXMM      *float64   // Width in mm (nullable)
    SizeYMM      *float64   // Height in mm (nullable)
    SizeZMM      *float64   // Depth in mm (nullable)
}
```

Reference: `https://help.inspector-cloud.com/docs/api/backend/methods/v1.5/catalog/sku`

#### Visit
```go
type Visit struct {
    ID          int
    Shop        int       // Shop ID
    Agent       string    // Merchandiser/agent name
    StartedDate time.Time // Visit start time
    Latitude    float64   // GPS latitude
    Longitude   float64   // GPS longitude
}
```

## API Workflow

### Typical Recognition Flow

1. **Upload Image(s)**
   ```go
   img, err := client.Image.UploadByURL(ctx, "https://example.com/shelf.jpg")
   ```

2. **Create Visit (Optional)**
   ```go
   visit, err := client.Visit.AddVisit(ctx)
   ```

3. **Trigger Recognition**
   ```go
   recResp, err := client.Recognize.Recognize(ctx, inspector.RecognizeRequest{
       Images:      []int{img.ID},
       ReportTypes: []string{inspector.ReportTypeFACING_COUNT},
       Visit:       visit.ID,
       Webhook:     "https://myapp.com/webhook",
   })
   ```

4. **Poll for Report** (if no webhook)
   ```go
   reportID := recResp.Reports[inspector.ReportTypeFACING_COUNT]
   
   for {
       report, err := client.Report.GetReport(ctx, reportID)
       if report.Status == inspector.ReportStatusREADY {
           break
       }
       time.Sleep(2 * time.Second)
   }
   ```

5. **Parse Report**
   ```go
   facingCount, err := client.Report.ToFacingCount(report.Json)
   ```

### Asynchronous Processing

- Recognition is **always asynchronous**
- `Recognize()` returns immediately with report IDs
- Reports must be polled until `Status == READY`
- Alternative: Provide webhook URL for push notification

#### Polling Helper

Use `WaitForReport` to poll report status with defaults (2s interval, 60s timeout) and optional callbacks.

```go
report, err := client.Report.WaitForReport(ctx, reportID, &inspector.ReportWaitOptions{
    Interval: 2 * time.Second,
    Timeout:  60 * time.Second,
    OnProgress: func(report *inspector.Report) {
        log.Printf("report status: %s", report.Status)
    },
})
if err != nil {
    return err
}
```

### Webhook Integration

When a webhook URL is provided:
1. Recognition completes asynchronously
2. Inspector Cloud POSTs to webhook URL when ready
3. Payload format: `WebhookReports` struct
4. Parse with: `ParseWebhookReports(requestBody)`

## Pagination

Most list endpoints (including `/sku`) return standard pagination fields:

- `count`: total number of objects
- `next`: URL for the next page (nullable)
- `previous`: URL for the previous page (nullable)
- `results`: list of objects

Reference: `https://help.inspector-cloud.com/docs/api/backend/methods/pagination`
SKU endpoint: `https://help.inspector-cloud.com/docs/api/backend/methods/v1.5/catalog/sku`

## Testing Strategy

### Framework
- **Library:** `github.com/stretchr/testify/assert`
- **Mocking:** `net/http/httptest` for HTTP server mocking

### Test Organization

- One test file per source file: `service.go` → `service_test.go`
- Test function naming: `TestService_Method`
- Use `t.Run()` for sub-tests and variations

### Test Pattern Template

```go
func TestService_Method(t *testing.T) {
    // Create mock HTTP server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request method, path, headers
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "/endpoint/", r.URL.Path)
        
        // Validate request body
        body, err := ioutil.ReadAll(r.Body)
        assert.NoError(t, err)
        
        // Return mock response
        fmt.Fprintln(w, `{"id": 123, "status": "success"}`)
    }))
    defer ts.Close()
    
    // Create client with test server
    client, err := NewClient(ClintConf{
        Instance: ts.URL,
        APIKey:   "test-key",
    })
    assert.NoError(t, err)
    
    // Execute method
    result, err := client.Service.Method(context.Background(), params)
    
    // Assert results
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test Data Management

- Simple data: Inline JSON strings
- Complex data: `testdata/` directory
  - `webhook_reports.json` - Webhook payload example
  - `REALOGRAM_1_5.json` - Realogram report structure

### Coverage Requirements

✅ **Currently Tested:**
- ImageService: UploadByURL
- RecognizeService: Recognize
- ReportService: GetReport, ToPriceTags, ToFacingCount, ToRealogram, ParseWebhookReports
- SkuService: ToSku

❌ **Not Tested:**
- Client initialization and configuration
- VisitService: AddVisit
- RecognizeService: RecognitionError
- Error handling edge cases
- Context cancellation behavior

## Build and Deployment

### Dependencies

**External:**
- `github.com/germangorelkin/http-client` v0.7.0 - Custom HTTP client wrapper
- `github.com/mitchellh/mapstructure` v1.1.2 - Dynamic type conversion
- `github.com/stretchr/testify` v1.7.0 - Testing framework

**Standard Library:**
- `context` - Request context management
- `encoding/json` - JSON marshaling
- `fmt` - String formatting and error wrapping
- `time` - Timestamps and durations
- `net/http` - HTTP client configuration

### Build Commands

This is a library, not an application:

```bash
# Verify library builds
go build ./...

# Install dependencies
go mod download
go mod tidy

# Build CLI example (for testing)
go build -o cli cmd/cli/main.go
```

### Testing Commands

```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# With coverage
go test ./... -cover

# Specific package
go test ./inspector -v

# Single test
go test ./inspector -run TestImageService_UploadByURL -v
```

### CI/CD

**GitHub Actions** (`.github/workflows/test.yml`):
- Triggers: Push, Pull Request
- Platforms: Ubuntu, macOS, Windows
- Go Version: 1.24 (sourced from go.mod)
- Command: `go test ./... -v`

**No Additional Tooling:**
- No Makefile
- No linting (golangci-lint, etc.)
- No code formatting enforcement
- No pre-commit hooks
- No code generation

## Environment Configuration

### Required Environment Variables (for CLI examples)

```bash
export API_KEY="your_api_key_here"
export INSTANCE="https://instance.inspector-cloud.ru/api/v1.5/"
```

### Configuration Notes

- API key authentication (no OAuth/JWT)
- Instance URL is customer-specific
- API version is part of URL path (`/api/v1.5/`)
- No configuration file support (environment variables only)

## Known Limitations and TODOs

### Current Limitations

1. **Direct File Upload Buffers In Memory**
   - `Upload(ctx, reader, filename)` uses multipart/form-data
   - The multipart helper buffers data in memory (avoid very large uploads)

2. **No Automatic Report Polling**
   - ✅ `WaitForReport(ctx, reportID, opts)` helper available
   - Supports timeout, interval, optional backoff, and progress callbacks

3. **Pagination Helpers**
   - `GetSKU()` returns single page
   - `IterateSKU()` provides automatic pagination with iterator pattern
   - `GetAllSKU()` fetches all pages automatically
   - Includes safeguards against infinite loops

4. **Limited Error Context**
   - HTTP errors wrapped but not typed
   - No error codes or retry classification
   - Consider: Custom error types for 4xx vs 5xx

5. **Type Name Typo**
   - `ClintConf` should be `ClientConf`
   - Breaking change to fix (requires major version bump)

### Future Enhancements

**Priority: High**
- Document direct file upload method
- Add polling helper with timeout/retry

**Priority: Medium**
- Add request/response logging hooks
- Support custom HTTP middleware
- Add metrics/instrumentation hooks

**Priority: Low**
- Add mock client for testing
- Add helper for batch image upload
- Add report caching layer

## Version History

**Current:** v1.1.0

**Tags:** v0.5.0 → v1.1.0 (18 versions)

## Support and Resources

- **Documentation:** README.md
- **Examples:** cmd/cli/main.go (mostly commented out)
- **HTTP Examples:** cmd/cli/examples.http
- **License:** MIT License
