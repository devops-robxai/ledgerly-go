package httpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devops-robxai/ledgerly-go/internal/billing"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	root := httpRepoRoot(t)
	store := billing.NewStore()
	billing.Seed(store)
	srv, err := NewFromFS(
		store,
		os.DirFS(filepath.Join(root, "web", "templates")),
		os.DirFS(filepath.Join(root, "web", "static")),
		os.DirFS(filepath.Join(root, "runbooks")),
	)
	if err != nil {
		t.Fatal(err)
	}
	return srv
}

func TestRunbooksIndexRedirectsTo101(t *testing.T) {
	srv := testServer(t)
	for _, path := range []string{"/runbooks", "/runbooks/", "/workflows", "/analysis"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Errorf("%s: status %d, want %d", path, rec.Code, http.StatusFound)
		}
		if loc := rec.Header().Get("Location"); loc != "/runbooks/101" {
			t.Errorf("%s: Location %q, want /runbooks/101", path, loc)
		}
	}
}

func TestRunbooks101RendersCopyableCards(t *testing.T) {
	srv := testServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/runbooks/101", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		"Runbooks",
		"Before (09:22)",
		"Modes (09:24)",
		"Rules (09:27)",
		"Skills (09:31)",
		"Subagent and MCP (09:34)",
		"Hooks (09:36)",
		"Prove it (09:37)",
		"CLI and CI (09:39)",
		"If the agent stalls",
		"Do not open Multitask",
		"/ask List the Copilot customization files in this repo",
		"/plan Map each Copilot customization file",
		"/debug A rule was saved as .cursor/rules/go-money.md",
		"SUGGESTED_CREDIT_API_VERSION",
		"TestSuggestedCreditUsesV2Cap",
		"Starter $49",
		"Scale $249",
		"/migrate-to-skills",
		"/invoice-memo Write the memo line for invoice inv_1048",
		"@demo-reset",
		"Run in the terminal",
		"curl https://cursor.com/install -fsS | bash",
		"cursor-agent -p --force",
		"--output-format json",
		"There is no image flag",
		"demo/screenshots/inv-1048-scale.png",
		"./scripts/apply-migration-fallback.sh",
		"Paste in Cursor",
		"data-copy",
		"runbooks/101.md",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("101 page missing %q", want)
		}
	}
}

func TestRunbooksUnknownTrack404(t *testing.T) {
	srv := testServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/runbooks/999", nil)
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d, want 404", rec.Code)
	}
}

func httpRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "web", "templates", "layout.html")); err == nil {
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
