---
name: workflow-feature-development
description: Use when turning a feature request or non-trivial task into a vetted, ready-to-build plan in Go. A team of three specialists (strong-middle Go engineer, senior Go engineer, middle-strong DB architect) analyze independently, the senior consolidates, the same trio runs a reverse (adversarial) review, the final implementation plan goes to the user for approval, then the same trio implements it. Use for features that touch logic + data + API, schema/migration changes, or cross-cutting work — not for one-line fixes.
category: workflow
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (layering, DI, storage, testing)
  - skills/methodology/00-canonical-full.md
related:
  - skills/methodology/02-layered-architecture.md
  - skills/methodology/06-storage.md
  - skills/methodology/12-testing.md
---

# Feature Development Flow

## Intent

Turn a feature request or a non-trivial task into an implementation plan that has already
survived three independent expert viewpoints and an adversarial review — *before* a single
line of production code is written — and then ship it with the same team that vetted it.

The flow simulates a small, senior engineering pod and runs it as a deterministic pipeline:

```
                 ┌──────────────── identical brief ────────────────┐
                 ▼                      ▼                            ▼
        Strong-Middle Go        Senior Go engineer          DB architect
        engineer (impl)         (architecture)              (data model)
                 │                      │                            │
                 └──────────┬───────────┴────────────┬──────────────┘
                            ▼  PHASE 2: consolidate (senior owns the merge)
                       draft plan
                            ▼  PHASE 3: reverse review (same trio attacks the draft)
                  approve / change / block  +  cross-domain checks
                            ▼  PHASE 4: senior folds findings → FINAL PLAN
                   ───────  USER APPROVAL GATE  ───────   (nothing built before sign-off)
                            ▼  PHASE 5: same trio implements, task-by-task, peer-reviewed
                        shipped feature
```

## Context

Most "design" steps fail in one of two ways. A **committee-first** process (everyone in one
room from minute zero) collapses into the loudest voice or the first idea on the whiteboard —
anchoring and groupthink quietly delete most of the solution space. A **single-author**
process is coherent but inherits one person's blind spots, especially across the
logic/data boundary where Go service bugs concentrate.

This flow fixes both. **Independent-first** (Phase 1) forces three different lenses to explore
the problem in isolation, so divergence — and the blind spots it exposes — surfaces before any
merge. A **single consolidator** (Phase 2) then keeps the plan coherent and *owned* rather than
a lowest-common-denominator compromise. A **reverse review** (Phase 3) specifically hunts the
two failure modes of any synthesis: good ideas dropped during the merge, and new defects
introduced by it. The **human gate** (Phase 4) keeps authority with the user. And running the
**same trio** through implementation (Phase 5) means the people who vetted the plan own its
execution — no handoff loss.

**A note on simulated independence.** When this flow is run by a single model wearing three
hats, the roles are not three independent minds — they share weights and priors, so some blind
spots are correlated and "independence" is weaker than it sounds. What the structure still buys
you is real and worth the ceremony: *distinct lenses* (different prompts pull genuinely different
concerns to the front) and *output isolation* (no role anchors on another's text before forming
its own view). Treat the gain as lens-diversity plus anti-anchoring, not as true cognitive
independence — and lean on the reverse review and the human gate to catch what shared priors miss.

## The trio (roles)

The three roles are fixed. Each has one dominant lens; together they cover the surface area
where Go service features actually break.

### 1. Strong-Middle Go Engineer — the implementation lens

Owns "how do we actually build this in *this* codebase." Pragmatic, code-first.

- Concrete package/file layout; which existing types and interfaces to reuse or extend.
- Idiomatic Go: `context.Context` propagation, error wrapping (`%w`), zero-value usability,
  goroutine lifecycle and cancellation, avoiding premature abstraction.
- Edge cases, input validation, partial-failure behavior, and what the tests must cover.
- A realistic task breakdown with effort signal (S/M/L) and the riskiest unknown called out.

### 2. Senior Go Engineer — the architecture lens (and plan owner)

Owns the shape of the change and is the **consolidator** and **final-plan author**.

- API/interface design, package boundaries, dependency direction, backward compatibility.
- Trade-offs and the decisions behind them; alignment with the service methodology
  (layering, DI via constructors, no global state).
- Cross-cutting concerns: observability, security, rollout/rollback, failure modes, SLOs.
- Decomposition into ordered tasks with owners, dependencies, and acceptance criteria.

### 3. Middle-Strong DB Architect — the data lens

Owns everything below the repository interface.

- Data model: tables/collections, normalization vs. deliberate denormalization, keys.
- Migrations: forward **and** rollback; online/lock-safe strategy for large tables.
- Indexes matched to the actual access patterns; query shapes and their cost.
- Transactions, isolation level, contention/locking, consistency guarantees, retention/volume.

