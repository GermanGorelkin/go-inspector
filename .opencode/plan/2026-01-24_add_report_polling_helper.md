# Task: Add Report Polling Helper

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Consumers must hand-roll polling loops to wait for asynchronous reports, leading to duplicated, error-prone code without shared defaults.

## Proposed Solution

Provide a `WaitForReport` helper with configurable `WaitOptions` (interval, timeout, backoff, logger) and thorough tests plus documentation.

## Detailed Steps

1. [ ] Step 1: Define `WaitOptions` and defaults
   - Files: `inspector/report.go`, `inspector/constants.go`
   - Changes: capture polling interval, timeout, optional backoff, and logging hook.

2. [ ] Step 2: Implement `WaitForReport`
   - Files: `inspector/report.go`
   - Changes: loop with context awareness, evaluate status, expose progress callbacks.

3. [ ] Step 3: Add table-driven tests
   - Files: `inspector/report_test.go`
   - Changes: mock report responses, cover timeout, cancellation, and success cases.

4. [ ] Step 4: Document usage
   - Files: `README.md`, `specs/spec.md`
   - Changes: provide example code and guidance on defaults.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestReportService_WaitForReport`
- [ ] Integration tests: Optional staging validation with slow reports
- [ ] Manual testing: CLI workflow demonstrating polling helper

## Open Questions

1. Should the helper expose logging hooks or rely solely on `WaitOptions` callbacks?
2. What default timeout best balances UX with rate limits (e.g., 60s vs 120s)?

## Risks and Edge Cases

- Long-running polls could exceed rate limits; need jitter/backoff.
- Callers might expect continuous polling even after context deadlines.

## Rollback Strategy

Delete `WaitForReport`, remove new constants, and revert documentation.
