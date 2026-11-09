package runbooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBetssonDemoAssets(t *testing.T) {
	root := repoRoot(t)
	required := []string{
		"AGENTS.md",
		".github/copilot-instructions.md",
		".github/instructions/go-money.instructions.md",
		".github/instructions/billing-tests.instructions.md",
		".github/prompts/explain-invoice.prompt.md",
		".github/agents/billing-reviewer.md",
		".github/skills/invoice-memo/SKILL.md",
		".claude/skills/dispute-note/SKILL.md",
		".vscode/mcp.json",
		".github/workflows/cursor-agent-demo.yml",
		"demo/screenshots/inv-1048-scale.png",
		"demo/migration-fallback/.cursor/rules/go-money.mdc",
		"demo/migration-fallback/.cursor/rules/suggested-credit-seam.mdc",
		"demo/migration-fallback/.cursor/rules/demo-reset.mdc",
		"demo/migration-fallback/.cursor/skills/explain-invoice/SKILL.md",
		"demo/migration-fallback/.cursor/agents/billing-reviewer.md",
		"demo/migration-fallback/.cursor/mcp.json",
		"demo/migration-fallback/.cursor/hooks.json",
		"demo/migration-fallback/.cursor/hooks/format-go.sh",
		"demo/migration-fallback/.cursor/hooks/guard-secret-print.sh",
		"scripts/apply-migration-fallback.sh",
		".cursor/rules/ledgerly.mdc",
	}
	for _, rel := range required {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Errorf("missing %s: %v", rel, err)
		}
	}

	ledgerly, err := os.ReadFile(filepath.Join(root, ".cursor/rules/ledgerly.mdc"))
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{"alwaysApply: true", "Scale — $249 (24900 cents)", "SUGGESTED_CREDIT_API_VERSION", "dsp_1043"} {
		if !strings.Contains(string(ledgerly), needle) {
			t.Errorf("ledgerly.mdc lost %q", needle)
		}
	}

	png, err := os.ReadFile(filepath.Join(root, "demo/screenshots/inv-1048-scale.png"))
	if err != nil {
		t.Fatal(err)
	}
	if len(png) < 8 || string(png[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatal("screenshot is not a PNG")
	}

	wf, err := os.ReadFile(filepath.Join(root, ".github/workflows/cursor-agent-demo.yml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(wf)
	for _, needle := range []string{
		"curl https://cursor.com/install -fsS | bash",
		"cursor-agent -p --force",
		"--output-format json",
		"CURSOR_API_KEY: ${{ secrets.CURSOR_API_KEY }}",
		"workflow_dispatch",
		"demo/screenshots/inv-1048-scale.png",
		`data.get("result")`,
	} {
		if !strings.Contains(text, needle) {
			t.Errorf("workflow missing %q", needle)
		}
	}
	for _, banned := range []string{"--image", "--yolo", "pull_request:"} {
		if strings.Contains(text, banned) {
			t.Errorf("workflow should not contain %q", banned)
		}
	}

	mcp, err := os.ReadFile(filepath.Join(root, ".vscode/mcp.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(mcp), "ghp_") || strings.Contains(string(mcp), "sk-") {
		t.Fatal("mcp config contains a secret-shaped token")
	}
	if !strings.Contains(string(mcp), "${env:LEDGERLY_MCP_TOKEN}") {
		t.Fatal("copilot mcp config should interpolate LEDGERLY_MCP_TOKEN")
	}
	cursorMCP, err := os.ReadFile(filepath.Join(root, "demo/migration-fallback/.cursor/mcp.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cursorMCP), `"mcpServers"`) || !strings.Contains(string(cursorMCP), "${env:LEDGERLY_MCP_TOKEN}") {
		t.Fatalf("cursor mcp.json shape:\n%s", cursorMCP)
	}

	hooks, err := os.ReadFile(filepath.Join(root, "demo/migration-fallback/.cursor/hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(hooks), `"version": 1`) || !strings.Contains(string(hooks), "afterFileEdit") || !strings.Contains(string(hooks), "preToolUse") {
		t.Fatalf("hooks.json:\n%s", hooks)
	}
}
