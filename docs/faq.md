# Agenticons FAQ

## What is Agenticons?

Agenticons is a Codex skill package that routes explicit delegation requests to named custom subagents. It includes one skill file, `SKILL.md`, plus custom agent specs in `.codex/agents/*.toml`.

## Does Agenticons spawn subagents automatically?

No. Agenticons dispatches only when the user explicitly asks for agenticons, subagents, delegation, parallel execution, or model-tier routing.

## How do I opt out after mentioning Agenticons?

Use one of these phrases in the request:

```text
no subagents
do not use subagents
handle locally
do this yourself
do not use agenticons
```

## Where do the files go?

For repo-local installation, run:

```bash
./scripts/install.sh --target /path/to/your-repo
```

The installer copies `SKILL.md` to `.agents/skills/agenticons/SKILL.md` and `.codex/agents/*.toml` to `.codex/agents/` in the target repository.

For global installation, run:

```bash
./scripts/install.sh --global
```

The installer copies `SKILL.md` to `~/.agents/skills/agenticons/SKILL.md` and `.codex/agents/*.toml` to `~/.codex/agents/` for the current user.

## Can I install directly from GitHub?

Yes:

```bash
curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --target /path/to/your-repo
curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --global
```

Remote install uses `curl` or `wget` to fetch package files from GitHub.

## What if files already exist?

The installer is idempotent when existing files match. If an existing file differs, the installer stops unless you pass `--force`.

```bash
./scripts/install.sh --target /path/to/your-repo --force
./scripts/install.sh --global --force
```

Preview writes with:

```bash
./scripts/install.sh --target /path/to/your-repo --dry-run
./scripts/install.sh --global --dry-run
```

## Which agent should I ask for?

| Task | Agent |
|---|---|
| Challenge a high-leverage technical direction before planning | `advisor` |
| Explain recurring behavior and identify system leverage points | `systems_thinker` |
| Plan a complex feature or migration after the direction is accepted | `planner` |
| Implement normal code changes | `coding_worker` |
| Make a small localized fix | `fast_coding_worker` |
| Quick lookup or reconnaissance without editing | `helper_worker` |
| Find the root cause of a hard or intermittent failure | `forensic_analyst` |
| Audit docs for drift/accuracy, obsolete low-value material, and AI-confusing content (report only) | `doc_reviewer` |
| Apply confirmed documentation fixes or removals | `fast_coding_worker` or `coding_worker` |
| Review code for correctness and regressions | `reviewer` |
| Review auth, payments, data loss, or production-risk changes (one change) | `reviewer` |
| Exercise a feature or release end-to-end before shipping | `qa_engineer` |
| Find edge cases or missing test coverage nobody considered | `edge_case_analyst` |
| Review instrumentation, design SLOs and burn alerts, or control telemetry cost | `observability_engineer` |
| Audit the security posture of the whole system: threat model, authn/authz, secrets, supply chain, CI/CD (report only) | `security_auditor` |

## When should I use `advisor` instead of `planner`?

Use `advisor` to decide whether a technical direction is sound. Use it for
multiple credible approaches, cross-system boundaries, expensive reversals,
conflicting recommendations, or trajectory drift.

Use `planner` after the parent accepts the direction. The planner converts that
direction into an implementation sequence, affected-file map, risk list, and
verification strategy.

Do not run `advisor` as a mandatory gate for routine work.

## When should I use `systems_thinker`?

Use `systems_thinker` when a behavior repeats over time and local fixes do not hold. It examines feedback loops, stocks and flows, delays, constraints, information flows, and relevant leverage points.

Use `advisor` when you need a technical direction. Use `forensic_analyst` when you need evidence for a hard failure's cause.

Do not use `systems_thinker` as a mandatory gate or for an isolated event without behavior-over-time evidence.

## When should I use `observability_engineer`?

Use `observability_engineer` for instrumentation and wide-event review, OpenTelemetry strategy, SLO and burn-alert design, alert-fatigue cleanup, telemetry sampling and pipeline cost, and observability for CI/CD, frontend/mobile, or LLM applications. Its practice follows Observability Engineering, 2nd edition. It may create scratch scripts or configs to demonstrate a recommendation, but production instrumentation edits go to `coding_worker` or `fast_coding_worker`.

Use `forensic_analyst` to prove the cause of one hard failure; `observability_engineer` designs the telemetry that makes such debugging possible. Use `qa_engineer` to exercise a feature end-to-end; `observability_engineer` judges what its telemetry reveals. Observability governance questions — business case, build versus buy, vendor selection — belong to `advisor`.

## When should I use `security_auditor`?

Use `security_auditor` for a defensive audit of the system as a whole: threat model, authentication and session handling, authorization and tenant isolation, input handling, secrets, data protection, dependency and supply-chain exposure, and CI/CD hardening. It audits against OWASP ASVS 5.0, confirms findings with minimal evidence, and returns a severity-ranked report. It is read-only; fixes go to `coding_worker` or `fast_coding_worker`.

Use `reviewer` for the security of one change. Use `edge_case_analyst` to enumerate hostile inputs as test cases; `security_auditor` models the attacker and the trust boundaries those inputs cross. Use `forensic_analyst` to investigate a suspected breach; `security_auditor` audits before the incident. Security governance questions — compliance programs, tooling and vendor selection — belong to `advisor`.

## Can I use a different model for an Agenticons role?

No. Use the model configured in the role's `.codex/agents/*.toml` file and shown in the package docs. Agenticons should not substitute unlisted models or providers at dispatch time; changing a role's model requires updating the agent spec and docs.

## Who orchestrates Agenticons subagents?

The parent agent is the orchestrator and DRA (Directly Responsible Agent). DRA means the parent remains accountable for the project outcome: it selects subagents, assigns scope, sequences work, resolves conflicts, verifies results, decides what subagent output to accept, and owns the final response. Subagents do not delegate or route; they return advisory findings/results to the parent.

## Can multiple workers edit files in parallel?

Yes, but each writable worker should receive disjoint file or feature ownership. The parent agent should coordinate the work and consolidate results.

## What should a subagent prompt include?

Include the concrete task, scope, constraints, ownership, and expected output. For writable work, tell the agent which files or feature area it owns and what it must avoid.

## How do I limit fan-out?

Set global agent limits in `.codex/config.toml` or `~/.codex/config.toml`:

```toml
[agents]
max_threads = 6
max_depth = 1
```

## How do I validate package changes?

Run:

```bash
go run ./scripts/validate_package.go
go test ./...
go vet ./...
```

## What does the validator check?

The validator checks that every agent TOML file is parseable, required fields are present, agent names match filenames, sandbox modes and reasoning efforts are supported, `README.md`, `SKILL.md`, `docs/design.md`, and `docs/faq.md` mention every configured agent, the model tables in `README.md` and `docs/design.md` match the TOML specs, the `docs/design.md` Sandbox column matches each agent's `sandbox_mode`, and the agent list in `scripts/install.sh` matches the agent files.

## Why is my new agent failing validation?

Check that the TOML file has all required fields, the `name` matches the filename, the sandbox mode is `read-only` or `workspace-write`, the reasoning effort is a supported value, `README.md`, `SKILL.md`, `docs/design.md`, and `docs/faq.md` mention the new agent (including its model in the README and design tables), and the agent is listed in `scripts/install.sh`.
