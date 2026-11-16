package hooks

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".cursor", "hooks.json")); err == nil {
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

func hookPath(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(repoRoot(t), ".cursor", "hooks", name)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("missing hook %s: %v", name, err)
	}
	return p
}

func runHook(t *testing.T, script string, payload any) (string, string, int) {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(script)
	cmd.Dir = repoRoot(t)
	cmd.Stdin = bytes.NewReader(body)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("run %s: %v", script, err)
		}
	}
	return stdout.String(), stderr.String(), code
}

func permissionOf(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &out); err != nil {
		t.Fatalf("hook stdout is not JSON: %q: %v", stdout, err)
	}
	return out
}

func TestGuardShellDecisions(t *testing.T) {
	script := hookPath(t, "guard-shell.sh")
	cases := []struct {
		name, command, want string
	}{
		{"allow test", "go test ./internal/billing/", "allow"},
		{"deny rm -rf", "rm -rf /tmp/ledgerly-demo-scratch", "deny"},
		{"deny rm -fr", "rm -fr /tmp/ledgerly-demo-scratch", "deny"},
		{"deny force push", "git push --force origin betsson-201", "deny"},
		{"deny push -f", "git push -f origin betsson-201", "deny"},
		{"deny drop database", "psql -c 'DROP DATABASE ledgerly'", "deny"},
		{"deny db drop", "echo db drop", "deny"},
		{"ask git push", "git push --dry-run origin betsson-201", "ask"},
		{"ask rm", "rm /tmp/ledgerly-demo-scratch", "ask"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runHook(t, script, map[string]any{
				"command": tc.command,
				"cwd":     "/tmp",
				"sandbox": false,
			})
			if code != 0 {
				t.Fatalf("exit %d stderr=%s stdout=%s", code, stderr, stdout)
			}
			got := permissionOf(t, stdout)
			if got["permission"] != tc.want {
				t.Fatalf("permission=%v want %s (stdout=%s)", got["permission"], tc.want, stdout)
			}
		})
	}
}

func TestGuardMCPAllowlist(t *testing.T) {
	script := hookPath(t, "guard-mcp.sh")
	cases := []struct {
		name, server, want string
	}{
		{"allow linear", "linear", "allow"},
		{"deny slack", "slack", "deny"},
		{"deny github", "github", "deny"},
		{"deny missing name", "", "deny"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]any{
				"tool_name":  "list_issues",
				"tool_input": "{}",
			}
			if tc.server != "" {
				payload["mcp_server_name"] = tc.server
			}
			stdout, stderr, code := runHook(t, script, payload)
			if code != 0 {
				t.Fatalf("exit %d stderr=%s stdout=%s", code, stderr, stdout)
			}
			got := permissionOf(t, stdout)
			if got["permission"] != tc.want {
				t.Fatalf("permission=%v want %s (stdout=%s)", got["permission"], tc.want, stdout)
			}
		})
	}
}

func TestFormatGoAfterFileEdit(t *testing.T) {
	script := hookPath(t, "format-go.sh")
	dir := t.TempDir()
	path := filepath.Join(dir, "messy.go")
	original := "package messy\nfunc  Add( a,b int )int{return a+b}\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := runHook(t, script, map[string]any{
		"file_path": path,
		"edits":     []map[string]string{{"old_string": "", "new_string": original}},
	})
	if code != 0 {
		t.Fatalf("exit %d stderr=%s stdout=%s", code, stderr, stdout)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "package messy\n\nfunc Add(a, b int) int { return a + b }\n"
	if string(got) != want {
		t.Fatalf("gofmt result:\n%s\nwant:\n%s", got, want)
	}

	other := filepath.Join(dir, "notes.md")
	if err := os.WriteFile(other, []byte("# hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, stderr, code = runHook(t, script, map[string]any{"file_path": other, "edits": []any{}})
	if code != 0 {
		t.Fatalf("non-go edit should be a no-op, exit %d stderr=%s", code, stderr)
	}
	md, err := os.ReadFile(other)
	if err != nil {
		t.Fatal(err)
	}
	if string(md) != "# hi\n" {
		t.Fatalf("markdown was changed: %q", md)
	}
}

func TestHooksJSONSchema(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), ".cursor", "hooks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]any
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg["version"] != float64(1) {
		t.Fatalf("version=%v", cfg["version"])
	}
	hooks, _ := cfg["hooks"].(map[string]any)
	for _, name := range []string{"beforeShellExecution", "afterFileEdit", "beforeMCPExecution"} {
		list, ok := hooks[name].([]any)
		if !ok || len(list) == 0 {
			t.Fatalf("missing hook list %s", name)
		}
		first, _ := list[0].(map[string]any)
		cmd, _ := first["command"].(string)
		if !strings.HasPrefix(cmd, ".cursor/hooks/") {
			t.Fatalf("%s command should be project-relative, got %q", name, cmd)
		}
	}
	mcp := hooks["beforeMCPExecution"].([]any)[0].(map[string]any)
	if mcp["failClosed"] != true {
		t.Fatalf("beforeMCPExecution should set failClosed, got %#v", mcp["failClosed"])
	}
	allow, err := os.ReadFile(filepath.Join(repoRoot(t), ".cursor", "hooks", "mcp-allowlist.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(allow), `"linear"`) {
		t.Fatalf("allowlist should include linear: %s", allow)
	}
}
