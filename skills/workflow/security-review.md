---
name: workflow-security-review
description: Use when security-reviewing Go code — a service, a diff, an auth/session flow, or data-handling code. A team of four specialists (senior Go engineer, strong-middle Go engineer, middle-strong DB architect, senior SecOps) audits independently, the senior Go engineer consolidates one deduplicated findings report, the same team runs a reverse analysis (cull false positives, verify exploitability, surface missed attack chains), and a prioritized remediation plan goes to the user for approval — then the same team fixes the vulnerabilities. Use for security audits, pre-release hardening, and post-incident review; not for a one-line obvious fix.
category: workflow
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (transport, storage, validation, errors)
  - skills/methodology/00-canonical-full.md
related:
  - skills/workflow/feature-development.md
  - skills/methodology/07-http-transport.md
  - skills/methodology/10-validation.md
  - skills/methodology/11-error-handling.md
---

# Security Code Review

## Intent

Turn "is this code safe?" into a **prioritized, exploitability-checked remediation plan** —
produced by four independent security lenses and hardened by an adversarial reverse pass that
kills false positives and surfaces attack chains — *before* anyone starts patching, and then fix
the confirmed vulnerabilities with the same team.

It is the [Feature Development Flow](skills/workflow/feature-development.md) retargeted at
security: same independent-first → consolidate → reverse → human gate → execute shape, but the
unit of work is a **finding** (a vulnerability) instead of a task, and the reverse pass is tuned
to the two failure modes of security review — **noise** (false positives) and **blind spots**
(missed chains).

```
                 ┌──────────────── identical review brief ────────────────┐
                 ▼              ▼                    ▼                      ▼
        Senior Go        Strong-Middle Go     DB architect          Senior SecOps
        (Go security)    (impl security)      (data security)       (threat/ops security)
                 │              │                    │                      │
                 └──────┬───────┴──────────┬─────────┴──────────┬──────────┘
                        ▼  PHASE 2: consolidate (senior Go owns the merge; dedup + attack chains)
                  one findings report
                        ▼  PHASE 3: reverse analysis (same team attacks the report)
              cull false positives · verify exploitability · find missed chains
                        ▼  PHASE 4: senior folds findings → REMEDIATION PLAN (severity-ranked)
                 ───────  USER APPROVAL GATE  ───────   (nothing patched before sign-off)
                        ▼  PHASE 5: same team fixes each finding, peer-verified + regression-tested
                     hardened code
```

## Context

