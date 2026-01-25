# Task: Introduce SKU Pagination Iterator

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Clients manually issue paginated SKU requests, risking infinite loops and inconsistent page handling.

## Proposed Solution

Create an iterator or helper (`IterateSKU`/`GetAllSKU`) that abstracts pagination, guards against loops, and offers ergonomic accessors.

## Detailed Steps

1. [ ] Step 1: Design iterator API and safeguards
   - Files: `inspector/sku.go`
   - Changes: define struct/function shape, maximum pages, optional callbacks.

2. [ ] Step 2: Implement pagination logic
   - Files: `inspector/sku.go`
   - Changes: fetch pages, update cursors, detect repeated tokens.

3. [ ] Step 3: Write multi-page tests
   - Files: `inspector/sku_test.go`
   - Changes: simulate paginated responses, cover empty and error cases.

4. [ ] Step 4: Document iterator usage
   - Files: `README.md`, `specs/spec.md`
   - Changes: example loop and warnings about rate limits.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestSkuService_Iterate`
- [ ] Integration tests: Optional real pagination validation
- [ ] Manual testing: CLI script fetching all SKUs

## Open Questions

1. Which pagination fields are canonical (page/limit vs cursor)?
2. Do we expose a channel-based iterator or synchronous slice helper?

## Risks and Edge Cases

- API changes to pagination schema could break assumptions.
- Infinite loop safeguards must not accidentally truncate valid datasets.

## Rollback Strategy

Remove iterator types/functions and revert documentation updates.