> The roles are deliberately *overlapping at the edges* — the Go engineer cares about how the
> data layer feels from the caller's side; the DB architect cares whether the Go plan's access
> patterns are serviceable. That overlap is what makes the reverse review (Phase 3) sharp.

## The phases

### Phase 0 — Frame the task (shared brief)

Write one brief and hand it **identically** to all three roles. Identical input is what makes
their outputs independent and therefore comparable. Garbage-in dominates this entire pipeline —
every later phase reasons *from the brief alone* — so hold it to a bar:

- The request in one or two sentences, plus the *why* (user/business outcome).
- **Non-goals** — what this change explicitly will *not* do (the cheapest scope-control you have).
- Hard constraints with **concrete numbers**: compatibility, deadlines, data volume,
  latency/SLO (e.g. "p99 < 50ms", not "fast"), security/compliance.
- **Testable acceptance criteria** — prefer Given/When/Then, not adjectives.
- The relevant slice of the codebase by **exact paths/symbols/current schema** — enough to reason
  without guessing, no more.

**Convergence pre-check.** A vague brief produces three confidently-different misreadings that
consolidation then launders into a decisive-looking (wrong) plan. So after Phase 1, compare the
three `understanding` statements: if they disagree on *what is being built*, the brief failed —
fix it and re-run Phase 1 before consolidating. (The orchestration surfaces this as a
`convergenceCheck` field on the draft.)

### Phase 1 — Independent analysis (fan-out)

Each role produces a structured proposal *from the brief alone*. **No role sees another's
output.** Each proposal contains: problem understanding, the role's domain plan, risks, open
questions, and a domain task breakdown. Divergence here is a feature — it is the cheapest
blind-spot detector you have.

**Ground it in the real codebase.** The implementation and data roles must read the actual files
the brief references (read-only) rather than reason about an imagined codebase — an
"implementation lens" that never looks at the code is fiction. They propose against what exists:
real types, real schema, real call sites.

### Phase 2 — Consolidation (senior Go)

The senior merges the three proposals into **one** draft plan:

- **Reconcile conflicts explicitly.** When the DB architect wants a new table and the Go
  engineer wants to reuse an existing one, the senior decides and records the rationale.
  Conflicts are resolved on the record, never silently dropped.
- **De-duplicate and fill gaps** the individual lenses missed at their seams.
- Produce a single **ordered task list**: each task has an owner role, dependencies, and
  acceptance criteria.

### Phase 3 — Reverse review (the same trio)

"Reverse" means the pipeline runs backward: instead of *producing*, each specialist now
*attacks* the consolidated draft. Each role returns a verdict — **approve / approve-with-changes
/ block** — plus concrete, actionable findings, looking for:

1. **Consolidation loss** — anything from their Phase-1 proposal that the merge dropped or distorted.
2. **Merge-introduced defects** — problems that exist only because the parts were combined.
3. **Cross-domain check** — each role audits one neighbor's area:
   - Go engineer ⇒ is the data plan ergonomic and correct from the calling code?
   - DB architect ⇒ are the Go plan's access patterns actually serviceable at the storage layer?
   - Senior ⇒ do both still satisfy the methodology and the acceptance criteria?

A single `block` verdict means the draft is not done — the senior addresses it and the trio
re-reviews. Default to skepticism: if a reviewer is unsure a concern is real, they raise it.

### Phase 4 — Final plan + approval gate

The senior folds the reverse-review findings into the **final implementation plan** — the
canonical artifact (see template below). This plan is presented to the **user for approval**.
**Nothing is implemented before explicit sign-off.** If the user requests changes, the size of
the change decides the loop-back: edits to the **data model, public API, or task dependencies**
force a fresh **Phase 3 re-review** (they can introduce new cross-domain defects); cosmetic or
scoping-only edits need only a **Phase 2 re-consolidate**. The orchestration exposes this as an
`amend` mode so an edited plan is actually re-vetted rather than silently hand-patched.

### Phase 5 — Implementation (the same trio)

After approval, the trio implements the plan **task by task**, in plan order so dependencies are
honored:

- Each task is routed to its **owner role**.
- After each task: a two-stage peer review — **spec compliance first, then code quality** —
  performed by a *different* role than the implementer.
- **The review gates the task.** A failed spec or quality check is not just recorded — it sends
  the task back to its owner for bounded rework. A task that still fails after the rework budget
  is **flagged as not-done and escalated**, never silently shipped.
- **Tests are evidence, not assertion.** The implementer runs `go test -race` and captures the
  exact command and pass/fail output; the reviewer verifies against that captured evidence rather
  than trusting the prose summary. ("Race-safe" is only as wide as the tests you actually run —
  see `skills/methodology/12-testing.md` for the bar.)