A single reviewer or a single SAST scanner fails security review in two opposite directions at
once. It **misses** whole classes of issues outside its lens (a Go-memory reviewer doesn't see a
broad SQL grant; a scanner doesn't see a broken authorization model), and it **drowns** real
findings in false positives until nobody trusts the report.

Four independent lenses (Phase 1) widen coverage across the classes where Go services actually get
breached. A single consolidator (Phase 2) produces one deduplicated, severity-ranked report and —
critically for security — looks for **attack chains**, where several individually-minor findings
compose into a critical one. The **reverse analysis** (Phase 3) is the noise filter and the
blind-spot check: the same specialists now try to *refute* each finding (kill false positives,
right-size severity, prove or disprove exploitability) and to find what the merge missed. The
**human gate** (Phase 4) keeps the risk-acceptance decision with the owner. And the **same team**
fixes the findings (Phase 5), so the people who understood the vulnerability close it.

> The same honest caveat as the feature flow applies: run by one model, the four roles are not
> four independent minds — they share priors, so some blind spots correlate. The value is
> lens-diversity (different prompts surface different vulnerability classes) plus output isolation
> (no role anchors on another's findings first). Lean on the reverse pass, real exploit reasoning
> grounded in the actual code, and the human gate to catch what shared priors miss. This flow
> complements — it does not replace — `govulncheck`, `go vet`, fuzzing, and a real pentest.

## The team (roles)

Four fixed roles. Their lenses overlap at the edges on purpose — that overlap is what makes the
reverse analysis sharp.

### 1. Senior Go Engineer — Go-language security (and report owner)

Owns Go-specific vulnerability classes and is the **consolidator** and **plan owner**.

- crypto misuse (`math/rand` where `crypto/rand` is required, non-constant-time comparison of
  secrets, weak/again-used nonces, home-grown crypto), TLS config, JWT/session handling.
- injection through Go APIs (`os/exec` with a shell, `text/template` vs `html/template`,
  `database/sql` string-built queries), path traversal, SSRF, unsafe deserialization.
- concurrency bugs that are security bugs (data races on auth state, TOCTOU), `unsafe`,
  integer overflow, resource exhaustion / unbounded goroutines & allocations (DoS).
- secrets leaking through errors/logs/panics; error wrapping that exposes internals.

### 2. Strong-Middle Go Engineer — implementation-level security

Owns "is each boundary actually checked," concrete and code-first.

- per-endpoint/handler authentication **and** authorization checks (the missing `if !authorized`).
- input validation and output encoding at every trust boundary; mass-assignment / overposting.
- HTTP hardening: security headers, cookie flags (`HttpOnly`/`Secure`/`SameSite`), CORS, CSRF,
  rate limiting, request size limits, timeouts.
- dependency risk (`govulncheck`), error responses that leak stack traces / internal detail.

### 3. Middle-Strong DB Architect — data-layer security

Owns everything below the repository interface.

- SQL/NoSQL injection (parameterization, not string concatenation), ORM query injection.
- least privilege (DB roles/grants), row-level security, tenant isolation, `IDOR` at the query.
- PII/secrets at rest (encryption, column-level), data retention, audit logging of sensitive access.
- migrations that widen exposure, connection-string/credential handling, backups.

### 4. Senior SecOps — threat, architecture & operational security

Owns the system view and the attacker's view.

- threat model & trust boundaries; authn/authz **architecture** (not just per-handler checks).
- secrets management (env/vault, rotation), transport security, supply-chain & build integrity.
- container/deploy hardening, least privilege at the infra edge, egress control.
- detection: security logging, monitoring, alerting, and incident-response readiness; blast radius.
- maps findings to **OWASP Top 10 / CWE** and judges real-world exploitability and likelihood.

## The phases

### Phase 0 — Frame the review (shared brief)

One brief, handed **identically** to all four roles. Hold it to a bar — a vague scope produces a
vague audit:

- **Scope**: exactly what is under review (paths/packages/the diff), and what is explicitly out.
- **Threat model**: who the attacker is (anonymous internet, authenticated tenant, insider) and
  what they're after.
- **Assets & trust boundaries**: what must be protected (PII, money, credentials), where untrusted
  input crosses into trusted code.
- **Context**: how the code is deployed/exposed, the relevant data model, and any compliance bar
  (OWASP ASVS level, PCI, GDPR) findings must be judged against.

### Phase 1 — Independent security analysis (fan-out)

Each role audits *from the brief alone* and **grounds findings in the actual code** (reads the real
files — an audit of an imagined codebase is theatre). **No role sees another's findings.** Each
returns vulnerabilities with: title, **severity** (Critical/High/Medium/Low/Info) + **CWE**, exact
location, a concrete **exploit/impact**, a proposed remediation, and a **confidence** level. False
positives are expensive downstream, so each finding states *why it is exploitable here*, not just
that a risky API appears. Where the code compiles, **run the tooling and fold it in** — `govulncheck`
(dependency CVEs + reachability) and `go vet` — so the audit is grounded in real signal, not only
reading. And hunt **negative space** as hard as present bugs: a *missing* authorization check, a
*missing* rate limit, *missing* input validation. Absences are where real breaches live and are the
easiest thing for an independent reader to miss.

### Phase 2 — Consolidation (senior Go)

The senior merges the four reports into **one** findings report:

- **Deduplicate** the same vulnerability reported by multiple lenses into one finding (keep the
  highest-confidence write-up; record which lenses raised it).
- **Reconcile severity conflicts explicitly** (record the disagreement and the decided severity
  with rationale — never silently drop or downgrade).
- **Surface attack chains**: link findings that, combined, raise the effective severity (e.g. an
  open redirect + a token in a URL → account takeover). A chain is itself a finding.
- Output a single severity-ranked report; assign each finding an **owner** role for the eventual fix.

### Phase 3 — Reverse analysis (the same team)

"Reverse" means the team now *attacks the report itself*. Each role returns a verdict
(**approve / approve-with-changes / block**) plus:

1. **False positives** — findings that are not actually reachable/exploitable in this code; argue
   for removal or downgrade. This is the noise filter; a security report nobody trusts is useless.
2. **Severity adjustments** — right-size over- and under-rated findings against the threat model.
3. **Consolidation loss** — real findings from their Phase-1 report that the merge dropped.
4. **Missed attack chains / blind spots** — exploit paths visible only when findings are combined.
5. **Exploitability check** — for each High/Critical, state whether a realistic exploit path exists
   given the trust boundaries; an unexploitable "critical" is downgraded.
6. **Coverage** — name any endpoint, query/sink, or trust boundary in scope that was *not* examined.
   This pass culls false positives, so it must consciously guard against the opposite bias —
   under-reporting. A gap here is a finding: "we didn't look at X."

A `block` means the report is not done (e.g. a Critical is unverified, or a likely false positive
is still rated Critical) — the senior revises and the team re-reviews. Default to skepticism on
*existence* (does the bug reproduce?) and on *severity* (is it really that bad here?).

### Phase 4 — Remediation plan + approval gate

The senior folds the reverse analysis into the **remediation plan**: confirmed findings only, each
with severity/CWE, impact, the fix, an owner, dependencies, and a **verification** step (how we'll
prove it's closed). The plan is ordered by risk (Critical first) and states **residual risk** and
any **accepted risks**. It goes to the **user for approval** — the risk-acceptance and
prioritization call is theirs. Nothing is patched before sign-off. User edits loop back via
`amend` (a severity change or a new finding forces a fresh Phase-3 pass; a wording/priority tweak
needs only re-consolidation).

### Phase 5 — Remediation (the same team)

After approval, the team fixes findings **in priority order** (Critical → High → …):

- Each finding is routed to its **owner** role.
- After each fix, a **peer** role (not the fixer) verifies two things **independently** — by
  inspecting the diff and *re-running the test itself*, not trusting the fixer's self-report: the
  **vulnerability class is actually closed** (not just the reported payload masked) and **no
  regression / no new issue** was introduced. A failed verification sends the fix back for bounded
  rework; a fix that never passes is **flagged and escalated**, never silently shipped.
- **Verify before persist**: a fix is committed *only after* it passes independent verification; a
  rejected fix's working-tree changes are discarded, so an unverified (possibly-insecure) change
  never lands on the branch.
- Each fix ships with a **regression/security test** that fails on the vulnerable code and passes
  on the fix (proof, not assertion); `go test -race` for anything concurrency-related.
- A final whole-review pass confirms the fixes don't interact badly and the residual risk is
  acceptable.

This composes with `subagent-driven-development` and the [Feature Development Flow](skills/workflow/feature-development.md).

## Severity & CWE

Rate every finding **Critical / High / Medium / Low / Info** by *exploitability × impact in this
threat model* (not by how scary the API name sounds), and tag a **CWE** (e.g. CWE-89 SQL
injection, CWE-79 XSS, CWE-200 info exposure, CWE-352 CSRF, CWE-862 missing authorization,
CWE-326 weak crypto). Severity drives Phase-4 ordering; a finding with no realistic exploit path in
scope is Info or dropped.

## When to use

- A security audit of a service, a sensitive subsystem (auth, payments, multi-tenant data), or a
  diff that touches a trust boundary.
- Pre-release hardening, or a periodic review of code handling untrusted input or secrets.
- Post-incident review: find the class of bug, not just the one that fired.

## When NOT to use

- A single obvious fix (a hard-coded secret to rotate, a missing `Secure` cookie flag) — just fix it.
- A pure dependency bump or a finding already triaged — run `govulncheck` and patch.
- Greenfield design with no code yet — use the [Feature Development Flow](skills/workflow/feature-development.md)
  and bake security in as you build.

## Remediation-plan template

```markdown
# <Target> — Security Remediation Plan

**Scope:** <what was reviewed>   **Threat model:** <attacker + assets>
**Owner:** Senior Go engineer

## Severity summary
Critical: N · High: N · Medium: N · Low: N · Info: N

## Findings (ordered by risk)
| ID | Severity | CWE | Title | Location | Owner | Verification |
|----|----------|-----|-------|----------|-------|--------------|
| F1 | Critical | CWE-89 | SQL injection in … | repo/user.go:NN | db-arch | test injects `' OR 1=1`, expect parameter-bound no-op |

For each: impact, the exploit, the fix, dependencies.

## Attack chains
<finding combinations that raise effective severity, and the combined rating>

## Residual & accepted risk
<what remains after remediation; risks the owner explicitly accepts and why>

## Test strategy
<security/regression tests proving each fix; what must pass before release>

## Open questions
<anything needing a human risk decision — surfaced, not hidden>
```

## Runnable orchestration

Ships as a runnable Claude Code workflow: [`workflows/security-review.mjs`](workflows/security-review.mjs).
Because a background workflow cannot pause for input, the approval gate splits the run across
invocations. The `brief` may be a plain string or an object with
`{ scope, threatModel, assets, compliance, codebase }`. Review mode **refuses an empty brief**, and
the script tolerates `args` arriving as a JSON string.

```text
# 1. Audit and produce the vetted remediation plan (Phases 1–4).
Workflow({ scriptPath: "workflows/security-review.mjs",
           args: { mode: "review", brief: { /* the Phase-0 brief — string or object */ } } })

# 2a. If you edit the plan (re-rate a finding, accept a risk), re-vet it before fixing.
Workflow({ scriptPath: "workflows/security-review.mjs",
           args: { mode: "amend", plan: { /* edited plan */ }, changeRequest: "what changed & why" } })

# 2b. After you approve, remediate (Phase 5) with the same team.
Workflow({ scriptPath: "workflows/security-review.mjs",
           args: { mode: "remediate", plan: { /* the approved plan */ } } })
```

Phase mapping: `parallel()` fan-outs for the independent analysis and the reverse analysis; a
single consolidator `agent()` for the merge and the remediation plan; remediation walks findings
**in priority order**, each through a gated fix→peer-verify rework loop. See the script header for
the exact contract.

## Gotchas

- **A report that's all noise.** False positives are the #1 way a security review gets ignored.
  The reverse pass MUST try to refute each finding and downgrade what isn't exploitable here.
- **A report that's all single findings.** The expensive misses are *chains*. Consolidation and
  the reverse pass must look at findings together, not just individually.
- **Severity theatre.** Rating by scary API name, not by exploitability in the actual threat model,
  inflates Critical counts and buries the real ones. Right-size against the brief.
- **Masking the symptom.** Phase-5 verification checks the *vulnerability* is closed, not that the
  one reported payload now fails. Add a test that fails on the vulnerable code.
- **Skipping the gate.** Risk acceptance and fix-ordering are the owner's call. Don't auto-patch.
- **Treating this as a substitute for tooling.** It complements `govulncheck`, fuzzing, and a
  human pentest; it does not replace them.

## See also

- workflows/security-review.mjs — the runnable orchestration of this flow.
- skills/workflow/feature-development.md — the same method, retargeted at building features.
- skills/methodology/10-validation.md — input validation at the boundary.
- skills/methodology/11-error-handling.md — errors that don't leak internals.
