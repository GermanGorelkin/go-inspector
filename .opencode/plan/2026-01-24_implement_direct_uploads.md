# Task: Implement Direct File Uploads

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Direct file uploads remain unimplemented in `ImageService`, forcing every client to self-host assets before invoking the API.

## Proposed Solution

Introduce a multipart-powered `Upload(ctx context.Context, r io.Reader, filename string)` helper, backed by `http-client`, with regression tests and documentation updates.

## Detailed Steps

1. [ ] Step 1: Confirm or extend multipart support within `http-client`
   - Files: `inspector/image.go`, dependency `github.com/germangorelkin/http-client`
   - Changes: add helper utilities if the transport lacks multipart helpers.

2. [ ] Step 2: Implement `ImageService.Upload`
   - Files: `inspector/image.go`
   - Changes: construct multipart bodies, stream from `io.Reader`, wrap errors with context.

3. [ ] Step 3: Add regression tests
   - Files: `inspector/image_test.go`, `inspector/testdata/`
   - Changes: `httptest` server validating headers and multipart payloads.

4. [ ] Step 4: Document usage and caveats
   - Files: `README.md`, `specs/spec.md`
   - Changes: include examples, note size limits and expected content types.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestImageService_Upload`
- [ ] Integration tests: Upload sample files against staging
- [ ] Manual testing: CLI sample performing direct upload

## Open Questions

1. Does `github.com/germangorelkin/http-client` expose multipart helpers or do we extend it?
2. Are there maximum file sizes or timeouts that the client should enforce proactively?

## Risks and Edge Cases

- Large files could exhaust memory if not streamed efficiently.
- Missing multipart support might require broader transport changes.

## Rollback Strategy

Remove the new method, revert documentation updates, and drop any dependency changes.
