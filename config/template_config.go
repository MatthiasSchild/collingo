package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Well-known template / export kinds (JSON template.kind).
const (
	TemplateKindVueI18n = "vue-i18n"
	TemplateKindFlutter = "flutter"
	TemplateKindJSON    = "json"
	TemplateKindI18next = "i18next"
)

// TemplateConfig is the workspace template object: kind plus optional export settings.
type TemplateConfig struct {
	Kind       string `json:"kind,omitempty"`
	Directory  string `json:"directory,omitempty"`
	OutputFile string `json:"outputFile,omitempty"`
	Formatted  bool   `json:"formatted,omitempty"`
}

func (t *TemplateConfig) HasKind() bool {
	return t != nil && t.Kind != ""
}

// EnvDescription returns lines for collingo env (no trailing newline).
func (t *TemplateConfig) EnvDescription() string {
	if t == nil || !t.HasKind() {
		return "(none)"
	}
	var b strings.Builder
	b.WriteString(t.Kind)
	if t.Directory != "" {
		fmt.Fprintf(&b, "\n    directory: %s", t.Directory)
	}
	if t.OutputFile != "" {
		fmt.Fprintf(&b, "\n    outputFile: %s", t.OutputFile)
	}
	if t.Formatted {
		b.WriteString("\n    formatted: true")
	}
	return b.String()
}

var (
	ErrNoTemplateForExport             = errors.New("no workspace template configured; run collingo init or set template in .collingo.json")
	ErrTemplateExportDirectoryRequired = errors.New("template.directory is required for this export kind; set it in .collingo.json or use collingo export <subcommand> -d …")
	ErrTemplateExportUnknownKind       = errors.New("unknown template kind in .collingo.json")
)

// ValidateForDefaultExport checks that workspace template can run bare collingo export.
func (t *TemplateConfig) ValidateForDefaultExport() error {
	if t == nil || !t.HasKind() {
		return ErrNoTemplateForExport
	}
	switch t.Kind {
	case TemplateKindFlutter, TemplateKindI18next:
		if strings.TrimSpace(t.Directory) == "" {
			return ErrTemplateExportDirectoryRequired
		}
	case TemplateKindVueI18n, TemplateKindJSON:
		// directory not used; outputFile may be empty (stdout)
	default:
		return fmt.Errorf("%w: %q", ErrTemplateExportUnknownKind, t.Kind)
	}
	return nil
}

// ResolvePathUnderWorkspace joins a path relative to workspace root. userPath must be relative
// and must not escape the workspace.
func ResolvePathUnderWorkspace(workspaceRoot, userPath string) (string, error) {
	userPath = strings.TrimSpace(userPath)
	if userPath == "" {
		return "", nil
	}
	if filepath.IsAbs(userPath) {
		return "", fmt.Errorf("path must be relative to the workspace root: %q", userPath)
	}
	clean := filepath.Clean(userPath)
	full := filepath.Join(workspaceRoot, clean)
	rootClean := filepath.Clean(workspaceRoot)
	rel, err := filepath.Rel(rootClean, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes workspace: %q", userPath)
	}
	return full, nil
}
