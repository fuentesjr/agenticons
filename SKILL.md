---
name: codex-dispatch
description: Selects and spawns named Codex custom subagents for planning, implementation, review, documentation review, and helper work. Use for tasks that need delegation, parallel execution, code changes, review, investigation, documentation drift checks, or model-tier routing.
---

# Codex Dispatch

Use this skill as a lean dispatcher. Your job is to select and spawn the right named custom subagent for the task, then consolidate results.

## Hard rules

- Default to dispatching through subagents, unless the user explicitly opts out.
- Honor escape hatches exactly. If the user says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use codex-dispatch`, do not spawn subagents for that request.
- Dispatch to named subagents. Do not merely recommend models.
- Use the exact subagent names from `.codex/agents/*.toml`:
  - `planner`
  - `coding_worker`
  - `fast_coding_worker`
  - `helper_worker`
  - `doc_reviewer`
  - `reviewer`
  - `premium_reviewer`
- Standard review always uses `reviewer`.
- Documentation drift review uses `doc_reviewer`.
- Use `premium_reviewer` only for rare, very high-stakes/high-cost review.
- Keep orchestration shallow: parent coordinates; subagents do focused work.
- Avoid recursive delegation unless the user explicitly asks for it.
- Prefer one subagent for simple work, two to four subagents for parallel review or independent implementation slices.
- Every spawned subagent must receive a concrete task, scope, expected output, and constraints.
- Assign disjoint file or feature ownership to implementation agents when more than one worker may edit files.
- Ask subagents for evidence: file paths, symbols, commands run, test results, and unresolved risks.
- Do not let implementation agents silently expand scope.

## Routing table

| Situation | Spawn |
|---|---|
| Complex feature, vague request, architectural tradeoff, phased plan | `planner` |
| Normal code implementation, multi-file edit, bug fix, refactor | `coding_worker` |
| Small localized edit, mechanical change, simple failing test | `fast_coding_worker` |
| Read-only lookup, repo reconnaissance, dependency/API check, test triage, docs/help task | `helper_worker` |
| Documentation correctness, stale docs, README/API drift, changelog/release-note coverage | `doc_reviewer` |
| Normal correctness/security/maintainability review | `reviewer` |
| Rare high-stakes review: security boundary, data loss, payments, auth, migration, production incident, irreversible or expensive decision | `premium_reviewer` |

## Dispatch patterns

### Plan then implement

1. Spawn `planner` to produce a short implementation plan, risks, files likely touched, and verification strategy.
2. Spawn `coding_worker` with the accepted plan.
3. Spawn `reviewer` to review the diff and test evidence.

### Fast fix

1. Spawn `fast_coding_worker` for a tightly scoped fix.
2. Spawn `reviewer` only if the change affects behavior, tests, security, data, or public API.

### Parallel review

Spawn one or more `reviewer` agents with separate review angles only when useful, for example:
- correctness and edge cases
- security and permissions
- tests and regressions

Use `premium_reviewer` instead of `reviewer` only when the change is genuinely high-stakes.

### Documentation drift review

Spawn `doc_reviewer` when code behavior, public APIs, installation steps, commands, configuration, examples, or user-facing workflows may have changed. Ask it to compare implementation and docs, identify stale or missing documentation, and recommend exact doc updates.

### Investigation before editing

Spawn `helper_worker` first when the correct files, runtime path, API behavior, or failure mode are unclear. Then hand the evidence to `coding_worker` or `fast_coding_worker`.

## Subagent prompt template

When spawning a subagent, include:

```text
Agent: <name>

Task:
<single concrete task>

Scope:
<files, feature area, branch, issue, or PR range>

Constraints:
<what not to change, style rules, compatibility requirements, performance/security limits>

Ownership:
<files or feature area this subagent owns; note any files/areas it must avoid>

Expected output:
<plan, patch summary, review findings, test results, evidence, open questions>
```

## Consolidation

After subagents finish:

- Summarize each subagent result briefly.
- Call out conflicts between agents.
- Identify the recommended next action.
- Include verification status and remaining risks.
