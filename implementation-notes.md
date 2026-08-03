# Implementation notes: systems_thinker role

Scope contract: Add a read-only `systems_thinker` role that applies Meadows-style systems thinking to recurring sociotechnical dynamics and integrates it into dispatch, installation, documentation, and behavioral verification.

## Decisions

- Full maintainer path: this change adds a public role and a new agent spec.
- Role contract: `gpt-5.6-sol`, `xhigh`, `read-only`; recurring dynamics and leverage points, not technical direction or incident forensics.
- Nicknames: `Systems Thinker Atlas`, `Systems Thinker Delta`, and `Systems Thinker Echo` keep the role explicit and follow the package convention.
- Use Meadows' concepts pragmatically. Require evidence, a defensible system map, reversible intervention experiments, and explicit learning loops.
- Detect semantic drift with a documented smoke rubric and a blind forward test before release.

## Order

1. Add the role spec and capture the package validator's expected registration failure.
2. Align dispatcher, docs, and installer.
3. Run static checks, then add the approved global per-file symlink.
4. Forward-test the role, run fresh review, resolve findings, and re-run checks.

## Constraints

- Preserve unrelated untracked `.gitattributes`, `.gitignore`, and `.trk/` content.
- Do not alter existing role models or the global agent directory topology.
- Leave the unrelated missing global `advisor` symlink unchanged.
- Do not commit, push, publish, or open a pull request without separate approval.

## Verification

- Baseline before edits: package validator, Go tests, Go vet, and Skill Creator validation passed for the ten-role package.
- Red proof after adding the role spec: package validation failed because `README.md` did not mention `systems_thinker`.
- `go run ./scripts/validate_package.go` — `Validated 11 agent specs, 4 docs, and the install script.`
- `go test ./...` — `ok agenticons/scripts (cached)`; `go vet ./...` exited successfully with no output.
- Skill Creator `quick_validate.py /Users/sal/Projects/skills/agenticons` — `Skill is valid!`
- `scripts/install.sh --target . --dry-run` listed `.codex/agents/systems_thinker.toml` as unchanged.
- `git diff --check` and documentation-link checks passed.
- `readlink /Users/sal/.codex/agents/systems_thinker.toml` returned `/Users/sal/Projects/skills/agenticons/.codex/agents/systems_thinker.toml`, and the target resolved.
- The current runtime could not spawn the newly registered custom type because role discovery is session-scoped. A fresh-context `gpt-5.6-sol` fallback read only the role spec and raw evidence packet; its report passed every documented positive-case rubric item.
- Fresh standard review found no role, installer, permissions, or public-contract blocker. Its stale-notes finding is fixed by this section.

## Remaining verification

- Start a new Codex session and spawn the registered `systems_thinker` role once to confirm runtime discovery.
- Use an isolated-event negative control for that first spawn. Passing means the role declines to manufacture a systems model and names the missing behavior-over-time evidence.
