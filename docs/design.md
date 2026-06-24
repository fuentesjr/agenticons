# Agenticons Design

## Goal

Agenticons provides a small, explicit delegation layer for Codex. It gives users a fixed, named set of subagents for planning, implementation, review, documentation review, investigation, and QA verification without turning the parent agent into a heavy workflow engine.

## Non-Goals

- Agenticons does not dispatch unless the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.
- Agenticons does not replace the parent agent's judgment. The parent still owns orchestration, consolidation, and final response quality.
- Agenticons does not define project-specific engineering policy. Repository instructions and user requests remain authoritative.
- Agenticons does not provide a full reference manual for every Codex feature.

## Package Layout

```text
agenticons/
  SKILL.md
  README.md
  LICENSE
  go.mod
  go.sum
  docs/
    design.md
    faq.md
  .codex/
    agents/
      planner.toml
      coding_worker.toml
      fast_coding_worker.toml
      helper_worker.toml
      forensic_analyst.toml
      doc_reviewer.toml
      reviewer.toml
      qa_engineer.toml
  .github/
    workflows/
      validate.yml
  scripts/
    install.sh
    validate_package.go
    validate_package_test.go
```

The package installs into a target repository or into the current user's global Codex locations with `scripts/install.sh`. Repo-local installation places `SKILL.md` under `.agents/skills/agenticons/` and agent specs under `.codex/agents/`. Global installation places them under `~/.agents/skills/agenticons/` and `~/.codex/agents/`.

## Dispatch Contract

Dispatch is opt-in. The skill should activate when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.

Escape hatches take precedence. If the user says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use agenticons`, the parent agent handles the task directly.

Model routing is part of the package contract. The parent agent must use only the model identifiers configured in `.codex/agents/*.toml` and documented in the Agent Roles table; it must not override a role to an unlisted model or provider at dispatch time. Model changes should happen by updating the relevant agent spec and docs.

Fixing one model per role also keeps each subagent cache-coherent by construction: a role never switches models mid-task, so its prompt-prefix cache is never invalidated by a routing change. Dynamic auto-routers add cache-aware machinery to recover this property; a fixed roster has it for free.

## Parent Orchestration Contract

The parent agent is the orchestrator and DRA (Directly Responsible Agent). DRA means the parent remains accountable for the project outcome rather than handing accountability to subagents. It remains responsible for:

- selecting the right subagent
- assigning concrete scope and constraints
- sequencing work and deciding when to stop or continue
- assigning disjoint ownership for parallel writable work
- resolving conflicts between subagent outputs
- verifying results and deciding which findings or patches to accept
- treating subagent output as advisory until the parent accepts it
- consolidating results, conflicts, verification, and remaining risks into the final response

Subagents must not delegate or route. They return findings/results to the parent orchestrator, who remains accountable for the project outcome.

## Agent Roles

| Agent | Sandbox | Model | Responsibility |
|---|---|---:|---|
| `planner` | `read-only` | `gpt-5.5` | Architecture, decomposition, sequencing, risk analysis |
| `coding_worker` | `workspace-write` | `gpt-5.3-codex` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `workspace-write` | `gpt-5.3-codex-spark` | Small localized edits and quick fixes |
| `helper_worker` | `read-only` | `gpt-5.4-mini` | Quick lookup, repo reconnaissance, evidence gathering |
| `forensic_analyst` | `read-only` | `gpt-5.5` | Deep root-cause investigation, intermittent and cross-system failures, forensic reports |
| `doc_reviewer` | `read-only` | `gpt-5.4-mini` | Documentation correctness and drift review |
| `reviewer` | `read-only` | `gpt-5.5` | Standard correctness, security, maintainability, regression review |
| `qa_engineer` | `workspace-write` | `gpt-5.5` | Exploratory QA verification: exercises changes end-to-end, probes regressions, performance, and user-facing rough edges |

## Agent Spec Contract

Each `.codex/agents/*.toml` file must define:

| Field | Purpose |
|---|---|
| `name` | Spawnable agent identifier. Must match the filename without `.toml`. |
| `description` | Short role summary. |
| `model` | Model assigned to the role. |
| `model_reasoning_effort` | Reasoning effort assigned to the role. |
| `sandbox_mode` | `read-only` or `workspace-write`. |
| `nickname_candidates` | Human-readable nicknames. |
| `developer_instructions` | Role-specific behavior and output contract. |

## Orchestration Model

Agenticons keeps orchestration shallow. The parent agent delegates bounded subtasks, receives findings/results back from subagents, and then synthesizes the result.

Common patterns:

- Plan then implement: `planner` -> `coding_worker` -> `reviewer`
- Fast fix: `fast_coding_worker`, with `reviewer` when behavior or public API changes
- Investigation before editing: `helper_worker` -> `coding_worker` or `fast_coding_worker`
- Deep root-cause investigation: `forensic_analyst` -> `coding_worker` once a cause is confirmed; the parent saves the accepted report to a file when the user requests it
- Documentation drift review: `doc_reviewer`
- High-stakes or security-sensitive review: `reviewer`
- Exploratory QA verification: `qa_engineer` after a feature lands or before a release; it exercises the change rather than reading it, and the parent routes confirmed findings to `coding_worker` or `fast_coding_worker`

Parallel writable work should use disjoint ownership. Parallel review work should use distinct review angles such as correctness, security, and regression risk.

## Validation

`scripts/validate_package.go` protects the package contract before publishing or installation. It checks:

- every top-level `.codex/agents/*.toml` file is valid TOML
- required agent fields are present and non-blank
- agent names match filenames and are unique
- sandbox modes and model reasoning efforts are supported values
- nickname candidates are non-empty
- `README.md`, `SKILL.md`, and `docs/design.md` mention every configured agent
- `README.md` lists every agent TOML file path
- `README.md` and `docs/design.md` document each agent with its configured model on one line
- `SKILL.md`'s exact dispatch list matches the agent files
- `scripts/install.sh`'s agent list matches the agent files
- deprecated project identifiers do not remain in primary docs

Run validation and tests with:

```bash
go run ./scripts/validate_package.go
go test ./...
go vet ./...
```

`.github/workflows/validate.yml` runs all three commands on every push and pull request.

## Installation Script

`scripts/install.sh` is the primary distribution path. It supports:

- `--target <repo>` to choose the repository to install into
- `--global` to install for the current user under `~/.agents` and `~/.codex`
- `--dry-run` to preview writes
- `--force` to overwrite differing files
- `--ref <git-ref>` for remote installs from a specific Git ref

The script works from a local checkout. It also works through a raw GitHub pipe, where it downloads `SKILL.md` and `.codex/agents/*.toml` from the selected ref.

## Change Policy

When adding, removing, or renaming an agent, update the TOML file, `SKILL.md`, `README.md`, the agent list in `scripts/install.sh`, and any relevant docs in the same change. Run the validator and tests before publishing.

Model, sandbox, and role changes should be deliberate because they affect the delegation contract users rely on.
