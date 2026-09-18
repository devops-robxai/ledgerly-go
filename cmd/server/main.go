package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/robbarnes/ledgerly/internal/billing"
	httpserver "github.com/robbarnes/ledgerly/internal/http"
)

const addr = ":43173"

func main() {
	store := billing.NewStore()
	billing.Seed(store)

	root, err := findRepoRoot()
	if err != nil {
		log.Fatalf("repo root: %v", err)
	}
	tmplFS := os.DirFS(filepath.Join(root, "web", "templates"))
	staticFS := os.DirFS(filepath.Join(root, "web", "static"))

	srv, err := httpserver.NewFromFS(store, tmplFS, staticFS)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	fmt.Printf("Ledgerly listening on http://localhost%s\n", addr)
	fmt.Printf("  dispute seam: %s → shows v1 credit for dsp_1043\n", billing.SuggestedCreditPath("dsp_1043"))
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}

func findRepoRoot() (string, error) {
	// Prefer cwd when run as `go run ./cmd/server` from module root.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	candidates := []string{cwd}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Dir(exe), filepath.Join(filepath.Dir(exe), "..", ".."))
	}
	for _, c := range candidates {
		c = filepath.Clean(c)
		if hasWeb(c) {
			return c, nil
		}
		// walk up a few levels
		dir := c
		for i := 0; i < 6; i++ {
			if hasWeb(dir) {
				return dir, nil
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return "", fmt.Errorf("could not find web/templates (run from module root)")
}

func hasWeb(dir string) bool {
	_, err := fs.Stat(os.DirFS(dir), "web/templates/dashboard.html")
	return err == nil
}
