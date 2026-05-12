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

Run:

```bash
./scripts/install.sh --target /path/to/your-repo
```

The installer copies `SKILL.md` to `.agents/skills/agenticons/SKILL.md` and `.codex/agents/*.toml` to `.codex/agents/` in the target repository.

## Can I install directly from GitHub?

Yes:

```bash
curl -fsSL https://raw.githubusercontent.com/fuentesjr/agenticons/main/scripts/install.sh | sh -s -- --target /path/to/your-repo
```

Remote install uses `curl` or `wget` to fetch package files from GitHub.

## What if files already exist?

The installer is idempotent when existing files match. If an existing file differs, the installer stops unless you pass `--force`.

```bash
./scripts/install.sh --target /path/to/your-repo --force
```

Preview writes with:

```bash
./scripts/install.sh --target /path/to/your-repo --dry-run
```

## Which agent should I ask for?

| Task | Agent |
|---|---|
| Plan a complex feature or migration | `planner` |
| Implement normal code changes | `coding_worker` |
| Make a small localized fix | `fast_coding_worker` |
| Investigate without editing | `helper_worker` |
| Check documentation drift | `doc_reviewer` |
| Review code for correctness and regressions | `reviewer` |
| Review auth, payments, data loss, or production-risk changes | `premium_reviewer` |

## When should I use `premium_reviewer`?

Use `premium_reviewer` only for rare high-stakes work: auth/session boundaries, payments, migrations, data loss risk, privacy risk, production incidents, or expensive irreversible decisions.

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

The validator checks that every agent TOML file is parseable, required fields are present, agent names match filenames, sandbox modes are supported, and `README.md` and `SKILL.md` mention every configured agent.

## Why is my new agent failing validation?

Check that the TOML file has all required fields, the `name` matches the filename, the sandbox mode is `read-only` or `workspace-write`, and both `README.md` and `SKILL.md` mention the new agent.
