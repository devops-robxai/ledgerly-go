#!/usr/bin/env bash
# beforeShellExecution: allow ordinary commands, ask on git push / rm, deny
# destructive ones (rm -rf, git push --force, database drops).
# stdin: {"command":"...","cwd":"...","sandbox":false}
# stdout: {"permission":"allow"|"deny"|"ask","user_message":"...","agent_message":"..."}
set -u
input=$(cat || true)
printf '%s' "$input" | python3 -c '
import json, re, sys

raw = sys.stdin.read()
try:
    data = json.loads(raw) if raw.strip() else {}
except Exception:
    data = {}

command = data.get("command") if isinstance(data, dict) else None
if not isinstance(command, str):
    command = ""
text = command.strip()
lower = text.lower()

def reply(permission, user_message, agent_message=None):
    out = {"permission": permission, "user_message": user_message}
    if agent_message:
        out["agent_message"] = agent_message
    print(json.dumps(out))

if not text:
    reply("deny", "Shell hook: command was missing from the event.",
          "The beforeShellExecution event had no command. Do not retry.")
    raise SystemExit(0)

# Deny: recursive force-delete, force push, and database drops.
if re.search(r"\brm\s+(-[^\s]*f[^\s]*r|-[^\s]*r[^\s]*f|--recursive\s+--force|--force\s+--recursive)\b", lower):
    reply("deny", "Blocked rm -rf. The project shell hook does not allow recursive force deletes.",
          "Do not run rm -rf. The project beforeShellExecution hook denied it.")
    raise SystemExit(0)

if re.search(r"\bgit\s+push\b", lower) and re.search(r"(--force-with-lease|--force|-f)(\s|$)", lower):
    reply("deny", "Blocked git push --force. The project shell hook does not allow force pushes.",
          "Do not force-push. The project beforeShellExecution hook denied it.")
    raise SystemExit(0)

if re.search(r"\bdrop\s+(database|schema|table|db)\b", lower) or re.search(r"\bdb\s+drop\b", lower):
    reply("deny", "Blocked a database drop. The project shell hook does not allow drop commands.",
          "Do not drop a database or table. The project beforeShellExecution hook denied it.")
    raise SystemExit(0)

# Ask: a normal git push, or rm without -rf, so the operator can refuse in the UI.
if re.search(r"\bgit\s+push\b", lower):
    reply("ask", "The project shell hook wants confirmation before git push. Deny this if you are in the live demo.")
    raise SystemExit(0)

if re.search(r"(^|[;&|]\s*)rm\s+", lower):
    reply("ask", "The project shell hook wants confirmation before rm.")
    raise SystemExit(0)

reply("allow", "")
'
exit 0
