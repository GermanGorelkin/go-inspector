# Task: Fix VisitService.AddVisit Request Body

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

`VisitService.AddVisit` previously omitted the empty request body required by Inspector Cloud, causing payloads to be dropped silently.

## Proposed Solution

Always send an explicit empty JSON object in AddVisit requests, ensure regression coverage, and document the behavior for contributors.

## Detailed Steps

1. [x] Step 1: Update `VisitService.AddVisit` to send `{}`
   - Files: `inspector/visit.go`
   - Changes: force an empty buffer for POST calls that do not accept data.

2. [x] Step 2: Strengthen AddVisit unit tests
   - Files: `inspector/visit_test.go`
   - Changes: assert the body contents and verify the handler receives the payload.

3. [x] Step 3: Document AddVisit request expectations
   - Files: `README.md`
   - Changes: mention the empty body requirement for visit creation.

## Testing Strategy

- [x] Unit tests: `go test ./inspector -run TestVisitService_AddVisit`
- [ ] Integration tests: Validate against Inspector Cloud staging when available
- [ ] Manual testing: CLI sample flow for visit creation

## Open Questions

1. None.

## Risks and Edge Cases

- Minimal risk; future API updates could change payload expectations.

## Rollback Strategy

Revert the commits touching `inspector/visit.go`, `inspector/visit_test.go`, and related documentation.
