#!/usr/bin/env bash
# beforeMCPExecution: allow only MCP servers listed in mcp-allowlist.json.
# Missing or unknown mcp_server_name is a deny.
# stdin includes tool_name, tool_input, mcp_server_name, and url or command.
# stdout: {"permission":"allow"|"deny","user_message":"...","agent_message":"..."}
set -u
HERE="$(cd "$(dirname "$0")" && pwd)"
ALLOWLIST="$HERE/mcp-allowlist.json"
export LEDGERLY_MCP_ALLOWLIST="$ALLOWLIST"
input=$(cat || true)
printf '%s' "$input" | python3 -c '
import json, os, sys

raw = sys.stdin.read()
try:
    data = json.loads(raw) if raw.strip() else {}
except Exception:
    data = {}
if not isinstance(data, dict):
    data = {}

def reply(permission, user_message, agent_message=None):
    out = {"permission": permission, "user_message": user_message}
    if agent_message:
        out["agent_message"] = agent_message
    print(json.dumps(out))

path = os.environ.get("LEDGERLY_MCP_ALLOWLIST", "")
try:
    allow = json.loads(open(path, encoding="utf-8").read())
    servers = allow.get("servers") or []
    allowed = {str(name).lower() for name in servers if isinstance(name, str) and name.strip()}
except Exception:
    reply("deny", "MCP hook: allowlist could not be read.",
          "The beforeMCPExecution hook could not read mcp-allowlist.json. Do not retry.")
    raise SystemExit(0)

name = data.get("mcp_server_name")
if not isinstance(name, str) or not name.strip():
    reply("deny", "Blocked an MCP call with no server name. The project allowlist only permits listed servers.",
          "mcp_server_name was missing. Treat unknown MCP servers as denied. Do not retry.")
    raise SystemExit(0)

key = name.strip().lower()
if key not in allowed:
    reply("deny",
          "Blocked MCP server %r. The project allowlist permits: %s." % (name, ", ".join(sorted(allowed)) or "(none)"),
          "MCP server %r is not on the project allowlist. Do not call it again." % name)
    raise SystemExit(0)

reply("allow", "")
'
exit 0
