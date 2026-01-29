# Task: Refactor CLI to Individual Binary Examples

**Date:** 2026-01-29  
**Status:** Completed

## Problem Statement

The current `cmd/cli/` directory contains a single `main.go` file with mostly commented-out code demonstrating various library features. This approach has several issues:

1. Hard to discover and understand individual features
2. Commented code is not tested or maintained
3. Single monolithic file doesn't scale well
4. Not idiomatic for Go library examples

**Goal:** Create individual, self-contained CLI examples in `examples/` directory, each demonstrating a specific feature of the `go-inspector` library. Remove the old `cmd/cli/` directory after migration.

## Proposed Solution

Create an `examples/` directory with individual subdirectories, each containing a standalone `main.go` that demonstrates one specific library feature. Use standard library `flag` package for argument parsing and environment variables for credentials.

### Directory Structure

```
examples/
├── README.md                    # Overview of all examples
├── upload-url/
│   └── main.go                  # Upload image by URL
├── upload-file/
│   └── main.go                  # Upload image from local file
├── recognize/
│   └── main.go                  # Full recognition workflow
├── get-report/
│   └── main.go                  # Get and parse existing report
├── wait-report/
│   └── main.go                  # Poll report until ready
├── sku-list/
│   └── main.go                  # List SKUs with pagination
├── sku-all/
│   └── main.go                  # Fetch all SKUs using iterator
├── visit-create/
│   └── main.go                  # Create a new visit
├── recognition-error/
│   └── main.go                  # Report recognition error
└── full-workflow/
    └── main.go                  # Complete end-to-end example
```

## Detailed Steps

### Phase 1: Setup and Infrastructure

1. [x] **Step 1.1:** Create `examples/` directory structure
   - Files: `examples/README.md`
   - Changes: Create overview documentation explaining all examples, environment setup, and common patterns

2. [x] **Step 1.2:** Create shared helper for client initialization
   - Files: Common pattern documented in README (no shared code to avoid import complications)
   - Changes: Document standard client init pattern using `API_KEY` and `INSTANCE` environment variables

### Phase 2: Image Service Examples

3. [x] **Step 2.1:** Create `examples/upload-url/main.go`
   - Service: `ImageService.UploadByURL`
   - Flags: `-url` (required) - Image URL to upload
   - Output: JSON with image ID, dimensions, URL

4. [x] **Step 2.2:** Create `examples/upload-file/main.go`
   - Service: `ImageService.Upload`
   - Flags: `-file` (required) - Path to local image file
   - Output: JSON with image ID, dimensions, URL

### Phase 3: Recognition Service Examples

5. [x] **Step 3.1:** Create `examples/recognize/main.go`
   - Service: `RecognizeService.Recognize`
   - Flags:
     - `-images` (required) - Comma-separated image IDs
     - `-types` (optional, default: "FACING_COUNT") - Comma-separated report types
     - `-webhook` (optional) - Webhook URL
     - `-visit` (optional) - Visit ID
   - Output: JSON with recognition ID and report IDs

6. [x] **Step 3.2:** Create `examples/recognition-error/main.go`
   - Service: `RecognizeService.RecognitionError`
   - Flags:
     - `-images` (required) - Comma-separated image IDs
     - `-sku` (required) - SKU ID that was mistakenly recognized
     - `-scene` (required) - Scene UUID
     - `-message` (optional) - Error description
   - Output: JSON with recognition error ID

### Phase 4: Report Service Examples

7. [x] **Step 4.1:** Create `examples/get-report/main.go`
   - Service: `ReportService.GetReport`, `To*` converters
   - Flags:
     - `-id` (required) - Report ID
     - `-type` (optional) - Report type for parsing (FACING_COUNT, PRICE_TAGS, REALOGRAM)
   - Output: JSON report data (raw or parsed based on type)

8. [x] **Step 4.2:** Create `examples/wait-report/main.go`
   - Service: `ReportService.WaitForReport`
   - Flags:
     - `-id` (required) - Report ID
     - `-interval` (optional, default: 2s) - Polling interval
     - `-timeout` (optional, default: 60s) - Overall timeout
     - `-type` (optional) - Report type for parsing
   - Output: JSON report data with progress logging

### Phase 5: SKU Service Examples

9. [x] **Step 5.1:** Create `examples/sku-list/main.go`
   - Service: `SkuService.GetSKU`
   - Flags:
     - `-offset` (optional, default: 0) - Pagination offset
     - `-limit` (optional, default: 10) - Items per page
   - Output: JSON with pagination info and SKU list

