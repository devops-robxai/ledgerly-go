#!/usr/bin/env bash
# preToolUse: allow normal work. Refuse an event that looks like it prints a demo secret.
# TO CONFIRM: stdin JSON shape, and that stdout {"permission":"allow"|"deny"} is the decision.
set -u
input=$(cat || true)
printf '%s' "$input" | python3 -c '
import json, sys
raw = sys.stdin.read()
text = raw.lower()
secret_names = ("cursor_api_key", "ledgerly_mcp_token")
printers = ("printenv", "echo ", "echo\t", "env", "printf ")
bad = any(name in text for name in secret_names) and any(word in text for word in printers)
if bad:
    print(json.dumps({
        "permission": "deny",
        "user_message": "Refusing a command that would print a demo secret.",
    }))
else:
    print(json.dumps({"permission": "allow"}))
'
exit 0
