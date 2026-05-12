---
name: agenticons
description: Selects and spawns named Codex custom subagents for planning, implementation, review, documentation review, and helper work. Use when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.
---

# Agenticons

Use this skill as a lean dispatcher when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing. Your job is to select and spawn the right named custom subagent for the task, then consolidate results.

## Hard rules

- Dispatch through subagents only when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.
- Honor escape hatches exactly. If the user says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use agenticons`, do not spawn subagents for that request.
- Dispatch to named subagents. Do not merely recommend models.
- Use the exact subagent names from the installed Agenticons specs (`.codex/agents/*.toml` for repo-local installs or `~/.codex/agents/*.toml` for global installs):
  - `planner`
  - `coding_worker`
  - `fast_coding_worker`
  - `helper_worker`
  - `doc_reviewer`
  - `reviewer`
  - `premium_reviewer`
- Standard review always uses `reviewer`.
- Documentation drift review uses `doc_reviewer`.
- Never spawn `premium_reviewer` unless the user explicitly asks for `premium_reviewer` or premium review in the current request. High-stakes context alone is not enough.
- If the user asks to use premium review only if needed, spawn `reviewer` first and ask it to say whether premium review is justified.
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
| Explicitly requested premium review for a rare high-stakes case: security boundary, data loss, payments, auth, migration, production incident, irreversible or expensive decision | `premium_reviewer` |

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

Use `reviewer` for normal review, including security-sensitive review. If the user explicitly asks for premium review in the current request, use `premium_reviewer` only for rare high-stakes cases. If the user asks for premium only if needed, use `reviewer` first and report whether escalation is justified.

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
