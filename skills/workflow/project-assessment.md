---
name: workflow-project-assessment
description: Use FIRST on an existing Go project — point go-skills at a codebase you already have and get a grounded health assessment + a prioritized improvement roadmap, where every finding cites the exact go-skills methodology chapter or pattern to apply. A team of five (senior Go engineer, strong-middle Go engineer, middle-strong DB architect, test/reliability engineer, senior SecOps) assesses independently, the senior Go engineer consolidates one maturity report, the same team runs a reverse analysis (cull generic nitpicks, check coverage), and the roadmap goes to the user for approval — then improvements are executed via the Feature Development and Security Review flows. Read-only; it never modifies your code.
category: workflow
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (the conformance baseline)
  - skills/methodology/00-canonical-full.md
related:
  - skills/workflow/feature-development.md
  - skills/workflow/security-review.md
  - skills/methodology/01-principles-and-layout.md
  - skills/methodology/12-testing.md
---

# Project Assessment

## Intent

Point go-skills at an **existing** Go codebase and answer "where does this project stand, and what
should we fix first?" — as a grounded **maturity assessment + a prioritized roadmap** in which every
finding names the exact go-skills [methodology](METHODOLOGY.md) chapter or [pattern](PATTERNS.md) to
apply. It is the **entry point** to the whole catalog: instead of reading 14 chapters and 52 patterns
and guessing what's relevant, you get a personalized, ranked action list that links straight to the
skills you actually need.

It completes the lifecycle the other two workflow skills cover:

```
   ┌─────────────────────┐   ┌──────────────────────┐   ┌────────────────────┐
   │  Project Assessment  │ → │ Feature Development   │   │  Security Review   │
   │  understand what you │   │ build what's next     │   │  keep it safe      │
   │  already have  (HERE)│   │                       │   │                    │
   └─────────────────────┘   └──────────────────────┘   └────────────────────┘
            the roadmap from HERE hands off to the other two for execution
```

It is **read-only** — it assesses and recommends; it never edits your code. Fixes are executed,
after you approve the roadmap, via [Feature Development Flow](skills/workflow/feature-development.md)
(refactors/features) and [Security Code Review](skills/workflow/security-review.md) (the security items).

```
                 ┌──────────────── the codebase + your goals ──────────────────┐
          ▼            ▼              ▼                ▼                ▼
   Senior Go     Strong-Middle Go   DB architect   Test/Reliability  SecOps (triage)
   (arch +       (idioms, re-       (data layer)   (tests, ops,      (top risks →
   methodology)  invented patterns)                deps, CI)          Security Review)
          │            │              │                │                │
          └─────┬──────┴──────┬───────┴────────┬───────┴────────┬───────┘
                ▼  PHASE 2: consolidate (senior owns it; one maturity report, scored)
          health report
                ▼  PHASE 3: reverse analysis (same team attacks the report)
        cull generic nitpicks · right-size priority · check coverage
                ▼  PHASE 4: senior folds → PRIORITIZED ROADMAP (each item → a go-skill)
         ───────  USER APPROVAL GATE  ───────
                ▼  hand off to Feature Development / Security Review to execute
```

## Context

A reference library is only useful if you know *which page applies to you*. A team dropping go-skills
onto an existing service doesn't need all 52 patterns — it needs the three that are reinvented badly
in *their* code, the methodology chapter their layering violates, and the one missing piece (graceful
shutdown, structured errors, a test seam) that's hurting them now. Reading the catalog cover-to-cover
to find that is the wrong order.

This flow inverts it: five lenses read the *actual* code and map it onto the catalog. **Independent
assessment** (Phase 1) widens what gets noticed across architecture, idioms, data, reliability, and
security. **Consolidation** (Phase 2) produces one scored maturity report instead of five overlapping
opinions. The **reverse analysis** (Phase 3) is the noise filter — the failure mode of any automated
project review is a flood of generic "add more tests / rename this / extract that" nitpicks, so the
team must cull anything not grounded in the code and a real go-skills recommendation, and check what
was *not* looked at. The **human gate** (Phase 4) keeps prioritization with the owner. Execution is
deliberately *handed off*, not done here — assessment and change are different jobs.

> Same honest caveat as the sibling flows: one model wearing five hats is not five independent minds;
> the win is lens-diversity + output isolation, hardened by the reverse pass and grounded in the real
> code. This complements `go vet`, `golangci-lint`, `govulncheck`, and human judgment — it does not
> replace them.

## The team (roles)

Five fixed lenses, chosen to cover the dimensions a Go service's health actually breaks along.

### 1. Senior Go Engineer — architecture & methodology conformance (report owner)

Owns the big picture and is the **consolidator** and **roadmap author**.

- Layering & dependency direction vs the [methodology](skills/methodology/02-layered-architecture.md);
  package boundaries, `cmd/` vs `internal/` layout, DI via constructors vs global state.
- `context.Context` propagation, graceful shutdown, configuration, bootstrap.
- Which architectural patterns are missing, present-but-misapplied, or reinvented.

### 2. Strong-Middle Go Engineer — idioms & code-level health

