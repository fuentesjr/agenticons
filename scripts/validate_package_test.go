package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunValidPackageFromNestedDirectory(t *testing.T) {
	root := newPackageTree(t)
	nestedDir := filepath.Join(root, "scripts")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("mkdir nested dir: %v", err)
	}
	chdir(t, nestedDir)

	if err := run(); err != nil {
		t.Fatalf("run() error = %v, want nil", err)
	}
}

func TestRunRejectsPackageDrift(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(t *testing.T, root string)
		wantErr string
	}{
		{
			name: "missing agent field",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, agentPath(root), `model = "gpt-5.4-mini"`+"\n", "")
			},
			wantErr: ".codex/agents/helper_worker.toml missing fields: model",
		},
		{
			name: "agent filename and name mismatch",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, agentPath(root), `name = "helper_worker"`, `name = "reviewer"`)
			},
			wantErr: `.codex/agents/helper_worker.toml has name "reviewer", expected "helper_worker"`,
		},
		{
			name: "unsupported sandbox mode",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, agentPath(root), `sandbox_mode = "read-only"`, `sandbox_mode = "danger-full-access"`)
			},
			wantErr: `.codex/agents/helper_worker.toml has unsupported sandbox_mode "danger-full-access"`,
		},
		{
			name: "blank nickname",
			mutate: func(t *testing.T, root string) {
				replaceInFile(t, agentPath(root), `nickname_candidates = ["Help Atlas"]`, `nickname_candidates = ["Help Atlas", ""]`)
			},
			wantErr: ".codex/agents/helper_worker.toml nickname_candidates[1] must not be blank",
		},
		{
			name: "missing readme agent file path",
			mutate: func(t *testing.T, root string) {
				writeFile(t, filepath.Join(root, "README.md"), "# Agenticons\n\n`helper_worker`\n")
			},
			wantErr: "README.md does not mention .codex/agents/helper_worker.toml",
		},
		{
			name: "missing skill exact agent list",
			mutate: func(t *testing.T, root string) {
				writeFile(t, filepath.Join(root, "SKILL.md"), "# Agenticons\n\nUse `helper_worker`.\n")
			},
			wantErr: "SKILL.md exact subagent list missing: helper_worker",
		},
		{
			name: "unknown skill exact agent list entry",
			mutate: func(t *testing.T, root string) {
				writeFile(t, filepath.Join(root, "SKILL.md"), "# Agenticons\n\n  - `helper_worker`\n  - `ghost_worker`\n")
			},
			wantErr: "SKILL.md exact subagent list references missing agent files: ghost_worker",
		},
		{
			name: "deprecated project identifier",
			mutate: func(t *testing.T, root string) {
				writeFile(t, filepath.Join(root, "README.md"), "# Agenticons\n\ncodex-dispatch\n`helper_worker`\n`.codex/agents/helper_worker.toml`\n")
			},
			wantErr: `README.md still mentions deprecated project identifier "codex-dispatch"`,
		},
		{
			name: "missing doc",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "SKILL.md")); err != nil {
					t.Fatalf("remove SKILL.md: %v", err)
				}
			},
			wantErr: "could not find repository root containing README.md, SKILL.md, and .codex/agents",
		},
		{
			name: "no agent files",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(agentPath(root)); err != nil {
					t.Fatalf("remove agent file: %v", err)
				}
			},
			wantErr: "no agent TOML files found in .codex/agents",
		},
		{
			name: "invalid toml",
			mutate: func(t *testing.T, root string) {
				writeFile(t, agentPath(root), "name =")
			},
			wantErr: ".codex/agents/helper_worker.toml is invalid TOML:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newPackageTree(t)
			tt.mutate(t, root)
			chdir(t, root)

			err := run()
			if err == nil {
				t.Fatal("run() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("run() error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func newPackageTree(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), validReadme())
	writeFile(t, filepath.Join(root, "SKILL.md"), validSkill())
	writeFile(t, agentPath(root), validAgentTOML())
	return root
}

func validReadme() string {
	return strings.Join([]string{
		"# Agenticons",
		"",
		"- `.codex/agents/helper_worker.toml`",
		"- `helper_worker`",
		"",
	}, "\n")
}

func validSkill() string {
	return strings.Join([]string{
		"# Agenticons",
		"",
		"- Use the exact subagent names from `.codex/agents/*.toml`:",
		"  - `helper_worker`",
		"",
	}, "\n")
}

func validAgentTOML() string {
	return strings.Join([]string{
		`name = "helper_worker"`,
		`description = "Read-only helper for investigation and evidence gathering."`,
		`model = "gpt-5.4-mini"`,
		`model_reasoning_effort = "medium"`,
		`sandbox_mode = "read-only"`,
		`nickname_candidates = ["Help Atlas"]`,
		`developer_instructions = """`,
		`Investigate and gather evidence. Do not edit files.`,
		`"""`,
		"",
	}, "\n")
}

func agentPath(root string) string {
	return filepath.Join(root, agentsDir, "helper_worker.toml")
}

func replaceInFile(t *testing.T, path, old, replacement string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	text := string(data)
	if !strings.Contains(text, old) {
		t.Fatalf("%s does not contain %q", path, old)
	}

	writeFile(t, path, strings.Replace(text, old, replacement, 1))
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory %s: %v", previousDir, err)
		}
	})
}
