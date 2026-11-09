package runbooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// betsson101 is the 10 Nov 2026 Copilot-to-Grok Build block in runbooks/101.md.
// Titles and the prompt examples below must stay in this order.
var betsson101 = []struct {
	section string
	beats   []string
}{
	{
		section: "Before (09:22)",
		beats:   []string{"Open Ledgerly", "Copilot files on this branch"},
	},
	{
		section: "Modes (09:24)",
		beats:   []string{"Which modes", "Ask", "Plan", "Debug", "Back to Agent"},
	},
	{
		section: "Rules (09:27)",
		beats:   []string{"Migrate instructions", "Four rule types"},
	},
	{
		section: "Skills (09:31)",
		beats:   []string{"Copy the GitHub skill", "Prompt file becomes a skill"},
	},
	{
		section: "Subagent and MCP (09:34)",
		beats:   []string{"Custom agent", "MCP config"},
	},
	{
		section: "Hooks (09:36)",
		beats:   []string{"Optional hooks"},
	},
	{
		section: "Prove it (09:37)",
		beats:   []string{"Run the skill", "Use the manual rule"},
	},
	{
		section: "CLI and CI (09:39)",
		beats:   []string{"Install the CLI", "Read a screenshot", "The Actions workflow"},
	},
	{
		section: "If the agent stalls",
		beats:   []string{"Apply the fallback"},
	},
}

func TestParse101SourceOfTruth(t *testing.T) {
	md, err := os.ReadFile(filepath.Join(repoRoot(t), "runbooks", "101.md"))
	if err != nil {
		t.Fatal(err)
	}
	track, err := Parse(string(md))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(track.Title, "101") {
		t.Fatalf("title %q, want 101", track.Title)
	}
	if !strings.Contains(track.Intro, "43173") {
		t.Fatalf("intro should mention port 43173: %q", track.Intro)
	}

	if len(track.Sections) != len(betsson101) {
		t.Fatalf("sections %d, want %d", len(track.Sections), len(betsson101))
	}
	for i, sec := range track.Sections {
		want := betsson101[i]
		if sec.Title != want.section {
			t.Errorf("section[%d] title %q, want %q", i, sec.Title, want.section)
		}
		if len(sec.Beats) != len(want.beats) {
			t.Errorf("section %q beats %d, want %d", sec.Title, len(sec.Beats), len(want.beats))
		}
		for j, beat := range sec.Beats {
			if j < len(want.beats) && beat.Title != want.beats[j] {
				t.Errorf("section %q beat[%d] %q, want %q", sec.Title, j, beat.Title, want.beats[j])
			}
			if len(beat.Prompts) != len(beat.PromptKinds) {
				t.Errorf("section %q beat %q prompts %d kinds %d", sec.Title, beat.Title, len(beat.Prompts), len(beat.PromptKinds))
			}
		}
	}
	if hasSectionContaining(track, "Reset") {
		t.Fatal("do not invent a Reset section; reset docs live on the script")
	}
	if hasSectionContaining(track, "Multitask") {
		t.Fatal("do not add a Multitask section")
	}

	wantPrompts := []string{
		"/ask List the Copilot customization files in this repo and say what each one is for. Do not edit files.",
		"/plan Map each Copilot customization file to the Grok Build file that should replace it. Do not create or edit files. Do not change SUGGESTED_CREDIT_API_VERSION.",
		"/debug A rule was saved as .cursor/rules/go-money.md and Grok Build is not applying it. Find out why. Do not change SUGGESTED_CREDIT_API_VERSION. Do not edit tests.",
		"Migrate Copilot instructions to Grok Build rules.",
		"Copy .github/skills/ to .cursor/skills/. Do not edit Go code.",
		"/migrate-to-skills Turn .github/prompts/explain-invoice.prompt.md into a skill at .cursor/skills/explain-invoice/SKILL.md. Do not create a command. Do not edit Go code.",
		"Create .cursor/agents/billing-reviewer.md from .github/agents/billing-reviewer.md.",
		"Create .cursor/mcp.json from .vscode/mcp.json. Use the mcpServers key.",
		`Add .cursor/hooks.json with {"version":1,"hooks":{...}}.`,
		"/invoice-memo Write the memo line for invoice inv_1048. Do not edit files.",
		"@demo-reset What do I run between workshop groups, and which suggested-credit version must stay in place?",
		"curl https://cursor.com/install -fsS | bash",
		`cursor-agent -p --force "Read the image at demo/screenshots/inv-1048-scale.png. Reply with one sentence that includes the plan name and the invoice amount in USD. Do not edit files." --output-format json`,
		"./scripts/apply-migration-fallback.sh",
	}
	for _, want := range wantPrompts {
		if !trackHasPrompt(track, want) {
			t.Errorf("missing prompt %q", want)
		}
	}

	modes := findBeat(track, "Which modes")
	if modes == nil || !strings.Contains(modes.Detail, "Do not open Multitask") {
		t.Fatal("modes card should tell the presenter not to open Multitask")
	}
	if strings.Contains(trackText(track), "--image") {
		t.Fatal("do not document an image flag")
	}

	shot := findBeat(track, "Read a screenshot")
	if shot == nil {
		t.Fatal("missing Read a screenshot")
	}
	if len(shot.Prompts) != 1 || shot.PromptKinds[0] != "terminal" {
		t.Fatalf("screenshot card should be one terminal command, got %#v %#v", shot.Prompts, shot.PromptKinds)
	}
	if !strings.Contains(shot.Detail, "There is no image flag") {
		t.Errorf("screenshot card should say there is no image flag, got %q", shot.Detail)
	}

	migrate := findBeat(track, "Migrate instructions")
	if migrate == nil {
		t.Fatal("missing Migrate instructions")
	}
	body := migrate.Prompts[0]
	for _, needle := range []string{
		".cursor/rules/ledgerly.mdc",
		"go-money.mdc",
		"suggested-credit-seam.mdc",
		"demo-reset.mdc",
		"SUGGESTED_CREDIT_API_VERSION",
		"TestSuggestedCreditUsesV2Cap",
		"Starter $49 (4900 cents)",
		"Growth $99 (9900 cents)",
		"Scale $249 (24900 cents)",
		"FormatUSD",
	} {
		if !strings.Contains(body, needle) {
			t.Errorf("migrate prompt missing %q", needle)
		}
	}

	if trackHasPrompt(track, "Fix the failing test.") {
		t.Error("this session does not fix the planted test")
	}
	if trackHasPrompt(track, "Redact the customer email") {
		t.Error("email redact cards are not part of this session")
	}
	full := trackText(track)
	for _, banned := range []string{".claude/", ".codex/", "dispute-note", "Leave the Claude skill"} {
		if strings.Contains(full, banned) {
			t.Errorf("runbook still mentions %q", banned)
		}
	}
}

