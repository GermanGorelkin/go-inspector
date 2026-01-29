# Task: Add Makefile Workflow

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

Contributors currently run commands manually, lacking a standardized Makefile to encapsulate fmt, lint, build, and test flows.

## Proposed Solution

Create a Makefile with common targets, document usage, and ignore generated artifacts like coverage profiles.

## Detailed Steps

1. [x] Step 1: Author Makefile targets (`fmt`, `lint`, `test`, `build`, `coverage`)
   - Files: `Makefile`
   - Changes: ensure cross-platform compatibility and configurable variables.

2. [x] Step 2: Ignore generated artifacts
   - Files: `.gitignore`
   - Changes: add coverage outputs or binaries built via Makefile.

3. [x] Step 3: Document workflow
   - Files: `README.md`
   - Changes: add contributor instructions referencing Make targets.

## Testing Strategy

- [x] Manual testing: run each Make target locally
- [ ] CI: optionally invoke Make targets in workflows
- [ ] No automated tests required

## Open Questions

1. Should Make targets wrap `golangci-lint` and `govulncheck` or remain minimal?

## Risks and Edge Cases

- GNU make extensions may not work on BSD make; keep commands portable.

## Rollback Strategy

Remove the Makefile and related documentation updates.