Owns "is this idiomatic Go," code-first.

- Reinvented wheels that map to a catalog pattern (a hand-rolled worker pool →
  [Bounded Parallelism](skills/concurrency/bounded-parallelism.md); ad-hoc construction →
  [Functional Options](skills/idiom/functional-options.md) / [Builder](skills/creational/builder.md)).
- Error wrapping, zero-value usability, goroutine lifecycle/leaks, premature abstraction, naming.
- Concrete smells with the specific skill that fixes each.

### 3. Middle-Strong DB Architect — data layer

Owns everything below the repository interface.

- Schema & access patterns, indexing, migrations (forward + rollback), transactions/isolation.
- Query construction (parameterization), N+1s, connection pooling, data lifecycle/retention.
- Maps to [Storage](skills/methodology/06-storage.md) and the data-relevant patterns.

### 4. Test/Reliability Engineer — tests, observability & operations

Owns "can you change this safely and run it in production."

- Test coverage **and design** (table-driven, race-safe, real seams vs mocks-everywhere) vs
  [Testing](skills/methodology/12-testing.md); flakiness, missing `-race`.
- Observability (structured [logging](skills/methodology/09-logging.md), metrics, tracing), graceful
  shutdown, health/readiness, timeouts/retries/[stability patterns](PATTERNS.md).
- Build/deploy/CI ([Build & Deploy](skills/methodology/13-build-and-deploy.md)), dependency health
  (`govulncheck`, stale/abandoned deps, module hygiene).

### 5. Senior SecOps — security triage (hands off, does not duplicate)

Owns a **shallow** first pass — top security risks and posture — then **delegates depth**.

- Surface-level: secrets in code/config, obvious injection/authz gaps, transport, dependency CVEs.
- Output is a triage verdict + the recommendation to run the full
  [Security Code Review](skills/workflow/security-review.md) when the risk warrants. It does **not**
  reproduce that flow's depth.

## The phases

### Phase 0 — Frame the assessment (the brief)

One brief, handed **identically** to all five roles:

