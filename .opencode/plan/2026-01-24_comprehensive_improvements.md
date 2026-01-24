# Task: Comprehensive Project Improvements

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

The go-inspector SDK has accumulated critical bugs, missing functionality, and weak infrastructure that collectively block reliable adoption. Visit creation silently drops payloads, direct file uploads are stubbed, complex workflows (report polling, pagination) force every client to reimplement helpers, and errors remain untyped. Documentation, constants, and validation are inconsistent, while linting, Makefile tooling, and CI coverage lag behind modern contributor expectations.

## Proposed Solution

Deliver a staged improvement program covering four priority areas: unblock critical bugs, ship high-impact helpers, elevate code quality, and modernize tooling/CI. Each stage introduces targeted fixes with accompanying tests, documentation, and guardrails so future contributors can iterate safely. Optional enhancements (mock client, rate limiting, observability hooks) remain on deck once core stability returns.

## Detailed Steps

1. [x] Step 1: Fix `VisitService.AddVisit` request body
   - Files: `inspector/visit.go`, `inspector/visit_test.go`
   - Changes: always send empty body required by API, adjust tests and docs.

2. [ ] Step 2: Implement direct file uploads
   - Files: `inspector/image.go`, `inspector/image_test.go`, `README.md`
   - Changes: add `Upload(ctx context.Context, r io.Reader, filename string)`, extend `http-client` if multipart support is missing, add tests and docs.

3. [ ] Step 3: Add report polling helper
   - Files: `inspector/report.go`, `inspector/report_test.go`, `README.md`
   - Changes: implement `WaitForReport` with `WaitOptions`, optional backoff/logging, plus mock-based tests and usage example.

4. [ ] Step 4: Introduce SKU pagination iterator
   - Files: `inspector/sku.go`, `inspector/sku_test.go`, `README.md`
   - Changes: create `GetAllSKU` or `IterateSKU`, guard against infinite loops, test multi-page flows, document behavior.

5. [ ] Step 5: Define typed error system
   - Files: `inspector/errors.go`, `inspector/errors_test.go`, service files, `.opencode/spec.md`
   - Changes: add `ErrorType` enum, wrap HTTP failures, expose helpers, document error handling.

6. [ ] Step 6: Add request `Validate()` methods
   - Files: `inspector/recognize.go`, `inspector/visit.go`, `inspector/image.go`, related tests
   - Changes: enforce required fields before API calls, add negative tests, update existing coverage.

7. [ ] Step 7: Improve package godoc
   - Files: all exported `.go` files in `inspector/`
   - Changes: add package comment, document public types/methods, include inline examples, verify via `go doc`.

8. [ ] Step 8: Centralize shared constants
   - Files: `inspector/constants.go`, service files, constant tests
   - Changes: declare default polling values and endpoint strings, replace literals across services, ensure tests cover defaults.

9. [ ] Step 9: Configure golangci-lint
   - Files: `.golangci.yml`, `.github/workflows/test.yml`, `README.md`
   - Changes: add lint config, wire CI job, resolve reported issues, document badge.

10. [ ] Step 10: Add Makefile workflow
    - Files: `Makefile`, `README.md`, `.gitignore`
    - Changes: define common targets (fmt, lint, test, build, coverage), ignore coverage artifacts, document usage.

11. [ ] Step 11: Expand CI/CD coverage
    - Files: `.github/workflows/test.yml`
    - Changes: add coverage upload, lint job, `go mod tidy` consistency check, `govulncheck` security scan.

12. [ ] Step 12: Optional pre-commit hooks
    - Files: `.pre-commit-config.yaml`, `CONTRIBUTING.md`
    - Changes: configure Go hooks, describe installation steps, consider simple shell fallback.

13. [ ] Step 15: Optional metrics/observability hooks
    - Files: service files introducing hook interfaces
    - Changes: add `MetricsHook`, invoke callbacks on request lifecycle, document Prometheus integration ideas.

## Testing Strategy

- [ ] Unit tests: `go test ./...`
- [ ] Integration tests: manual verification against Inspector Cloud staging (if available)
- [ ] Manual testing: CLI sample flows for uploads, polling, pagination

## Open Questions

1. Should typed errors or other breaking API changes trigger a v2.0.0 release?
2. Is a staging Inspector Cloud environment available for validating multipart uploads?
3. What level of backward compatibility is mandatory for existing clients?

## Risks and Edge Cases

- Breaking API changes may require semantic versioning coordination; plan for v2 if compatibility cannot be preserved.
- `http-client` might lack multipart support; be ready to extend it or introduce an alternative transport.
- Large refactors can destabilize existing tests; update suites in lockstep to avoid regressions.
- Added complexity (polling, iterators, hooks) must remain opt-in to keep the public API approachable.

## Rollback Strategy

Each phase resides in its own branch or commit series. Revert the relevant commits or drop the feature branch to roll back, ensuring CI remains green before merging future work.

## Additional Resources

- `.opencode/spec.md`, `AGENTS.md`, `README.md`
- Inspector Cloud API documentation (link pending)
- https://go.dev/doc/effective_go
- https://semver.org/
