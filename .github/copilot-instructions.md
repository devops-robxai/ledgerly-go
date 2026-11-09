# Ledgerly (Copilot)

Follow `AGENTS.md`. This file is the repo-wide Copilot instruction.

Voice is precise billing ops. Operator Avery Quinn. Customers use synthetic `.example` addresses only. Demo clock is 23 August 2026.

Catalog, read from `internal/billing/plans.go` and formatted with `FormatUSD`:

- Starter — $49 (4900 cents)
- Growth — $99 (9900 cents)
- Scale — $249 (24900 cents)

Do not invent a fourth plan.

Suggested credit ships on v1. `SUGGESTED_CREDIT_API_VERSION` is in `internal/billing/suggested_credit.go`. Dispute `dsp_1043` claims $400 against Scale invoice `inv_1048` ($249). v1 returns the raw claim. v2 caps at $249. Keep both `/api/v1` and `/api/v2` suggested-credit routes. Do not edit `TestSuggestedCreditUsesV2Cap` to get a green test.

Path rules live in `.github/instructions/`. The explain-invoice prompt file is `.github/prompts/explain-invoice.prompt.md`. The custom agent is `.github/agents/billing-reviewer.md`. Skills used with Copilot in this repo are `.github/skills/invoice-memo/` and `.claude/skills/dispute-note/`. MCP for VS Code Copilot is `.vscode/mcp.json`.
