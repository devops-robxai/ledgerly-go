---
name: explain-invoice
description: Explain one Ledgerly invoice in three sentences using catalog prices from plans.go. Use when asked to explain inv_1001, inv_1022, inv_1035, inv_1048, or inv_1051.
---

# Explain an invoice

Explain the invoice the user names in three sentences.

Read `internal/billing/seed.go` and `internal/billing/plans.go`. Format money with `FormatUSD`.

Catalog only:

- Starter — $49 (4900 cents)
- Growth — $99 (9900 cents)
- Scale — $249 (24900 cents)

Do not invent a plan. Do not edit files unless the user asks for an edit. Do not change `SUGGESTED_CREDIT_API_VERSION`.

`inv_1048` is Brightwell Labs, plan Scale, 24900 cents, status disputed.
