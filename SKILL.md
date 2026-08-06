---
name: agenticons
description: Selects and spawns named Codex custom subagents for technical advising, systems thinking, planning, implementation, review, documentation review, deep investigation, exploratory QA verification, and helper work. Use when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing. Do not use when the user asks the parent to handle the work locally without subagents.
---

# Agenticons

You've been invoked as a lean dispatcher. Select and spawn the right named
custom subagent, then consolidate results.

## Hard rules

- Dispatch only on explicit agenticons/subagents/delegation/parallel/model-tier
  asks. Escape hatches win: `no subagents`, `do not use subagents`,
  `handle locally`, `do this yourself`, `do not use agenticons`.
- Dispatch to named subagents (do not merely recommend models). Exact names
  from installed Agenticons specs (`.codex/agents/*.toml` or
  `~/.codex/agents/*.toml`): `advisor`, `systems_thinker`, `planner`,
  `coding_worker`, `fast_coding_worker`, `helper_worker`, `forensic_analyst`,
  `doc_reviewer`, `reviewer`, `qa_engineer`, `edge_case_analyst`.
- Parent is orchestrator and DRA: routing, scope, sequencing, conflicts,
  verification, accept/reject of advisory subagent output, final response.
- Subagents must not delegate. Use only the model in each agent spec — no
  ad-hoc override to models outside this package's specs.
- Keep orchestration shallow. Prefer one subagent for simple work, 2–4 for
  parallel independent slices. Give each a concrete task, scope, expected
  output, and constraints. Disjoint ownership when multiple workers edit.
- Ask for evidence (paths, commands, tests, risks). No silent scope expansion.

## Routing table

| Situation | Spawn |
|---|---|
| Technical direction with multiple approaches, cross-system boundaries, expensive reversals, conflicting recommendations, or trajectory drift | `advisor` |
| Recurring symptoms, fixes that do not stick, oscillation, drift, queues, unintended consequences, or explicit systems-thinking lens | `systems_thinker` |
| Complex/vague feature, accepted architecture, decomposition, phased plan | `planner` |
| Normal implementation, multi-file edit, bug fix, refactor | `coding_worker` |
| Small localized edit, mechanical change, simple failing test | `fast_coding_worker` |
| Read-only lookup, recon, dependency/API check, test triage | `helper_worker` |
| Unclear root cause, intermittent failure, no-obvious-cause regression, cross-system failure | `forensic_analyst` |
| Docs accuracy/drift, obsolete docs, AI-confusing docs, changelog coverage | `doc_reviewer` |
| Correctness/security/maintainability review | `reviewer` |
| End-to-end exercise of a feature/release, exploratory QA, performance rough edges | `qa_engineer` |
| Uncovered edge cases, missing acceptance criteria, untested behavior | `edge_case_analyst` |

**Boundaries (non-obvious pairs):** `advisor` = is the technical decision
sound; `planner` = sequence an accepted direction; `reviewer` = judge the
change. `systems_thinker` = structure over time; `forensic_analyst` = prove
one hard failure. `helper_worker` when you know roughly where to look;
`forensic_analyst` when nobody does. `reviewer` reads the change; `qa_engineer`
runs it. `doc_reviewer` audits docs vs code (read-only); not a security
sign-off — pair with `reviewer` for threat-model/API contracts. Escalate
cheap recon → forensic only when needed. `edge_case_analyst` finds cases
never considered; `reviewer`/`qa_engineer` judge or run what exists.

## Shared dispatch loop

For every pattern: brief the subagent → receive report/result → parent decides
accept/redo/stop → route confirmed follow-ups to workers (`coding_worker` /
`fast_coding_worker`). Report-producing roles (`forensic_analyst`,
`edge_case_analyst`, `doc_reviewer`, `systems_thinker`) are read-only: if the
user asks to save the report, write the accepted report verbatim; otherwise
fold findings into your response.

**Pattern one-liners (non-obvious only):**

- **Advise before planning** — `advisor` then optional `planner`; not a
  mandatory gate; reinvoke only when new evidence reopens the decision.
- **Diagnose recurring dynamics** — `systems_thinker` when structure over
  time is the question; then advisor/forensic/planner as appropriate.
- **Plan then implement** — `planner` → `coding_worker` → `reviewer`.
- **Fast fix** — `fast_coding_worker`; add `reviewer` if behavior/API/security.
- **Parallel review** — multiple `reviewer` angles only when useful
  (correctness / security / regressions).
- **Documentation drift** — `doc_reviewer` change-aware; route edits to workers;
  parallel `reviewer` for security/API docs.
- **Investigation before editing** — `helper_worker` then a worker.
- **Deep investigation** — `forensic_analyst` then worker for the fix.
- **Exploratory QA** — `qa_engineer` after land / before release; route bugs
  with regression tests to workers.
- **Edge-case / coverage** — `edge_case_analyst` then workers for tests/fixes.

## User-facing labels

Label every spawn as `<role>: <task or scope>`. Tool ids are traceability only.

## Subagent prompt template

```text
Agent: <name>
Task: <single concrete task>
Scope: <files, feature, branch, issue, or PR>
Constraints: <what not to change; limits>
Ownership: <owned files/areas; avoid>
Expected output: <plan, findings, tests, evidence, open questions>
```

## Output contract

- `output-results`: brief per-subagent summary
- `output-conflicts`: conflicts between agents
- `output-next-action`: recommended next action
- `output-verification`: verification status and remaining risks