10. [x] **Step 5.2:** Create `examples/sku-all/main.go`
    - Service: `SkuService.GetAllSKU` or `IterateSKU`
    - Flags:
      - `-page-size` (optional, default: 100) - Items per page
      - `-max` (optional) - Maximum total items to fetch
    - Output: JSON array of all SKUs with progress logging

### Phase 6: Visit Service Example

11. [x] **Step 6.1:** Create `examples/visit-create/main.go`
    - Service: `VisitService.AddVisit`
    - Flags: None (server assigns defaults)
    - Output: JSON with visit ID and details

### Phase 7: Full Workflow Example

12. [x] **Step 7.1:** Create `examples/full-workflow/main.go`
    - Services: All services combined
    - Flags:
      - `-url` or `-file` (one required) - Image source
      - `-types` (optional, default: "FACING_COUNT") - Report types
      - `-wait` (optional, default: true) - Wait for reports
    - Output: Complete workflow with all steps logged
    - Steps: Upload → Create Visit → Recognize → Wait → Parse Reports

### Phase 8: Cleanup

13. [x] **Step 8.1:** Update project documentation
    - Files: `README.md`, `specs/spec.md`
    - Changes: Update references from `cmd/cli/` to `examples/`

14. [x] **Step 8.2:** Remove old `cmd/cli/` directory
    - Files: Remove `cmd/cli/main.go`, `cmd/cli/examples.http`, `cmd/cli/run.sh`, `cmd/cli/test.txt`, `cmd/cli/main`
    - Note: Keep `examples.http` content - move to `examples/` or repository root

15. [x] **Step 8.3:** Update Makefile (if exists)
    - Files: `Makefile`
    - Changes: Update build targets for examples

## Testing Strategy

### Build Verification
- [x] Each example should compile without errors: `go build ./examples/...`

### Integration Tests (Manual)
- [x] Test each example with real API (requires `API_KEY` and `INSTANCE`)
- [x] Verify output format is valid JSON
- [x] Verify error messages are helpful

### Documentation Tests
- [x] Verify README instructions work
- [x] Verify all flags are documented

## Example Code Pattern

Each example will follow this consistent pattern:

```go
package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/germangorelkin/go-inspector/inspector"
)

func main() {
    // Define flags
    var flagName = flag.String("flag", "default", "Flag description")
    flag.Parse()

    // Validate required flags
    if *flagName == "" {
        flag.Usage()
        os.Exit(1)
    }

    // Get credentials from environment
    apiKey := os.Getenv("API_KEY")
    instance := os.Getenv("INSTANCE")
    if apiKey == "" || instance == "" {
        log.Fatal("API_KEY and INSTANCE environment variables must be set")
    }

    // Create client
    client, err := inspector.NewClient(inspector.ClientConf{
        APIKey:   apiKey,
        Instance: instance,
        Timeout:  30 * time.Second,
    })
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Execute operation
    ctx := context.Background()
    result, err := client.Service.Method(ctx, ...)
    if err != nil {
        log.Fatalf("Failed to execute: %v", err)
    }

    // Output JSON result
    output, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(output))
}
```

## Open Questions

None - all questions were clarified before planning.

## Risks and Edge Cases

### Risk 1: API Rate Limiting
- **Description:** Examples that make multiple API calls (e.g., `full-workflow`, `sku-all`) may hit rate limits
- **Mitigation:** Add appropriate delays between calls, document rate limit behavior

### Risk 2: Large SKU Catalogs
- **Description:** `sku-all` example may take very long or run out of memory for large catalogs
- **Mitigation:** Add `-max` flag to limit total items, use streaming output if needed

### Risk 3: Long-Running Operations
- **Description:** `wait-report` and `full-workflow` may run for extended periods
- **Mitigation:** Document expected durations, use context timeouts, show progress

### Edge Case 1: Invalid Image URLs
- **Description:** `upload-url` with invalid/inaccessible URLs
- **Handling:** Let API error propagate with clear message

### Edge Case 2: Report Never Ready
- **Description:** Recognition may fail or hang indefinitely
- **Handling:** `wait-report` has timeout flag, will exit with error

## Rollback Strategy

If issues arise:
1. Keep `cmd/cli/` directory until all examples are verified
2. Examples are additive - can be removed individually without breaking library
3. Git history preserves old code if needed

## Dependencies

No new dependencies required - all examples use:
- Standard library (`flag`, `encoding/json`, `os`, `context`, `time`, `log`, `fmt`)
- Existing project dependency: `github.com/germangorelkin/go-inspector/inspector`

## Success Criteria

1. All examples compile: `go build ./examples/...`
2. All examples have `-help` flag working
3. Each example demonstrates one clear feature
4. `examples/README.md` provides clear usage instructions
5. Old `cmd/cli/` directory removed
6. Project documentation updated