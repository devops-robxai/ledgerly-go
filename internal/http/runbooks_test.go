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
		"/debug the failing test",
		"SUGGESTED_CREDIT_API_VERSION",
		"Paste in Cursor",
		"data-copy",
		"Ask",
		"Plan",
		"cmd/server/main.go",
		"internal/billing/plans.go",
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
