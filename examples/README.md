# go-inspector Examples

This directory contains standalone CLI examples demonstrating various features of the `go-inspector` library. Each example is a separate binary that showcases a specific API functionality.

## Prerequisites

- Go 1.24 or later
- Inspector Cloud API credentials (API key and instance URL)
- Network access to your Inspector Cloud instance

## Environment Setup

All examples require two environment variables:

```bash
export API_KEY="your_api_key_here"
export INSTANCE="https://your-instance.inspector-cloud.ru/api/v1.5/"
```

You can also create a `.env` file (not tracked in git):

```bash
# .env
API_KEY=your_api_key_here
INSTANCE=https://your-instance.inspector-cloud.ru/api/v1.5/
```

Then source it before running examples:

```bash
source .env
```

## Building Examples

Build all examples:

```bash
go build ./examples/...
```

Build a specific example:

```bash
go build -o upload-url ./examples/upload-url
```

## Running Examples

All examples support the `-help` flag to display usage information:

```bash
go run ./examples/upload-url/main.go -help
```



## Available Examples

### Image Upload Examples

#### `upload-url` - Upload Image by URL
Upload an image to Inspector Cloud from a publicly accessible URL.

```bash
go run ./examples/upload-url/main.go -url "https://example.com/shelf.jpg"
```

**Output:** JSON with image ID, dimensions, and URL

#### `upload-file` - Upload Image from Local File
Upload an image to Inspector Cloud from your local filesystem.

```bash
go run ./examples/upload-file/main.go -file "/path/to/image.jpg"
```

**Output:** JSON with image ID, dimensions, and URL

### Recognition Examples

#### `recognize` - Trigger Recognition
Start an asynchronous recognition job on uploaded images.

```bash
go run ./examples/recognize/main.go -images "12345,12346" -types "FACING_COUNT,PRICE_TAGS"
```

**Flags:**
- `-images` (required) - Comma-separated image IDs
- `-types` (optional, default: "FACING_COUNT") - Comma-separated report types
- `-webhook` (optional) - Webhook URL for async notifications
- `-visit` (optional) - Visit ID to associate with recognition

**Output:** JSON with recognition ID and report IDs

**Available Report Types:**
- `FACING_COUNT` - Count product facings
- `PRICE_TAGS` - Price tag recognition
- `REALOGRAM` - Visual shelf layout with annotations
- `SHARE_OF_SPAC` - Share of space analysis
- `MHL_COMPLIANCE` - Minimum handling level compliance
- `PLANOGRAM_COMPLIANCE` - Planogram compliance check

#### `recognition-error` - Report Recognition Error
Submit feedback about incorrect recognition results.

```bash
go run ./examples/recognition-error/main.go \
  -images "12345" \
  -sku 789 \
  -scene "abc-123-def" \
  -message "Product was misidentified"
```

**Flags:**
- `-images` (required) - Comma-separated image IDs
- `-sku` (required) - SKU ID that was mistakenly recognized
- `-scene` (required) - Scene UUID from recognition
- `-message` (optional) - Error description

**Output:** JSON with recognition error ID

### Report Examples

#### `get-report` - Get Report
Retrieve a report by ID and optionally parse it to a specific type.

```bash
go run ./examples/get-report/main.go -id 12345 -type "FACING_COUNT"
```

**Flags:**
- `-id` (required) - Report ID
- `-type` (optional) - Report type for parsing (FACING_COUNT, PRICE_TAGS, REALOGRAM)

**Output:** JSON report data (raw or parsed based on type)

#### `wait-report` - Wait for Report
Poll a report until it's ready, with configurable timeout and interval.

```bash
go run ./examples/wait-report/main.go -id 12345 -interval 3s -timeout 120s -type "PRICE_TAGS"
```

**Flags:**
- `-id` (required) - Report ID
- `-interval` (optional, default: 2s) - Polling interval
- `-timeout` (optional, default: 60s) - Overall timeout
- `-type` (optional) - Report type for parsing

**Output:** JSON report data with progress logging to stderr

### SKU Catalog Examples

