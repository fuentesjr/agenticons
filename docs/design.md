# Agenticons Design

## Goal

Agenticons provides a small, explicit delegation layer for Codex. It gives users named subagents for technical advising, systems thinking, planning, implementation, review, documentation review, investigation, QA verification, and observability engineering.

The fixed roster keeps routing explicit without turning the parent agent into a workflow engine.

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
      advisor.toml
      systems_thinker.toml
      planner.toml
      coding_worker.toml
      fast_coding_worker.toml
      helper_worker.toml
      forensic_analyst.toml
      doc_reviewer.toml
      reviewer.toml
      qa_engineer.toml
      edge_case_analyst.toml
      observability_engineer.toml
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
| `advisor` | `read-only` | `gpt-5.6-sol` | Principal-level technical direction, tradeoffs, reversibility, and decision consistency |
| `systems_thinker` | `read-only` | `gpt-5.6-sol` | Recurring sociotechnical dynamics, feedback loops, leverage points, and intervention learning loops |
| `planner` | `read-only` | `gpt-5.6-sol` | Implementation decomposition, sequencing, risk analysis, and verification |
| `coding_worker` | `workspace-write` | `gpt-5.6-terra` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `workspace-write` | `gpt-5.6-luna` | Small localized edits and quick fixes |
| `helper_worker` | `read-only` | `gpt-5.6-luna` | Quick lookup, repo reconnaissance, evidence gathering |
| `forensic_analyst` | `read-only` | `gpt-5.6-sol` | Deep root-cause investigation, intermittent and cross-system failures, forensic reports |
| `doc_reviewer` | `read-only` | `gpt-5.6-luna` | Doc accuracy/drift, prune obsolete low-value docs, flag AI-confusing content; saveable report |
| `reviewer` | `read-only` | `gpt-5.6-sol` | Standard correctness, security, maintainability, regression review |
| `qa_engineer` | `workspace-write` | `gpt-5.6-terra` | Exploratory QA verification: exercises changes end-to-end, probes regressions, performance, and user-facing rough edges |
| `edge_case_analyst` | `read-only` | `gpt-5.6-sol` | Edge-case and coverage-gap discovery: finds unconsidered cases and specifies expected behavior and test cases |
| `observability_engineer` | `workspace-write` | `gpt-5.6-terra` | Instrumentation and wide-event review, OpenTelemetry strategy, SLO and burn-alert design, telemetry sampling and pipeline cost; scratch artifacts only, production edits routed to workers |

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

- Advise before planning: `advisor` challenges a high-leverage technical decision; the parent accepts or rejects the direction before invoking `planner`
- Diagnose recurring dynamics: `systems_thinker` models the structure behind persistent behavior; the parent accepts a leverage point before routing further work
- Plan then implement: `planner` -> `coding_worker` -> `reviewer`
- Fast fix: `fast_coding_worker`, with `reviewer` when behavior or public API changes
- Investigation before editing: `helper_worker` -> `coding_worker` or `fast_coding_worker`
- Deep root-cause investigation: `forensic_analyst` -> `coding_worker` once a cause is confirmed; the parent saves the accepted report to a file when the user requests it
- Documentation drift review: `doc_reviewer` returns a saveable report covering drift, obsolete docs to remove, and critical AI-confusing content; the parent routes confirmed updates and removals to `fast_coding_worker` or `coding_worker`, and pairs with `reviewer` when security or public API contract docs need code review as well
- High-stakes or security-sensitive review: `reviewer`
- Exploratory QA verification: `qa_engineer` after a feature lands or before a release; it exercises the change rather than reading it, and the parent routes confirmed findings to `coding_worker` or `fast_coding_worker`
- Edge-case and coverage analysis: `edge_case_analyst` returns a report of uncovered cases with proposed specs and concrete test cases; the parent saves the report when the user requests it and routes confirmed cases to `coding_worker` or `fast_coding_worker`
- Observability design and review: `observability_engineer` returns a saveable report on instrumentation, SLO/alerting design, and telemetry cost; it may create scratch scripts or configs to demonstrate a recommendation but must not modify existing source, and the parent routes confirmed instrumentation edits to `coding_worker` or `fast_coding_worker`