- **The codebase**: repo path / packages in scope (and what's out — vendored code, generated files).
- **Goals & pain**: why now — onboarding, a rewrite decision, recurring incidents, scaling, an audit.
- **Constraints**: Go version, deploy target, team size, what *can't* change.
- **Depth**: a quick triage vs a deep audit, and any subsystem to prioritize.

### Phase 1 — Independent assessment (fan-out)

Each role reads the **actual code** (read-only) and returns, for its dimension: a **maturity grade**
(see below), concrete findings, and for *each finding the specific go-skills chapter/pattern that
applies* — a finding with no catalog mapping and no code evidence is noise and is dropped. No role
sees another's assessment. Where the module builds, each role also **runs the deterministic tooling
and folds it in** — `go build`, `go vet`, `go test`, `govulncheck` — so grades rest on real signal,
not eyeballing; a dependency-CVE or "dependency health" finding must cite a `govulncheck` result.

### Phase 2 — Consolidation (senior Go)

The senior merges the five assessments into **one** maturity report: deduplicate overlapping findings,
reconcile grade disagreements explicitly (record the decided grade + rationale), compute an **overall
maturity** and per-dimension grades, and assign each finding an **owner** and a **priority** (P0–P3
by impact × effort). Cross-cutting "this keeps biting you" themes are called out, not just point issues.

### Phase 3 — Reverse analysis (the same team)

The team attacks the report. Each role returns a verdict (approve / approve-with-changes / block) plus:

1. **Nitpick cull** — remove generic advice not grounded in *this* code or not tied to a real
   go-skills recommendation. This is the noise filter; an assessment that's 80% boilerplate gets ignored.
2. **Priority/grade right-sizing** — is each item's impact and effort honest? Demote the cosmetic,
   promote the load-bearing.
3. **Consolidation loss** — real findings from their Phase-1 assessment the merge dropped.
4. **Coverage** — name any package, layer, or concern in scope that was *not* actually examined; a
   gap here is itself a finding ("we didn't look at X").

Crucially the reverse pass **re-opens the cited code** — it verifies each kept-or-culled finding's
evidence against the real source, not just the report text. A finding whose evidence can't be
confirmed in the code is culled; this is how confident fabrications die.

A `block` means the report is not done (e.g. a P0 is unjustified, or coverage missed a core package).

### Phase 4 — Prioritized roadmap + approval gate

The senior folds the reverse analysis into the **roadmap**: an ordered list (P0 first) where each item
has the problem, the **go-skill to apply**, an owner, an effort estimate, the expected payoff, and the
**execution route** — which workflow runs it: [Feature Development](skills/workflow/feature-development.md)
for refactors/features, [Security Review](skills/workflow/security-review.md) for the security items.
It states the overall maturity, the biggest risk, and a suggested first sprint. It goes to the **user
for approval** — prioritization is theirs. Nothing is changed here. Edits loop back via `amend`.

### Phase 5 — Hand off to execute (not done here)

After approval, each roadmap item is executed by the workflow it was routed to — assessment and change
are separate jobs, so this flow stops at the approved roadmap. (Composes with
`subagent-driven-development` once a feature/fix is in flight.)

## Maturity model

Grade each dimension and the project overall on a simple, defensible scale — by *evidence in the code*,
not vibes:

- **1 — Ad hoc**: no consistent approach; the dimension is a frequent source of bugs/risk.
- **2 — Emerging**: some good practice, applied inconsistently.
- **3 — Solid**: consistent, idiomatic, matches the methodology in the common case.
- **4 — Strong**: consistent + handles the hard cases (failure modes, edge cases, scale).
- **5 — Exemplary**: a reference others should copy.

Dimensions: architecture & layering · idioms & API design · data layer · concurrency & resource
lifecycle · error handling & validation · testing · observability & ops · build/deploy & deps ·
security posture. The roadmap targets the lowest grades with the highest impact first.

**Overall maturity is the minimum of the load-bearing dimensions** (architecture, data, testing,
security) — a single critical 1 is not averaged away by strengths elsewhere. Every grade must cite
concrete code evidence; a grade with no evidence is not a grade.

## When to use

- **First contact** with an existing Go project — onboarding, inheriting a service, a "should we
  rewrite?" decision, or pre-scaling/pre-audit triage.
- A periodic health check to measure drift from the methodology over time.
- Deciding *which* go-skills to apply next — this flow chooses them for you, grounded in your code.

## When NOT to use

- Greenfield with no code yet — use the [Feature Development Flow](skills/workflow/feature-development.md)
  and bake quality in as you build.
- A specific known task or a single security concern — go straight to the relevant workflow.
- A pure mechanical lint pass — run `golangci-lint`/`go vet`; this flow is for judgment, not style.

## Roadmap template

```markdown
# <Project> — Assessment & Roadmap

**Scope:** <packages reviewed>   **Goal:** <why now>   **Overall maturity: N/5**
**Owner:** Senior Go engineer

## Maturity by dimension
| Dimension | Grade | One-line justification (with code evidence) |
|-----------|-------|---------------------------------------------|
| Architecture & layering | 2/5 | transport imports the DB driver directly — §2 violation |
| Testing | 3/5 | table-driven + race in core; handlers untested |

## Roadmap (priority order)
| # | Priority | Finding | Apply go-skill | Owner | Effort | Run via |
|---|----------|---------|----------------|-------|--------|---------|
| 1 | P0 | hand-rolled unbounded worker pool (OOM risk) | [Bounded Parallelism] | go-mid | M | feature-development |
| 2 | P0 | string-built SQL in repo layer | [Security Review] | secops | M | security-review |

## Cross-cutting themes
<the "this keeps biting you" issues a roadmap of point-fixes would miss>

## Suggested first sprint
<the 3–5 items to do first and why>

## Coverage & open questions
<what was NOT examined; anything needing a human decision>
```

## Runnable orchestration

Ships as a runnable Claude Code workflow: [`workflows/project-assessment.mjs`](workflows/project-assessment.mjs).
Read-only. The `brief` may be a plain string or an object with
`{ codebase, goals, constraints, depth }`. Assess mode **refuses an empty brief**, and the script
tolerates `args` arriving as a JSON string.

```text
# 1. Assess and produce the roadmap (Phases 1–4).
Workflow({ scriptPath: "workflows/project-assessment.mjs",
           args: { mode: "assess", brief: { /* the Phase-0 brief — string or object */ } } })

# 2. If you re-prioritize or scope-change the roadmap, re-vet it.
Workflow({ scriptPath: "workflows/project-assessment.mjs",
           args: { mode: "amend", plan: { /* edited roadmap */ }, changeRequest: "what changed & why" } })

# 3. Execute approved items via the other workflows (not this one):
#    feature-development.mjs for refactors/features · security-review.mjs for security items.
```

Phase mapping: `parallel()` fan-outs for the independent assessment and the reverse analysis; a single
consolidator `agent()` for the merge and the roadmap. Execution is intentionally delegated — there is
no remediate mode here. See the script header for the exact contract.

## Gotchas

- **Generic-nitpick soup.** The #1 way a project assessment gets ignored is boilerplate advice. Every
  finding must be grounded in the actual code AND tied to a specific go-skill; the reverse pass culls
  the rest.
- **A report with no priorities.** Twenty equal findings is not a roadmap. Force P0–P3 by impact ×
  effort and name the first sprint.
- **Findings that don't link to a skill.** The whole value is "apply *this* page next." A finding with
  no catalog mapping is either out of scope or under-researched.
- **Duplicating the deep flows.** Security depth lives in [Security Review](skills/workflow/security-review.md);
  this flow triages and hands off. Don't re-run it here.
- **Assessing the whole monorepo at once.** Scope it; on a large codebase, sample the core packages and
  state in coverage what you did not read.

## See also

- workflows/project-assessment.mjs — the runnable orchestration of this flow.
- skills/workflow/feature-development.md — executes the build/refactor roadmap items.
- skills/workflow/security-review.md — executes the security roadmap items.
- METHODOLOGY.md / PATTERNS.md — the catalog every finding maps onto.
