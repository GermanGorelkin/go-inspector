# Task: Centralize Shared Constants

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Polling defaults, endpoint paths, and other repeated literals are scattered across services, making updates error-prone.

## Proposed Solution

Create `inspector/constants.go` to house shared defaults (endpoints, timeouts, report statuses) and replace literals across the codebase.

## Detailed Steps

1. [ ] Step 1: Catalog duplicated literals
   - Files: service files under `inspector/`
   - Changes: list endpoints, default intervals, header names.

2. [ ] Step 2: Add `constants.go` with exported values
   - Files: `inspector/constants.go`
   - Changes: document each constant and provide sensible defaults.

3. [ ] Step 3: Refactor services to use constants
   - Files: affected service files
   - Changes: replace inline strings/numbers, update tests if needed.

4. [ ] Step 4: Update tests referencing literals
   - Files: relevant `*_test.go`
   - Changes: import constants to keep assertions in sync.

## Testing Strategy

- [ ] Unit tests: `go test ./...`
- [ ] Manual testing: sanity-check CLI flows touching endpoints
- [ ] No dedicated integration tests required

## Open Questions

1. Which constants remain internal vs exported for SDK consumers?

## Risks and Edge Cases

- Over-exposing constants could freeze internal APIs.

## Rollback Strategy

Revert `constants.go` and restore previous literals.
