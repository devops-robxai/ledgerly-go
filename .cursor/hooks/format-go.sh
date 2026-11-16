#!/usr/bin/env bash
# afterFileEdit: run gofmt -w on the edited file when it is a .go file.
# stdin: {"file_path":"<absolute path>","edits":[{"old_string":"...","new_string":"..."}]}
set -u
input=$(cat || true)
path=$(printf '%s' "$input" | python3 -c '
import json, sys
raw = sys.stdin.read()
try:
    data = json.loads(raw) if raw.strip() else {}
except Exception:
    raise SystemExit(0)
if not isinstance(data, dict):
    raise SystemExit(0)
path = data.get("file_path") or ""
if isinstance(path, str) and path.endswith(".go"):
    print(path)
' 2>/dev/null || true)
if [[ -n "${path}" && -f "${path}" ]] && command -v gofmt >/dev/null; then
  gofmt -w "${path}" || true
fi
exit 0
