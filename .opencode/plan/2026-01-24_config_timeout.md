# Task: Add Timeout Config and Rename ClientConf

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

The client configuration struct (`ClintConf`) lacks an explicit timeout option and still carries a long-standing typo in its name. Users need the ability to customize HTTP timeouts without supplying a full custom `http.Client`, and the exported configuration type should have the correct name while preserving backward compatibility.

## Proposed Solution

Introduce a `Timeout time.Duration` field in the configuration, defaulting to 30 seconds when unspecified. Rename the struct to `ClientConf` and provide a type alias `type ClintConf = ClientConf` to avoid breaking existing consumers. Update `NewClient` to honor the new timeout when constructing the default HTTP client, refresh documentation/spec references, and add tests if practical.

## Detailed Steps

1. [x] Step 1: Update configuration struct and client initialization
   - Files: `inspector/client.go`
   - Changes: Rename struct to `ClientConf`, add `Timeout time.Duration`, ensure default 30s timeout, and keep backward-compatible alias.

2. [x] Step 2: Update documentation/spec references
   - Files: `.opencode/spec.md`, `README.md` (if needed)
   - Changes: Reflect new struct name, timeout option, and default description.

3. [ ] Step 3: Add/adjust tests (if applicable)
   - Files: `inspector/client_test.go` (new or existing)
   - Changes: Introduce tests covering timeout handling and struct rename usage.

## Testing Strategy

- [ ] Unit tests: Add coverage for timeout configuration in `inspector/client_test.go`.
- [ ] Integration tests: Not required (library change only).
- [ ] Manual testing: Not required.

## Open Questions

1. Should we keep both `ClientConf` and `ClintConf` exported via alias, or is a hard rename acceptable?
2. Do we need to expose the timeout default via a constant for users to reference?

## Risks and Edge Cases

- Risk 1: Breaking API change if alias not provided. Mitigation: Add `type ClintConf = ClientConf` and keep existing docs noting alias.
- Edge case 1: Users supply custom `http.Client` and timeout; ensure new logic does not override provided client.

## Rollback Strategy

Revert changes to `inspector/client.go` and documentation, removing the timeout field and alias, to restore previous behavior.
