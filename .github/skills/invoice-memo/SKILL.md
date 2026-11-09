---
name: invoice-memo
description: Draft a one-line Ledgerly invoice memo from the catalog. Use when asked for a memo for inv_1001, inv_1022, inv_1035, inv_1048, or inv_1051.
---

# Invoice memo

Read the invoice in `internal/billing/seed.go` and the price in `internal/billing/plans.go`. Format the amount with `FormatUSD`.

Catalog only:

- Starter — $49 (4900 cents)
- Growth — $99 (9900 cents)
- Scale — $249 (24900 cents)

Reply with one line: invoice id, customer, plan, formatted amount, status. Do not invent a plan. Do not edit files unless the user asks for an edit.

Current seed row for the demo: `inv_1048` is Brightwell Labs, plan Scale, 24900 cents, status disputed.
