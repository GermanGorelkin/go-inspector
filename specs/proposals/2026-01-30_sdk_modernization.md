# Proposal: SDK Modernization and Refinement

**Date:** 2026-01-30  
**Status:** Draft / Proposed  
**Author:** AI Assistant  
**Target Version:** v1.2.0

## 1. Problem Statement

While the current SDK (v1.1.0) is functional and follows a solid service-oriented architecture, it has several areas that could be improved by leveraging modern Go features (v1.24), improving reliability in URL handling, and ensuring consistency across documentation and code.

Key issues:
- Use of `any` in pagination and converters reduces type safety.
- Manual query string building is prone to errors.
- Inconsistent use of `mapstructure.Decode` vs `WeakDecode`.
- Stale comments and documentation typos.
- Rigid webhook report structures.

## 2. Proposed Changes

### 2.1 Type-Safe Pagination with Generics (Go 1.24)
Refactor the `Pagination` struct to use generics, allowing users to receive typed results without manual casting or dynamic decoding.

**Current:**
```go
type Pagination struct {
    Results any `json:"results"`
}
```

**Proposed:**
```go
type Pagination[T any] struct {
    Count    int     `json:"count"`
    Next     *string `json:"next,omitempty"`
    Previous *string `json:"previous,omitempty"`
    Results  []T     `json:"results"`
}
```

### 2.2 Robust Query Parameter Handling
Replace manual string formatting for URL parameters with `url.Values` to ensure proper encoding and prevent malformed URLs.

**Affected Method:** `SkuService.GetSKU`

**Implementation:**
```go
q := req.URL.Query()
q.Set("limit", strconv.Itoa(limit))
q.Set("offset", strconv.Itoa(offset))
req.URL.RawQuery = q.Encode()
```

### 2.3 Standardization on `WeakDecode`
Standardize all dynamic decoding to use `mapstructure.WeakDecode`. This is safer for API responses where the backend might occasionally return a stringified number or vice-versa.

**Affected Files:** `sku.go`, `report.go`.

### 2.4 Flexible Webhook Reports
Refactor `WebhookReports` to handle dynamic report types more gracefully, preventing errors when new report types are added to the IC API.

**Proposed:**
```go
type WebhookReports struct {
    ID      int            `json:"id"`
    Display int            `json:"display"`
    Reports map[string]any `json:"reports"` // Flexible mapping
}
```

### 2.5 Documentation & Comment Cleanup
- Fix "copy-paste" comments in `visit.go`, `sku.go`, and `recognize.go` where structures are described as "ReportPriceTagsJson".
- Update `AGENTS.md` to reflect that the `Upload` method (multipart) is actually implemented.

## 3. Implementation Plan

### Phase 1: Core Refactoring
1. [ ] Update `Pagination` in `client.go` to use generics.
2. [ ] Update `GetSKU` in `sku.go` to use `Pagination[Sku]` and `url.Values`.
3. [ ] Switch `SkuService.ToSku` to use `WeakDecode`.

### Phase 2: Services & Models
1. [ ] Update `WebhookReports` in `report.go` to use a map for reports.
2. [ ] Correct all misleading comments in service files.

### Phase 3: Infrastructure & Tests
1. [ ] Add `golangci-lint` configuration and GitHub Action.
2. [ ] Implement missing tests for `VisitService` and `RecognitionError`.
3. [ ] Add test cases for `WaitForReport` context cancellation.

## 4. Testing Strategy

- **Unit Tests:** Update existing tests to support the new generic `Pagination` type.
- **Integration Tests:** Add `VisitService` tests using `httptest`.
- **Regression:** Ensure `ClintConf` alias still works and provides backward compatibility.

## 5. Backward Compatibility

- **Breaking Changes:** Changing `Pagination` to `Pagination[T]` and changing `WebhookReports` structure are breaking changes for users relying on the old field types.
- **Decision:** Given these improvements significantly improve the DX, we should target this for a minor version (v1.2.0) or major (v2.0.0) if strictly following SemVer.

## 6. Open Questions

1. Should we provide helper methods to convert the new `map[string]any` in `WebhookReports` to specific types, similar to `ToFacingCount`?
2. Do we want to support Go versions older than 1.24? (Currently, project is 1.24).