#!/usr/bin/env bash
# afterFileEdit: format a Go file when the event names one.
# TO CONFIRM: stdin is a JSON event. This script tries common path fields.
set -u
input=$(cat || true)
path=$(printf '%s' "$input" | python3 -c '
import json, sys
raw = sys.stdin.read()
try:
    data = json.loads(raw)
except Exception:
    sys.exit(0)

found = []

def take(value):
    if isinstance(value, str):
        found.append(value)

def walk(node):
    if isinstance(node, dict):
        for key in ("file_path", "path", "filePath", "uri", "file"):
            take(node.get(key))
        for value in node.values():
            walk(value)
    elif isinstance(node, list):
        for value in node:
            walk(value)

walk(data)
for candidate in found:
    if candidate.endswith(".go"):
        print(candidate)
        break
' 2>/dev/null || true)
if [[ -n "${path}" && -f "${path}" ]] && command -v gofmt >/dev/null; then
  gofmt -w "${path}" || true
fi
exit 0
