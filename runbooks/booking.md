# Booking.com: Grok Build for engineers

Thu 5 Nov 2026, 08:00 UK (16:00 Shanghai). Remote. Demos only. App: **Ledgerly**, port **43173**.

## Pre-flight (tonight)

### Team and Auto

Sign in with the **cursor.sh account on the ADM/FE team**. Check that Auto shows **Cost, Balance and Intelligence**.

### Start the app

**Open a new agent.** Local agent. Check the dashboard, invoices and disputes pages load in the Browser view.

```text
Start the application.
```

```text
/rename-chat Start app
```

### Reset the planted state

Afterwards, dsp_1043 shows **$400.00** and `TestSuggestedCreditUsesV2Cap` fails.

```text
Run ./scripts/reset-demo-state.sh and then run go test ./internal/billing/ -run TestSuggestedCreditUsesV2Cap. Report the result only. Do not fix anything.
```

### Clean up after rehearsal

Delete `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and any saved plan file. Keep `.cursor/rules/ledgerly.mdc`.

### Screen

Light mode, larger fonts, left panel collapsed, notifications muted.

## Morning of

### Final checks

App running in the Browser view. Reset run. No Demo 4 rule or skill left over. Model on **Auto**, mode on **Ask**. Join at 07:55.

## Demo 1: The agent loop: Ask, Plan, Agent, Debug, and Auto

### Show Ledgerly

Show http://localhost:43173 in the Browser view. Explain what the app does.

### Replay the Start app chat

Open the **Start app** chat. Explain the loop: understand, plan, act, verify.

### Ask: understand

**Open a new agent.** Show the four modes. Ask mode. Explain that Ask reads but never edits.

```text
/ask What does Ledgerly do? Name the main packages for invoices, disputes and the HTTP layer in 5 bullets. Do not edit files.
```

### Plan: map the approach

**Open a new agent.** Plan mode. Show http://localhost:43173/invoices/inv_1048 first: the customer email is read-only. Explain reviewing the plan before any code changes.

```text
/plan I want a new feature to update the customer email in the invoice detail customer card. Don’t implement email validation.
```

### Agent: build the plan

Click **Build** on the plan. When it finishes, reload inv_1048 and show the email edit.

### Pick a model

Open the model picker. Explain starting with the default and moving up only if it misses.

### Auto: Cost, Balance, Intelligence

Click **Auto**, then **Balance**. Explain the three modes and how each bills.

### Debug: investigate

**Open a new agent.** Debug mode. Show the top hypothesis. Don't fix it here.

```text
/debug the failing test. Give me your top hypotheses for the root cause. Do not fix it yet.
```

```text
/rename-chat debug-failing-test
```

## Demo 2: Steer, review, approve, rewind

### Make a small, safe change

**Open a new agent.** Agent mode, Auto.

```text
/rename-chat Steering demo
```

```text
Add a second example to the doc comment on FormatUSD in internal/billing/money.go showing a negative amount. Do not change any logic.
```

### Review the diff

Open the diff. Show it's one comment and no logic changed.

### Run mode

Show Settings > Agents > Executions & Approvals > Run Mode. Explain approvals and allowlists.

### Approve, then rewind

Keep the change. Click **restore checkpoint** on the prompt. Show the comment is gone from `money.go`.

## Demo 3: Same fix, fewer tokens

### Show the bug

Open http://localhost:43173/disputes/dsp_1043. It shows **$400.00** against a $249 invoice. Background: the fix is `SUGGESTED_CREDIT_API_VERSION` in `internal/billing/suggested_credit.go`.

### Explain context cost

Show **Context Usage** on the debug-failing-test chat. Explain that every turn re-sends the whole chat.

### One specific prompt, fresh agent

**Open a new agent.** Agent mode, Auto. Explain referencing the Debug chat, and naming the file, the constraint and what done looks like.

```text
@debug-failing-test Fix the failing test by switching the suggested-credit client to the v2 API in internal/billing/suggested_credit.go. Keep both /api/v1 and /api/v2 routes. Do not edit the test. Done when go test ./... passes. Continue to the end, do not wait for my approval.
```

### Verify

Reload dsp_1043, which now shows **$249.00**. Show the one-line diff and the Context Usage.

### Habits

Explain starting fresh at around 60% context, summarisation near the limit, and one agent per piece of work.

### Reset

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```

## Demo 4a: Rules

### Show the existing rule

Open `.cursor/rules/ledgerly.mdc`. Explain project rules and always-apply.

### Baseline, without the new rule

**Open a new agent.** Ask mode.

```text
/ask How would you add a 2.5% late fee to an invoice total in Ledgerly? Show the Go code you would write. Do not edit files.
```

### Create the rule

**Open a new agent.** Show the new file `.cursor/rules/money-rounding.mdc` and its apply type.

```text
/create-rule Money in Ledgerly is int cents. Percentages are basis points (2.5% = 250 bps), applied with integer maths and an explicit round-half-up helper. Never use float64 for currency.
```

### Test the rule

**Open a new agent.** Ask mode. Compare with the baseline: basis points and integer rounding, not `float64`.

```text
/ask How would you add a 2.5% late fee to an invoice total in Ledgerly? Show the Go code you would write. Do not edit files.
```

```text
Which rules and skills did you apply to that answer? Name each file.
```

## Demo 4b: Skills

### Create a skill

**Open a new agent.** Explain what a skill is. Show `.cursor/skills/release-note/SKILL.md`.

```text
/create-skill Create a project skill called release-note. It writes a customer-facing release note for a Ledgerly change: a one-line title, what changed, who it affects, and no internal file names or code. Keep it under 80 words.
```

### Use the skill

**Open a new agent.** Show the output follows the skill's shape.

```text
/release-note Suggested credit on disputes now uses the v2 API and is capped at the customer's plan price.
```

```text
Which rules and skills did you apply to that answer? Name each file.
```

### What loads before you type

Open the **Context Usage** breakdown. Explain what's loaded up front: rules, skill descriptions, MCP tools.

## After the session

### Clean up

Run the reset. Delete `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and any saved plan file.

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```
