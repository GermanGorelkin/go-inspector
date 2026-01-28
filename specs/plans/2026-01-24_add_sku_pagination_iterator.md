# Task: Introduce SKU Pagination Iterator

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

Clients manually issue paginated SKU requests, risking infinite loops and inconsistent page handling.

## Proposed Solution

Create an iterator or helper (`IterateSKU`/`GetAllSKU`) that abstracts pagination, guards against loops, and offers ergonomic accessors.

## Detailed Steps

1. [x] Step 1: Design iterator API and safeguards
   - Files: `inspector/sku.go`
   - Changes: Added `SKUIterator` struct with `IterateSKU()` and `GetAllSKU()` methods, infinite loop safeguards (max 1000 pages, duplicate offset detection).

2. [x] Step 2: Implement pagination logic
   - Files: `inspector/sku.go`
   - Changes: Implemented `Next()` method with automatic offset increment, last page detection (`len(skus) < pageSize || pag.Next == nil`), error wrapping.

3. [x] Step 3: Write multi-page tests
   - Files: `inspector/sku_test.go`
   - Changes: Added `TestSkuService_IterateSKU` and `TestSkuService_GetAllSKU` with multi-page simulation, error cases, and empty result handling.

4. [x] Step 4: Document iterator usage
   - Files: `README.md`, `specs/spec.md`
   - Changes: Updated README with iterator examples, updated service table, updated spec.md to reflect pagination helpers are now implemented.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestSkuService_Iterate`
- [ ] Integration tests: Optional real pagination validation
- [ ] Manual testing: CLI script fetching all SKUs

## Open Questions

1. Which pagination fields are canonical (page/limit vs cursor)?
   - **Answer:** The API uses offset/limit pagination as seen in `GetSKU()`. The `Pagination` struct has `Next`/`Previous` URL fields, but the iterator uses offset/limit internally.

2. Do we expose a channel-based iterator or synchronous slice helper?
   - **Answer:** Implemented both: `IterateSKU()` returns an iterator for manual control, `GetAllSKU()` is a synchronous convenience method that returns all SKUs.

## Risks and Edge Cases

- API changes to pagination schema could break assumptions.
- Infinite loop safeguards must not accidentally truncate valid datasets.

## Rollback Strategy

Remove iterator types/functions and revert documentation updates.
