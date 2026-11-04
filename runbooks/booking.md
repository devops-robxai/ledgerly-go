# Booking.com: Grok Build for engineers

Thu 5 Nov 2026, 08:00 UK (16:00 Shanghai). Remote. 47 min of slides and demos, then 10–12 min of Q&A. Topics: skills, model selection (including Auto), token optimisation, context window management. Always call the app **Ledgerly**. Port **43173**. Review every result.

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

Restores the v1 suggested-credit client. Afterwards, http://localhost:43173/disputes/dsp_1043 should show **$400.00** (wrong) and `TestSuggestedCreditUsesV2Cap` should fail.

```text
Run ./scripts/reset-demo-state.sh and then run go test ./internal/billing/ -run TestSuggestedCreditUsesV2Cap. Report the result only. Do not fix anything.
```

### Clean up after rehearsal

Delete anything you create in rehearsal so you can create it live: `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and any plan file such as `fix-dispute-cap.md`. Keep `.cursor/rules/ledgerly.mdc`, the repo's own always-apply rule.

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

## Morning of (07:30–07:55)

### Final checks

Check the app still runs in the Browser view (re-run **Start the application.** in the Start app chat if not). Run the reset prompt again. Check no Demo 4 rule or skill is left over. Model picker on **Auto**, mode on **Ask**. Join at 07:55 and check with Thomas whether to record.

## Slides 1–10 (7 min, 08:00–08:07)

### Present

Keep slides 1, 3 and 4 short. Spend the time on 5–10 (levers, the agent loop, modes, models and Auto). End on: "Let's see that loop in Grok Build."

## Demo 1: The agent loop, modes and Auto (6 min, 08:07–08:13)

### Show Ledgerly

Show http://localhost:43173 in the Browser view. "This is Ledgerly, our demo app. It handles invoices and disputes."

### Replay the Start app chat

Open the **Start app** chat and scroll it. "I asked one thing: start the application. It read the project, made a plan, ran commands and checked the app loaded. Understand, plan, act, verify."

### The four modes

**Open a new agent.** Open the mode picker. Ask is for understanding and doesn't edit files. Plan writes the steps, Agent does the work, and Debug investigates. You can also type `/ask`, `/plan` or `/debug` at the start of a prompt.

```text
/ask What does Ledgerly do? Name the main packages for invoices, disputes and the HTTP layer in 5 bullets. Do not edit files.
```

### Pick a model

Open the model picker. "Start with your default. Move up only if it misses. Routine work, a faster model. Ambiguous work, a higher-effort model."

### Auto: Cost, Balance, Intelligence

Click **Auto**, then **Balance**. Cost keeps token spend lowest and bills like Grok or Composer. Balance weighs intelligence, speed and cost. Intelligence goes to the most capable models. Balance and Intelligence bill at the list price of the model they pick. No prices or benchmark ranks from memory.

```text
/ask Which file in internal/billing would I open first to understand how disputes work? One sentence. Do not edit files.
```

### If it goes wrong

No Balance: "Your admin can turn these on for your team," and go back to slide 10. Ask is slow: keep talking through the picker. App not running: skip the Browser view and show the Start app chat only.

## Slides 11–12 (2 min, 08:13–08:15)

### Present

How do I work with an agent? Then steer, review, approve. End on: "Let me show you review, approve and rewind."

## Demo 2: Steer, review, approve, rewind (4 min, 08:15–08:19)

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

## Slides 13–16 (5 min, 08:19–08:24)

### Present

Same result, fewer tokens. Most time on slide 14 (every turn re-sends the whole chat). Then decide earlier, re-send less, and specific prompts. End on: "Let's fix a real bug the low-token way."

## Demo 3: Ask, Plan, build fresh (7 min, 08:24–08:31)

### Show the bug

Open http://localhost:43173/disputes/dsp_1043. It shows **$400.00** suggested credit against a $249 Scale invoice. "This is wrong. Let's find out why, without guessing." Background, not to read out: the client calls the v1 API (raw claim). v2 caps at the plan price. The fix is `SUGGESTED_CREDIT_API_VERSION` in `internal/billing/suggested_credit.go`.

### Ask (fresh chat)

**Open a new agent.** Ask mode. Then point to **Context Usage** below the chat: "Every turn re-sends the whole chat, so this only goes up."

```text
/ask Why does dispute dsp_1043 show $400.00 suggested credit instead of $249.00? Name the files and the API version in 3 bullets. Do not edit files.
```

```text
/rename-chat dispute-cap-ask
```

### Plan (fresh chat)

**Open a new agent.** Plan mode. "I'm pointing at the Ask chat so it doesn't search again. Reviewing the plan is the cheapest moment to change my mind." Save the plan as `fix-dispute-cap.md`, or use whatever name it offers in the next prompt.

```text
/plan @dispute-cap-ask Fix the suggested credit so the client uses the v2 API. Keep both /api/v1 and /api/v2 routes. Do not edit the failing test. Done when go test ./... passes. Stop after the plan for my review.
```

### Build (fresh chat)

**Open a new agent.** Agent mode. "Fresh agent, clear plan, clear definition of done."

```text
Implement @fix-dispute-cap.md. Change only the files it names. Done when go test ./... passes. Continue to the end, do not wait for my approval.
```

### Verify and compare

Reload dsp_1043 and it now shows **$249.00**. Read the diff. Point to Context Usage in the build chat: "Look how little context this used." If staged, open **Attempt 1 Long**: "Same bug, three vague turns, far more tokens." Habits: start fresh at around 60% context. Near the limit, the chat is summarised and detail is dropped.

### Reset

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```

