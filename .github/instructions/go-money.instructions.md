---
applyTo: "internal/billing/**/*.go,internal/http/**/*.go"
---

# Go money

When you edit these files:

- Read plan amounts from `internal/billing/plans.go`. Do not hard-code a new price.
- Catalog only: Starter $49 (4900 cents), Growth $99 (9900 cents), Scale $249 (24900 cents).
- Format money with `FormatUSD` in `internal/billing/money.go`.
- Customers stay synthetic `.example` addresses. Operator Avery Quinn. Demo clock 23 August 2026.
- Do not invent a fourth plan.