- A final whole-implementation review before the feature is called done.

This phase composes directly with `subagent-driven-development`: one fresh implementer per task,
gated two-stage review between tasks.

## When to use

- The change touches **more than one layer** — typically logic **and** data **and** transport.
- Any **schema or migration** change, especially on large or hot tables.
- Cross-cutting work: a new bounded context, an integration with an external system, a change
  to a core domain invariant.
- High-stakes work where a wrong plan is expensive to unwind (billing, auth, data retention).

## When NOT to use

- One-line fixes, copy changes, dependency bumps, mechanical refactors — the overhead dwarfs
  the value. Just do them.
- Pure exploration/spikes where the goal is to *learn*, not to ship — use a single throwaway pass.
- Trivially-scoped tasks fully owned by one role with no data or API impact.

## Final-plan template

```markdown
# <Feature> — Implementation Plan

**Goal:** <one sentence — the user/business outcome>
**Owner:** Senior Go engineer

## Architecture
<the shape of the change: packages, interfaces, dependency direction, trade-offs decided>

## Data model & migrations
<tables/columns, indexes matched to access patterns, forward + rollback migration, lock safety>

## Tasks (ordered)
| # | Task | Owner | Depends on | Acceptance criteria |
|---|------|-------|-----------|---------------------|
| 1 | ... | go-mid | — | ... |
| 2 | ... | db-arch | 1 | ... |

## Risks & mitigations
<each open risk, its blast radius, and how we de-risk it>

## Rollout & rollback
<feature flag / phased rollout / how to revert safely, incl. the migration>

## Test strategy
<unit, integration, race; what must pass before merge>

## Open questions
<anything still needing a human decision — surfaced, not hidden>
```

## Runnable orchestration

This flow ships as a runnable Claude Code workflow: [`workflows/feature-development.mjs`](workflows/feature-development.mjs).
Because a background workflow cannot pause for input, the approval gate splits the run across
invocations. The `brief` may be a plain string or an object with
`{ request, why, constraints, acceptance, codebase }`. Plan mode **refuses an empty brief** —
without one the roles invent scope from the environment — and the script tolerates `args`
arriving as a JSON string.

```text
# 1. Produce the vetted plan (Phases 1–4). Returns the final plan for you to approve.
Workflow({ scriptPath: "workflows/feature-development.mjs",
           args: { mode: "plan", brief: { /* the Phase-0 brief — string or object */ } } })

# 2a. If you edit the returned plan, re-vet the edit before building (real Phase-4 loop-back).
Workflow({ scriptPath: "workflows/feature-development.mjs",
           args: { mode: "amend", plan: { /* edited plan */ }, changeRequest: "what changed & why" } })

# 2b. After you approve, implement it (Phase 5) with the same trio.
Workflow({ scriptPath: "workflows/feature-development.mjs",
           args: { mode: "implement", plan: { /* the approved plan */ } } })
```

The script maps each phase to a Workflow primitive: `parallel()` fan-outs for the independent
analysis and the reverse review; a single consolidator `agent()` for the merge and the final plan;
and, in implementation, tasks run **in plan order** (to honor `dependsOn`), each through a gated
implement→peer-review rework loop. See the script header for the exact contract.

## Gotchas

- **Letting the roles talk in Phase 1.** The moment they share context the independence — and
  the blind-spot detection — is gone. Keep the briefs identical and the analyses isolated.
- **Silent consolidation.** If the merge drops a Phase-1 idea without saying why, the reverse
  review can't catch it as *loss* — it just looks absent. Record every reconciled conflict.
- **A reverse review that only rubber-stamps.** Skepticism is a process requirement — each
  reviewer must actually *attempt to refute* the draft. But don't turn that into a quota of
  findings: manufacturing objections to look diligent is its own failure mode (performative
  review). The right check on an all-approve round of a non-trivial plan is *"did each reviewer
  cite a specific risk they considered and dismissed?"* — not an assumption that approval is
  suspect. A small, well-scoped plan can simply be correct.
- **Skipping the gate "to save time."** The gate is the cheapest place to catch a wrong plan.
  Implementing before sign-off trades a five-minute review for a multi-day unwind.
- **Using it for trivial work.** Three analyses + a review for a typo is theater. Match the
  ceremony to the stakes (see *When NOT to use*).

## See also

- workflows/feature-development.mjs — the runnable orchestration of this flow.
- skills/methodology/02-layered-architecture.md — the architecture the senior role aligns to.
- skills/methodology/06-storage.md — the data layer the DB architect role owns.
- skills/methodology/12-testing.md — the race-safe test bar Phase 5 holds to.
