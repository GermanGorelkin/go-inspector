# Task: Configure golangci-lint

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

The project lacks automated linting, allowing stylistic and correctness issues to slip into CI.

## Proposed Solution

Introduce a `golangci-lint` configuration, wire it into GitHub Actions, and resolve reported issues.

## Detailed Steps

1. [ ] Step 1: Create `.golangci.yml`
   - Files: `.golangci.yml`
   - Changes: enable relevant linters aligned with project style.

2. [ ] Step 2: Update CI workflow
   - Files: `.github/workflows/test.yml`
   - Changes: add lint job or extend existing pipeline.

3. [ ] Step 3: Fix lint findings
   - Files: codebase-wide
   - Changes: address formatting, style, and potential bugs.

4. [ ] Step 4: Document lint workflow
   - Files: `README.md`
   - Changes: describe how to run lint locally.

## Testing Strategy

- [ ] Lint: `golangci-lint run`
- [ ] CI: ensure workflow passes on supported OS targets
- [ ] Manual testing: none required beyond verifying instructions

## Open Questions

1. Which linters should be mandatory vs optional (e.g., `gosimple`, `staticcheck`)?

## Risks and Edge Cases

- Aggressive linters could block contributions; need balanced configuration.

## Rollback Strategy

Remove `.golangci.yml` and CI lint steps if necessary.
