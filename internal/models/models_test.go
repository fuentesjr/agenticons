package models

import "testing"

const catalogJSON = `{
  "models": [
    {"slug": "gpt-6-sol", "visibility": "list", "supported_reasoning_levels": [{"effort": "high"}]},
    {"slug": "gpt-6.1-sol", "visibility": "list", "supported_reasoning_levels": [{"effort": "high"}]},
    {"slug": "gpt-6.2-sol", "visibility": "hide", "supported_reasoning_levels": [{"effort": "high"}]},
    {"slug": "gpt-7-sol", "visibility": "list", "supported_reasoning_levels": [{"effort": "high"}]},
    {"slug": "gpt-6-luna", "visibility": "list", "supported_reasoning_levels": [{"effort": "low"}]},
    {"slug": "gpt-6.1-luna", "visibility": "list", "supported_reasoning_levels": [{"effort": "high"}]},
    {"slug": "gpt-5.6-terra", "visibility": "list", "supported_reasoning_levels": [{"effort": "medium"}]},
    {"slug": "gpt-5.10-terra", "visibility": "list", "supported_reasoning_levels": [{"effort": "medium"}]}
  ]
}`

func TestLatest(t *testing.T) {
	catalog, err := ParseCatalog([]byte(catalogJSON))
	if err != nil {
		t.Fatalf("ParseCatalog() error = %v", err)
	}

	tests := []struct {
		name, slug, effort, want string
	}{
		{"point release in same family", "gpt-6-sol", "high", "gpt-6.1-sol"},
		{"hidden release is skipped", "gpt-6.1-sol", "high", "gpt-6.1-sol"},
		{"release lacking the effort is skipped", "gpt-6-luna", "low", "gpt-6-luna"},
		{"other variants and major versions are not upgrades", "gpt-6-astra", "high", "gpt-6-astra"},
		{"point versions compare numerically", "gpt-5.6-terra", "medium", "gpt-5.10-terra"},
		{"unrecognized slug is unchanged", "codex-auto-review", "high", "codex-auto-review"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := catalog.Latest(tt.slug, tt.effort); got != tt.want {
				t.Errorf("Latest(%q, %q) = %q, want %q", tt.slug, tt.effort, got, tt.want)
			}
		})
	}
}
