export const meta = {
  name: 'project-assessment',
  description: 'Team-of-five read-only assessment of an existing Go project against the go-skills catalog: independent assessment → consolidate (maturity report) → reverse analysis (cull nitpicks, check coverage) → prioritized roadmap (approval gate) → hand off to feature-development / security-review.',
  whenToUse: 'First contact with an existing Go codebase — onboarding, a rewrite decision, pre-scaling/pre-audit triage, or choosing which go-skills to apply next. Not for greenfield or a single known task.',
  phases: [
    { title: 'Frame' },
    { title: 'Independent assessment' },
    { title: 'Consolidate' },
    { title: 'Reverse analysis' },
    { title: 'Roadmap' },
  ],
}

// =============================================================================
// Project Assessment — runnable orchestration of skills/workflow/project-assessment.md
//
// NOTE: runs under the Claude Code Workflow runtime (injects agent()/parallel()/
// pipeline()/phase()/log()/args; wraps the body in an async function). NOT meant to
// parse via `node --check` standalone.
//
// READ-ONLY: this flow assesses and recommends; it never edits code. Execution of the
// approved roadmap is handed off to feature-development.mjs / security-review.mjs.
//
//   mode: "assess"  — Phases 1–4. Returns the maturity report + prioritized roadmap.
//        args: { mode: "assess", brief: <string | brief object> }
//
//   mode: "amend"   — re-vet an edited roadmap (re-prioritized, scope changed).
//        args: { mode: "amend", plan: <edited roadmap>, changeRequest: <what changed/why> }
//
// The brief (Phase 0) may be a plain string or an object with any of:
//   { codebase, goals, constraints, depth }
// =============================================================================

const MAX_REVIEW_ROUNDS = 2 // reverse-analysis re-runs before remaining blocks become open questions

// ---- The fixed team ---------------------------------------------------------
const ROLES = [
  {
    key: 'go-sr',
    title: 'Senior Go Engineer',
    lens:
      'Architecture & methodology conformance (and report owner). Layering and dependency direction ' +
      'vs the service methodology; cmd/ vs internal/ layout; DI via constructors vs global state; ' +
      'context.Context propagation; graceful shutdown; configuration & bootstrap. Which architectural ' +
      'patterns are missing, present-but-misapplied, or reinvented.',
  },
  {
    key: 'go-mid',
    title: 'Strong-Middle Go Engineer',
    lens:
      'Idioms & code-level health, code-first. Reinvented wheels that map to a catalog pattern ' +
      '(hand-rolled worker pool → Bounded Parallelism; ad-hoc construction → Functional Options / ' +
      'Builder); error wrapping, zero-value usability, goroutine lifecycle/leaks, premature ' +
      'abstraction, naming. Concrete smells, each tied to the specific go-skill that fixes it.',
  },
  {
    key: 'db-arch',
    title: 'Middle-Strong DB Architect',
    lens:
      'Data layer below the repository interface. Schema & access patterns, indexing, migrations ' +
      '(forward + rollback), transactions/isolation, query construction (parameterization), N+1s, ' +
      'connection pooling, data lifecycle/retention. Maps to the Storage chapter and data patterns.',
  },
  {
    key: 'test-rel',
    title: 'Test/Reliability Engineer',
    lens:
      'Tests, observability & operations — can you change this safely and run it in production. Test ' +
      'coverage AND design (table-driven, race-safe, real seams vs mocks-everywhere), flakiness, ' +
      'missing -race; observability (structured logging, metrics, tracing), graceful shutdown, ' +
      'health/readiness, timeouts/retries/stability patterns; build/deploy/CI and dependency health ' +
      '(govulncheck, stale/abandoned deps, module hygiene).',
  },
  {
    key: 'secops',
    title: 'Senior SecOps',
    lens:
      'Security TRIAGE only — a shallow first pass that hands off depth. Secrets in code/config, ' +
      'obvious injection/authz gaps, transport security, dependency CVEs, overall posture. Output is ' +
      'a triage verdict plus a recommendation to run the full Security Code Review when warranted. Do ' +
      'NOT reproduce that flow\'s depth here.',
  },
]
const roleByKey = (k) => ROLES.find((r) => r.key === k)
const senior = roleByKey('go-sr')

const PRIORITIES = ['P0', 'P1', 'P2', 'P3']
const PRIO_RANK = { P0: 0, P1: 1, P2: 2, P3: 3 }
const RUN_VIA = ['feature-development', 'security-review']

