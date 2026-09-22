# Ledgerly (Go)

Go-native billing-ops demo for **Cursor 101** workshops (SLB Europe). Server-rendered HTML with planted seams for the in-app 101 track (Ask → Plan → Agent → Debug, then steer and govern).

Fictional B2B billing. Operator **Avery Quinn**. Catalog: Starter **$49**, Growth **$99**, Scale **$249**. Demo clock frozen at **23 August 2026**. Synthetic `.example` customers only.

Go alternative to Rosemary’s JS Ledgerly (`joatmon08x/ce-field-demos`).

## Run

Go 1.22+.

```bash
cd ~/Developer/ledgerly   # or this repo root
go run ./cmd/server
```

Open **http://localhost:43173**. Workshop prompts: **http://localhost:43173/runbooks/101**.

```bash
go build -o /tmp/ledgerly ./cmd/server
go test ./...
```

On a clean tree, `go test ./...` is **1 failed** — `TestSuggestedCreditUsesV2Cap` — by design.

## Planted bug (the 101 Debug beat)

| Seam | Shipped state | Correct state |
| --- | --- | --- |
| Suggested-credit client | `SUGGESTED_CREDIT_API_VERSION = "v1"` in `internal/billing/suggested_credit.go` | `"v2"` |
| `GET /api/v1/disputes/dsp_1043/suggested-credit` | **40000** cents ($400.00) — raw claim | deprecated |
| `GET /api/v2/disputes/dsp_1043/suggested-credit` | **24900** cents ($249.00) — Scale cap | correct |
| Dispute HTML `/disputes/dsp_1043` | Calls v1 path → shows **$400.00** | should show **$249.00** |
| `TestSuggestedCreditUsesV2Cap` | **Fails** until client uses v2 | passes after flip |

Do **not** edit the failing test to get green. Flip the client to v2 (and keep both API routes).

## App map

| Path | Role |
| --- | --- |
| `/` | Dashboard KPIs + lists |
| `/invoices`, `/invoices/{id}` | Invoice list / detail (customer email is read-only) |
| `POST /invoices/{id}/email` | Intentionally absent — 101 Plan→Agent builds this live |
| `/disputes`, `/disputes/{id}` | Dispute list / detail + suggested credit |
| `/runbooks`, `/runbooks/101` | In-app 101 prompt cards (copy-paste). `/workflows` and `/analysis` redirect here |
| `/api/v1\|v2/disputes/{id}/suggested-credit` | JSON suggested-credit APIs |

## Workshop notes

1. Copy-paste prompts in the product UI: [http://localhost:43173/runbooks/101](http://localhost:43173/runbooks/101). Cards parse from [`runbooks/101.md`](runbooks/101.md), which mirrors Rosemary’s 101 track (Go paths only where her prompts cite JS/npm).
2. Reset planted seams after a demo: `./scripts/reset-demo-state.sh`. That restores the v1 suggested-credit client **and** the absent invoice email form (strips the form / write path if a practice agent added them). Plan / Build assume the email feature is missing; Debug / Fix assume v1.
3. Invoice email beat uses `inv_1048` (Brightwell Labs). The update form is intentionally missing on a clean tree so Plan→Agent can implement it.
4. Do not push to GitHub unless asked; local `git init` is fine.

## Module

`github.com/devops-robxai/ledgerly-go`
