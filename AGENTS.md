# Agent Guidelines for go-inspector

This document provides guidelines for AI coding agents working on the go-inspector codebase. Read this carefully before making any changes.

## Quick Reference

**Project:** Go SDK for Inspector Cloud API (retail image recognition service)  
**Package:** `inspector`  
**Go Version:** 1.24  
**Architecture:** Service-oriented client library with 5 domain services

## Essential Documentation

**📋 Full Specification:** `specs/spec.md` - Complete project architecture, data models, API workflows, testing strategy, and known limitations. **Read this first** before making any changes.

## Build, Test, and Lint Commands

### Building
```bash
# Verify library builds (this is a library, not an app)
go build ./...

# Build CLI example
go build -o cli cmd/cli/main.go

# Tidy dependencies
go mod tidy
```

### Testing
```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# With coverage
go test ./... -cover

# Test specific package
go test ./inspector -v

# Run single test
go test ./inspector -run TestImageService_UploadByURL -v
```

### Linting
```bash
# No linting configured yet - use standard tools
go vet ./...
gofmt -l .
```

## Code Style Guidelines

### 1. Import Organization

**Always use this order with blank line separator:**

```go
package inspector

import (
    // Standard library imports (alphabetical)
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    // Blank line separator
    
    // Third-party imports (alphabetical)
    "github.com/germangorelkin/http-client"
    "github.com/mitchellh/mapstructure"
)
```

### 2. Error Handling Pattern

**Use consistent error wrapping with context:**

```go
// Pattern: fmt.Errorf("failed to [action] [details]:%w", contextInfo, err)

// ✅ Good examples:
return nil, fmt.Errorf("failed to build http-client:%w", err)
return nil, fmt.Errorf("failed to NewRequest(POST, recognize/, %v):%w", rr, err)
return nil, fmt.Errorf("failed to Do with Request(POST, uploads/, %v):%w", body, err)

// ❌ Bad examples:
return nil, err                                    // No context
return nil, fmt.Errorf("error: %v", err)          // No action description, no %w
return nil, errors.New("request failed")          // Loses original error
```

**Rules:**
- Always use `%w` verb for error wrapping (enables `errors.Is()` and `errors.As()`)
- Prefix with "failed to [action]"
- Include relevant context (method name, endpoint, request data)
- Never return raw errors without wrapping

### 3. Naming Conventions

**Types:**
- Exported: `PascalCase` (Client, RecognizeService, Image)
- Services: `XxxService` suffix (ImageService, SkuService)
- Request/Response: `XxxRequest`, `XxxResponse` suffix
- JSON reports: `ReportXxxJson` suffix
- Private: `camelCase`

**Functions:**
- Exported: `PascalCase` (NewClient, Recognize, GetReport)
- Converters: `ToXxx` prefix (ToPriceTags, ToFacingCount, ToSku)
- Parsers: `ParseXxx` prefix (ParseWebhookReports)
- Private: `camelCase`

**Variables:**
- Local: `camelCase`, short names in small scopes (`c`, `err`, `b`)
- Descriptive in larger scopes (`report`, `priceTags`, `facingCount`)

**Constants:**
- Exported: `SCREAMING_SNAKE_CASE` (ReportTypeFACING_COUNT, ReportStatusREADY)

### 4. Type Definitions

**Struct Tags:**
```go
type Image struct {
    ID          int       `json:"id"`                    // Basic tag
    URL         string    `json:"url,omitempty"`         // Optional field
    CreatedDate time.Time `json:"created_date"`          // Snake case for API
}

type Sku struct {
    SizeXMM *float64 `json:"size_x_mm,omitempty" mapstructure:"size_x_mm"` // Multiple tags
}
```

**Rules:**
- JSON tags use `snake_case` (API convention)
- Use `omitempty` for optional fields
- Use `mapstructure` tags for dynamic conversion
- Use pointers (`*string`, `*int`, `*float64`) for nullable fields

### 5. Documentation

**Godoc format:**
```go
// Service provides access to the [Feature] functions in the IC API.
type Service struct {}

// MethodName does [action] and returns [result].
// Context ctx is used for cancellation and timeout.
func (srv *Service) MethodName(ctx context.Context, ...) (..., error) {}
```

**Inline field comments:**
```go
type RecognizeRequest struct {
    Images      []int    `json:"images"`       // list of IC image IDs
    ReportTypes []string `json:"report_types"` // list of reports to be generated
    Webhook     string   `json:"webhook"`      // webhook URL for async notification
}
```

### 6. Context Usage

**All API methods must accept context.Context as first parameter:**

```go
// ✅ Good:
func (srv *ImageService) UploadByURL(ctx context.Context, url string) (Image, error)
func (srv *RecognizeService) Recognize(ctx context.Context, rr RecognizeRequest) (*RecognizeResponse, error)

// ❌ Bad:
func (srv *ImageService) UploadByURL(url string) (Image, error)  // Missing context
```