#### `sku-list` - List SKUs with Pagination
Retrieve a paginated list of SKUs from the catalog.

```bash
go run ./examples/sku-list/main.go -offset 0 -limit 10
```

**Flags:**
- `-offset` (optional, default: 0) - Pagination offset
- `-limit` (optional, default: 10) - Items per page

**Output:** JSON with pagination info and SKU list

#### `sku-all` - Fetch All SKUs
Retrieve all SKUs using automatic pagination.

```bash
go run ./examples/sku-all/main.go -page-size 100 -max 1000
```

**Flags:**
- `-page-size` (optional, default: 100) - Items per page
- `-max` (optional) - Maximum total items to fetch (omit for all)

**Output:** JSON array of all SKUs with progress logging to stderr

**Warning:** For large catalogs, this may take significant time and memory.

### Visit Example

#### `visit-create` - Create Visit
Create a new visit record (server assigns default values).

```bash
go run ./examples/visit-create/main.go
```

**Flags:** None (server assigns defaults)

**Output:** JSON with visit ID and details

### Complete Workflow Example

#### `full-workflow` - End-to-End Recognition
Complete workflow from image upload through report retrieval.

```bash
# Upload by URL
go run ./examples/full-workflow/main.go -url "https://example.com/shelf.jpg"

# Upload from file
go run ./examples/full-workflow/main.go -file "/path/to/image.jpg"

# Custom report types
go run ./examples/full-workflow/main.go -url "https://example.com/shelf.jpg" -types "FACING_COUNT,PRICE_TAGS"

# Skip waiting for reports
go run ./examples/full-workflow/main.go -url "https://example.com/shelf.jpg" -wait=false
```

**Flags:**
- `-url` or `-file` (one required) - Image source
- `-types` (optional, default: "FACING_COUNT") - Comma-separated report types
- `-wait` (optional, default: true) - Wait for reports to complete

**Output:** Complete workflow with all steps logged

**Steps:**
1. Upload image (by URL or file)
2. Create visit
3. Trigger recognition
4. Wait for reports (if `-wait=true`)
5. Parse and display reports

## Common Patterns

### Client Initialization

All examples follow this standard pattern:

```go
apiKey := os.Getenv("API_KEY")
instance := os.Getenv("INSTANCE")
if apiKey == "" || instance == "" {
    log.Fatal("API_KEY and INSTANCE environment variables must be set")
}

client, err := inspector.NewClient(inspector.ClientConf{
    APIKey:   apiKey,
    Instance: instance,
    Timeout:  30 * time.Second,
})
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

### Error Handling

All examples use consistent error handling:

```go
result, err := client.Service.Method(ctx, params)
if err != nil {
    log.Fatalf("Failed to execute: %v", err)
}
```

### JSON Output

All examples output results as formatted JSON:

```go
output, _ := json.MarshalIndent(result, "", "  ")
fmt.Println(string(output))
```

## Troubleshooting

### Authentication Errors

```
Failed to create client: authentication failed
```

**Solution:** Verify `API_KEY` and `INSTANCE` environment variables are set correctly.

### Network Errors

```
Failed to execute: connection refused
```

**Solution:** Check network connectivity and verify the `INSTANCE` URL is accessible.

### Timeout Errors

```
Failed to execute: context deadline exceeded
```

**Solution:** Increase timeout in client configuration or use longer `-timeout` flag for polling operations.

### Report Not Ready

```
Failed to wait for report: status ERROR
```

**Solution:** Recognition may have failed. Check image quality, format, and API usage limits.

## Additional Resources

- **Main Documentation:** [../README.md](../README.md)
- **API Specification:** [../specs/spec.md](../specs/spec.md)
- **Agent Guidelines:** [../AGENTS.md](../AGENTS.md)
- **Inspector Cloud Docs:** https://help.inspector-cloud.com/docs/api/backend/

## Contributing

When adding new examples:

1. Follow the standard code pattern shown above
2. Use the `flag` package for argument parsing
3. Accept credentials via environment variables
4. Output results as JSON
5. Include helpful error messages
6. Update this README with usage instructions