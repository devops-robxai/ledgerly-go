# Booking.com: Grok Build for engineers

Thu 5 Nov 2026, 08:00 UK (16:00 Shanghai). Remote. Slides and demos, then Q&A. Topics: skills, model selection (including Auto), token optimisation, context window management. Always call the app **Ledgerly**. Port **43173**. Review every result.

## Pre-flight (tonight)

### Team and Auto

Sign in with the **cursor.sh account on the ADM/FE team**. Open the model picker, click **Auto**, and check that **Cost, Balance and Intelligence** all show. If Balance is missing, an admin has not turned on the extra Auto modes for this team, so explain the three modes from slide 10 instead.

### Start the app

**Open a new agent.** Local agent, not a Cloud Agent. It should run `go run ./cmd/server` and open http://localhost:43173 in the Browser view. Check that the dashboard, invoices and disputes pages load.

```text
Start the application.
```

```text
/rename-chat Start app
```

### Reset the planted state

Restores the v1 suggested-credit client and removes the invoice email form, so Plan and Agent can build it live in Demo 1. Afterwards, http://localhost:43173/disputes/dsp_1043 should show **$400.00** (wrong) and `TestSuggestedCreditUsesV2Cap` should fail.

```text
Run ./scripts/reset-demo-state.sh and then run go test ./internal/billing/ -run TestSuggestedCreditUsesV2Cap. Report the result only. Do not fix anything.
```

### Clean up after rehearsal

Delete anything you create in rehearsal so you can create it live: `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and any plan file Demo 1 saved. Keep `.cursor/rules/ledgerly.mdc`, the repo's own always-apply rule.

### Optional: stage Attempt 1 Long

**Open a new agent.** Agent mode, Auto. Send these three prompts one after another, letting each finish. Then name the chat, note its Context Usage reading, and run the reset again.

```text
Fix the dispute credit.
```

```text
That's not it. The credit still looks wrong on disputes. Look again.
```

```text
Still wrong. Check which API the credit comes from, fix it, and run the tests.
```

```text
/rename-chat Attempt 1 Long
```

### Screen

Light mode. Larger editor and chat font. Collapse the left panel. Mute notifications. Slides and Grok Build on the same shared screen.

## Morning of

### Final checks

Check the app still runs in the Browser view (re-run **Start the application.** in the Start app chat if not). Run the reset prompt again. Check no Demo 4 rule or skill is left over. Model picker on **Auto**, mode on **Ask**. Join at 07:55 and check with Thomas whether to record.

## Demo 1: The agent loop: Ask, Plan, Agent, Debug, and Auto

### Show Ledgerly

Show http://localhost:43173 in the Browser view. "This is Ledgerly, our demo app. It handles invoices and disputes."

### Replay the Start app chat

Open the **Start app** chat and scroll it. "I asked one thing: start the application. It read the project, made a plan, ran commands and checked the app loaded. Understand, plan, act, verify. Each step has a mode."

### Ask: understand

**Open a new agent.** Open the mode picker and show the four modes. Ask mode. "Ask is for understanding. It reads, but never edits." Point out that no diff appears.

```text
/ask What does Ledgerly do? Name the main packages for invoices, disputes and the HTTP layer in 5 bullets. Do not edit files.
```

### Plan: map the approach

**Open a new agent.** Plan mode. Show http://localhost:43173/invoices/inv_1048 first: the customer email in the Customer card is read-only. "Plan writes the steps before any code changes. I review it, and change my mind here, where it's cheap."

```text
/plan I want a new feature to update the customer email in the invoice detail customer card. Don’t implement email validation.
```

### Agent: build the plan

Click **Build** on the plan (or switch to Agent mode with shift-tab and ask it to implement the plan). While it runs, move on to the model picker. When it finishes, reload inv_1048 and show the email edit in the Customer card.

### Pick a model

Open the model picker. "Start with your default. Move up only if it misses. Routine work, a faster model. Ambiguous work, a higher-effort model."

### Auto: Cost, Balance, Intelligence

Click **Auto**, then **Balance**. Cost keeps token spend lowest and bills like Grok or Composer. Balance weighs intelligence, speed and cost. Intelligence goes to the most capable models. Balance and Intelligence bill at the list price of the model they pick. No prices or benchmark ranks from memory.

### Debug: verify and investigate

**Open a new agent.** Debug mode. "There's a failing test in this repo. Debug mode investigates and gives me hypotheses on the root cause before anything changes." Read out the top hypothesis, and don't fix it here. "Hold that thought. We'll fix it properly in the tokens section."

```text
/debug the failing test. Give me your top hypotheses for the root cause. Do not fix it yet.
```

```text
/rename-chat debug-failing-test
```

### If it goes wrong

No Balance: "Your admin can turn these on for your team," and go back to slide 10. The build is slow: carry on with the model picker and come back to inv_1048 at the end. The build fails: show the plan and the diff, and say "that's what Debug is for." Debug starts fixing: stop it with the Stop button and say "I'll fix this properly later." Running long: skip the browser check of the email feature.

## Demo 2: Steer, review, approve, rewind

### Make a small, safe change

**Open a new agent.** Agent mode, Auto. Don't touch the suggested-credit seam here, because Demo 3 fixes it.

```text
/rename-chat Steering demo
```

```text
Add a second example to the doc comment on FormatUSD in internal/billing/money.go showing a negative amount. Do not change any logic.
```

### Review the diff

Open the diff. "Green is added, red is removed. It's one comment, no logic changed. Always read the diff before you approve."

### Run mode

Show the run mode setting (Settings > Agents > Executions & Approvals > Run Mode). "Start with the agent asking you. Loosen it once you trust it, for example with an allowlist for `go test`."

### Approve, then rewind

Keep the change. Then scroll up to the prompt and click the **restore checkpoint** icon next to it. Show `money.go` again, and the comment is gone. "Checkpoints undo an agent's change without touching git."

### If it goes wrong

It changes logic as well: "This is why we read the diff." Reject and rewind. No restore icon: hover over the prompt. Comment still there after rewinding: run the reset prompt before Demo 3. Running long: mention run mode aloud only.

## Demo 3: Same fix, fewer tokens

### Call back to Demo 1

"You've already seen the low-token pattern. In Demo 1 we asked, planned once and gave each job its own agent. That's slide 15. Now let's look at what it costs when you don't, and fix the bug Debug found with one specific prompt." Don't re-run Ask or Plan here.

### Show the bug

Open http://localhost:43173/disputes/dsp_1043. It shows **$400.00** suggested credit against a $249 Scale invoice. "This is the failing test Debug looked at." Background, not to read out: the client calls the v1 API (raw claim), v2 caps at the plan price, and the fix is `SUGGESTED_CREDIT_API_VERSION` in `internal/billing/suggested_credit.go`.

### The expensive way: Attempt 1 Long

Open the **Attempt 1 Long** chat you staged last night and point to its **Context Usage** below the chat. "Same bug, three vague turns. Every turn re-sent everything before it, so this number only went up." Not staged: point back to slide 14 instead, and say "each turn adds to the bar."

### The cheap way: one specific prompt, fresh agent

**Open a new agent.** Agent mode, Auto. "I'm pointing at the Debug chat, so it doesn't search again. I'm naming the file, the constraint and what done looks like. That's slide 16."

```text
@debug-failing-test Fix the failing test by switching the suggested-credit client to the v2 API in internal/billing/suggested_credit.go. Keep both /api/v1 and /api/v2 routes. Do not edit the test. Done when go test ./... passes. Continue to the end, do not wait for my approval.
```

### Verify and compare

Reload dsp_1043, which now shows **$249.00**. Read the diff, which should be one line. Then compare the Context Usage in this chat with Attempt 1 Long: "Same result, a fraction of the context."

### Habits

"Start fresh at around 60% context. That's a habit, not a setting. Near the limit, the chat is summarised and detail gets dropped. One agent per piece of work."

### Reset

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```