func TestParseEmptyDetailBeatStillRenders(t *testing.T) {
	track, err := Parse("# 101\n\n## How do I work with an agent?\n\n### Redact (partial)\n\n```text\nRedact the customer email in the UI.\n```\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(track.Sections) != 1 || len(track.Sections[0].Beats) != 1 {
		t.Fatalf("got %+v", track)
	}
	b := track.Sections[0].Beats[0]
	if b.Title != "Redact (partial)" || b.Detail != "" || len(b.Prompts) != 1 {
		t.Fatalf("got %+v", b)
	}
	if len(b.PromptKinds) != 1 || b.PromptKinds[0] != "cursor" {
		t.Fatalf("default fence kind %+v", b.PromptKinds)
	}
}

func TestTerminalFenceKind(t *testing.T) {
	track, err := Parse("# 101\n\n## CLI and CI (09:39)\n\n### Install the CLI\n\n```bash\ncurl https://cursor.com/install -fsS | bash\n```\n")
	if err != nil {
		t.Fatal(err)
	}
	b := track.Sections[0].Beats[0]
	if len(b.PromptKinds) != 1 || b.PromptKinds[0] != "terminal" {
		t.Fatalf("got %+v", b.PromptKinds)
	}
}

func TestSafeFileName(t *testing.T) {
	name, ok := SafeFileName("101")
	if !ok || name != "101.md" {
		t.Fatalf("got %q %v", name, ok)
	}
	if _, ok := SafeFileName("../etc/passwd"); ok {
		t.Fatal("rejected path traversal")
	}
	if _, ok := SafeFileName("101/extra"); ok {
		t.Fatal("rejected slash")
	}
}

func TestInlineHTML(t *testing.T) {
	got := string(InlineHTML("Port **43173**. Use `go test`."))
	if !strings.Contains(got, "<strong>43173</strong>") {
		t.Fatalf("bold: %s", got)
	}
	if !strings.Contains(got, "<code>go test</code>") {
		t.Fatalf("code: %s", got)
	}
}

func findBeat(track *Track, title string) *Beat {
	for i := range track.Sections {
		for j := range track.Sections[i].Beats {
			if track.Sections[i].Beats[j].Title == title {
				return &track.Sections[i].Beats[j]
			}
		}
	}
	return nil
}

func trackText(track *Track) string {
	var b strings.Builder
	b.WriteString(track.Intro)
	for _, sec := range track.Sections {
		b.WriteString(sec.Title)
		for _, beat := range sec.Beats {
			b.WriteString(beat.Detail)
			for _, p := range beat.Prompts {
				b.WriteString(p)
			}
		}
	}
	return b.String()
}

func trackHasPrompt(track *Track, substr string) bool {
	for _, sec := range track.Sections {
		for _, beat := range sec.Beats {
			for _, p := range beat.Prompts {
				if strings.Contains(p, substr) {
					return true
				}
			}
		}
	}
	return false
}

func hasSectionContaining(track *Track, substr string) bool {
	for _, sec := range track.Sections {
		if strings.Contains(strings.ToLower(sec.Title), strings.ToLower(substr)) {
			return true
		}
	}
	return false
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "runbooks", "101.md")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("repo root not found")
	return ""
}
