// validate_package verifies that the agenticons package is internally
// consistent before it is published or installed into another repository.
//
// The checks intentionally focus on package contract drift:
//   - every custom agent TOML file is parseable and has the required fields
//   - each agent's declared name matches its filename
//   - sandbox modes and model reasoning efforts are supported values
//   - README.md, SKILL.md, and docs/design.md mention every configured agent
//   - README.md and docs/design.md document each agent with its configured model
//   - SKILL.md's exact dispatch list matches the files in .codex/agents
//   - scripts/install.sh's agent list matches the files in .codex/agents
//
// Keeping those rules in code makes documentation updates harder to forget
// when a role is added, renamed, or removed.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	agentsDir         = ".codex/agents"
	installScriptPath = "scripts/install.sh"
)

var (
	// docsToValidate are the public package docs that must stay aligned with
	// the concrete agent files.
	docsToValidate = []string{"README.md", "SKILL.md", "docs/design.md"}

	// modelTableDocs are the docs that publish the role/model contract as a
	// table. Each must keep an agent's name and configured model on one line
	// so model changes in TOML cannot leave the tables stale.
	modelTableDocs = []string{"README.md", "docs/design.md"}

	// requiredAgentFields mirrors the supported agent TOML contract. The
	// validator checks these with TOML metadata so missing fields are caught
	// even when Go would otherwise decode them as zero values.
	requiredAgentFields = []string{
		"name",
		"description",
		"model",
		"model_reasoning_effort",
		"sandbox_mode",
		"nickname_candidates",
		"developer_instructions",
	}

	// validSandboxModes lists the only sandbox values this package currently
	// ships. New modes should be added here deliberately when Codex supports
	// and the package documents them.
	validSandboxModes = map[string]struct{}{
		"read-only":       {},
		"workspace-write": {},
	}

	// validModelReasoningEfforts lists the reasoning effort values Codex
	// accepts, so a typo in an agent spec fails validation instead of
	// shipping.
	validModelReasoningEfforts = map[string]struct{}{
		"minimal": {},
		"low":     {},
		"medium":  {},
		"high":    {},
		"xhigh":   {},
	}

	// installAgentsLineRE matches the single-quoted agents list in
	// scripts/install.sh. The installer iterates this list, so it must stay
	// in sync with the agent files or installs silently drop roles.
	installAgentsLineRE = regexp.MustCompile(`(?m)^agents='([^']*)'$`)

	// deprecatedProjectIdentifiers catch incomplete project renames in the two
	// primary docs. Codex remains valid product wording, but the old package
	// name should no longer appear in public instructions.
	deprecatedProjectIdentifiers = []string{
		"codex-dispatch-package",
		"codex-dispatch",
		"Codex Dispatch Package",
		"Codex Dispatch",
	}

	// skillAgentLineRE matches the exact two-space bullet list in SKILL.md's
	// "Use the exact subagent names" section. This intentionally mirrors the
	// documented contract rather than scraping every inline code span.
	skillAgentLineRE = regexp.MustCompile(`(?m)^  - ` + "`" + `([a-z0-9_]+)` + "`" + `$`)
)

// agentSpec is the subset of each .codex/agents/*.toml file that agenticons
// promises to ship and validate.
type agentSpec struct {
	Name                  string   `toml:"name"`
	Description           string   `toml:"description"`
	Model                 string   `toml:"model"`
	ModelReasoningEffort  string   `toml:"model_reasoning_effort"`
	SandboxMode           string   `toml:"sandbox_mode"`
	NicknameCandidates    []string `toml:"nickname_candidates"`
	DeveloperInstructions string   `toml:"developer_instructions"`
}

// agentFile keeps parsed TOML together with its repository-relative path so
// downstream validation can report useful errors and check README references.
type agentFile struct {
	spec    agentSpec
	relPath string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	root, err := findRepositoryRoot()
	if err != nil {
		return err
	}

	agents, err := loadAgents(root)
	if err != nil {
		return err
	}

	if err := validateDocs(root, agents); err != nil {
		return err
	}

	if err := validateInstallScript(root, agents); err != nil {
		return err
	}

	fmt.Printf("Validated %d agent specs, %d docs, and the install script.\n", len(agents), len(docsToValidate))
	return nil
}

// findRepositoryRoot walks upward from the current working directory until it
// finds the files that define an agenticons source tree. This lets the command
// work from the repo root or from a nested directory such as scripts/.
func findRepositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	for {
		if fileExists(filepath.Join(dir, "README.md")) &&
			fileExists(filepath.Join(dir, "SKILL.md")) &&
			dirExists(filepath.Join(dir, agentsDir)) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find repository root containing README.md, SKILL.md, and .codex/agents")
		}
		dir = parent
	}
}