### If it goes wrong

@debug-failing-test doesn't resolve: pick the chat from the @ menu, or leave it out, because the prompt names the file anyway. Tests still fail: "That's the verify step doing its job." Use the prompt below, and if it's still red, reset and move on. Running long: skip the Attempt 1 Long comparison and use slide 14.

```text
Tests are failing. Read the failure, fix only the cause, and run go test ./... again.
```

## Demo 4a: Rules

### Show the existing rule

Open `.cursor/rules/ledgerly.mdc`. "Rules are .mdc files in the repo, so the whole team gets them. This one is always-apply, so it loads into every chat."

### Baseline, without the new rule

**Open a new agent.** Ask mode. "This is our 'without' run. Watch how it handles the percentage."

```text
/ask How would you add a 2.5% late fee to an invoice total in Ledgerly? Show the Go code you would write. Do not edit files.
```

### Create the rule

**Open a new agent.** Use /create-rule in the project chat, which makes a project rule. The New button under Customize > Rules makes a User rule, which is just for you. Open the new file in `.cursor/rules/` (expected: `money-rounding.mdc`) and show its apply type: always, by file pattern, or when relevant.

```text
/create-rule Money in Ledgerly is int cents. Percentages are basis points (2.5% = 250 bps), applied with integer maths and an explicit round-half-up helper. Never use float64 for currency.
```

### Test the rule

**Open a new agent.** Ask mode. Send the same baseline prompt. Compare the two: basis points and integer rounding, not `float64`.

```text
/ask How would you add a 2.5% late fee to an invoice total in Ledgerly? Show the Go code you would write. Do not edit files.
```

### Verify it was used

"Ask it, then check the code yourself."

```text
Which rules and skills did you apply to that answer? Name each file.
```

### If it goes wrong

The baseline already uses integer cents: "The code already shows the pattern. The rule makes it hold everywhere," and focus on the rounding. The rule isn't listed: set its apply type to always for the demo and re-run in a fresh agent.

## Demo 4b: Skills

### What a skill is

"A skill is a how-to: a name, a description, then the steps. Only the name and description load up front. The body loads when it's used."

### Create a skill

**Open a new agent.** Then show the new folder `.cursor/skills/release-note/` and open its `SKILL.md`.

```text
/create-skill Create a project skill called release-note. It writes a customer-facing release note for a Ledgerly change: a one-line title, what changed, who it affects, and no internal file names or code. Keep it under 80 words.
```

### Use the skill

**Open a new agent.** Point out the title, the sections and that there are no file names. That's the skill's shape.

```text
/release-note Suggested credit on disputes now uses the v2 API and is capped at the customer's plan price.
```

```text
Which rules and skills did you apply to that answer? Name each file.
```

### What loads before you type

Open the **Context Usage** breakdown. "Before I type anything, some context is already loaded: AGENTS.md, always-apply rules, skill names and descriptions, and MCP tool names. Keep those short, because they're in every chat."

```text
/ask Which instructions in our rules could a Go test or lint check enforce instead? Name up to 3, and for each, say whether a test already covers it. Do not edit files.
```

### If it goes wrong

/release-note isn't recognised: use the prompt below. Running long: skip the last Ask prompt.

```text
Use the release-note skill to write a release note: Suggested credit on disputes now uses the v2 API and is capped at the customer's plan price.
```

## After the session

### Clean up

Run the reset prompt. Delete `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and any plan file Demo 1 saved. Send Thomas the resources from slide 20 and any open questions.

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```
