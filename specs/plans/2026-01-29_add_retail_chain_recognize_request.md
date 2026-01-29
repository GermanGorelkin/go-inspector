# Task: Add retail_chain to recognition request

**Date:** 2026-01-29  
**Status:** Completed

## Problem Statement

The recognition request model lacks the `retail_chain` field required by the API. We need to add it to the request struct so clients can pass the retail chain identifier.

## Proposed Solution

Extend `RecognizeRequest` with a `RetailChain` string field, wired to JSON as `retail_chain`, and update documentation/tests as needed.

## Detailed Steps

1. [x] Step 1: Add field to request model
   - Files: `inspector/recognize.go`
   - Changes: Add `RetailChain string` with `json:"retail_chain,omitempty"` and inline comment.

2. [x] Step 2: Update spec and docs (if applicable)
   - Files: `specs/spec.md`, `README.md` (if they document the request)
   - Changes: Add `retail_chain` to the recognition request documentation.

3. [x] Step 3: Update tests (if any coverage exists for recognize request)
   - Files: `inspector/recognize_test.go`
   - Changes: Ensure JSON serialization includes `retail_chain` when set.

## Testing Strategy

- [x] Unit tests: Add/extend tests for recognize request JSON encoding (`go test ./... -v`).
- [ ] Manual testing: Validate sample request payloads in docs.

## Open Questions

1. Should `retail_chain` be optional (`omitempty`) or always included? — Optional (`omitempty`).
2. Do we want to update any public examples in `README.md` to include `retail_chain`? — Yes, updated.

## Risks and Edge Cases

- Risk: Breaking changes if existing tests assert exact JSON without the new field.
- Edge case: Empty string value should be omitted if `omitempty` is used.

## Rollback Strategy

Revert the added field and any documentation/test updates.