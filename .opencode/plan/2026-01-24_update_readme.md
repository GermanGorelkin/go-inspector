# Task: Refresh README documentation

**Date:** 2026-01-24  
**Status:** Completed

## Problem Statement

The current `README.md` still reflects an early beta version of the SDK. It uses outdated examples (no context usage, deprecated Upload method), lacks explanations of the service-oriented structure, omits requirements/testing instructions, and does not mention the CLI example or async recognition workflow. New contributors struggle to understand capabilities and setup steps.

## Proposed Solution

Rewrite the README around the canonical information from `.opencode/spec.md` and `AGENTS.md`. Provide a concise overview, installation guide, end-to-end quickstart with context-aware snippets, service/architecture summary, CLI usage, and development/testing workflow. Highlight gotchas (asynchronous reports, missing direct upload) and link to deeper resources.

## Detailed Steps

1. [x] Step 1: Audit existing docs for canonical facts
   - Files: `.opencode/spec.md`, `AGENTS.md`, `cmd/cli/main.go`
   - Changes: none (information gathering)

2. [x] Step 2: Restructure README header and overview
   - Files: `README.md`
   - Changes: update title, project status/version, add badges/description

3. [x] Step 3: Add requirements, install, and configuration sections
   - Files: `README.md`
   - Changes: mention Go 1.24, module path, environment variables, CLI build instructions

4. [x] Step 4: Rewrite quickstart and workflow examples
   - Files: `README.md`
   - Changes: provide code snippets for client init with context, upload by URL, recognition flow, report parsing, SKU pagination, visit creation, webhook parsing

5. [x] Step 5: Document architecture, services, report types, and CLI usage
   - Files: `README.md`
   - Changes: add service overview table/list, report type descriptions, CLI commands and expected output

6. [x] Step 6: Add development/testing guidance and resource links
   - Files: `README.md`
   - Changes: include build/test commands, references to `.opencode/spec.md`, AGENTS guide, license info, issue reporting instructions

## Testing Strategy

- [x] Proofread updated README  
- [ ] Optional: run `markdownlint` manually if available (not required)  
- [x] Verify code snippets compile conceptually (Go vet mentally)

## Open Questions

1. Should we include screenshots or diagrams? (default: no)  
2. Should CLI usage highlight sample output? (default: brief textual description)

## Risks and Edge Cases

- Risk: README diverges from code if future changes occur. Mitigation: link to spec and emphasize canonical sources.  
- Risk: Overly long README. Mitigation: keep sections concise, link out for deep dives.  
- Edge case: Direct file upload still unimplemented—call out explicitly to avoid confusion.

## Rollback Strategy

Revert changes to `README.md` and delete this plan file if necessary.
