# Agenticons

<img width="1672" height="941" alt="agenticons-readme" src="https://github.com/user-attachments/assets/733b5678-8ad6-491e-b521-2b65ec17f31f" />


*Pronounced uh-JEN-tih-conz — like Decepticons, not "agent icons."*

Agenticons is a Codex skill package for explicit, named subagent delegation. Its roles cover technical advising, systems thinking, planning, implementation, review, documentation, investigation, QA verification, observability engineering, and security auditing.

## Why This Exists

Codex can delegate work to subagents, but repeated manual routing is easy to make inconsistent. Agenticons packages a small routing skill and named agent specs.

Users request delegation once, then route work through fixed roles such as `advisor`, `systems_thinker`, `planner`, `coding_worker`, or `reviewer`.

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

Model routing is fixed by the installed Agenticons agent specs. Spawn each role with the model and reasoning effort listed in its `.codex/agents/*.toml` file and in the table below; do not substitute unlisted models or providers at dispatch time.

The roster prioritizes quality, then speed. It assigns Astra to architecture, systems analysis, planning, difficult investigations, review, edge-case analysis, and security audits. It assigns Sol to substantial implementation and verification, Terra to bounded evidence gathering, and Luna to small mechanical changes. Higher reasoning effort is reserved for work whose analytical scope needs it.

The parent agent is the orchestrator and DRA (Directly Responsible Agent). That means the parent remains accountable for the project outcome: subagents do not delegate or route; they return advisory findings/results to the parent, who owns sequencing, verification, conflict resolution, deciding what to accept, and the final response.

Standard review always routes to `reviewer`, including security-sensitive and high-stakes review requests. A security audit of the system as a whole, rather than of one change, routes to `security_auditor`.

| Subagent | Agent file | Model | Effort | Use |
|---|---|---:|---:|---|
| `advisor` | `.codex/agents/advisor.toml` | `gpt-6-astra` | `xhigh` | Principal-level technical direction, tradeoffs, reversibility, and decision consistency |
| `systems_thinker` | `.codex/agents/systems_thinker.toml` | `gpt-6-astra` | `xhigh` | Recurring sociotechnical dynamics, feedback loops, leverage points, and intervention learning loops |
| `planner` | `.codex/agents/planner.toml` | `gpt-6-astra` | `high` | Implementation decomposition, sequencing, risk analysis, and verification |
| `coding_worker` | `.codex/agents/coding_worker.toml` | `gpt-5.6-sol` | `high` | Normal implementation, bug fixes, refactors |
| `fast_coding_worker` | `.codex/agents/fast_coding_worker.toml` | `gpt-5.6-luna` | `low` | Small localized edits and quick fixes |
| `helper_worker` | `.codex/agents/helper_worker.toml` | `gpt-5.6-terra` | `medium` | Read-only lookup, repo reconnaissance, docs/API lookup |
| `forensic_analyst` | `.codex/agents/forensic_analyst.toml` | `gpt-6-astra` | `max` | Deep root-cause investigation, intermittent and cross-system failures, forensic reports |
| `doc_reviewer` | `.codex/agents/doc_reviewer.toml` | `gpt-5.6-sol` | `high` | Doc accuracy/drift, prune obsolete low-value docs, flag AI-confusing content |
| `reviewer` | `.codex/agents/reviewer.toml` | `gpt-6-astra` | `high` | Standard code review |
| `qa_engineer` | `.codex/agents/qa_engineer.toml` | `gpt-5.6-sol` | `high` | Exploratory QA: exercises changes end-to-end, probes regressions, performance, and rough edges |
| `edge_case_analyst` | `.codex/agents/edge_case_analyst.toml` | `gpt-6-astra` | `xhigh` | Find unconsidered edge cases and design gaps; specify expected behavior and concrete test cases |
| `observability_engineer` | `.codex/agents/observability_engineer.toml` | `gpt-5.6-sol` | `high` | Instrumentation and wide-event review, OpenTelemetry strategy, SLO and burn-alert design, telemetry sampling/pipeline cost, observability for CI/CD, frontend, and LLM apps |
| `security_auditor` | `.codex/agents/security_auditor.toml` | `gpt-6-astra` | `xhigh` | Defensive security posture audit: threat model, ASVS 5.0-grounded findings ranked by severity, dependency/supply-chain and CI/CD hardening; read-only, fixes routed to workers |

Challenge technical direction before planning:

```text
Use agenticons.

We need to choose between extending our monolith and splitting billing into a
service. Spawn advisor to assess the boundary, tradeoffs, reversibility, and
operational blast radius. Do not produce an implementation plan yet.
```

Analyze a recurring system behavior:

```text
Use agenticons.

Review queues clear after each release push, then grow again within two weeks.
Adding reviewers helps briefly, but lead time and rework keep returning to their
previous levels. Spawn systems_thinker to model the recurring dynamics, identify
supported leverage points, and propose measurable intervention experiments.
Do not edit files or produce an implementation plan.
```

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
configuration, API docs, agent instructions, and release notes against the
implementation. Flag anything that could confuse AI agents as critical. Recommend
removing obsolete low-value docs, not only rewriting them.
Save the accepted report to docs/reviews/doc-drift.md.
Then apply confirmed small doc fixes and removals with fast_coding_worker.
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

Find uncovered edge cases before writing code:

```text
Use agenticons.

We are adding coupon stacking to checkout. Spawn edge_case_analyst to find the
edge cases and missing acceptance criteria we have not considered (expired
coupons, conflicting discounts, rounding, currency, concurrency) and return a
report with a concrete test case for each. Do not edit code.
```

Review instrumentation and alerting before a launch:

```text
Use agenticons.

We are taking the payments service to GA next month. Spawn observability_engineer
to review its instrumentation and alerting: event width and trace coverage against
OpenTelemetry conventions, SLO and burn-alert design for checkout, threshold
alerts we should delete, and expected telemetry volume with a sampling strategy.
Return a report; route any instrumentation edits back to me for coding_worker.
```

Audit security posture before a release:

```text
Use agenticons.

We are opening the tenant API to external customers next quarter. Spawn
security_auditor for a pre-release audit of the whole service: threat model,
authentication and tenant isolation, secrets handling, dependency and CI/CD
exposure. Confirm findings with minimal evidence and rank them by severity.
Do not edit code; route fixes back to me for coding_worker.
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