// loadAgents reads, decodes, and validates every top-level TOML file in
// .codex/agents. It returns agents sorted by name so all downstream checks and
// error messages are deterministic.
func loadAgents(root string) ([]agentFile, error) {
	pattern := filepath.Join(root, agentsDir, "*.toml")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob %s: %w", pattern, err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no agent TOML files found in %s", agentsDir)
	}

	sort.Strings(paths)
	seenNames := map[string]string{}
	agents := make([]agentFile, 0, len(paths))

	for _, path := range paths {
		agent, err := decodeAgent(root, path)
		if err != nil {
			return nil, err
		}

		if previous, ok := seenNames[agent.spec.Name]; ok {
			return nil, fmt.Errorf("duplicate agent name %q in %s and %s", agent.spec.Name, previous, agent.relPath)
		}
		seenNames[agent.spec.Name] = agent.relPath
		agents = append(agents, agent)
	}

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].spec.Name < agents[j].spec.Name
	})
	return agents, nil
}

// decodeAgent performs per-file validation while the TOML metadata is still
// available. That keeps "missing field" errors distinct from blank or invalid
// values, which makes package maintenance faster.
func decodeAgent(root, path string) (agentFile, error) {
	relPath := relativePath(root, path)

	var spec agentSpec
	metadata, err := toml.DecodeFile(path, &spec)
	if err != nil {
		return agentFile{}, fmt.Errorf("%s is invalid TOML: %w", relPath, err)
	}

	var missing []string
	for _, field := range requiredAgentFields {
		if !metadata.IsDefined(field) {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return agentFile{}, fmt.Errorf("%s missing fields: %s", relPath, strings.Join(missing, ", "))
	}

	if err := validateAgentSpec(spec, relPath, filepath.Base(path)); err != nil {
		return agentFile{}, err
	}

	return agentFile{spec: spec, relPath: relPath}, nil
}

// validateAgentSpec checks semantic rules that TOML decoding cannot express:
// filename/name alignment, non-empty required strings, allowed sandbox modes,
// and at least one usable nickname.
func validateAgentSpec(spec agentSpec, relPath, filename string) error {
	expectedName := strings.TrimSuffix(filename, filepath.Ext(filename))
	if spec.Name != expectedName {
		return fmt.Errorf("%s has name %q, expected %q", relPath, spec.Name, expectedName)
	}

	requiredStrings := map[string]string{
		"name":                   spec.Name,
		"description":            spec.Description,
		"model":                  spec.Model,
		"model_reasoning_effort": spec.ModelReasoningEffort,
		"sandbox_mode":           spec.SandboxMode,
		"developer_instructions": spec.DeveloperInstructions,
	}
	for field, value := range requiredStrings {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s field %q must not be blank", relPath, field)
		}
	}

	if _, ok := validSandboxModes[spec.SandboxMode]; !ok {
		return fmt.Errorf("%s has unsupported sandbox_mode %q", relPath, spec.SandboxMode)
	}

	if _, ok := validModelReasoningEfforts[spec.ModelReasoningEffort]; !ok {
		return fmt.Errorf("%s has unsupported model_reasoning_effort %q", relPath, spec.ModelReasoningEffort)
	}

	if len(spec.NicknameCandidates) == 0 {
		return fmt.Errorf("%s must define at least one nickname candidate", relPath)
	}
	for index, nickname := range spec.NicknameCandidates {
		if strings.TrimSpace(nickname) == "" {
			return fmt.Errorf("%s nickname_candidates[%d] must not be blank", relPath, index)
		}
	}

	return nil
}

// validateDocs keeps the package's human-facing documentation aligned with
// the machine-readable agent definitions.
func validateDocs(root string, agents []agentFile) error {
	docTexts := map[string]string{}
	for _, docPath := range docsToValidate {
		text, err := readDoc(root, docPath)
		if err != nil {
			return err
		}
		docTexts[docPath] = text

		if err := validateProjectRename(docPath, text); err != nil {
			return err
		}
		if err := validateDocMentionsAgents(docPath, text, agents); err != nil {
			return err
		}
	}

	if err := validateReadmeMentionsAgentFiles(docTexts["README.md"], agents); err != nil {
		return err
	}
	if err := validateSkillExactAgentList(docTexts["SKILL.md"], agents); err != nil {
		return err
	}

	for _, docPath := range modelTableDocs {
		if err := validateDocModelRows(docPath, docTexts[docPath], agents); err != nil {
			return err
		}
	}

	return nil
}

// validateDocModelRows enforces the model routing contract in the docs that
// publish it: each agent's backticked name must share a line with its
// configured model, so editing a TOML model without updating the tables fails.
func validateDocModelRows(docPath, text string, agents []agentFile) error {
	for _, agent := range agents {
		if !hasAgentModelLine(text, agent.spec) {
			return fmt.Errorf("%s does not document `%s` with model `%s` on one line", docPath, agent.spec.Name, agent.spec.Model)
		}
	}
	return nil
}

