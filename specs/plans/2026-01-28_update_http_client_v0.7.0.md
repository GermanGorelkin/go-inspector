# Task: Update http-client dependency to v0.7.0

**Date:** 2026-01-28  
**Status:** Completed

## Problem Statement

The `github.com/germangorelkin/http-client` dependency is currently at v0.5.0 and needs to be updated to v0.7.0 to access new features including multipart file upload support and retry interceptor.

## Proposed Solution

Simple dependency version bump. Analysis confirms no breaking API changes between v0.5.0 and v0.7.0.

## Version Comparison

**Current:** v0.5.0  
**Target:** v0.7.0

### Changes between v0.5.0 and v0.7.0

1. **Multipart file upload support** (v0.7.0)
   - `NewMultipartForm()`, `AddFile()`, `PostMultipart()`
   - Enables future implementation of direct file upload in go-inspector

2. **Retry interceptor** (v0.6.0)
   - `NewRetryInterceptor()`, `DefaultRetryInterceptor()`
   - Automatic retry with exponential backoff for transient errors

3. **Documentation improvements**
   - README with usage examples
   - GoDoc comments on all exported entities

4. **Go 1.24 compatibility**
   - Updated go.mod to Go 1.24
   - Replaced `interface{}` with `any` (backward compatible)

### API Compatibility Check

All functions used by go-inspector are present in v0.7.0:
- `New()` ✓
- `WithBaseURL()` ✓
- `WithAuthorization()` ✓
- `WithInterceptor()` ✓
- `ResponseInterceptor` ✓
- `DumpInterceptor` ✓
- `AddInterceptor()` ✓

## Detailed Steps

1. [x] Step 1: Update go.mod
   - Files: `go.mod`
   - Changes: Change version from `v0.5.0` to `v0.7.0` on line 6

2. [x] Step 2: Run go mod tidy
   - Files: `go.sum`
   - Changes: Update checksums for new version

3. [x] Step 3: Run tests
   - Command: `go test ./... -v`
   - Verify all existing tests pass

4. [x] Step 4: Verify build
   - Command: `go build ./...`
   - Ensure library compiles without errors

## Testing Strategy

- [x] Unit tests: Run existing test suite - no new tests needed for dependency update
- [x] Manual testing: Verify `go build ./...` succeeds

## Open Questions

None - straightforward version bump with confirmed backward compatibility.

## Risks and Edge Cases

- **Risk:** None identified
- **Mitigation:** All API functions used are confirmed present in v0.7.0
- **Note:** The new multipart upload feature can be used to implement the TODO `Upload()` method in future work

## Rollback Strategy

Revert go.mod and go.sum to previous state:
```bash
git checkout go.mod go.sum
```

## Additional Notes

- `specs/spec.md` line 358 and `AGENTS.md` line 480 reference v0.4.0 - these are outdated documentation and should be updated separately if needed