**Benefits:**
- Request cancellation
- Timeout propagation
- Deadline enforcement
- Request-scoped values

## Testing Requirements

### Test Framework
- Use `github.com/stretchr/testify/assert` for assertions
- Use `net/http/httptest` for HTTP mocking
- Use `context.Background()` for test contexts

### Test File Organization
```bash
# Pattern: service.go → service_test.go
inspector/
├── client.go
├── image.go
├── image_test.go        # Tests for image.go
├── recognize.go
├── recognize_test.go    # Tests for recognize.go
```

### Test Naming
```go
// Format: TestService_Method
func TestImageService_UploadByURL(t *testing.T) { }
func TestRecognizeService_Recognize(t *testing.T) { }

// For variations, use t.Run():
func TestReportService_ToRealogram(t *testing.T) {
    t.Run("valid report", func(t *testing.T) { })
    t.Run("empty report", func(t *testing.T) { })
}
```

### Test Template
```go
func TestService_Method(t *testing.T) {
    // Setup mock HTTP server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request
        assert.Equal(t, "POST", r.Method)
        assert.Equal(t, "/endpoint/", r.URL.Path)
        
        // Read and validate body
        body, err := ioutil.ReadAll(r.Body)
        assert.NoError(t, err)
        
        // Return mock response
        fmt.Fprintln(w, `{"id": 123}`)
    }))
    defer ts.Close()
    
    // Create test client
    client, err := NewClient(ClintConf{
        Instance: ts.URL,
        APIKey:   "test-key",
    })
    assert.NoError(t, err)
    
    // Test the method
    result, err := client.Service.Method(context.Background(), params)
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test Data
- Simple: Inline JSON strings
- Complex: `testdata/filename.json`

### Assertions
```go
// ✅ Good:
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.Len(t, result, 2)
assert.NotNil(t, ptr)

// ❌ Avoid:
if err != nil {
    t.Fatal(err)  // Use assert.NoError instead
}
```

## Architecture Patterns

### Service-Oriented Structure

```
Client (orchestrator)
├── ImageService      → Image uploads
├── RecognizeService  → Recognition triggers
├── ReportService     → Report retrieval and parsing
├── SkuService        → SKU data management
└── VisitService      → Visit creation
```

**Rules:**
- Each service has its own file: `service.go`
- All services share the same HTTP client instance
- Services are created in `NewClient()` constructor
- Services should not depend on each other

### Type Conversion Pattern

Use `ToXxx()` methods for parsing dynamic API responses:

```go
// Get report as interface{}
report, err := client.Report.GetReport(ctx, reportID)

// Convert to specific type
priceTags, err := client.Report.ToPriceTags(report.Json)
facingCount, err := client.Report.ToFacingCount(report.Json)
realogram, err := client.Report.ToRealogram(report.Json)
```

## Common Pitfalls

### Known Issues

**⚠️ Type Name Typo:**
```go
type ClintConf struct { }  // Note: "Clint" not "Client"
```
This is a known typo in the public API. Do NOT fix without major version bump.

**⚠️ TODO - File Upload:**
```go
// Direct file upload is NOT implemented - only URL upload works
func (srv *ImageService) Upload(r io.Reader, filename string) (*Image, error) {
    // TODO
}
```

**⚠️ Asynchronous Recognition:**
- `Recognize()` returns immediately with report IDs
- Reports must be polled until `Status == READY`
- No built-in polling helper (manual loop required)

### API Quirks

1. **API uses snake_case, Go uses camelCase:**
   ```go
   type Image struct {
       CreatedDate time.Time `json:"created_date"`  // API: created_date, Go: CreatedDate
   }
   ```

2. **Nullable fields use pointers:**
   ```go
   type Sku struct {
       EAN13 *string `json:"ean13,omitempty"`  // Can be null in API
   }
   ```

3. **Report processing is async:**
   ```go
   // Must poll until ready
   for {
       report, _ := client.Report.GetReport(ctx, reportID)
       if report.Status == inspector.ReportStatusREADY {
           break
       }
       time.Sleep(2 * time.Second)
   }
   ```

## Development Workflow for New Changes

### Planning Phase

When starting any new feature, bug fix, or significant change, **ALWAYS** follow this workflow:

#### 1. Create Detailed Plan

**Create a plan file:** `.opencode/plan/YYYY-MM-DD_task_name.md`

**Plan must include:**
- Problem statement and context
- Proposed solution with detailed steps
- Files to be modified/created
- Testing strategy
- Potential risks and edge cases
- Rollback strategy (if applicable)

**Example filename:**
```bash
.opencode/plan/2026-01-24_add_report_polling_helper.md
.opencode/plan/2026-01-24_fix_timeout_handling.md
.opencode/plan/2026-01-24_implement_file_upload.md
```

#### 2. Ask Clarifying Questions

**Before implementing, ask about:**
- Ambiguous requirements
- Design decisions with multiple valid approaches
- Breaking changes that affect public API
- Test coverage expectations
- Performance implications
- Backward compatibility concerns

**Example questions:**
```
"Should the polling helper return an error after N attempts or continue indefinitely?"
"Do we need to support both sync and async polling modes?"
"Should this be a breaking change (v2.0.0) or can we maintain backward compatibility?"
```

#### 3. Iterative Development

**Development MUST be iterative with confirmation:**

1. **Implement Step 1** → Request confirmation
2. **After confirmation** → Implement Step 2 → Request confirmation
3. **After confirmation** → Implement Step 3 → Request confirmation
4. Continue until complete

**Benefits:**
- Catch issues early
- Adjust approach based on feedback
- Ensure alignment with expectations
- Maintain control over large changes

**Example workflow:**
```
Agent: "I've completed Step 1: Added the polling helper function. 
        Please review before I proceed to Step 2 (adding tests)."
        
