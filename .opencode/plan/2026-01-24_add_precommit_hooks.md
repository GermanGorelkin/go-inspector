# Task: Optional Pre-commit Hooks

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

Contributors lack automated local checks before committing, leading to avoidable CI failures.

## Proposed Solution

Provide an optional `.pre-commit-config.yaml` with Go-focused hooks and document installation steps in `CONTRIBUTING.md`.

## Detailed Steps

1. [ ] Step 1: Define pre-commit configuration
   - Files: `.pre-commit-config.yaml`
   - Changes: include gofmt, golangci-lint, go test, and shell fallback hooks.

2. [ ] Step 2: Document setup instructions
   - Files: `CONTRIBUTING.md`, `README.md`
   - Changes: explain installation and optional nature.

3. [ ] Step 3: Provide helper script (optional)
   - Files: `hack/` scripts if needed
   - Changes: wrap gofmt/lint/test for contributors without pre-commit.

## Testing Strategy

- [ ] Manual testing: run `pre-commit run --all-files`
- [ ] No automated tests required
- [ ] Optional CI job to ensure config stays valid

## Open Questions

1. Should hooks run automatically via Makefile for users without pre-commit installed?

## Risks and Edge Cases

- Overly heavy hooks could slow contributors; keep optional.

## Rollback Strategy

Remove `.pre-commit-config.yaml` and revert documentation references.
