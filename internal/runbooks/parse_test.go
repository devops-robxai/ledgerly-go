package runbooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// rosemary101 is the Go 101 track mirrored from Rosemary's lib/runbooks/beats/101.ts.
// Titles and prompt examples must stay in this order; adapt only Go-specific notes.
var rosemary101 = []struct {
	section string
	beats   []string
}{
	{
		section: "What is Grok Build?",
		beats: []string{
			"Ask",
			"Plan",
			"Build in Agent mode",
			"Debug",
			"Check the models",
			"Plan to fix the bug",
		},
	},
	{
		section: "How do I work with an agent?",
		beats: []string{
			"Run Mode Allowlist",
			"Verify the email feature",
			"Redact (partial)",
			"Stop the prompt",
			"Interrupt and steer",
			"Continue to the end",
			"Review diffs",
			"Restore from a checkpoint",
		},
	},
	{
		section: "How do I govern my agent?",
		beats: []string{
			"Create a user rule",
			"Test the rule",
			"Create a user skill",
			"Test the skill",
			"Canvas",
			"MCP / Figma",
		},
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

	if len(track.Sections) != len(rosemary101) {
		t.Fatalf("sections %d, want %d", len(track.Sections), len(rosemary101))
	}
	for i, sec := range track.Sections {
		want := rosemary101[i]
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
		}
	}
	if hasSectionContaining(track, "Reset") {
		t.Fatal("do not invent a Reset section; reset docs live on the script")
	}

	wantPrompts := []string{
		"/ask Tell me what this application does in 3 sentences",
		"/plan I want a new feature to update the customer email in the invoice detail customer card. Don’t implement email validation.",
		"/debug the failing test",
		"Fix the failing test.",
		"Redact the customer email in the UI. The first two characters and domain are plaintext. When I click to type in the box, clear it and save the new email.",
		"Redact the customer email in the UI. Show it in plaintext when I click the box to edit it. Stop every time you change a file for me to review.",
		"Continue to the end, do not wait for my approval.",
		"/create-rule Preserve the invoice view. Do not rename, restyle, or rearrange invoice screens unless the user names the **exact** new copy (or a specific layout change). This is a personal rule.",
		`Change "Line Items" in the UI to something else.`,
		"/create-skill Use domain-driven design to break down the domains in this application and match it to available APIs or data schemas. This is a personal skill.",
		"Use domain-driven design on this application. Do not edit files.",
		"Create a canvas explaining what we did today.",
		"Create three slides in Figma Slides outlining how I used Grok Build to develop a new feature. I want to use this as part of my demo showcase.",
	}
	for _, want := range wantPrompts {
		if !trackHasPrompt(track, want) {
			t.Errorf("missing prompt %q", want)
		}
	}

	// Go-only presenter note on the same Fix card — not a different beat or prompt.
	fix := findBeat(track, "Plan to fix the bug")
	if fix == nil {
		t.Fatal("missing Plan to fix the bug")
	}
	if !strings.Contains(fix.Detail, "SUGGESTED_CREDIT_API_VERSION") ||
		!strings.Contains(fix.Detail, "internal/billing/suggested_credit.go") {
		t.Errorf("fix beat should keep a presenter note for the Go constant, got %q", fix.Detail)
	}
	if !strings.Contains(fix.Detail, "Use shift-tab to toggle to Agent mode.") {
		t.Errorf("fix beat should keep Rosemary's detail, got %q", fix.Detail)
	}
	if len(fix.Prompts) != 1 || fix.Prompts[0] != "Fix the failing test." {
		t.Errorf("fix prompt %q, want exact card example", fix.Prompts)
	}

	allow := findBeat(track, "Run Mode Allowlist")
	if allow == nil {
		t.Fatal("missing Run Mode Allowlist")
	}
	if !strings.Contains(allow.Detail, "Settings > Agents > Executions & Approvals > Run Mode > Allowlist") {
		t.Errorf("allowlist should keep Rosemary's Settings path, got %q", allow.Detail)
	}
	if !strings.Contains(allow.Detail, "go test") || !strings.Contains(allow.Detail, "go run") {
		t.Errorf("allowlist should mention go test / go run, got %q", allow.Detail)
	}

	// Empty-detail steering/govern beats still parse as cards with their example.
	for _, title := range []string{"Redact (partial)", "Continue to the end", "Test the rule", "Test the skill"} {
		b := findBeat(track, title)
		if b == nil || len(b.Prompts) == 0 {
			t.Errorf("%q should parse as a card with a prompt", title)
		}
	}

	if trackHasPrompt(track, "/create-rule Future code must never call /api/v1/disputes") {
		t.Error("old v1 guardrail prompt should not appear")
	}
	if trackHasPrompt(track, "./scripts/reset-demo-state.sh") {
		t.Error("reset script is not a 101 card prompt")
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
