# Agenticons

Agenticons is a Codex skill package for users who want explicit, named subagent delegation for planning, implementation, review, documentation review, and investigation work.

## Why This Exists

Codex can delegate work to subagents, but repeated manual routing is easy to make inconsistent. Agenticons packages a small routing skill plus named agent specs so a user can ask for delegation once and get predictable roles such as `planner`, `coding_worker`, `reviewer`, or `doc_reviewer`.

```text
Use agenticons.

Review this branch against main for correctness, security, regressions, and missing tests.
Use the standard review agent.
```

## Install

```bash
TARGET_REPO=/path/to/your-repo
mkdir -p "$TARGET_REPO/.agents/skills/agenticons" "$TARGET_REPO/.codex/agents"
cp SKILL.md "$TARGET_REPO/.agents/skills/agenticons/SKILL.md"
cp .codex/agents/*.toml "$TARGET_REPO/.codex/agents/"
```

## Quick Start

Paste this into Codex from the target repository after installation:

```text
Use agenticons.

Review this branch against main for correctness, security, regressions, and missing tests.
Use the standard review agent.
```

Expected routing:

```text
reviewer
```

## Usage

Agenticons dispatches only when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing. If a request says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use agenticons`, the parent agent should handle it directly.

| Subagent | Agent file | Model | Use |
|---|---|---:|---|
| `planner` | `.codex/agents/planner.toml` | `gpt-5.5` | Architecture, decomposition, sequencing, risk analysis |
| `coding_worker` | `.codex/agents/coding_worker.toml` | `gpt-5.3-codex` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `.codex/agents/fast_coding_worker.toml` | `gpt-5.3-codex-spark` | Small localized edits and quick fixes |
| `helper_worker` | `.codex/agents/helper_worker.toml` | `gpt-5.4-mini` | Read-only investigation, lookup, repo reconnaissance, docs/API help |
| `doc_reviewer` | `.codex/agents/doc_reviewer.toml` | `gpt-5.4-mini` | Documentation correctness, stale docs, doc drift |
| `reviewer` | `.codex/agents/reviewer.toml` | `gpt-5.5` | Standard code review |
| `premium_reviewer` | `.codex/agents/premium_reviewer.toml` | `gpt-5.5-pro` | Rare high-stakes/high-cost review only |

Plan, implement, and review:

```text
Use agenticons.

Build tenant-level API keys with creation, rotation, revocation, and audit logging.
Spawn planner first to produce a short implementation plan and risk list.
Then spawn coding_worker to implement the accepted plan.
Then spawn reviewer for standard review.
```

Investigate before editing:

```text
Use agenticons.

The checkout flow sometimes double-charges in staging.
Do not edit code yet. Spawn a helper to trace the likely code path and identify evidence.
After that, decide whether to spawn coding_worker.
```

Review documentation drift:

```text
Use agenticons.

Review this branch for documentation drift. Check README, setup steps, CLI examples,
configuration, API docs, and release notes against the implementation changes.
```

Run a high-stakes review:

```text
Use agenticons.

Review the auth/session migration before release.
This affects login, session invalidation, and production user access.
Use premium review only if it meets the high-stakes threshold.
```

For more detail, see [docs/faq.md](docs/faq.md) and [docs/design.md](docs/design.md).

## Configuration

| File | Setting | Purpose |
|---|---|---|
| `.codex/config.toml` or `~/.codex/config.toml` | `[agents].max_threads` | Caps total subagent fan-out. |
| `.codex/config.toml` or `~/.codex/config.toml` | `[agents].max_depth` | Caps recursive delegation depth. |

```toml
[agents]
max_threads = 6
max_depth = 1
```

## Requirements

| Requirement | Version or value |
|---|---|
| Codex | Custom skills and custom subagents enabled |
| Go | 1.22 or newer, for validation and tests |
| Installed skill path | `.agents/skills/agenticons/SKILL.md` |
| Installed agent path | `.codex/agents/*.toml` |

## Contributing / Development

```bash
go run ./scripts/validate_package.go
go test ./...
go vet ./...
```

The validator parses every agent TOML file and checks that `README.md` and `SKILL.md` mention every configured agent. It is also run by `.github/workflows/validate.yml`.

## License

No license has been specified.
