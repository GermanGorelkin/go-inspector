# Task: Fix SKU Iterator Pagination

**Date:** 2026-01-28  
**Status:** In Progress

## Problem Statement

`SKUIterator.Next()` increments by the requested page size and stops when the returned page is smaller than `pageSize`. If the API caps `limit` or returns short pages while still providing `next`, the iterator can stop early or skip records.

## Proposed Solution

Use the server-provided `Pagination.Next` URL as the continuation signal and advance the iterator offset based on the `next` URL when available, falling back to the actual returned count as a safety.

## Detailed Steps

1. [x] Step 1: Update iterator pagination logic
   - Files: `inspector/sku.go`
   - Changes: Use `pag.Next` to determine continuation; parse next offset when present; fallback to advancing by returned length when needed.

2. [x] Step 2: Add/adjust tests for short pages with next
   - Files: `inspector/sku_test.go`
   - Changes: Add a test where the server returns a short page with a valid `next` URL to ensure iteration continues without skipping.

3. [ ] Step 3: Verify tests
   - Files: none
   - Changes: Run focused tests for SKU iterator pagination.

## Testing Strategy

- [ ] Unit tests: `go test ./inspector -run TestSkuService_IterateSKU -v`

## Open Questions

1. Should `pag.Next` be treated as the canonical continuation signal even when the page is short?

## Risks and Edge Cases

- If `pag.Next` is present but does not include a parseable `offset`, fallback behavior must not loop indefinitely.
- Empty pages with `next` could indicate server-side filtering or eventual consistency.

## Rollback Strategy

Revert the iterator logic changes and remove the added tests.
