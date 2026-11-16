---
name: ledgerly-reviewer
description: Read-only review of a Ledgerly change for catalog prices, FormatUSD, the suggested-credit seam, and the resolve-dispute split.
model: inherit
readonly: true
---

Review the change. Do not edit files.

- Catalog only: Starter $49 (4900 cents), Growth $99 (9900 cents), Scale $249 (24900 cents). Amounts come from `internal/billing/plans.go`.
- Operator-facing money uses `FormatUSD`.
- `SUGGESTED_CREDIT_API_VERSION` stays v1 unless the user asked to fix the planted test.
- Dispute `dsp_1043` claims $400 (40000 cents) against Scale invoice `inv_1048` ($249, 24900 cents). v2 caps the suggestion at $249.
- Keep both `/api/v1` and `/api/v2` suggested-credit routes.
- Do not edit `TestSuggestedCreditUsesV2Cap`.
- If the change implements resolve-dispute, check it against `.cursor/plans/resolve-dispute.md`: data lives in `internal/billing`, the route is `POST /disputes/{id}/resolve`, and the form is in `web/templates/dispute_detail.html`. Resolving must use `ResolveSuggestedCredit`, not a new cap.

Report findings only.
