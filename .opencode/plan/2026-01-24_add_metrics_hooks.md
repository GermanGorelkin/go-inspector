# Task: Optional Metrics and Observability Hooks

**Date:** 2026-01-24  
**Status:** Planning

## Problem Statement

There is no way to observe SDK request lifecycle events, limiting insight into retries, latency, and failures.

## Proposed Solution

Add optional hook interfaces (e.g., `MetricsHook`) invoked around request execution, enabling Prometheus or custom telemetry integrations.

## Detailed Steps

1. [ ] Step 1: Design hook interfaces
   - Files: relevant service files, new hook definitions
   - Changes: define events (before request, after response, error) and payloads.

2. [ ] Step 2: Wire hooks into services
   - Files: service files
   - Changes: invoke hooks without affecting existing behavior; ensure nil-safe usage.

3. [ ] Step 3: Document integration patterns
   - Files: `README.md`, `specs/spec.md`
   - Changes: describe sample Prometheus integration and performance considerations.

## Testing Strategy

- [ ] Unit tests: ensure hooks fire under success/error paths
- [ ] Manual testing: sample hook implementation verifying event order
- [ ] No integration tests required initially

## Open Questions

1. Should hooks run synchronously (blocking) or asynchronously to avoid latency impact?

## Risks and Edge Cases

- Hook panics must not crash SDK consumers; need recovery strategy.

## Rollback Strategy

Remove hook interfaces and call sites.
