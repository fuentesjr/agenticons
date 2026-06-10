---
name: agenticons
description: Selects and spawns named Codex custom subagents for planning, implementation, review, documentation review, deep investigation, and helper work. Use when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.
---

# Agenticons

Use this skill as a lean dispatcher when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing. Your job is to select and spawn the right named custom subagent for the task, then consolidate results.

## Hard rules

- Dispatch through subagents only when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.
- Honor escape hatches exactly. If the user says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use agenticons`, do not spawn subagents for that request.
- Dispatch to named subagents. Do not merely recommend models.
- The parent agent is the orchestrator and DRA (Directly Responsible Agent). As DRA, the parent remains accountable for the project outcome: it owns routing, scope, sequencing, conflict resolution, verification, deciding whether to accept subagent output, and the final response.
- Subagents must not delegate or route. They return findings/results to the parent orchestrator; their output is advisory until the parent accepts it.
- Use only the model assigned in each installed Agenticons agent spec. Do not override a spawned Agenticons subagent to any model or provider not listed in this package's `.codex/agents/*.toml` files.
- Use the exact subagent names from the installed Agenticons specs (`.codex/agents/*.toml` for repo-local installs or `~/.codex/agents/*.toml` for global installs):
  - `planner`
  - `coding_worker`
  - `fast_coding_worker`
  - `helper_worker`
  - `forensic_analyst`
  - `doc_reviewer`
  - `reviewer`
- Standard review always uses `reviewer`.
- Documentation drift review uses `doc_reviewer`.
- Deep root-cause investigation uses `forensic_analyst`.
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
| Read-only lookup, repo reconnaissance, dependency/API check, test triage, docs/API lookup | `helper_worker` |
| Unclear root cause, intermittent or hard-to-reproduce failure, regression with no obvious cause, cross-system failure | `forensic_analyst` |
| Documentation correctness, stale docs, README/API drift, changelog/release-note coverage | `doc_reviewer` |
| Normal correctness/security/maintainability review | `reviewer` |

Rule of thumb for investigations: if you already know roughly where to look, use `helper_worker`; if the question is why something fails and nobody knows, use `forensic_analyst`.

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

Use `reviewer` for normal review, including security-sensitive review.

### Documentation drift review

Spawn `doc_reviewer` when code behavior, public APIs, installation steps, commands, configuration, examples, or user-facing workflows may have changed. Ask it to compare implementation and docs, identify stale or missing documentation, and recommend exact doc updates.

### Investigation before editing

Spawn `helper_worker` first when you know roughly where to look but need the correct files, runtime path, or API behavior confirmed. Then hand the evidence to `coding_worker` or `fast_coding_worker`.

### Deep investigation

Spawn `forensic_analyst` when the question is why something fails and nobody knows: unclear root cause, intermittent or hard-to-reproduce failures, regressions with no obvious cause, or failures spanning systems.

1. Give it the symptom, scope, and all evidence gathered so far.
2. It returns a standalone Markdown report: summary, symptom and scope, evidence, ranked hypotheses, causal analysis, reproduction, recommended next steps, and open questions.
3. If the user asks to save the report, write the accepted report verbatim to the requested path. The analyst is read-only and does not write files.
4. Hand the confirmed cause to `coding_worker` or `fast_coding_worker` for the fix.

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