User: "Looks good, but change the timeout to 30s instead of 60s"

Agent: "Updated timeout to 30s. I've completed Step 2: Added unit tests.
        Please review before I proceed to Step 3 (updating documentation)."
```

#### 4. Track Progress

- Immediately mark each plan step or tracked task as completed once it’s done so reviewers can see progress without rereading the whole plan.
- Keep unchecked items accurate; if the scope changes, update the checklist instead of leaving stale entries.

### Plan Template

Use this template for `.opencode/plan/YYYY-MM-DD_task_name.md`:

```markdown
# Task: [Brief Description]

**Date:** YYYY-MM-DD  
**Status:** Planning | In Progress | Completed | Cancelled

## Problem Statement

[Describe what needs to be done and why]

## Proposed Solution

[High-level approach]

## Detailed Steps

1. [ ] Step 1: [Description]
   - Files: `file1.go`, `file2.go`
   - Changes: [Brief description]
   
2. [ ] Step 2: [Description]
   - Files: `file3_test.go`
   - Changes: [Brief description]

3. [ ] Step 3: [Description]
   - Files: `README.md`
   - Changes: [Brief description]

## Testing Strategy

- [ ] Unit tests: [Description]
- [ ] Integration tests: [Description]
- [ ] Manual testing: [Description]

## Open Questions

1. [Question 1]?
2. [Question 2]?

## Risks and Edge Cases

- Risk 1: [Description and mitigation]
- Edge case 1: [Description and handling]

## Rollback Strategy

[How to undo this change if needed]
```

## Dependencies

**External:**
- `github.com/germangorelkin/http-client` v0.7.0 - Custom HTTP wrapper
- `github.com/mitchellh/mapstructure` v1.1.2 - Dynamic type conversion
- `github.com/stretchr/testify` v1.7.0 - Testing assertions

**When adding dependencies:**
1. Justify the need (avoid unnecessary dependencies)
2. Check license compatibility (MIT preferred)
3. Verify maintenance status (recent commits, active issues)
4. Update go.mod and go.sum: `go mod tidy`

## CI/CD

**GitHub Actions** runs on every push and PR:
- Platforms: Ubuntu, macOS, Windows
- Go version: 1.24 (sourced from go.mod)
- Command: `go test ./... -v`

**Before pushing:**
```bash
# Ensure tests pass
go test ./... -v

# Check for issues
go vet ./...

# Format code
gofmt -w .

# Verify build
go build ./...
```

## Common Tasks

### Adding a New Service

1. Create `newservice.go` in `inspector/` package
2. Define service struct: `type NewService struct { c *Client }`
3. Add methods with context: `func (srv *NewService) Method(ctx context.Context, ...) (..., error)`
4. Add service to Client struct in `client.go`
5. Initialize in `NewClient()` constructor: `c.New = &NewService{c: c}`
6. Create `newservice_test.go` with test coverage
7. Update `specs/spec.md` with new service documentation

### Adding a New Report Type

1. Add constant in `report.go`: `const ReportTypeNEW_TYPE = "NEW_TYPE"`
2. Define struct: `type ReportNewTypeJson struct { ... }`
3. Add converter: `func (srv *ReportService) ToNewType(v interface{}) ([]ReportNewTypeJson, error)`
4. Add test in `report_test.go`: `func TestReportService_ToNewType(t *testing.T)`
5. Add example JSON in `testdata/new_type.json`
6. Update `specs/spec.md` with new report type

### Fixing a Bug

1. Write failing test that reproduces the bug
2. Fix the bug
3. Verify test passes
4. Add regression test if needed
5. Update documentation if behavior changed

## Resources

- **Full Spec:** `specs/spec.md` - Complete architecture and data models
- **API Docs:** README.md - Usage examples
- **HTTP Examples:** cmd/cli/examples.http - Raw API requests
- **Test Data:** inspector/testdata/ - JSON fixtures

## Questions?

When uncertain about:
- Architecture decisions → Check `specs/spec.md`
- Code style → Check existing files for patterns
- Testing approach → Check existing `*_test.go` files
- API behavior → Check `README.md` and `examples.http`

**When in doubt, ask the user before implementing.**
