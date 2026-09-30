# Choosing Models for Agent Roles

Use this guide when proposing a model or reasoning-effort change. Optimize for
**quality first, then speed**. Cost does not decide between candidates.
The owner approves role assignments; this guide supplies evidence for that
decision. Current assignments remain in [the design role table](design.md#agent-roles)
and `.codex/agents/*.toml`.

## Define the Decision

Compare a complete configuration: model version, reasoning effort, agent harness
and version, role prompt, tools, permissions, and execution limits. A result for
Astra max does not establish Astra low's performance. Sol 6 and Sol 6.1 are
different candidates.

Before searching, write down:

- The role and its typical tasks.
- What a correct result must contain, including evidence and boundary discipline.
- Failures that would disqualify a candidate, such as fabricated findings or
  unauthorized production edits.
- The incumbent and proposed configurations.
- Whether the decision concerns a parent orchestrator or a bounded subagent.

## Choose Evidence That Matches the Role

Use the most relevant evidence available, in this order: controlled local role
comparisons, external benchmarks matching the task and configuration, broader
evaluations, then vendor guidance. A weak local test need not outweigh a strong
external evaluation; record coverage and limitations for both.

| Source | Best use | Limits to preserve |
|---|---|---|
| [Artificial Analysis Coding Agent Index](https://artificialanalysis.ai/agents/coding-agents) | Implementation, terminal workflows, repository understanding; compare evaluated configurations and active task runtime | Inspect the component scores and harness. A coding composite does not directly measure security review or architecture advice. |
| [Arena Agent Code](https://arena.ai/leaderboard/agent/code) | Interactive orchestration, accepting corrections, tool reliability, and user feedback during coding | Confirm the exact model and effort are listed. Arena's harness and user task mix differ from local Codex. |
| [Artificial Analysis OpenAI provider comparison](https://artificialanalysis.ai/providers/openai) | Broad capability screening and API response characteristics | Aggregate intelligence and 500-token response times cannot substitute for coding-agent results or task completion time. |
| [OpenAI model selection](https://developers.openai.com/api/docs/guides/model-selection) | Supported use cases and a starting shortlist | Vendor positioning is not a controlled comparison of effort settings. |
| [OpenAI reasoning guidance](https://developers.openai.com/api/docs/guides/reasoning#reasoning-effort) | Understanding effort controls and supported settings | More effort is not evidence of better measured performance on a particular role. |

Artificial Analysis's coding index equally weights DeepSWE v1.1, Terminal-Bench
4.0, and SWE-Atlas-QnA, using three attempts per task. Its reported task time
excludes environment startup and verifier or judge time. See its
[benchmark description](https://artificialanalysis.ai/agents/coding-agents).

Arena estimates orchestrator effects from randomized model assignments and
signals including user feedback, correction handling, command recovery, and
tool hallucination. Its improvement percentages are relative effects against
a baseline, not raw task-success percentages. See its
[methodology](https://arena.ai/blog/agent-arena-methodology).

## Match Each Role to a Quality Check

The following are proposed evaluation criteria, not claims about a model's
strengths. Use external results to shortlist candidates; use representative
role tasks to resolve gaps.

| Role | Most relevant evidence | Local acceptance criteria |
|---|---|---|
| Parent orchestrator, outside the installed roster | Arena agent signals and complete local workflows | Preserves requirements, handles corrections, delegates within scope, verifies accepted work, completes the task |
| `coding_worker` | Coding-agent implementation and terminal results | Correct patch, relevant checks pass, existing conventions followed, limited rework |
| `fast_coding_worker` | Local small edits; coding benchmarks as background | Exact requested edit, no collateral changes, low time to accepted result |
| `helper_worker` | Repository Q&A and retrieval tasks | Correct locations and citations, useful coverage, no invented behavior |
| `advisor` | Local design decisions; broad reasoning as background | Identifies real constraints and alternatives, explains reversibility, preserves owner authority |
| `systems_thinker` | Local recurring-problem cases | Supported causal explanation, useful intervention, explicit detection path |
| `planner` | Local plans checked against completed implementations | Correct dependencies, executable steps, risks and verification accounted for |
| `forensic_analyst` | Local incident and debugging cases; terminal evidence | Reproduces or explains the failure, distinguishes evidence from hypotheses, identifies the cause |
| `reviewer` | Local changes with known defects | Finds consequential defects with evidence; low false-positive burden |
| `edge_case_analyst` | Local specifications and defect cases | Finds distinct missing cases and defines expected behavior and meaningful checks |
| `doc_reviewer` | Repository Q&A and local documentation drift | Correctly traces claims to implementation and identifies consequential omissions |
| `qa_engineer` | Local browser and application tasks | Reproduces failures, covers important paths, supplies actionable evidence |
| `observability_engineer` | Local instrumentation and incident cases | Relates telemetry to user symptoms, proposes useful detection, respects edit boundaries |
| `security_auditor` | Local audits with known vulnerabilities | Finds required vulnerabilities, confirms findings safely, respects trust and role boundaries |

For systems, observability, and security roles, also use the behavioral checks
already defined in [the design document](design.md). Do not infer audit quality
from patch-generation scores alone.

## Read Results Without Overstating Them

1. Record the source URL, access date, benchmark version, exact configuration,
   component scores, sample size, uncertainty, and timing definition. Mark
   information the source does not disclose as unknown.
2. Separate observations, interpretations, and proposals. A measured score is
   an observation; assigning a role from that score is a proposal.
3. Compare within the same evaluation. Do not average Arena percentages with
   Artificial Analysis scores or compare scores across benchmark versions.
4. Check uncertainty and rounding. A one-point lead is not automatically a
   meaningful advantage. Overlapping intervals alone do not prove equivalence
   or settle a pairwise significance test.
5. Treat an absent configuration as missing evidence. Do not substitute an
   older release or infer low-effort results from max-effort results.
6. Inspect task-specific scores before the composite. Identify which abilities
   the role actually needs and which important failures the benchmark misses.
7. Compare speed only after quality. Prefer elapsed time to an accepted result,
   including retries, verification, and human correction. Keep output throughput,
   initial latency, active agent runtime, and total workflow time distinct.

If sources disagree, explain differences in task mix, harness, effort, and
measurement before changing the recommendation. Preserve unresolved uncertainty
instead of declaring whichever leaderboard was read last the winner.

## Resolve Gaps With a Bounded Local Comparison

Propose a small trial when external evidence does not answer the role decision.
Agree on the tasks and acceptance criteria before running it. Include routine
work, difficult work, and a case that previously exposed a failure.

Run candidates from the same repository state with the same inputs, role prompt,
tools, permissions, and completion criteria. Use fresh sessions and multiple
attempts where variability could change the decision. Record client versions,
service tier, environment differences, and whether cache conditions were comparable.

Score outputs without model labels where practical. Record correctness,
completeness, evidence quality, boundary violations, false positives, human
corrections, and time to accepted completion. Keep failures in the results.

Reject candidates with critical failures. Choose the stronger quality result;
use speed to choose between candidates whose quality is practically equivalent
under the agreed criteria. If evidence is inconclusive, retain the incumbent
and name the additional evidence needed. This is a decision aid, not a claim
that a small trial establishes statistical superiority.

## Evidence Snapshot: 2026-09-30

These observations seed a shortlist. They do not approve configuration changes.
Refresh the linked results before using this snapshot for a later decision.

Artificial Analysis reported these Codex variants:

| Configuration | Coding index | DeepSWE | Terminal-Bench | SWE-Atlas-QnA | Average active time |
|---|---:|---:|---:|---:|---:|
| Sol 6.1 xhigh | 63 | 73% | 55% | 61% | 15.5 min |
| Astra max | 62 | 68% | 56% | 62% | 29.4 min |
| Sol 6.1 medium | 61 | 72% | 52% | 61% | 10.9 min |
| Sol 6.1 high | 60 | 71% | 50% | 60% | 13.3 min |
| Sol 6.1 max | 60 | 70% | 53% | 58% | 24.4 min |

Source: [Codex model variants](https://artificialanalysis.ai/agents/coding-agents/comparisons/codex-vs-kimi-code-cli#model-variants).
The listed Codex variants did not include Astra low or medium. Sol medium
outscored high in this run; increasing effort did not produce monotonic gains.

Arena's September 29 table placed Astra max second and Sol 6 max fourth.
It did not list Sol 6.1 or Astra low. These results support considering Astra
max for orchestration but do not compare the original three candidates.
Source: [Arena Agent Code](https://arena.ai/leaderboard/agent/code).

**Proposed shortlist:** Sol 6.1 medium and xhigh for implementation; Astra max
as an orchestration candidate. Retain Astra low in an implementation trial as
the incumbent. Sol high remains eligible, but these coding results provide no
measured advantage over medium. Specialized review and advisory roles need
their own evidence. No local comparison was run for this snapshot.

## Record the Decision

### Approved Assignment: 2026-09-30

The owner approved applying the recommendations in the accompanying discussion.
The evidence snapshot above remains a historical record of the shortlist.

| Role | Previous configuration | Approved configuration |
|---|---|---|
| `coding_worker` | Astra low | Sol 6.1 xhigh |
| `fast_coding_worker` | Luna low | Sol 6.1 medium |
| `helper_worker` | Terra medium | Sol 6.1 medium, provisional trial |

The remaining specialist assignments are retained. Sol's coding-agent results
support the implementation choices; the helper assignment is an extrapolation
from repository-understanding results and needs local confirmation. No controlled
local role comparison has been completed. In particular, the external coding
comparison does not establish a direct win over the previous low-effort workers.

The recommended parent configuration is Astra max. The parent is outside the
installed roster and must be selected in the host client; these repository edits
do not change a running session's model.

Follow-up owner: the repository owner, assisted by the parent agent reviewing
results. Review after the first five completed tasks for each changed role,
before treating the helper trial as settled. Record correctness, supported
findings, scope violations, human rework, and elapsed time to accepted completion.
Reopen the choice immediately for fabricated evidence, unauthorized edits, or
a consequential correctness failure. Compare ordinary quality and timing with
recorded incumbent tasks where available; otherwise mark the baseline unknown.
If a regression is confirmed, propose restoring the previous configuration for
owner approval. The task sample is a regression screen, not proof of superiority.

### Decision Template

Copy this template into the change proposal. Keep evidence links and the owner's
decision together so later model releases do not erase the reasoning.

```text
Role and representative tasks:
Decision date and owner:
Incumbent model, effort, prompt, harness, and tools:
Candidate configurations:
Required quality and disqualifying failures:

Evidence receipts (URL or artifact, access date, benchmark version):
Observed results and uncertainty:
Configuration or workload mismatches:
Local comparison results, or why none was run:

Recommendation and alternatives:
Expected quality effect:
Expected speed effect and timing definition:
Unknowns and evidence that would reverse the recommendation:
Owner decision: pending / approved / rejected

Approved configuration changes:
Verification results:
Follow-up owner and review date or task count:
Signals to inspect and rollback condition:
```

## Apply and Revisit an Approved Choice

Follow [the package change policy](design.md#change-policy): update the agent
spec and both documented role tables together, then run the package checks.
The selection record does not authorize runtime overrides of a role's model.

Before adopting a new assignment, name who will inspect the first agreed batch
of completed tasks and when. Compare accepted-result quality, correction burden,
and completion time with the baseline. Revert or reopen the decision if the
recorded failure threshold is crossed. Review on that date or task count;
do not wait for a user complaint to discover a regression.

Reopen the evidence review after a relevant model release, harness change,
benchmark revision, or observed regression. This guide introduces no scheduled
automation. Any future automation needs a separate decision and detection plan.
