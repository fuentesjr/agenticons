# Agenticons
<img width="1672" height="941" alt="agenticons_ruby_banner" src="https://github.com/user-attachments/assets/91d20c4f-d80b-4069-934b-aa37f6d735c6" />



Agenticons is a Codex skill package for users who want explicit, named subagent delegation for planning, implementation, review, documentation review, investigation, and QA verification work.

## Why This Exists

Codex can delegate work to subagents, but repeated manual routing is easy to make inconsistent. Agenticons packages a small routing skill plus named agent specs so a user can ask for delegation once and route it through a fixed, explicit set of named roles such as `planner`, `coding_worker`, `reviewer`, or `doc_reviewer`.

```text
Use agenticons.

Review this branch against main for correctness, security, regressions, and missing tests.
Use the standard review agent.
```

## Install

Install into one repository:

```bash
./scripts/install.sh --target /path/to/your-repo
```

Or install globally for the current user so Agenticons is available from any repository:

```bash
./scripts/install.sh --global
```

Directly from GitHub:

```bash
curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --target /path/to/your-repo
curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --global
```

## Quick Start

Paste this into Codex from the target repository after repo-local installation, or from any repository after global installation:

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

Model routing is fixed by the installed Agenticons agent specs. Spawn each role with the model listed in its `.codex/agents/*.toml` file and in the table below; do not substitute unlisted models or providers at dispatch time.

The parent agent is the orchestrator and DRA (Directly Responsible Agent). That means the parent remains accountable for the project outcome: subagents do not delegate or route; they return advisory findings/results to the parent, who owns sequencing, verification, conflict resolution, deciding what to accept, and the final response.

Standard review always routes to `reviewer`, including security-sensitive and high-stakes review requests.

| Subagent | Agent file | Model | Use |
|---|---|---:|---|
| `planner` | `.codex/agents/planner.toml` | `gpt-5.5` | Architecture, decomposition, sequencing, risk analysis |
| `coding_worker` | `.codex/agents/coding_worker.toml` | `gpt-5.3-codex` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `.codex/agents/fast_coding_worker.toml` | `gpt-5.3-codex-spark` | Small localized edits and quick fixes |
| `helper_worker` | `.codex/agents/helper_worker.toml` | `gpt-5.4-mini` | Read-only lookup, repo reconnaissance, docs/API lookup |
| `forensic_analyst` | `.codex/agents/forensic_analyst.toml` | `gpt-5.5` | Deep root-cause investigation, intermittent and cross-system failures, forensic reports |
| `doc_reviewer` | `.codex/agents/doc_reviewer.toml` | `gpt-5.4-mini` | Documentation correctness, stale docs, doc drift |
| `reviewer` | `.codex/agents/reviewer.toml` | `gpt-5.5` | Standard code review |
| `qa_engineer` | `.codex/agents/qa_engineer.toml` | `gpt-5.5` | Exploratory QA: exercises changes end-to-end, probes regressions, performance, and rough edges |

Plan, implement, and review:

```text
Use agenticons.

Build tenant-level API keys with creation, rotation, revocation, and audit logging.
Spawn planner first to produce a short implementation plan and risk list.
Then spawn coding_worker to implement the accepted plan.
Then spawn reviewer for standard review.
```

Investigate a hard failure before editing:

```text
Use agenticons.

The checkout flow sometimes double-charges in staging and nobody knows why.
Do not edit code yet. Spawn forensic_analyst to find the root cause and report.
Save the accepted report to docs/forensics/double-charge.md.
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

Review the auth/session migration before release with reviewer.
This affects login, session invalidation, and production user access.
```

Exercise a release candidate:

```text
Use agenticons.

We are about to tag v2.4. Spawn qa_engineer to exercise the new rate limiter
end-to-end: build it, drive realistic traffic, probe restarts and config
reloads, and watch for latency regressions against v2.3. Report findings
before we tag.
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
| Shell | POSIX `sh` for installation |
| Remote install | `curl` or `wget` |
| Go | 1.22 or newer, for validation and tests |
| Installed skill path | Repo-local: `.agents/skills/agenticons/SKILL.md`; global: `~/.agents/skills/agenticons/SKILL.md` |
| Installed agent path | Repo-local: `.codex/agents/*.toml`; global: `~/.codex/agents/*.toml` |

## Contributing / Development

```bash
go run ./scripts/validate_package.go
go test ./...
go vet ./...
```

The validator parses every agent TOML file and checks that the package docs, the model tables, and the installer's agent list stay aligned with the agent specs. `.github/workflows/validate.yml` runs the validator, tests, and vet on every push and pull request.

## License

GNU General Public License v3.0. See `LICENSE`.
