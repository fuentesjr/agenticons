---
name: agenticons
description: Selects and spawns named Codex custom subagents for technical advising, planning, implementation, review, documentation review, deep investigation, exploratory QA verification, and helper work. Use when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing. Do not use when the user asks the parent to handle the work locally without subagents.
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
  - `advisor`
  - `planner`
  - `coding_worker`
  - `fast_coding_worker`
  - `helper_worker`
  - `forensic_analyst`
  - `doc_reviewer`
  - `reviewer`
  - `qa_engineer`
  - `edge_case_analyst`
- Principal-level technical direction advice uses `advisor`.
- Standard review always uses `reviewer`.
- Documentation drift review uses `doc_reviewer`.
- Deep root-cause investigation uses `forensic_analyst`.
- Exploratory QA verification uses `qa_engineer`.
- Edge-case and coverage analysis uses `edge_case_analyst`.
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
| Technical direction with multiple credible approaches, cross-system boundaries, expensive reversals, conflicting recommendations, or trajectory drift | `advisor` |
| Complex feature, vague request, accepted architecture, decomposition, or phased implementation plan | `planner` |
| Normal code implementation, multi-file edit, bug fix, refactor | `coding_worker` |
| Small localized edit, mechanical change, simple failing test | `fast_coding_worker` |
| Read-only lookup, repo reconnaissance, dependency/API check, test triage, docs/API lookup | `helper_worker` |
| Unclear root cause, intermittent or hard-to-reproduce failure, regression with no obvious cause, cross-system failure | `forensic_analyst` |
| Documentation accuracy/drift, obsolete low-value docs to remove, AI-confusing docs, changelog coverage | `doc_reviewer` |
| Normal correctness/security/maintainability review | `reviewer` |
| Feature or release that should be exercised end-to-end, change-aware exploratory testing, performance regression probing, rough-edge hunting | `qa_engineer` |
| Uncovered edge cases, missing acceptance criteria, design gaps, untested behavior, coverage holes | `edge_case_analyst` |

Rule of thumb for direction: `advisor` challenges whether the technical decision is sound; `planner` turns an accepted direction into an implementation sequence; `reviewer` judges the resulting change. Use `advisor` only when the decision has enough leverage to justify a separate advisory pass.

Rule of thumb for investigations: if you already know roughly where to look, use `helper_worker`; if the question is why something fails and nobody knows, use `forensic_analyst`.

Rule of thumb for verification: `reviewer` reads the change; `qa_engineer` runs it. Use `qa_engineer` when confidence requires executing the software, not just reviewing the diff.

Rule of thumb for documentation: use `helper_worker` to look up or summarize what existing docs say; use `doc_reviewer` to audit whether docs match implementation (drift, gaps, overpromises), to recommend removing obsolete low-value docs, and to flag anything that could confuse AI agents. Prefer delete/archive over leaving agent-misleading prose. `doc_reviewer` does not edit files.

Rule of thumb for escalation: start with the cheaper read-only agent (`helper_worker`) and escalate to `forensic_analyst` only when recon does not surface the cause. Match model tier to task difficulty; do not pay for capability the task does not need. For security/threat-model docs or precise public API contracts, pair `doc_reviewer` with `reviewer` on the related code — do not treat doc drift review alone as a security or correctness sign-off.

Rule of thumb for coverage: `reviewer` judges the change as written and `qa_engineer` runs it; `edge_case_analyst` asks what cases were never considered or tested and specifies them with concrete test cases.

## Dispatch patterns

### Advise before planning

Spawn `advisor` when implementation should wait for a principal-level technical decision.

1. Give it the decision, constraints, relevant evidence, inherited decisions, and credible alternatives already identified.
2. It returns a recommendation, tradeoffs, assumptions, reversibility analysis, and the evidence that would change its recommendation.
3. The parent decides whether to accept the direction.
4. If implementation is non-trivial, spawn `planner` with the accepted direction.

Do not make `advisor` a mandatory gate. Reinvoke it only when new evidence, changed constraints, conflicting recommendations, or trajectory drift reopens the decision.

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

