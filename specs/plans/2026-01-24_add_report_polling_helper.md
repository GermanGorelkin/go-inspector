# Task: Add Report Polling Helper

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

Consumers must hand-roll polling loops to wait for asynchronous reports, leading to duplicated, error-prone code without shared defaults.

## Proposed Solution

Provide a `WaitForReport` helper with configurable `WaitOptions` (interval, timeout, backoff, progress callback) and thorough tests plus documentation.

## Detailed Steps

1. [x] Step 1: Define `WaitOptions` and defaults
   - Files: `inspector/report.go`
   - Changes: capture polling interval, timeout, optional backoff, and progress callback.

2. [x] Step 2: Implement `WaitForReport`
   - Files: `inspector/report.go`
   - Changes: loop with context awareness, evaluate status, expose progress callbacks.

3. [x] Step 3: Add table-driven tests
   - Files: `inspector/report_test.go`
   - Changes: mock report responses, cover timeout, cancellation, and success cases.

4. [x] Step 4: Document usage
   - Files: `README.md`, `specs/spec.md`
   - Changes: provide example code and guidance on defaults.

## Testing Strategy

- [x] Unit tests: `go test ./inspector -run TestReportService_WaitForReport` (5 scenarios: ready, error, timeout, canceled, backoff)
- [ ] Integration tests: Optional staging validation with slow reports
- [ ] Manual testing: CLI workflow demonstrating polling helper

## Open Questions

1. Should the helper expose logging hooks or rely solely on `WaitOptions` callbacks? Answered: callbacks only.
2. What default timeout best balances UX with rate limits (e.g., 60s vs 120s)? Answered: 60s.

## Risks and Edge Cases

- Long-running polls could exceed rate limits; need jitter/backoff.
- Callers might expect continuous polling even after context deadlines.

## Rollback Strategy

Delete `WaitForReport`, remove new constants, and revert documentation.
