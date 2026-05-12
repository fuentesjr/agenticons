# Agenticons Design

## Goal

Agenticons provides a small, explicit delegation layer for Codex. It gives users predictable named subagents for planning, implementation, review, documentation review, and investigation without turning the parent agent into a heavy workflow engine.

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
  docs/
    design.md
    faq.md
  .codex/
    agents/
      planner.toml
      coding_worker.toml
      fast_coding_worker.toml
      helper_worker.toml
      doc_reviewer.toml
      reviewer.toml
      premium_reviewer.toml
  scripts/
    install.sh
    validate_package.go
```

The package installs into a target repository or into the current user's global Codex locations with `scripts/install.sh`. Repo-local installation places `SKILL.md` under `.agents/skills/agenticons/` and agent specs under `.codex/agents/`. Global installation places them under `~/.agents/skills/agenticons/` and `~/.codex/agents/`.

## Dispatch Contract

Dispatch is opt-in. The skill should activate when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.

Escape hatches take precedence. If the user says `no subagents`, `do not use subagents`, `handle locally`, `do this yourself`, or `do not use agenticons`, the parent agent handles the task directly.

Premium review is explicit opt-in. The parent agent must not spawn `premium_reviewer` unless the user asks for `premium_reviewer` or premium review in the current request. High-stakes context can justify recommending premium escalation, but it is not enough to spawn the expensive role. If the user asks for premium review only if needed, the parent should run `reviewer` first and have it report whether premium review is justified.

The parent agent remains responsible for:

- selecting the right subagent
- assigning concrete scope and constraints
- assigning disjoint ownership for parallel writable work
- avoiding recursive delegation unless requested
- consolidating results, conflicts, verification, and remaining risks

## Agent Roles

| Agent | Sandbox | Model | Responsibility |
|---|---|---:|---|
| `planner` | `read-only` | `gpt-5.5` | Architecture, decomposition, sequencing, risk analysis |
| `coding_worker` | `workspace-write` | `gpt-5.3-codex` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `workspace-write` | `gpt-5.3-codex-spark` | Small localized edits and quick fixes |
| `helper_worker` | `read-only` | `gpt-5.4-mini` | Investigation, lookup, repo reconnaissance, evidence gathering |
| `doc_reviewer` | `read-only` | `gpt-5.4-mini` | Documentation correctness and drift review |
| `reviewer` | `read-only` | `gpt-5.5` | Standard correctness, security, maintainability, regression review |
| `premium_reviewer` | `read-only` | `gpt-5.5-pro` | Explicitly requested rare high-stakes review for security, data, payments, migrations, and production risk |

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

Agenticons keeps orchestration shallow. The parent agent delegates bounded subtasks and then synthesizes the result.

Common patterns:

- Plan then implement: `planner` -> `coding_worker` -> `reviewer`
- Fast fix: `fast_coding_worker`, with `reviewer` when behavior or public API changes
- Investigation before editing: `helper_worker` -> `coding_worker` or `fast_coding_worker`
- Documentation drift review: `doc_reviewer`
- High-stakes review: `reviewer`, with a premium escalation recommendation when justified
- Explicit premium review: `premium_reviewer`

Parallel writable work should use disjoint ownership. Parallel review work should use distinct review angles such as correctness, security, and regression risk.

## Validation

`scripts/validate_package.go` protects the package contract before publishing or installation. It checks:

- every top-level `.codex/agents/*.toml` file is valid TOML
- required agent fields are present
- agent names match filenames
- sandbox modes are supported
- nickname candidates are non-empty
- `README.md` and `SKILL.md` mention every configured agent
- `SKILL.md`'s exact dispatch list matches the agent files
- deprecated project identifiers do not remain in primary docs

Run validation and tests with:

```bash
go run ./scripts/validate_package.go
go test ./...
go vet ./...
```

## Installation Script

`scripts/install.sh` is the primary distribution path. It supports:

- `--target <repo>` to choose the repository to install into
- `--global` to install for the current user under `~/.agents` and `~/.codex`
- `--dry-run` to preview writes
- `--force` to overwrite differing files
- `--ref <git-ref>` for remote installs from a specific Git ref

The script works from a local checkout. It also works through a raw GitHub pipe, where it downloads `SKILL.md` and `.codex/agents/*.toml` from the selected ref.

## Change Policy

When adding, removing, or renaming an agent, update the TOML file, `SKILL.md`, `README.md`, and any relevant docs in the same change. Run the validator and tests before publishing.

Model, sandbox, and role changes should be deliberate because they affect the delegation contract users rely on.