Spawn `doc_reviewer` when code behavior, public APIs, installation steps, commands, configuration, examples, user-facing workflows, or agent-facing instructions may have changed — or when docs may be obsolete, redundant, or confusing to AI agents.

1. Give it the change scope (branch, commits, PR, or feature) and the docs surfaces to check when known (README, install, CLI, config, API reference, examples, changelog, agent/skill instruction files).
2. It works change-aware: maps user- and agent-facing changes to docs, prefers implementation evidence over prose, treats AI-confusing docs as critical, recommends remove/archive for obsolete low-value material, and returns a standalone Markdown report (summary, scope, ranked findings with update/remove actions, prune candidates, current/verified, not covered, next steps).
3. If the user asks to save the report, write the accepted report verbatim to the requested path. The reviewer is read-only and does not write files.
4. Route confirmed doc updates and removals to `fast_coding_worker` for small localized edits/deletes, or `coding_worker` for larger multi-file rewrites or cleanups.
5. When the same change also needs code correctness or security review (especially threat-model or public API contract docs), spawn `reviewer` on the related code in parallel or after; do not substitute `doc_reviewer` for `reviewer`.

### Investigation before editing

Spawn `helper_worker` first when you know roughly where to look but need the correct files, runtime path, or API behavior confirmed. Then hand the evidence to `coding_worker` or `fast_coding_worker`.

### Deep investigation

Spawn `forensic_analyst` when the question is why something fails and nobody knows: unclear root cause, intermittent or hard-to-reproduce failures, regressions with no obvious cause, or failures spanning systems.

1. Give it the symptom, scope, and all evidence gathered so far.
2. It returns a standalone Markdown report: summary, symptom and scope, evidence, ranked hypotheses, causal analysis, reproduction, recommended next steps, and open questions.
3. If the user asks to save the report, write the accepted report verbatim to the requested path. The analyst is read-only and does not write files.
4. Hand the confirmed cause to `coding_worker` or `fast_coding_worker` for the fix.

### Exploratory QA verification

Spawn `qa_engineer` after a feature lands or before a release, when the change should be exercised rather than read: end-to-end behavior, regressions in adjacent features, performance changes, and rough edges a user would hit.

1. Give it the change scope (branch, commits, or feature), how to build and run the software, and any environment details it needs.
2. It inspects what changed, runs targeted scenarios, and returns a standalone Markdown verification report: verdict, scenarios executed, findings with reproduction steps, performance notes, rough edges, and what was not covered.
3. It must not modify existing source, tests, or docs; any scratch artifacts it creates are listed in the report.
4. Route confirmed findings to `coding_worker` or `fast_coding_worker`, including a regression test for each confirmed bug.

### Edge-case and coverage analysis

Spawn `edge_case_analyst` to find cases nobody considered and specify them: missing acceptance criteria, design edge cases, and untested behavior across input, failure, concurrency, and security dimensions.

1. Give it the feature, change scope, or area to analyze, plus where the existing tests and specs live so it can tell covered from uncovered.
2. It returns a standalone Markdown report: uncovered cases ranked by risk, each with scenario, expected behavior, a concrete test case, and coverage evidence, plus proposed acceptance criteria and open questions.
3. If the user asks to save the report, write the accepted report verbatim to the requested path. The analyst is read-only and does not write files.
4. Route confirmed cases to `coding_worker` or `fast_coding_worker` to add the tests and any fix.

## User-facing labels

- Maintain a meaningful local label for every spawned subagent in user-facing updates and summaries.
- Use labels in the form `<role>: <task or scope>`, for example `helper_worker: PackageDependencyPressure readiness review`.
- Treat tool-generated nicknames and agent ids as traceability metadata only. If needed, put them after the semantic label in parentheses; do not use them as the primary name.
- When consolidating results, refer to subagents by their semantic labels so readers can tell what each agent was responsible for without decoding generated nicknames.

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

## Output contract

After subagents finish:

- `output-results`: Summarize each subagent result briefly.
- `output-conflicts`: Call out conflicts between agents.
- `output-next-action`: Identify the recommended next action.
- `output-verification`: Include verification status and remaining risks.
