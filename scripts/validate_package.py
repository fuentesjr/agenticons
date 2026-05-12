#!/usr/bin/env python3
"""Validate the codex-dispatch package contract."""

from __future__ import annotations

import re
import sys
import tomllib
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
AGENTS_DIR = ROOT / ".codex" / "agents"
DOCS = [ROOT / "README.md", ROOT / "SKILL.md"]
REQUIRED_FIELDS = {
    "name",
    "description",
    "model",
    "model_reasoning_effort",
    "sandbox_mode",
    "nickname_candidates",
    "developer_instructions",
}
VALID_SANDBOXES = {"read-only", "workspace-write"}


def fail(message: str) -> None:
    print(f"ERROR: {message}", file=sys.stderr)
    raise SystemExit(1)


def load_agents() -> dict[str, Path]:
    if not AGENTS_DIR.is_dir():
        fail(f"missing agent directory: {AGENTS_DIR.relative_to(ROOT)}")

    agents: dict[str, Path] = {}
    for path in sorted(AGENTS_DIR.glob("*.toml")):
        try:
            data = tomllib.loads(path.read_text())
        except tomllib.TOMLDecodeError as exc:
            fail(f"{path.relative_to(ROOT)} is invalid TOML: {exc}")

        missing = REQUIRED_FIELDS - data.keys()
        if missing:
            fail(f"{path.relative_to(ROOT)} missing fields: {', '.join(sorted(missing))}")

        name = data["name"]
        if name != path.stem:
            fail(f"{path.relative_to(ROOT)} has name {name!r}, expected {path.stem!r}")
        if name in agents:
            fail(f"duplicate agent name: {name}")
        if data["sandbox_mode"] not in VALID_SANDBOXES:
            fail(f"{path.relative_to(ROOT)} has unsupported sandbox_mode {data['sandbox_mode']!r}")
        if not isinstance(data["nickname_candidates"], list) or not data["nickname_candidates"]:
            fail(f"{path.relative_to(ROOT)} must define at least one nickname candidate")

        agents[name] = path

    if not agents:
        fail("no agent TOML files found")
    return agents


def validate_docs(agents: dict[str, Path]) -> None:
    for doc in DOCS:
        if not doc.is_file():
            fail(f"missing doc: {doc.relative_to(ROOT)}")

        text = doc.read_text()
        for name in agents:
            if f"`{name}`" not in text:
                fail(f"{doc.relative_to(ROOT)} does not mention `{name}`")

    readme = (ROOT / "README.md").read_text()
    for name, path in agents.items():
        rel = path.relative_to(ROOT).as_posix()
        if f"`{rel}`" not in readme and rel not in readme:
            fail(f"README.md does not mention {rel}")

    skill = (ROOT / "SKILL.md").read_text()
    documented_names = set(re.findall(r"  - `([a-z0-9_]+)`", skill))
    missing_from_skill_list = set(agents) - documented_names
    if missing_from_skill_list:
        fail(
            "SKILL.md exact subagent list missing: "
            + ", ".join(sorted(missing_from_skill_list))
        )

    unknown_in_skill_list = documented_names - set(agents)
    if unknown_in_skill_list:
        fail(
            "SKILL.md exact subagent list references missing agent files: "
            + ", ".join(sorted(unknown_in_skill_list))
        )


def main() -> None:
    agents = load_agents()
    validate_docs(agents)
    print(f"Validated {len(agents)} agent specs and {len(DOCS)} docs.")


if __name__ == "__main__":
    main()