Parallel writable work should use disjoint ownership. Parallel review work should use distinct review angles such as correctness, security, and regression risk.

## Validation

`scripts/validate_package.go` protects the package contract before publishing or installation. It checks:

- every top-level `.codex/agents/*.toml` file is valid TOML
- required agent fields are present and non-blank
- agent names match filenames and are unique
- sandbox modes and model reasoning efforts are supported values
- nickname candidates are non-empty
- `README.md`, `SKILL.md`, `docs/design.md`, and `docs/faq.md` mention every configured agent
- `README.md` lists every agent TOML file path
- `README.md` and `docs/design.md` document each agent with its configured model on one line
- `docs/design.md` Sandbox column matches each agent's `sandbox_mode`
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

## Systems Thinker Behavioral Check

Run this check after changing the `systems_thinker` prompt or model. Complete it before release so semantic drift is detected within the same change.

1. Start a fresh Codex session from an empty temporary directory.
2. Provide a raw, dated evidence packet showing a managed roster diverging from its source over time.
3. Include roster snapshots, update-mechanism excerpts, and observed file types.
4. Do not include a suspected cause, proposed fix, this rubric, or the Agenticons repository.
5. Ask Agenticons to spawn `systems_thinker` for a read-only analysis of the recurring divergence.
6. Score the report after the subagent finishes.

The report passes when it:

- states whether the evidence supports a systems analysis
- separates observations, inferences, assumptions, and unknowns
- models only supported stocks, flows, feedback loops, constraints, and delays
- uses relevant leverage points without listing Meadows' full taxonomy by default
- proposes a reversible intervention with indicators, detection delay, review threshold, and adaptation path
- identifies the correct sibling role for technical decisions or unresolved failure causes

The report fails when it invents human motives, substitutes generic systems vocabulary for a causal model, or jumps directly to a one-off fix.

## Observability Engineer Behavioral Check

Run this check after changing the `observability_engineer` prompt or model. Complete it before release so semantic drift is detected within the same change.

Run two scenarios, each in its own fresh Codex session:

- Brownfield: a small service repository with partial, metrics-heavy instrumentation and a handful of static threshold alerts.
- Greenfield: a small service repository with no instrumentation at all.

For each scenario:

1. Ask Agenticons to spawn `observability_engineer` to review instrumentation and alerting before a release (brownfield) or to propose an instrumentation strategy (greenfield).
2. Do not include this rubric or the Agenticons repository in the session.
3. Score the report after the subagent finishes.

The report passes when it:

- assesses event width and trace coverage against OpenTelemetry conventions rather than proposing more dashboards or more metrics
- recommends unified wide-event storage for code the team owns rather than building out parallel metrics, logging, and tracing silos
- proposes SLO-based burn alerts tied to user-facing symptoms and names specific threshold alerts to delete (brownfield) or establishes SLOs before any threshold alerts exist (greenfield)
- addresses telemetry volume and sampling with a concrete strategy, not "collect everything"
- separates evidence-backed findings from assumptions and names what evidence would change the conclusion
- lists every scratch artifact it created and routes production instrumentation edits to the parent for a worker
- flags governance questions (build versus buy, vendor selection, budget ownership) as `advisor` territory instead of answering them

The report fails when it edits existing source, recommends tooling changes without evidence, treats three-pillars siloing as the target architecture, or produces generic monitoring advice unmoored from the repository's actual telemetry. In the greenfield scenario, the report additionally fails when it produces a boilerplate OpenTelemetry plan not grounded in the repository's frameworks and request paths, proposes instrumenting everything with no volume or sampling consideration, or selects a vendor or storage backend itself instead of flagging that choice for `advisor`.

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
