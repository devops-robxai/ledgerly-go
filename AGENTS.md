# Ledgerly

Fictional B2B billing. Precise billing-ops voice. Operator Avery Quinn. Synthetic `.example` customers only. Demo clock 23 August 2026.

Copilot and Grok Build both read this file.

## Catalog

Read amounts from `internal/billing/plans.go`. Format with `FormatUSD`.

- Starter — $49 (4900 cents)
- Growth — $99 (9900 cents)
- Scale — $249 (24900 cents)

Do not invent a fourth plan.

## Suggested credit

`SUGGESTED_CREDIT_API_VERSION` in `internal/billing/suggested_credit.go` ships as v1. Dispute `dsp_1043` claims $400 (40000 cents) against Scale invoice `inv_1048` ($249, 24900 cents). v1 returns the claim. v2 caps at the plan price. Keep both `/api/v1` and `/api/v2` suggested-credit routes. Do not edit `TestSuggestedCreditUsesV2Cap` unless the user asks to fix that planted test.

## App

Server port 43173. Workshop cards: `runbooks/101.md`.
