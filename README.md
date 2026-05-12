# Codex Dispatch Package

A lean Codex dispatch package that selects and spawns the right named custom subagent for each task.

It includes:

- `SKILL.md` — dispatch skill instructions
- `.codex/agents/planner.toml`
- `.codex/agents/reviewer.toml`
- `.codex/agents/coding_worker.toml`
- `.codex/agents/fast_coding_worker.toml`
- `.codex/agents/helper_worker.toml`
- `.codex/agents/premium_reviewer.toml`

## What this package does

This package turns Codex into a small dispatcher/coordinator. The parent agent should not merely recommend a model. It should spawn one or more named subagents:

| Subagent | Model | Use |
|---|---:|---|
| `planner` | `gpt-5.5` | Architecture, decomposition, sequencing, risk analysis |
| `reviewer` | `gpt-5.5` | Standard code review |
| `coding_worker` | `gpt-5.3-codex` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `gpt-5.3-codex-spark` | Small localized edits and quick fixes |
| `helper_worker` | `gpt-5.4-mini` | Read-only investigation, lookup, repo reconnaissance, docs/API help |
| `premium_reviewer` | `gpt-5.5-pro` | Rare high-stakes/high-cost review only |

## Installation

Copy the package contents into your repo root.

Expected layout:

```text
your-repo/
  SKILL.md
  README.md
  .codex/
    agents/
      planner.toml
      reviewer.toml
      coding_worker.toml
      fast_coding_worker.toml
      helper_worker.toml
      premium_reviewer.toml
```

If you want this as a Codex repo skill, place `SKILL.md` under your repository skill location instead, for example:

```text
your-repo/
  .agents/
    skills/
      codex-dispatch/
        SKILL.md
```

Keep `.codex/agents/*.toml` at the repo root so Codex can load the project-scoped custom agents.

## Recommended optional config

Codex supports global subagent limits under `[agents]`. Add this to `.codex/config.toml` or `~/.codex/config.toml` if you want conservative fan-out:

```toml
[agents]
max_threads = 6
max_depth = 1
```

## Practical examples

### 1. Plan a complex feature, then implement and review

```text
Use codex-dispatch.

Build tenant-level API keys with creation, rotation, revocation, and audit logging.

Spawn planner first to produce a short implementation plan and risk list.
Then spawn coding_worker to implement the accepted plan.
Then spawn reviewer for standard review.
```

Expected routing:

```text
planner -> coding_worker -> reviewer
```

### 2. Small bug fix

```text
Use codex-dispatch.

Fix the failing validation around blank display names.
This should be a small localized change. Spawn the right worker and include test evidence.
```

Expected routing:

```text
fast_coding_worker
optional reviewer if behavior changed
```

### 3. Unknown failure path

```text
Use codex-dispatch.

The checkout flow sometimes double-charges in staging.
Do not edit code yet. Spawn a helper to trace the likely code path and identify evidence.
After that, decide whether to spawn coding_worker.
```

Expected routing:

```text
helper_worker -> coding_worker -> reviewer
```

### 4. Standard PR review

```text
Use codex-dispatch.

Review this branch against main for correctness, security, regressions, and missing tests.
Use the standard review agent.
```

Expected routing:

```text
reviewer
```

### 5. High-stakes review

```text
Use codex-dispatch.

Review the auth/session migration before release.
This affects login, session invalidation, and production user access.
Use premium review only if it meets the high-stakes threshold.
```

Expected routing:

```text
premium_reviewer
```

Use `premium_reviewer` sparingly. It is for security boundaries, auth/session logic, payments, data loss risk, production migrations, irreversible changes, and expensive operational decisions.

### 6. Parallel focused review

```text
Use codex-dispatch.

Review this PR in parallel:
1. correctness and edge cases
2. security and permission boundaries
3. tests and regression risk

Spawn named review subagents and consolidate the findings.
```

Expected routing:

```text
reviewer x3, each with a distinct review angle
```

## Dispatch policy

Default to the cheapest capable subagent that preserves quality:

1. `helper_worker` for read-only discovery.
2. `fast_coding_worker` for small localized edits.
3. `coding_worker` for normal implementation.
4. `planner` when sequencing, architecture, or risk is the main problem.
5. `reviewer` for ordinary review.
6. `premium_reviewer` only for rare high-stakes review.

## Notes

- The agent TOML files pin each role to the model requested.
- Reviewers are read-only.
- Workers are workspace-write.
- The skill is intentionally lean; it exists to route and coordinate, not to impose a heavy process.