// hasAgentModelLine reports whether any single line mentions both the agent
// name and its model as backticked identifiers.
func hasAgentModelLine(text string, spec agentSpec) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "`"+spec.Name+"`") && strings.Contains(line, "`"+spec.Model+"`") {
			return true
		}
	}
	return false
}

// validateInstallScript keeps the installer's hardcoded agent list in sync
// with the agent files. Without this, a new role passes every doc check while
// installs silently skip it.
func validateInstallScript(root string, agents []agentFile) error {
	data, err := os.ReadFile(filepath.Join(root, installScriptPath))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing %s", installScriptPath)
		}
		return fmt.Errorf("read %s: %w", installScriptPath, err)
	}

	match := installAgentsLineRE.FindStringSubmatch(string(data))
	if match == nil {
		return fmt.Errorf("%s does not define an agents='...' list", installScriptPath)
	}

	listedNames := map[string]struct{}{}
	for _, name := range strings.Fields(match[1]) {
		listedNames[name] = struct{}{}
	}

	actualNames := map[string]struct{}{}
	for _, agent := range agents {
		actualNames[agent.spec.Name] = struct{}{}
	}

	missing := difference(actualNames, listedNames)
	if len(missing) > 0 {
		return fmt.Errorf("%s agent list missing: %s", installScriptPath, strings.Join(missing, ", "))
	}

	unknown := difference(listedNames, actualNames)
	if len(unknown) > 0 {
		return fmt.Errorf("%s agent list references missing agent files: %s", installScriptPath, strings.Join(unknown, ", "))
	}

	return nil
}

// validateProjectRename fails when the old package name survives in public
// docs. This catches partial renames before they reach installation examples.
func validateProjectRename(docPath, text string) error {
	for _, oldName := range deprecatedProjectIdentifiers {
		if strings.Contains(text, oldName) {
			return fmt.Errorf("%s still mentions deprecated project identifier %q", docPath, oldName)
		}
	}
	return nil
}

// readDoc loads one required documentation file and normalizes the common
// missing-file case into the same concise error style as the rest of the
// validator.
func readDoc(root, docPath string) (string, error) {
	fullPath := filepath.Join(root, docPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("missing doc: %s", docPath)
		}
		return "", fmt.Errorf("read %s: %w", docPath, err)
	}
	return string(data), nil
}

// validateDocMentionsAgents enforces that each public doc names every agent.
// The backtick requirement keeps the check targeted to documented identifiers
// instead of incidental prose.
func validateDocMentionsAgents(docPath, text string, agents []agentFile) error {
	for _, agent := range agents {
		if !strings.Contains(text, "`"+agent.spec.Name+"`") {
			return fmt.Errorf("%s does not mention `%s`", docPath, agent.spec.Name)
		}
	}
	return nil
}

// validateReadmeMentionsAgentFiles ensures installation docs list every
// concrete TOML file a user must copy into a target repository.
func validateReadmeMentionsAgentFiles(readme string, agents []agentFile) error {
	for _, agent := range agents {
		if !strings.Contains(readme, "`"+agent.relPath+"`") && !strings.Contains(readme, agent.relPath) {
			return fmt.Errorf("README.md does not mention %s", agent.relPath)
		}
	}
	return nil
}

// validateSkillExactAgentList compares SKILL.md's canonical bullet list with
// the actual agent files. This is stricter than validateDocMentionsAgents
// because the skill depends on an exact set of spawnable names.
func validateSkillExactAgentList(skill string, agents []agentFile) error {
	documentedNames := map[string]struct{}{}
	for _, match := range skillAgentLineRE.FindAllStringSubmatch(skill, -1) {
		documentedNames[match[1]] = struct{}{}
	}

	actualNames := map[string]struct{}{}
	for _, agent := range agents {
		actualNames[agent.spec.Name] = struct{}{}
	}

	missing := difference(actualNames, documentedNames)
	if len(missing) > 0 {
		return fmt.Errorf("SKILL.md exact subagent list missing: %s", strings.Join(missing, ", "))
	}

	unknown := difference(documentedNames, actualNames)
	if len(unknown) > 0 {
		return fmt.Errorf("SKILL.md exact subagent list references missing agent files: %s", strings.Join(unknown, ", "))
	}

	return nil
}

// difference returns sorted values that appear in left but not right, giving
// deterministic error messages for set comparisons.
func difference(left, right map[string]struct{}) []string {
	var values []string
	for value := range left {
		if _, ok := right[value]; !ok {
			values = append(values, value)
		}
	}
	sort.Strings(values)
	return values
}

// relativePath converts filesystem paths into slash-separated repository paths
// so errors and README checks match Markdown examples on every platform.
func relativePath(root, path string) string {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relPath)
}

// fileExists reports whether path is an existing regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// dirExists reports whether path is an existing directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