// ---- Schemas (structured agent output) --------------------------------------
const GRADE = {
  type: 'object',
  required: ['dimension', 'grade', 'justification'],
  properties: {
    dimension: { type: 'string' },
    grade: { type: 'integer', description: '1 ad-hoc, 2 emerging, 3 solid, 4 strong, 5 exemplary' },
    justification: { type: 'string', description: 'with code evidence' },
  },
}

const ASSESSMENT_SCHEMA = {
  type: 'object',
  required: ['summary', 'grades', 'findings'],
  properties: {
    role: { type: 'string', description: 'stamped by the orchestrator' },
    summary: { type: 'string' },
    grades: { type: 'array', items: GRADE },
    findings: {
      type: 'array',
      items: {
        type: 'object',
        required: ['title', 'location', 'evidence', 'recommendation', 'applySkill'],
        properties: {
          title: { type: 'string' },
          location: { type: 'string', description: 'file/package the finding is grounded in' },
          evidence: { type: 'string', description: 'what in the code shows this — not a generic claim' },
          recommendation: { type: 'string' },
          applySkill: { type: 'string', description: 'the go-skills chapter/pattern to apply' },
          priority: { type: 'string', enum: PRIORITIES },
          confidence: { type: 'string', enum: ['high', 'medium', 'low'] },
        },
      },
    },
    openQuestions: { type: 'array', items: { type: 'string' } },
  },
}

const ROADMAP_ITEM = {
  type: 'object',
  required: ['id', 'title', 'priority', 'applySkill', 'owner', 'runVia'],
  properties: {
    id: { type: 'string' },
    title: { type: 'string' },
    priority: { type: 'string', enum: PRIORITIES },
    location: { type: 'string' },
    evidence: { type: 'string' },
    applySkill: { type: 'string', description: 'the go-skills chapter/pattern to apply' },
    payoff: { type: 'string' },
    owner: { type: 'string', enum: ['go-sr', 'go-mid', 'db-arch', 'test-rel', 'secops'] },
    effort: { type: 'string', enum: ['S', 'M', 'L'] },
    runVia: { type: 'string', enum: RUN_VIA, description: 'which workflow executes this item' },
    sources: { type: 'array', items: { type: 'string' } },
  },
}

const REPORT_SCHEMA = {
  type: 'object',
  required: ['summary', 'overallMaturity', 'grades', 'conflictsResolved', 'findings'],
  properties: {
    summary: { type: 'string' },
    overallMaturity: { type: 'integer', description: '1–5; the MINIMUM of the load-bearing dimensions (architecture, data, testing, security) — do not average away a critical 1' },
    grades: { type: 'array', items: GRADE },
    crossCuttingThemes: { type: 'array', items: { type: 'string' } },
    conflictsResolved: {
      type: 'array',
      items: {
        type: 'object',
        required: ['conflict', 'decision', 'rationale'],
        properties: { conflict: { type: 'string' }, decision: { type: 'string' }, rationale: { type: 'string' } },
      },
    },
    findings: { type: 'array', items: ROADMAP_ITEM },
  },
}

const REVERSE_SCHEMA = {
  type: 'object',
  required: ['role', 'verdict', 'findings', 'crossDomainCheck'],
  properties: {
    role: { type: 'string' },
    verdict: { type: 'string', enum: ['approve', 'approve-with-changes', 'block'] },
    nitpicksToCull: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'why'],
        properties: { id: { type: 'string' }, why: { type: 'string' } },
      },
    },
    priorityAdjustments: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'from', 'to', 'why'],
        properties: { id: { type: 'string' }, from: { type: 'string' }, to: { type: 'string' }, why: { type: 'string' } },
      },
    },
    consolidationLoss: { type: 'array', items: { type: 'string' } },
    coverageGaps: { type: 'array', items: { type: 'string' }, description: 'packages/layers/concerns in scope not examined' },
    findings: {
      type: 'array',
      items: {
        type: 'object',
        required: ['severity', 'issue', 'fix'],
        properties: {
          severity: { type: 'string', enum: ['blocker', 'major', 'minor'] },
          issue: { type: 'string' },
          fix: { type: 'string' },
        },
      },
    },
    crossDomainCheck: { type: 'string' },
  },
}

