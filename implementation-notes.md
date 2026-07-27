# Implementation notes: doc_reviewer contract upgrade

## Decisions

- `doc_reviewer` stays on `gpt-5.6-luna` / `high` / `read-only`.
- Expanded charter: prune obsolete low-value docs; treat AI-confusing documentation as critical (agents load prose as instructions).
- Escalation is orchestration (pair with `reviewer` for security/API contracts), not model override — package forbids unlisted model substitution.
- Report shape mirrors forensic/edge-case/QA specialists so parents can save verbatim and route apply work.

## Scope

- `.codex/agents/doc_reviewer.toml` — description, boundaries, change-aware method, fixed report sections
- `SKILL.md` — docs rule of thumb, escalation note, numbered dispatch + apply path
- `README.md`, `docs/design.md`, `docs/faq.md` — wording + apply/save examples

## Verification

- `go run ./scripts/validate_package.go` — ok (9 agents, 4 docs, install script)
- `go test ./...` — ok
- `go vet ./...` — ok

---

# Implementation notes: principal advisor role

## Plan

- Add a read-only `advisor` agent using `gpt-5.6-sol` at `xhigh` effort.
- Keep its contract narrow: principal-level decision quality and technical
  trajectory, not implementation planning or code review.
- Update the skill roster, routing guidance, README, design reference, FAQ, and
  installer list.
- Prove package drift detection red before aligning those surfaces.
- Run package validation, Go tests, vet, and skill-lint. Defer live skill-tester
  evaluation until the active eval-infrastructure work completes.

## Scope boundaries

- Do not change existing model assignments. Sharpen only the planner boundary
  needed to distinguish accepted-direction planning from advisor judgment.
- Do not modify or stage unrelated untracked files under `.trk/`,
  `.gitattributes`, or `.gitignore`.
- Do not install the package. The user authorized committing and pushing the
  completed change.

## Verification

- Baseline `go test ./...` — ok.
- Red proof after adding the agent spec: package validator rejected the missing
  README registration.
- `go run ./scripts/validate_package.go` — ok (10 agents, 4 docs, install
  script).
- `go test ./...` — ok.
- `go vet ./...` — ok.
- `skill-lint` — 0 errors; 4 warnings: three package-layout orphans for
  repository-only scripts, plus the intentionally deferred missing eval suite.
- Live skill-tester certification — deferred by user direction until active
  eval-infrastructure work completes. A discarded partial run invoked
  Agenticons in 2/2 trials, but both were infrastructure-invalid when the grader
  subprocess failed; this is not certification evidence.
