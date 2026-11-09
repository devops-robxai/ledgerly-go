#!/usr/bin/env bash
# Copy the reviewed Copilot-to-Cursor after-state into the working tree.
# Does not modify .cursor/rules/ledgerly.mdc, AGENTS.md, or the Copilot source files.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="$ROOT/demo/migration-fallback/.cursor"
if [[ ! -d "$SRC/rules" ]]; then
  echo "missing $SRC" >&2
  exit 1
fi

mkdir -p "$ROOT/.cursor"
cp -R "$SRC/." "$ROOT/.cursor/"

mkdir -p "$ROOT/.cursor/skills"
rm -rf "$ROOT/.cursor/skills/invoice-memo"
cp -R "$ROOT/.github/skills/invoice-memo" "$ROOT/.cursor/skills/invoice-memo"

chmod +x "$ROOT/.cursor/hooks/"*.sh

if [[ ! -f "$ROOT/.cursor/rules/ledgerly.mdc" ]]; then
  echo "ledgerly.mdc disappeared; refusing to finish" >&2
  exit 1
fi

echo "Applied demo/migration-fallback."
echo "  rules: go-money.mdc, suggested-credit-seam.mdc, demo-reset.mdc (ledgerly.mdc kept)"
echo "  skills: invoice-memo (from .github/skills), explain-invoice"
echo "  subagent: .cursor/agents/billing-reviewer.md"
echo "  mcp: .cursor/mcp.json"
echo "  hooks: .cursor/hooks.json"