const PLAN_SCHEMA = {
  type: 'object',
  required: ['title', 'scope', 'overallMaturity', 'grades', 'findings', 'firstSprint', 'coverage', 'openQuestions'],
  properties: {
    title: { type: 'string' },
    scope: { type: 'string' },
    goal: { type: 'string' },
    overallMaturity: { type: 'integer', description: '1–5; the MINIMUM of the load-bearing dimensions — do not average away a critical 1' },
    grades: { type: 'array', items: GRADE },
    crossCuttingThemes: { type: 'array', items: { type: 'string' } },
    findings: { type: 'array', items: ROADMAP_ITEM },
    firstSprint: { type: 'array', items: { type: 'string' }, description: 'the 3–5 items to do first' },
    coverage: { type: 'string', description: 'what was and was NOT examined' },
    openQuestions: { type: 'array', items: { type: 'string' } },
  },
}

// ---- Helpers ----------------------------------------------------------------
function briefText(b) {
  if (b == null) return '(no brief provided)'
  if (typeof b === 'string') return b
  const parts = []
  if (b.codebase) parts.push(`Codebase / scope: ${JSON.stringify(b.codebase)}`)
  if (b.goals) parts.push(`Goals & pain: ${JSON.stringify(b.goals)}`)
  if (b.constraints) parts.push(`Constraints: ${JSON.stringify(b.constraints)}`)
  if (b.depth) parts.push(`Depth: ${JSON.stringify(b.depth)}`)
  return parts.length ? parts.join('\n') : JSON.stringify(b)
}

const J = (x) => JSON.stringify(x, null, 2)

function requirePlan(plan, mode) {
  if (!plan || !Array.isArray(plan.findings) || !plan.findings.length) {
    throw new Error(`${mode} mode requires args.plan with a non-empty findings[] array (the roadmap).`)
  }
}

