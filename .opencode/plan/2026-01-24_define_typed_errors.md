# Task: Define Typed Error System

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Errors across services are untyped, forcing string matching and hiding actionable context from SDK consumers.

## Proposed Solution

Introduce an `ErrorType` enum plus typed wrappers that capture HTTP status, API codes, and retry guidance, updating services and docs accordingly.

## Detailed Steps

1. [ ] Step 1: Add `inspector/errors.go` with core types
   - Files: `inspector/errors.go`
   - Changes: define `ErrorType`, `APIError`, helpers like `IsRateLimited`.

2. [ ] Step 2: Update services to wrap HTTP failures
   - Files: all service files
   - Changes: convert error paths to typed errors with `%w` context.

3. [ ] Step 3: Add regression tests
   - Files: `inspector/errors_test.go`, service tests
   - Changes: cover mapping logic and helper behavior.

4. [ ] Step 4: Document error handling guidance
   - Files: `specs/spec.md`, `README.md`
   - Changes: describe `ErrorType` enum and sample usage.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestErrorType`
- [ ] Integration tests: Optional staging checks for real error payloads
- [ ] Manual testing: CLI scenario demonstrating typed errors

## Open Questions

1. Does introducing typed errors require a v2.0.0 release per semver?
2. Should we expose raw response bodies on the error struct for debugging?

## Risks and Edge Cases

- Breaking behavior if callers relied on previous error strings.
- Mapping may drift if API adds new error codes.

## Rollback Strategy

Remove typed error definitions and revert service changes.
