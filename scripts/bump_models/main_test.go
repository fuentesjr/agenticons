package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBumpsAgentAndRoleTables(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		".codex/agents/doc_reviewer.toml": "name = \"doc_reviewer\"\nmodel = \"gpt-6-sol\"\nmodel_reasoning_effort = \"high\"\n",
		".codex/agents/planner.toml":      "name = \"planner\"\nmodel = \"gpt-6-astra\"\nmodel_reasoning_effort = \"high\"\n",
		"README.md":                       "| `doc_reviewer` | `.codex/agents/doc_reviewer.toml` | `gpt-6-sol` | `high` | Docs |\n| `planner` | `.codex/agents/planner.toml` | `gpt-6-astra` | `high` | Plans |\n",
		"docs/design.md":                  "| `doc_reviewer` | `read-only` | `gpt-6-sol` | `high` | Docs |\n| `planner` | `read-only` | `gpt-6-astra` | `high` | Plans |\n",
		"catalog.json": `{"models": [
  {"slug": "gpt-6.1-sol", "visibility": "list", "supported_reasoning_levels": [{"effort": "high"}]}
]}`,
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	if err := run("catalog.json", false); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	for name, want := range map[string]string{
		".codex/agents/doc_reviewer.toml": `model = "gpt-6.1-sol"`,
		".codex/agents/planner.toml":      `model = "gpt-6-astra"`,
		"README.md":                       "| `doc_reviewer` | `.codex/agents/doc_reviewer.toml` | `gpt-6.1-sol` |",
		"docs/design.md":                  "| `doc_reviewer` | `read-only` | `gpt-6.1-sol` |",
	} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), want) {
			t.Errorf("%s = %q, want it to contain %q", name, data, want)
		}
		if name != ".codex/agents/planner.toml" && strings.Contains(string(data), "gpt-6-sol") {
			t.Errorf("%s still mentions gpt-6-sol", name)
		}
	}
}
