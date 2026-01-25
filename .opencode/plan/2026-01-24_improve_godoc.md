# Task: Improve Package Godoc

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Godoc coverage is sparse, leaving exported types and services undocumented, which hurts discoverability.

## Proposed Solution

Add a package-level comment plus concise documentation for all exported types and methods, including inline examples where helpful.

## Detailed Steps

1. [ ] Step 1: Add package comment and high-level overview
   - Files: `inspector/doc.go` or existing root file
   - Changes: summarize SDK purpose and services.

2. [ ] Step 2: Document exported types and methods
   - Files: all exported `.go` files under `inspector/`
   - Changes: add Godoc comments aligning with style guide.

3. [ ] Step 3: Include code examples
   - Files: relevant service files or `_test.go`
   - Changes: add `Example` functions demonstrating workflows.

4. [ ] Step 4: Verify via `go doc`
   - Files: n/a
   - Changes: ensure formatting and references render correctly.

## Testing Strategy

- [ ] Documentation build: `go doc inspector`
- [ ] Manual testing: review rendered godoc locally or on pkg.go.dev preview
- [ ] No automated tests required

## Open Questions

1. Should we highlight advanced helpers (polling, iterators) in examples or keep them separate?

## Risks and Edge Cases

- Comments must stay synchronized with behavior to avoid misleading users.

## Rollback Strategy

Remove or adjust newly added comments if necessary.
