// bump_models moves each agent to the newest point release of its pinned
// model family (for example gpt-6-sol -> gpt-6.1-sol), using the model catalog
// reported by codex debug models. Run it from the repository root. It rewrites
// the agent TOML files and the role tables in README.md and docs/design.md;
// review and commit the diff, then run the package validator.
//
// Major versions and variant changes are never bumped; those stay deliberate
// decisions. A release that no longer accepts an agent's reasoning effort is
// skipped.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"agenticons/internal/models"

	"github.com/BurntSushi/toml"
)

const agentsDir = ".codex/agents"

// roleTableDocs keep one row per agent with its model on the same line.
var roleTableDocs = []string{"README.md", "docs/design.md"}

type agentSpec struct {
	Name                 string `toml:"name"`
	Model                string `toml:"model"`
	ModelReasoningEffort string `toml:"model_reasoning_effort"`
}

func main() {
	catalogPath := flag.String("catalog", "", "read the model catalog JSON from this file instead of codex debug models")
	dryRun := flag.Bool("n", false, "report available bumps without editing files")
	flag.Parse()

	if err := run(*catalogPath, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run(catalogPath string, dryRun bool) error {
	catalog, err := loadCatalog(catalogPath)
	if err != nil {
		return err
	}

	paths, err := filepath.Glob(filepath.Join(agentsDir, "*.toml"))
	if err != nil {
		return err
	}
	sort.Strings(paths)

	bumped := 0
	for _, path := range paths {
		var spec agentSpec
		if _, err := toml.DecodeFile(path, &spec); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		latest := catalog.Latest(spec.Model, spec.ModelReasoningEffort)
		if latest == spec.Model {
			continue
		}
		fmt.Printf("%s: %s -> %s\n", spec.Name, spec.Model, latest)
		bumped++
		if dryRun {
			continue
		}
		if err := bumpAgent(path, spec, latest); err != nil {
			return err
		}
	}

	switch {
	case bumped == 0:
		fmt.Println("All agents are on the newest point release of their model family.")
	case dryRun:
		fmt.Printf("%d agent(s) can be bumped; rerun without -n to apply.\n", bumped)
	default:
		fmt.Printf("Bumped %d agent(s). Review the diff, then run: go run ./scripts/validate_package.go\n", bumped)
	}
	return nil
}

func loadCatalog(path string) (models.Catalog, error) {
	if path == "" {
		return models.FetchCatalog()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return models.Catalog{}, err
	}
	return models.ParseCatalog(data)
}

// bumpAgent rewrites the model in the agent's TOML file and in its row of each
// role table. Edits are exact-line replacements so nothing else changes.
func bumpAgent(path string, spec agentSpec, latest string) error {
	oldLine := fmt.Sprintf("model = %q", spec.Model)
	newLine := fmt.Sprintf("model = %q", latest)
	if err := editLines(path, func(line string) string {
		if line == oldLine {
			return newLine
		}
		return line
	}); err != nil {
		return err
	}

	rowPrefix := fmt.Sprintf("| `%s` |", spec.Name)
	oldCell := fmt.Sprintf("| `%s` |", spec.Model)
	newCell := fmt.Sprintf("| `%s` |", latest)
	for _, doc := range roleTableDocs {
		if err := editLines(doc, func(line string) string {
			if strings.HasPrefix(line, rowPrefix) {
				return strings.Replace(line, oldCell, newCell, 1)
			}
			return line
		}); err != nil {
			return err
		}
	}
	return nil
}

func editLines(path string, edit func(string) string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	changed := false
	for i, line := range lines {
		if updated := edit(line); updated != line {
			lines[i] = updated
			changed = true
		}
	}
	if !changed {
		return fmt.Errorf("%s: no line to update", path)
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}