### If it goes wrong

Ask doesn't name the API version: use the prompt below. @dispute-cap-ask doesn't resolve: pick it from the @ menu or paste the Ask bullets. Tests still fail: "That's the verify step doing its job." Use the fix prompt below, and if it's still red at 08:31, reset and move on.

```text
/ask Which API version does SuggestedCreditPath in internal/billing/suggested_credit.go use? One line. Do not edit files.
```

```text
Tests are failing. Read the failure, fix only the cause, and run go test ./... again.
```

## Slides 17–18 (2 min, 08:31–08:33)

### Present

How do I govern my agent? Rules and skills. End on: "This next part is hands-on. Open your own project if you'd like to follow along." Paste the viewer prompts from the Viewer hands-on section into the meeting chat now.

## Demo 4a: Rules (5 min, 08:33–08:38)

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

## Demo 4b: Skills (4 min, 08:38–08:42)

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

## Viewer hands-on (3 min, 08:42–08:45)

### Your turn

"Open your own project in Grok Build. The prompts are in the chat. Take a couple of minutes, and drop your rule in the chat if you'd like." Tip: run the same task with and without the rule and compare. Viewers without a project can clone Ledgerly from slide 3.

### Step 1: Create a rule

```text
/create-rule <One thing every change in this repo must do, e.g. "All API handlers validate input with our shared schema helper.">
```

### Step 2: Create a skill

```text
/create-skill Create a project skill called <skill-name>. It <describes one repeatable task in this repo>, following the patterns already in the codebase.
```

### Step 3: Test it in a fresh agent

```text
/ask <A task that should trigger your rule or skill>. Show the code you would write. Do not edit files.
```

### Step 4: Check it was used

```text
Which rules and skills did you apply to that answer? Name each file.
```

## Wrap up (2 min, 08:45–08:47)

### Slides 19–21

Try it, resources (the workshops link), thank you. Open Q&A on slide 21.

## After the session

### Clean up

Run the reset prompt. Delete `.cursor/rules/money-rounding.mdc`, `.cursor/skills/release-note/` and `fix-dispute-cap.md`. Send Thomas the resources from slide 20 and any open questions.

```text
Run ./scripts/reset-demo-state.sh. Report the result only.
```
