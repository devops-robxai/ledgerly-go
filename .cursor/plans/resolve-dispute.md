# Resolve a dispute

Operators can resolve an open dispute by applying the suggested credit and marking the dispute resolved. Split the work across three subagents, each on its own worktree. Avoid conflicting edits: a subagent only changes the files listed in its section.

Do not change `SUGGESTED_CREDIT_API_VERSION`. It ships as v1. Do not edit `TestSuggestedCreditUsesV2Cap`. Keep both `/api/v1` and `/api/v2` suggested-credit routes. Catalog prices come from `internal/billing/plans.go` (Starter $49, Growth $99, Scale $249) and operator-facing amounts use `FormatUSD`.

When the three subagents finish, run the `ledgerly-reviewer` subagent on the combined result.

## Data

Files:

- `internal/billing/store.go`
- `internal/billing/models.go` (only if you need a field for the applied credit)
- `internal/billing/resolve_test.go` (new)

Add `Store.ResolveDispute(id string)`:

- Return an error if the dispute does not exist or is already `resolved`.
- Compute the credit with the existing `ResolveSuggestedCredit` helper. Do not reimplement v1 or v2, and do not flip the client constant.
- Set the dispute status to `resolved`.
- If you record the applied amount, store it on the dispute in cents.
- If the linked invoice status is `disputed`, set it to `open` (do not invent a new invoice status).

Cover `dsp_1043` in the new test: the shipped v1 client suggests $400 (40000 cents) against Scale invoice `inv_1048` ($249). Leave `TestSuggestedCreditUsesV2Cap` failing on a clean tree.

## API

Files:

- `internal/http/server.go`
- `internal/http/handlers.go`

Add `POST /disputes/{id}/resolve`. That exact path and method are the contract with the UI subagent.

- Call `Store.ResolveDispute`.
- On success, redirect to `/disputes/{id}`.
- On a missing dispute, return 404. If it is already resolved, return 409.
- Do not add new `/api/v1` or `/api/v2` routes for this, and do not change the suggested-credit handlers.

## UI

Files:

- `web/templates/dispute_detail.html`

On `/disputes/{id}`, if the dispute status is `open` or `needs_review`, show a form:

- `method="post"` `action="/disputes/{{.Dispute.ID}}/resolve"`
- A single submit button labelled "Resolve dispute".

If the dispute is `resolved`, show no form. If an applied credit is available, show it with `FormatUSD` (the `usd` template func). Do not edit other templates.

## Conflicts

- The route is `POST /disputes/{id}/resolve`. Do not pick a different path.
- Data owns `internal/billing`. API owns `internal/http`. UI owns `web/templates/dispute_detail.html`.
- Do not edit another subagent's files. Do not format or "tidy" files outside your list.
