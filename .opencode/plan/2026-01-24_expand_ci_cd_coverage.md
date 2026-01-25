# Task: Expand CI/CD Coverage

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Current CI only runs tests, lacking coverage uploads, linting, tidy checks, and security scanning.

## Proposed Solution

Enhance `.github/workflows/test.yml` to include coverage publishing, lint job, `go mod tidy` verification, and `govulncheck`.

## Detailed Steps

1. [ ] Step 1: Add coverage collection and upload
   - Files: `.github/workflows/test.yml`
   - Changes: run `go test ./... -coverprofile` and upload artifact or send to Codecov.

2. [ ] Step 2: Integrate lint job (depends on golangci-lint plan)
   - Files: `.github/workflows/test.yml`
   - Changes: add separate job or matrix entry.

3. [ ] Step 3: Add `go mod tidy` consistency check
   - Files: `.github/workflows/test.yml`
   - Changes: run tidy and fail if diff detected.

4. [ ] Step 4: Run `govulncheck`
   - Files: `.github/workflows/test.yml`
   - Changes: add security scanning step with caching.

## Testing Strategy

- [ ] CI: push branch to verify workflow runs on Ubuntu, macOS, Windows
- [ ] Manual testing: optional dry-run via `act`
- [ ] No unit tests required

## Open Questions

1. Which coverage service (Artifacts vs Codecov) should we adopt?
2. Should lint job block merges or be informational initially?

## Risks and Edge Cases

- Longer CI duration could slow feedback loops.
- `govulncheck` might flag false positives requiring triage.

## Rollback Strategy

Revert workflow changes or disable new jobs temporarily.
