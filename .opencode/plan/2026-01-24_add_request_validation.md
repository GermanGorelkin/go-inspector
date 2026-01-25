# Task: Add Request Validate Methods

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Requests lack client-side validation, allowing malformed payloads to reach the API and produce opaque failures.

## Proposed Solution

Introduce `Validate()` methods on request structs (recognize, visit, image) and enforce validation before HTTP execution, backed by negative tests.

## Detailed Steps

1. [ ] Step 1: Define validation rules per request
   - Files: `inspector/recognize.go`, `inspector/visit.go`, `inspector/image.go`
   - Changes: add `Validate()` implementations with descriptive errors.

2. [ ] Step 2: Invoke validation in service methods
   - Files: respective service files
   - Changes: call `Validate()` before building requests.

3. [ ] Step 3: Add positive/negative tests
   - Files: corresponding `*_test.go`
   - Changes: ensure invalid inputs fail fast with typed errors.

4. [ ] Step 4: Document validation behavior
   - Files: `README.md`, `specs/spec.md`
   - Changes: describe error types and example usage.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run Test*Validate`
- [ ] Integration tests: Optional staging confirmation for valid payloads
- [ ] Manual testing: CLI flow rejecting invalid input

## Open Questions

1. Should validation errors reuse the new typed error framework or plain `error`?
2. Do we allow partial success scenarios (e.g., some images missing URLs)?

## Risks and Edge Cases

- Overly strict validation could reject payloads the server accepts.
- Keeping validation rules in sync with API evolution requires diligence.

## Rollback Strategy

Remove `Validate()` methods and associated call sites/tests.