// =============================================================================
// MODE: assess — Phases 1–4, returns the maturity report + roadmap (approval gate)
// =============================================================================
async function runAssess(brief) {
  phase('Frame')
  const BRIEF = briefText(brief)
  log('Framed assessment brief; handing the identical brief to all five roles.')

  // --- Phase 1: independent assessment (fan-out; no role sees another's)
  phase('Independent assessment')
  const assessments = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are a ${r.title} assessing an EXISTING Go project against the go-skills catalog ` +
            `(the service methodology + 52 patterns).\n${r.lens}\n\n` +
            `Assess INDEPENDENTLY — you cannot see your teammates' assessments. READ THE ACTUAL CODE ` +
            `(read-only) the brief points at; an assessment of an imagined codebase is worthless. Where ` +
            `the module builds, also RUN the deterministic tooling and fold its output in — ` +
            `\`go build ./...\`, \`go vet ./...\`, \`go test ./...\`, and \`govulncheck ./...\`; any ` +
            `dependency-CVE or "dependency health" finding MUST cite a govulncheck result, not memory. ` +
            `For your dimension(s) give a MATURITY GRADE (1 ad-hoc … 5 exemplary) with code-grounded ` +
            `justification, and concrete findings. EVERY finding MUST (a) cite real code evidence ` +
            `(file/package) and (b) name the specific go-skills chapter or pattern to apply — a finding ` +
            `with no code evidence and no catalog mapping is generic nitpick noise; drop it. Do NOT ` +
            `edit any files.\n\n=== ASSESSMENT BRIEF ===\n${BRIEF}`,
          { schema: ASSESSMENT_SCHEMA, label: `assess:${r.key}`, phase: 'Independent assessment' },
        ),
      ),
    )
  )
    // parallel() preserves input order — stamp the role key by index BEFORE filtering.
    .map((p, i) => (p ? { ...p, role: ROLES[i].key } : null))
    .filter(Boolean)

  if (assessments.length < 2) {
    throw new Error('Independent assessment produced too few reports to consolidate.')
  }

  // --- Phase 2: consolidation (senior owns the merge)
  phase('Consolidate')
  let draft = await agent(
    `You are the ${senior.title} and you OWN this assessment. Consolidate the five independent ` +
      `assessments into ONE maturity report: deduplicate overlapping findings (record which lenses ` +
      `raised each in sources); reconcile grade disagreements EXPLICITLY (record the decided grade + ` +
      `rationale); compute an overallMaturity (1–5) and per-dimension grades; call out CROSS-CUTTING ` +
      `themes ("this keeps biting you"), not just point issues. Give every finding a stable id, an ` +
      `owner role (go-sr | go-mid | db-arch | test-rel | secops), a priority (P0 highest … P3), an ` +
      `effort (S/M/L), and a runVia route — "feature-development" for refactors/features, ` +
      `"security-review" for security items. Order findings by priority (P0 first).\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== INDEPENDENT ASSESSMENTS ===\n${J(assessments)}`,
    { schema: REPORT_SCHEMA, label: 'consolidate', phase: 'Consolidate' },
  )

  // --- Phase 3: reverse analysis (the same team attacks the report) — up to MAX_REVIEW_ROUNDS
  let reviews = []
  let round = 0
  while (true) {
    round++
    phase('Reverse analysis')
    reviews = (
      await parallel(
        ROLES.map((r) => () => {
          const own = assessments.find((p) => p.role === r.key) || null
          return agent(
            `You are the ${r.title}. Run a REVERSE ANALYSIS of the consolidated maturity report — ` +
              `attack the REPORT. Specifically: (1) NITPICK CULL — remove generic advice not grounded ` +
              `in THIS code or not tied to a real go-skills recommendation (the #1 way an assessment ` +
              `gets ignored is boilerplate); (2) PRIORITY/GRADE RIGHT-SIZING — is each item's impact × ` +
              `effort honest? demote cosmetic, promote load-bearing; (3) CONSOLIDATION LOSS — real ` +
              `findings from your own Phase-1 assessment the merge dropped; (4) COVERAGE — name any ` +
              `package/layer/concern in scope that was NOT actually examined (a gap is a finding). ` +
              `Do NOT just attack the text: for every finding you keep OR cull, RE-OPEN the cited file ` +
              `(read-only) and confirm the evidence actually appears in the source — cull any finding ` +
              `whose evidence you cannot confirm against the real code (that is how false findings die). ` +
              `Be skeptical and specific. Give a verdict (approve | approve-with-changes | block).\n\n` +
              `=== BRIEF ===\n${BRIEF}\n\n=== YOUR PHASE-1 ASSESSMENT ===\n${own ? J(own) : '(unavailable)'}\n\n` +
              `=== CONSOLIDATED REPORT ===\n${J(draft)}`,
            { schema: REVERSE_SCHEMA, label: `reverse:${r.key}:r${round}`, phase: 'Reverse analysis' },
          )
        }),
      )
    ).filter(Boolean)

    const blocking = reviews.filter((rv) => rv.verdict === 'block')
    if (!blocking.length || round >= MAX_REVIEW_ROUNDS) {
      if (blocking.length) {
        log(`Round ${round}: ${blocking.length} block(s) remain after the ${MAX_REVIEW_ROUNDS}-round cap — surfaced as open questions.`)
      }
      break
    }

    log(`Round ${round}: ${blocking.length} block verdict(s) — senior revises the report and the team re-reviews.`)
    phase('Consolidate')
    draft = await agent(
      `You are the ${senior.title}. The reverse analysis returned blocking findings on your report. ` +
        `Produce a revised report: cull confirmed nitpicks, apply justified priority/grade ` +
        `adjustments, restore dropped findings, and close coverage gaps (or record them). Keep the ` +
        `conflict log current.\n\n=== CURRENT REPORT ===\n${J(draft)}\n\n=== REVERSE-ANALYSIS FINDINGS ===\n${J(reviews)}`,
      { schema: REPORT_SCHEMA, label: `reconsolidate:r${round}`, phase: 'Consolidate' },
    )
  }

  // --- Phase 4: roadmap (senior folds reverse findings) → approval gate
  phase('Roadmap')
  const plan = await agent(
    `You are the ${senior.title}. Fold the reverse analysis into the FINAL prioritized roadmap. Each ` +
      `item: the problem (with code location/evidence), the go-skill to apply, an owner, an effort, the ` +
      `expected payoff, and the runVia route (feature-development | security-review). State scope, ` +
      `overallMaturity, per-dimension grades, cross-cutting themes, and a firstSprint (the 3–5 items to ` +
      `do first). Fill coverage with what was and was NOT examined. Any unresolved blocking concern ` +
      `needing a human decision goes into openQuestions. Order findings by priority (P0 first).\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== REPORT ===\n${J(draft)}\n\n=== REVERSE ANALYSIS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'roadmap', phase: 'Roadmap' },
  )

  // Deterministically surface any unresolved blocking findings as open questions.
  const unresolvedBlocks = reviews
    .filter((rv) => rv.verdict === 'block')
    .flatMap((rv) => (rv.findings || []).filter((f) => f.severity === 'blocker').map((f) => `[${rv.role}] ${f.issue}`))
  if (unresolvedBlocks.length) {
    plan.openQuestions = [...(plan.openQuestions || []), ...unresolvedBlocks]
  }

  // Deterministic groundedness backstop: the #1 gotcha is "grounded in code AND tied to a go-skill".
  // Drop any roadmap item that lacks an applySkill or any code anchor (evidence/location) — the
  // schema allows empty strings, so enforce it here rather than trust the model.
  const grounded = (f) =>
    f && f.applySkill && f.applySkill.trim() && ((f.evidence && f.evidence.trim()) || (f.location && f.location.trim()))
  if (Array.isArray(plan.findings)) {
    const dropped = plan.findings.filter((f) => !grounded(f))
    if (dropped.length) {
      plan.findings = plan.findings.filter(grounded)
      log(`Dropped ${dropped.length} ungrounded roadmap item(s) (missing applySkill or code anchor).`)
    }
  }

  return {
    mode: 'assess',
    plan,
    draft,
    reviews,
    assessments,
    approvalGate:
      'Present `plan` (maturity + roadmap) to the user for prioritization. This flow is read-only — ' +
      'nothing was changed. Execute approved items via feature-development.mjs (refactors/features) and ' +
      'security-review.mjs (security items). To re-vet user edits first, run { mode: "amend", plan, changeRequest }.',
  }
}

// =============================================================================
// MODE: amend — re-vet an edited roadmap, re-issue it
// =============================================================================
async function runAmend(plan, changeRequest) {
  requirePlan(plan, 'amend')
  const cr = changeRequest || '(no change request text supplied — review the roadmap as given)'

  phase('Reverse analysis')
  const reviews = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are the ${r.title}. The assessment roadmap below has been AMENDED. Change request: ${cr}\n\n` +
            `Reverse-analyze the amended roadmap from your lens: does the change drop a load-bearing ` +
            `finding, mis-prioritize, introduce a generic nitpick, or leave a coverage gap? Cross-check ` +
            `a neighbor's domain. Be skeptical; verdict + findings.\n\n=== AMENDED ROADMAP ===\n${J(plan)}`,
          { schema: REVERSE_SCHEMA, label: `amend-review:${r.key}`, phase: 'Reverse analysis' },
        ),
      ),
    )
  ).filter(Boolean)

  phase('Roadmap')
  const revised = await agent(
    `You are the ${senior.title}. Fold the amend reviews into a REVISED roadmap that honors the change ` +
      `request while keeping it sound (no load-bearing finding silently dropped or mis-prioritized). ` +
      `Surface any unresolved concern in openQuestions.\n\n` +
      `Change request: ${cr}\n\n=== CURRENT ROADMAP ===\n${J(plan)}\n\n=== REVIEWS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'amend-roadmap', phase: 'Roadmap' },
  )

  return {
    mode: 'amend',
    plan: revised,
    reviews,
    approvalGate: 'Present the revised `plan` to the user. Execute approved items via feature-development.mjs / security-review.mjs.',
  }
}

// ---- Entry point ------------------------------------------------------------
// Some runtimes deliver `args` as a JSON string rather than a parsed object — tolerate both,
// and treat a plain (non-JSON) string as the brief itself.
function parseArgs(a) {
  if (a == null) return {}
  if (typeof a === 'string') {
    const s = a.trim()
    if (!s) return {}
    try {
      const parsed = JSON.parse(s)
      return parsed && typeof parsed === 'object' ? parsed : { brief: a }
    } catch {
      return { brief: a }
    }
  }
  return a
}

const A = parseArgs(args)
const mode = A.mode || 'assess'
if (mode === 'amend') {
  return await runAmend(A.plan, A.changeRequest)
}

// assess mode — refuse an empty brief: without it the roles invent a project to assess.
const briefEmpty =
  A.brief == null ||
  (typeof A.brief === 'string' && !A.brief.trim()) ||
  (typeof A.brief === 'object' && !Array.isArray(A.brief) && Object.keys(A.brief).length === 0)
if (briefEmpty) {
  throw new Error(
    'assess mode requires a non-empty args.brief — without it the roles will invent a project. ' +
      'Pass { mode: "assess", brief: <string | { codebase, goals, constraints, depth }> }.',
  )
}
return await runAssess(A.brief)
