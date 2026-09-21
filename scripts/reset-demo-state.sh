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

# Strip the invoice email-update form and write path if a practice agent added them.
# The 101 Plan→Agent beat builds this live; a clean tree must show read-only email.
python3 - "$ROOT" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])


def write_if_changed(path: pathlib.Path, new: str, message: str) -> None:
    old = path.read_text()
    if old == new:
        print(f"{path.relative_to(root)}: {message} (already clean)")
        return
    path.write_text(new)
    print(f"{path.relative_to(root)}: {message}")


def strip_leading_comments(src: str, start: int) -> int:
    while start > 0:
        prev_nl = src.rfind("\n", 0, start - 1)
        line_start = 0 if prev_nl < 0 else prev_nl + 1
        line = src[line_start:start].strip()
        if line.startswith("//") or line == "":
            start = line_start
        else:
            break
    return start


def skip_string(src: str, i: int) -> int:
    quote = src[i]
    i += 1
    while i < len(src):
        if src[i] == "\\" and quote != "`":
            i += 2
            continue
        if src[i] == quote:
            return i + 1
        i += 1
    return i


def matching_brace_end(src: str, brace: int) -> int:
    depth = 0
    i = brace
    while i < len(src):
        ch = src[i]
        if ch in ('"', "'", "`"):
            i = skip_string(src, i)
            continue
        if ch == "/" and i + 1 < len(src) and src[i + 1] == "/":
            nl = src.find("\n", i)
            i = len(src) if nl < 0 else nl
            continue
        if ch == "/" and i + 1 < len(src) and src[i + 1] == "*":
            end = src.find("*/", i + 2)
            i = len(src) if end < 0 else end + 2
            continue
        if ch == "{":
            depth += 1
        elif ch == "}":
            depth -= 1
            if depth == 0:
                return i + 1
        i += 1
    return -1


def remove_func(src: str, signature: str) -> str:
    idx = src.find(signature)
    if idx < 0:
        return src
    start = src.rfind("\n", 0, idx)
    start = 0 if start < 0 else start + 1
    start = strip_leading_comments(src, start)
    brace = src.find("{", idx)
    if brace < 0:
        return src
    end = matching_brace_end(src, brace)
    if end < 0:
        return src
    if end < len(src) and src[end] == "\n":
        end += 1
    return src[:start] + src[end:]


def remove_if_with_needle(src: str, needle: str) -> str:
    idx = src.find(needle)
    if idx < 0:
        return src
    start = src.rfind("\n", 0, idx)
    start = 0 if start < 0 else start + 1
    # Walk back to the `if` that owns this condition (same or previous lines).
    look = start
    while look > 0 and "if " not in src[look:idx] and "if(" not in src[look:idx]:
        prev_nl = src.rfind("\n", 0, look - 1)
        look = 0 if prev_nl < 0 else prev_nl + 1
        if src[look:start].strip().startswith("//"):
            break
        start = look
    start = strip_leading_comments(src, start)
    brace = src.find("{", idx)
    if brace < 0:
        return src
    end = matching_brace_end(src, brace)
    if end < 0:
        return src
    while end < len(src) and src[end] in " \t":
        end += 1
    if end < len(src) and src[end] == "\n":
        end += 1
    return src[:start] + src[end:]


def strip_invoice_form(text: str) -> str:
    text = re.sub(r"\n[ \t]*<form\b.*?</form>[ \t]*\n?", "\n", text, flags=re.S | re.I)
    text = re.sub(
        r"\n[ \t]*<(?:label|input|button)\b[^>]*(?:email|Save)[^>]*>.*?(?:</(?:label|button)>|[ \t]*)",
        "",
        text,
        flags=re.S | re.I,
    )
    # Keep the customer card tight: no leftover blank line after a stripped form.
    text = re.sub(
        r'(<p>Email: <span class="mono" id="customer-email">\{\{\.Customer\.Email\}\}</span></p>)\n\n+',
        r"\1\n",
        text,
    )
    return text


# --- template: read-only customer email, no update form ---
tmpl = root / "web/templates/invoice_detail.html"
if tmpl.is_file():
    write_if_changed(tmpl, strip_invoice_form(tmpl.read_text()), "stripped email update form")
else:
    print("missing web/templates/invoice_detail.html", file=sys.stderr)
    sys.exit(1)

# --- handlers: no working POST /invoices/{id}/email ---
handlers = root / "internal/http/handlers.go"
if handlers.is_file():
    src = handlers.read_text()
    original = src
    for needle in (
        'HasSuffix(path, "/email")',
        'HasSuffix(r.URL.Path, "/email")',
        'HasSuffix(path, "/email/")',
        '"/email"',
    ):
        while needle in src and (
            "HasSuffix" in src[max(0, src.find(needle) - 80) : src.find(needle) + 40]
            or "TrimSuffix" in src[max(0, src.find(needle) - 80) : src.find(needle) + 40]
        ):
            nxt = remove_if_with_needle(src, needle)
            if nxt == src:
                break
            src = nxt
    src = remove_func(src, "UpdateCustomerEmail")
    # Drop email-updated flash wiring if the practice agent added it.
    src = re.sub(
        r'\n\tflash := ""\n\tif r\.URL\.Query\(\)\.Get\("saved"\) == "1" \{\n\t\tflash = "Customer email updated\."\n\t\}\n',
        "\n",
        src,
    )
    src = src.replace("Invoice: inv, Customer: cust, Flash: flash,", "Invoice: inv, Customer: cust,")
    src = re.sub(
        r'\n\ts\.mux\.HandleFunc\("/invoices/\{?id\}?/email".*\n',
        "\n",
        src,
    )
    if src != original:
        handlers.write_text(src)
        print("internal/http/handlers.go: removed invoice email write path")
    else:
        print("internal/http/handlers.go: no email write path (already clean)")
else:
    print("missing internal/http/handlers.go", file=sys.stderr)
    sys.exit(1)

# --- store: drop write method used only by the email POST ---
store = root / "internal/billing/store.go"
if store.is_file():
    src = store.read_text()
    new = remove_func(src, "func (s *Store) UpdateCustomerEmail")
    if new != src:
        store.write_text(new)
        print("internal/billing/store.go: removed UpdateCustomerEmail")
    else:
        print("internal/billing/store.go: UpdateCustomerEmail absent (already clean)")

# gofmt so a stripped practice tree still compiles cleanly
import shutil
import subprocess

gofmt = shutil.which("gofmt")
if gofmt:
    subprocess.run(
        [gofmt, "-w", str(handlers), str(store)],
        check=False,
    )
PY

echo "Planted state:"
echo "  - client: v1 (dsp_1043 UI → \$400.00)"
echo "  - TestSuggestedCreditUsesV2Cap should FAIL"
echo "  - both /api/v1 and /api/v2 suggested-credit routes preserved"
echo "  - invoice detail: customer email is read-only (no update form / POST)"
echo
echo "Quick check:"
echo "  go test ./internal/billing/ -run TestSuggestedCreditUsesV2Cap"
