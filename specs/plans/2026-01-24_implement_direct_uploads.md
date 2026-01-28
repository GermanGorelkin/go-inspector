# Task: Implement Direct File Uploads

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

Direct file uploads remain unimplemented in `ImageService`, forcing every client to self-host assets before invoking the API.

## Proposed Solution

Introduce a multipart-powered `Upload(ctx context.Context, r io.Reader, filename string)` helper, backed by `http-client`, with regression tests and documentation updates.

## Detailed Steps

1. [x] Step 1: Confirm or extend multipart support within `http-client`
   - Files: `inspector/image.go`, dependency `github.com/germangorelkin/http-client`
   - Changes: confirmed `MultipartForm` and `NewMultipartRequest` are available in v0.7.0 (buffers multipart body in memory).

2. [x] Step 2: Implement `ImageService.Upload`
   - Files: `inspector/image.go`
   - Changes: construct multipart bodies with field `file`, wrap errors with context, keep ctx-aware request handling.

3. [x] Step 3: Add regression tests
   - Files: `inspector/image_test.go`
   - Changes: `httptest` server validates multipart payload and auth header.

4. [x] Step 4: Document usage and caveats
   - Files: `README.md`, `specs/spec.md`
   - Changes: add direct upload example, note multipart buffering behavior.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestImageService_Upload -v`
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
