# Task: Expand Unit Tests

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

Current unit coverage skips several public API surfaces (`VisitService.AddVisit`, `RecognizeService.RecognitionError`, client initialization) and lacks guard-rail tests for error wrapping on converters. This leaves regressions undetected when request bodies change, when verbose configuration alters HTTP client wiring, or when JSON parsing fails.

## Proposed Solution

Add focused unit tests mirroring existing httptest patterns. Cover missing services, client construction branches, and negative parsing paths to ensure the SDK’s helper methods fail predictably and serialize requests correctly.

## Detailed Steps

1. [x] Step 1: Cover recognition error helper
   - Files: `inspector/recognize_test.go`
   - Changes: Add `TestRecognizeService_RecognitionError` verifying POST payload and response decoding.

2. [x] Step 2: Add Visit service test
   - Files: `inspector/visit_test.go`
   - Changes: New test ensuring `AddVisit` sends POST body and handles response.

3. [ ] Step 3: Test client initialization logic
   - Files: `inspector/client_test.go`
   - Changes: Verify default HTTP client timeout, injected client usage, verbose interceptor, and service wiring.

4. [x] Step 4: Strengthen conversion error handling tests
   - Files: `inspector/report_test.go`, `inspector/sku_test.go`
   - Changes: Add malformed JSON / decode failure cases asserting wrapped errors.

5. [x] Step 5: Run full unit suite
   - Files: n/a
   - Changes: Execute `go test ./...` and ensure green.

## Testing Strategy

- [x] Unit tests: `go test ./...`
- [ ] Integration tests: n/a
- [ ] Manual testing: n/a

## Open Questions

1. None (validated with user).
2. —

## Risks and Edge Cases

- Mock servers must match API paths to keep assertions reliable.
- Client tests should avoid mutating global HTTP state when swapping transports.
- Error-path tests must assert wrapped messages without leaking brittle formatting.

## Rollback Strategy

Revert the added tests/plan file via a single commit or branch reset; no external state changes.
