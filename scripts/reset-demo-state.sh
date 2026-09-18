#!/usr/bin/env bash
# Restore Ledgerly planted seams for the next Cursor 101 demo.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

FILE="$ROOT/internal/billing/suggested_credit.go"
if [[ ! -f "$FILE" ]]; then
  echo "missing $FILE" >&2
  exit 1
fi

# Ensure client is back on deprecated v1.
if grep -q 'SUGGESTED_CREDIT_API_VERSION = "v2"' "$FILE"; then
  if command -v sed >/dev/null; then
    sed -i.bak 's/SUGGESTED_CREDIT_API_VERSION = "v2"/SUGGESTED_CREDIT_API_VERSION = "v1"/' "$FILE"
    rm -f "${FILE}.bak"
  fi
  echo "restored SUGGESTED_CREDIT_API_VERSION to v1"
else
  echo "SUGGESTED_CREDIT_API_VERSION already v1 (or unexpected); left unchanged"
fi

# Drop live /create-rule leftover if present.
rm -f "$ROOT/.cursor/rules/suggested-credit-api-v2.mdc"

echo "Planted state:"
echo "  - client: v1 (dsp_1043 UI → \$400.00)"
echo "  - TestSuggestedCreditUsesV2Cap should FAIL"
echo "  - both /api/v1 and /api/v2 suggested-credit routes preserved"
echo
echo "Quick check:"
echo "  go test ./internal/billing/ -run TestSuggestedCreditUsesV2Cap"
