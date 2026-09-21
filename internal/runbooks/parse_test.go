package runbooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

	wantPrompts := []string{
		"/ask Tell me what this Ledgerly Go app does",
		"internal/billing/plans.go",
		"/plan I want a feature to update the customer email",
		"/debug the failing test",
		"SUGGESTED_CREDIT_API_VERSION",
		"/create-rule Future code must never call /api/v1/disputes",
		"./scripts/reset-demo-state.sh",
	}
	for _, want := range wantPrompts {
		if !trackHasPrompt(track, want) {
			t.Errorf("missing prompt containing %q", want)
		}
	}

	if !hasSectionTitle(track, "Ask") && !hasBeatTitle(track, "Ask") {
		t.Fatal("expected an Ask beat")
	}
	if !hasBeatTitle(track, "Plan") || !hasBeatTitle(track, "Debug") {
		t.Fatal("expected Plan and Debug beats")
	}
	if !hasSectionContaining(track, "Reset") {
		t.Fatal("expected a Reset section")
	}
	if !hasSectionContaining(track, "govern") {
		t.Fatal("expected a govern section")
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

func hasBeatTitle(track *Track, title string) bool {
	for _, sec := range track.Sections {
		for _, beat := range sec.Beats {
			if beat.Title == title {
				return true
			}
		}
	}
	return false
}

func hasSectionTitle(track *Track, title string) bool {
	for _, sec := range track.Sections {
		if sec.Title == title {
			return true
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
