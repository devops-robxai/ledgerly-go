# Ledgerly (Go)

Go-native billing-ops demo. Branch `betsson` is the Betsson Group 101 on Tue 10 Nov 2026, 09:00 London: migrate a GitHub Copilot setup to Cursor, then show Cursor CLI and a GitHub Actions workflow. In-app script: [`runbooks/101.md`](runbooks/101.md). The 20-minute block is 09:22–09:42.

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

Do **not** edit the failing test to get green. Flip the client to v2 (and keep both API routes). The Betsson migration demo does **not** flip this constant.

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

1. Copy-paste prompts: [http://localhost:43173/runbooks/101](http://localhost:43173/runbooks/101). Cards parse from [`runbooks/101.md`](runbooks/101.md).
2. Copilot files are already on `betsson` (`.github/`, `.claude/skills/`, `.vscode/mcp.json`, `AGENTS.md`). The live demo migrates them. If the agent stalls, `./scripts/apply-migration-fallback.sh` copies the reviewed after-state from `demo/migration-fallback/` and does not change `.cursor/rules/ledgerly.mdc`.
3. `./scripts/reset-demo-state.sh` restores the v1 suggested-credit client and strips an invoice email form if a practice session added one. This migration session does not add that form and does not flip v1.
4. CLI image: `demo/screenshots/inv-1048-scale.png`. Actions workflow: `.github/workflows/cursor-agent-demo.yml` (`workflow_dispatch`, secret `CURSOR_API_KEY`). No real secrets are committed.

## Module

`github.com/devops-robxai/ledgerly-go`
