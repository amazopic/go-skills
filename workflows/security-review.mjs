export const meta = {
  name: 'security-review',
  description: 'Team-of-four Go security review: independent audit → consolidate report (dedup + attack chains) → reverse analysis (cull false positives, verify exploitability) → remediation plan (approval gate) → remediation by the same team.',
  whenToUse: 'A security audit of a service, sensitive subsystem (auth/payments/multi-tenant data), or a diff crossing a trust boundary. Not for a one-line obvious fix.',
  phases: [
    { title: 'Frame' },
    { title: 'Independent analysis' },
    { title: 'Consolidate' },
    { title: 'Reverse analysis' },
    { title: 'Remediation plan' },
    { title: 'Remediation' },
  ],
}

// =============================================================================
// Security Code Review — runnable orchestration of skills/workflow/security-review.md
//
// NOTE: runs under the Claude Code Workflow runtime, which injects the globals
// agent()/parallel()/pipeline()/phase()/log()/args and wraps the body in an async
// function. NOT meant to parse via `node --check` standalone (top-level
// await/return are legal only under that wrapper).
//
// The human APPROVAL GATE splits the run across invocations:
//
//   mode: "review"    — Phases 1–4. Returns the vetted remediation plan to approve.
//        args: { mode: "review", brief: <string | brief object> }
//
//   mode: "amend"     — re-vet an edited plan (re-rate a finding, accept a risk).
//        args: { mode: "amend", plan: <edited plan>, changeRequest: <what changed/why> }
//
//   mode: "remediate" — Phase 5. Fixes the APPROVED findings with the same team.
//        args: { mode: "remediate", plan: <approved plan object> }
//
// The brief (Phase 0) may be a plain string or an object with any of:
//   { scope, threatModel, assets, compliance, codebase }
// =============================================================================

const MAX_REVIEW_ROUNDS = 2 // reverse-analysis re-runs before remaining blocks become open questions
const MAX_FIX_ATTEMPTS = 2  // fix→peer-verify reworks per finding before it is flagged failed

// ---- The fixed team ---------------------------------------------------------
const ROLES = [
  {
    key: 'go-sr',
    title: 'Senior Go Engineer',
    lens:
      'Go-language security (and report owner). crypto misuse (math/rand where crypto/rand is ' +
      'required, non-constant-time secret comparison, weak/reused nonces, home-grown crypto), TLS, ' +
      'JWT/session handling; injection via Go APIs (os/exec with a shell, text/template vs ' +
      'html/template, string-built database/sql queries), path traversal, SSRF, unsafe ' +
      'deserialization; security-relevant data races and TOCTOU, unsafe, integer overflow, ' +
      'resource exhaustion / unbounded goroutines (DoS); secrets leaking via errors/logs/panics.',
  },
  {
    key: 'go-mid',
    title: 'Strong-Middle Go Engineer',
    lens:
      'Implementation-level security, code-first. Per-endpoint authentication AND authorization ' +
      'checks (the missing "if !authorized"); input validation and output encoding at every trust ' +
      'boundary; mass-assignment/overposting; HTTP hardening (security headers, cookie ' +
      'HttpOnly/Secure/SameSite, CORS, CSRF, rate limiting, request size/timeouts); dependency risk ' +
      '(govulncheck); error responses leaking stack traces or internal detail.',
  },
  {
    key: 'db-arch',
    title: 'Middle-Strong DB Architect',
    lens:
      'Data-layer security below the repository interface. SQL/NoSQL/ORM injection (parameterization, ' +
      'not concatenation); least privilege (DB roles/grants), row-level security, tenant isolation, ' +
      'IDOR at the query; PII/secrets at rest (encryption, column-level), retention, audit logging of ' +
      'sensitive access; migrations that widen exposure; connection-string/credential handling, backups.',
  },
  {
    key: 'secops',
    title: 'Senior SecOps',
    lens:
      'Threat, architecture and operational security — the system view and the attacker view. Threat ' +
      'model & trust boundaries; authn/authz ARCHITECTURE (not just per-handler checks); secrets ' +
      'management (vault, rotation), transport security, supply-chain & build integrity; ' +
      'container/deploy hardening, least privilege at the edge, egress control; security logging, ' +
      'monitoring, alerting, incident-response readiness, blast radius. Maps findings to OWASP Top 10 / ' +
      'CWE and judges real-world exploitability and likelihood.',
  },
]
const roleByKey = (k) => ROLES.find((r) => r.key === k)
const senior = roleByKey('go-sr')
const validKeys = new Set(ROLES.map((r) => r.key))

const SEVERITIES = ['critical', 'high', 'medium', 'low', 'info']
const SEV_RANK = { critical: 0, high: 1, medium: 2, low: 3, info: 4 }

// ---- Schemas (structured agent output) --------------------------------------
const ANALYSIS_FINDING = {
  type: 'object',
  required: ['title', 'severity', 'location', 'exploit', 'remediation'],
  properties: {
    title: { type: 'string' },
    severity: { type: 'string', enum: SEVERITIES },
    cwe: { type: 'string', description: 'e.g. CWE-89' },
    location: { type: 'string', description: 'file:line or package/function' },
    description: { type: 'string' },
    exploit: { type: 'string', description: 'WHY it is exploitable here — concrete path, not just a risky API' },
    remediation: { type: 'string' },
    confidence: { type: 'string', enum: ['high', 'medium', 'low'] },
  },
}

const ANALYSIS_SCHEMA = {
  type: 'object',
  required: ['summary', 'findings'],
  properties: {
    role: { type: 'string', description: 'stamped by the orchestrator' },
    summary: { type: 'string' },
    findings: { type: 'array', items: ANALYSIS_FINDING },
    openQuestions: { type: 'array', items: { type: 'string' } },
  },
}

// A finding in the consolidated report / remediation plan: carries an owner + verification.
const PLAN_FINDING = {
  type: 'object',
  required: ['id', 'title', 'severity', 'location', 'impact', 'remediation', 'owner', 'verification'],
  properties: {
    id: { type: 'string' },
    title: { type: 'string' },
    severity: { type: 'string', enum: SEVERITIES },
    cwe: { type: 'string' },
    location: { type: 'string' },
    impact: { type: 'string' },
    exploit: { type: 'string' },
    remediation: { type: 'string' },
    owner: { type: 'string', enum: ['go-sr', 'go-mid', 'db-arch', 'secops'] },
    dependsOn: { type: 'array', items: { type: 'string' } },
    verification: { type: 'string', description: 'how we prove the vulnerability is closed (the test)' },
    sources: { type: 'array', items: { type: 'string' }, description: 'which lenses raised it' },
  },
}

const ATTACK_CHAIN = {
  type: 'object',
  required: ['chain', 'combinedSeverity'],
  properties: {
    chain: { type: 'string', description: 'the finding ids/titles that compose' },
    steps: { type: 'array', items: { type: 'string' } },
    combinedSeverity: { type: 'string', enum: SEVERITIES },
  },
}

const REPORT_SCHEMA = {
  type: 'object',
  required: ['summary', 'conflictsResolved', 'findings'],
  properties: {
    summary: { type: 'string' },
    dedupNotes: { type: 'string' },
    conflictsResolved: {
      type: 'array',
      items: {
        type: 'object',
        required: ['conflict', 'decision', 'rationale'],
        properties: { conflict: { type: 'string' }, decision: { type: 'string' }, rationale: { type: 'string' } },
      },
    },
    attackChains: { type: 'array', items: ATTACK_CHAIN },
    findings: { type: 'array', items: PLAN_FINDING },
  },
}

const REVERSE_SCHEMA = {
  type: 'object',
  required: ['role', 'verdict', 'findings', 'crossDomainCheck'],
  properties: {
    role: { type: 'string' },
    verdict: { type: 'string', enum: ['approve', 'approve-with-changes', 'block'] },
    falsePositives: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'why'],
        properties: { id: { type: 'string' }, why: { type: 'string' } },
      },
    },
    severityAdjustments: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'from', 'to', 'why'],
        properties: { id: { type: 'string' }, from: { type: 'string' }, to: { type: 'string' }, why: { type: 'string' } },
      },
    },
    consolidationLoss: { type: 'array', items: { type: 'string' } },
    missedChains: { type: 'array', items: { type: 'string' } },
    coverageGaps: { type: 'array', items: { type: 'string' }, description: 'endpoints/sinks/trust boundaries in scope that were NOT examined' },
    exploitabilityChecks: {
      type: 'array',
      items: {
        type: 'object',
        required: ['id', 'exploitable', 'note'],
        properties: { id: { type: 'string' }, exploitable: { type: 'boolean' }, note: { type: 'string' } },
      },
    },
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
  required: ['title', 'scope', 'summary', 'severitySummary', 'findings', 'residualRisk', 'testStrategy', 'openQuestions'],
  properties: {
    title: { type: 'string' },
    scope: { type: 'string' },
    threatModel: { type: 'string' },
    summary: { type: 'string' },
    severitySummary: {
      type: 'object',
      properties: {
        critical: { type: 'integer' },
        high: { type: 'integer' },
        medium: { type: 'integer' },
        low: { type: 'integer' },
        info: { type: 'integer' },
      },
    },
    findings: { type: 'array', items: PLAN_FINDING },
    attackChains: { type: 'array', items: ATTACK_CHAIN },
    residualRisk: { type: 'string' },
    acceptedRisk: { type: 'string' },
    testStrategy: { type: 'string' },
    openQuestions: { type: 'array', items: { type: 'string' } },
  },
}

const FIX_DONE_SCHEMA = {
  type: 'object',
  required: ['findingId', 'status', 'summary'],
  properties: {
    findingId: { type: 'string' },
    status: { type: 'string', enum: ['fixed', 'fixed-with-concerns', 'blocked'] },
    summary: { type: 'string' },
    files: { type: 'array', items: { type: 'string' } },
    tests: {
      type: 'object',
      description: 'the regression/security test proving the fix — fill when status is fixed / fixed-with-concerns',
      properties: {
        command: { type: 'string', description: 'e.g. go test -race ./...' },
        failsOnVulnerableCode: { type: 'boolean', description: 'does the test fail before the fix?' },
        passed: { type: 'boolean' },
        output: { type: 'string' },
      },
    },
    concerns: { type: 'array', items: { type: 'string' } },
  },
}

const FIX_REVIEW_SCHEMA = {
  type: 'object',
  required: ['findingId', 'vulnerabilityClosed', 'noRegression', 'notes'],
  properties: {
    findingId: { type: 'string' },
    vulnerabilityClosed: { type: 'boolean', description: 'the VULNERABILITY is closed, not just the one payload masked' },
    noRegression: { type: 'boolean' },
    notes: { type: 'string' },
  },
}

const FINAL_REVIEW_SCHEMA = {
  type: 'object',
  required: ['overallVerdict', 'residualRisks', 'notes'],
  properties: {
    overallVerdict: { type: 'string', enum: ['ship', 'ship-with-followups', 'block'] },
    residualRisks: { type: 'array', items: { type: 'string' } },
    notes: { type: 'string' },
  },
}

// Records what the orchestration did with a fix after the verify gate (verify-before-persist).
const PERSIST_SCHEMA = {
  type: 'object',
  required: ['findingId', 'action', 'detail'],
  properties: {
    findingId: { type: 'string' },
    action: { type: 'string', enum: ['committed', 'reverted', 'partial', 'none'] },
    detail: { type: 'string' },
  },
}

// ---- Helpers ----------------------------------------------------------------
function briefText(b) {
  if (b == null) return '(no brief provided)'
  if (typeof b === 'string') return b
  const parts = []
  if (b.scope) parts.push(`Scope: ${JSON.stringify(b.scope)}`)
  if (b.threatModel) parts.push(`Threat model: ${JSON.stringify(b.threatModel)}`)
  if (b.assets) parts.push(`Assets & trust boundaries: ${JSON.stringify(b.assets)}`)
  if (b.compliance) parts.push(`Compliance bar: ${JSON.stringify(b.compliance)}`)
  if (b.codebase) parts.push(`Codebase context: ${JSON.stringify(b.codebase)}`)
  return parts.length ? parts.join('\n') : JSON.stringify(b)
}

const J = (x) => JSON.stringify(x, null, 2)

function requirePlan(plan, mode, { checkOwners } = {}) {
  if (!plan || !Array.isArray(plan.findings) || !plan.findings.length) {
    throw new Error(`${mode} mode requires args.plan with a non-empty findings[] array.`)
  }
  if (checkOwners) {
    const ids = new Set(plan.findings.map((f) => f.id))
    for (const f of plan.findings) {
      if (!validKeys.has(f.owner)) {
        throw new Error(
          `${mode} mode: finding ${f.id || '(no id)'} has invalid owner "${f.owner}" — expected one of ${[...validKeys].join(', ')}.`,
        )
      }
      for (const dep of f.dependsOn || []) {
        if (!ids.has(dep)) {
          throw new Error(`${mode} mode: finding ${f.id || '(no id)'} dependsOn unknown finding "${dep}".`)
        }
      }
    }
  }
}

// =============================================================================
// MODE: review — Phases 1–4, returns the remediation plan for the approval gate
// =============================================================================
async function runReview(brief) {
  phase('Frame')
  const BRIEF = briefText(brief)
  log('Framed review brief; handing the identical brief to all four roles.')

  // --- Phase 1: independent security analysis (fan-out; no role sees another's findings)
  phase('Independent analysis')
  const reports = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are a ${r.title} performing a SECURITY review.\n` +
            `${r.lens}\n\n` +
            `Audit the target below INDEPENDENTLY — you cannot see your teammates' findings; ` +
            `divergence widens coverage. GROUND every finding in the ACTUAL code: read the real files ` +
            `the brief references (read-only) — an audit of an imagined codebase is worthless, and ` +
            `false positives are expensive. Where the code is a compilable module, RUN the tooling and ` +
            `treat its output as input: \`govulncheck ./...\` (dependency CVEs/reachability) and ` +
            `\`go vet ./...\`. Hunt NEGATIVE SPACE as hard as present bugs — MISSING authorization ` +
            `checks, MISSING rate limiting, MISSING input validation, absent security headers: absences ` +
            `are where real breaches hide and are the easiest to overlook. Do NOT write code or edit ` +
            `files. For each vulnerability give ` +
            `title, severity (critical/high/medium/low/info) + CWE, exact location, a CONCRETE exploit ` +
            `path that explains why it is exploitable HERE (not just that a risky API appears), a ` +
            `remediation, and your confidence. Rate severity by exploitability × impact in the threat ` +
            `model, not by how scary the API sounds.\n\n` +
            `=== REVIEW BRIEF ===\n${BRIEF}`,
          { schema: ANALYSIS_SCHEMA, label: `audit:${r.key}`, phase: 'Independent analysis' },
        ),
      ),
    )
  )
    // parallel() preserves input order — stamp the role key by index BEFORE filtering.
    .map((p, i) => (p ? { ...p, role: ROLES[i].key } : null))
    .filter(Boolean)

  if (reports.length < 2) {
    throw new Error('Independent analysis produced too few reports to consolidate.')
  }

  // --- Phase 2: consolidation (senior owns the merge)
  phase('Consolidate')
  let draft = await agent(
    `You are the ${senior.title} and you OWN this security report. Consolidate the four independent ` +
      `audits below into ONE report: DEDUPLICATE the same vulnerability reported by multiple lenses ` +
      `into a single finding (keep the highest-confidence write-up, record which lenses raised it in ` +
      `sources); reconcile every SEVERITY conflict EXPLICITLY (record the disagreement, your decided ` +
      `severity, and the rationale — never silently drop or downgrade); and SURFACE ATTACK CHAINS — ` +
      `combinations of individually-minor findings that together raise the effective severity (a chain ` +
      `is itself a finding). Give every finding a stable id, an owner role (go-sr | go-mid | db-arch | ` +
      `secops) for the fix, and a verification step. Order findings by severity (critical first).\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== INDEPENDENT AUDITS ===\n${J(reports)}`,
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
          const own = reports.find((p) => p.role === r.key) || null
          return agent(
            `You are the ${r.title}. Run a REVERSE ANALYSIS of the consolidated security report: ` +
              `attack the REPORT itself. Specifically: (1) FALSE POSITIVES — findings not actually ` +
              `reachable/exploitable in this code; argue removal or downgrade (this is the noise ` +
              `filter — a report nobody trusts is useless); (2) SEVERITY ADJUSTMENTS — right-size over- ` +
              `and under-rated findings against the threat model; (3) CONSOLIDATION LOSS — real ` +
              `findings from your own Phase-1 audit that the merge dropped; (4) MISSED CHAINS / blind ` +
              `spots visible only when findings are combined; (5) EXPLOITABILITY CHECK — for each ` +
              `High/Critical, state whether a realistic exploit path exists given the trust ` +
              `boundaries; (6) COVERAGE — name any endpoint, query/sink, or trust boundary in scope ` +
              `that was NOT examined (this pass culls false positives, so guard against the opposite ` +
              `bias of under-reporting). Default to skepticism on BOTH existence (does it reproduce?) ` +
              `and severity (is it really that bad here?). Give a verdict (approve | approve-with-changes | block).\n\n` +
              `=== BRIEF ===\n${BRIEF}\n\n=== YOUR PHASE-1 AUDIT ===\n${own ? J(own) : '(unavailable)'}\n\n` +
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
        `Produce a revised report: remove confirmed false positives, apply justified severity ` +
        `adjustments, restore dropped findings, and add missed attack chains. Keep the conflict log ` +
        `current.\n\n=== CURRENT REPORT ===\n${J(draft)}\n\n=== REVERSE-ANALYSIS FINDINGS ===\n${J(reviews)}`,
      { schema: REPORT_SCHEMA, label: `reconsolidate:r${round}`, phase: 'Consolidate' },
    )
  }

  // --- Phase 4: remediation plan (senior folds reverse findings) → approval gate
  phase('Remediation plan')
  const plan = await agent(
    `You are the ${senior.title}. Fold the reverse analysis into the FINAL remediation plan — ` +
      `CONFIRMED findings only (drop the false positives, apply the severity adjustments). Each finding ` +
      `keeps its id, severity+CWE, location, impact, the fix, an owner, dependencies, and a verification ` +
      `step (the test that proves it closed). Fill severitySummary with the per-severity counts. Order ` +
      `findings by risk (critical first). State residualRisk and any acceptedRisk. Any unresolved ` +
      `blocking concern needing a human risk decision goes into openQuestions (surfaced, not hidden).\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== REPORT ===\n${J(draft)}\n\n=== REVERSE ANALYSIS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'remediation-plan', phase: 'Remediation plan' },
  )

  // Deterministically surface any unresolved blocking findings as open questions.
  const unresolvedBlocks = reviews
    .filter((rv) => rv.verdict === 'block')
    .flatMap((rv) => (rv.findings || []).filter((f) => f.severity === 'blocker').map((f) => `[${rv.role}] ${f.issue}`))
  if (unresolvedBlocks.length) {
    plan.openQuestions = [...(plan.openQuestions || []), ...unresolvedBlocks]
  }

  return {
    mode: 'review',
    plan,
    draft,
    reviews,
    reports,
    approvalGate:
      'Present `plan` to the user for risk acceptance + prioritization. Nothing is patched yet. After ' +
      'sign-off, re-run with { mode: "remediate", plan }. To re-vet user edits first, run ' +
      '{ mode: "amend", plan, changeRequest }.',
  }
}

// =============================================================================
// MODE: amend — re-vet an edited remediation plan, re-issue it
// =============================================================================
async function runAmend(plan, changeRequest) {
  requirePlan(plan, 'amend')
  const cr = changeRequest || '(no change request text supplied — review the plan as given)'

  phase('Reverse analysis')
  const reviews = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are the ${r.title}. The approved security remediation plan below has been AMENDED. ` +
            `Change request: ${cr}\n\nReverse-analyze the amended plan from your lens: does the change ` +
            `introduce a false positive, an unsafe severity downgrade, a dropped finding, or a newly ` +
            `unguarded attack chain? Cross-check a neighbor's domain. Be skeptical; verdict + findings.\n\n` +
            `=== AMENDED PLAN ===\n${J(plan)}`,
          { schema: REVERSE_SCHEMA, label: `amend-review:${r.key}`, phase: 'Reverse analysis' },
        ),
      ),
    )
  ).filter(Boolean)

  phase('Remediation plan')
  const revised = await agent(
    `You are the ${senior.title}. Fold the amend reviews into a REVISED remediation plan that honors ` +
      `the change request while keeping the plan sound (no real vulnerability silently dropped or ` +
      `under-rated). Surface any unresolved concern in openQuestions.\n\n` +
      `Change request: ${cr}\n\n=== CURRENT PLAN ===\n${J(plan)}\n\n=== REVIEWS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'amend-remediation-plan', phase: 'Remediation plan' },
  )

  return {
    mode: 'amend',
    plan: revised,
    reviews,
    approvalGate: 'Present the revised `plan` to the user. After sign-off, run { mode: "remediate", plan }.',
  }
}

// =============================================================================
// MODE: remediate — Phase 5, the same team fixes the approved findings.
// Findings are fixed in PRIORITY order (critical → info). Each finding: owner
// fixes → a PEER role verifies the vulnerability is CLOSED (not masked) and no
// regression, gated with bounded rework. A finding that never passes is flagged.
// =============================================================================
async function runRemediate(plan) {
  requirePlan(plan, 'remediate', { checkOwners: true })
  phase('Remediation')

  // Priority order: critical first. Stable within a severity (preserve plan order).
  const ordered = plan.findings
    .map((f, i) => ({ f, i }))
    .sort((a, b) => (SEV_RANK[a.f.severity] ?? 9) - (SEV_RANK[b.f.severity] ?? 9) || a.i - b.i)
    .map((x) => x.f)

  log(`Remediating ${ordered.length} finding(s) in priority order, peer-verified (closed + no regression), up to ${MAX_FIX_ATTEMPTS} attempts each.`)

  const results = []
  for (let i = 0; i < ordered.length; i++) {
    const finding = ordered[i]
    const owner = roleByKey(finding.owner) // guaranteed by requirePlan(checkOwners)
    const candidates = ROLES.filter((r) => r.key !== owner.key)
    const peer = candidates[i % candidates.length] // rotate the verifier for even lens coverage
    const id = finding.id || `F${i + 1}`

    let attempt = 0
    let done = null
    let review = null
    let priorFeedback = ''
    while (attempt < MAX_FIX_ATTEMPTS) {
      attempt++
      done = await agent(
        `You are the ${owner.title}. Fix this confirmed vulnerability with an idiomatic, ` +
          `production-grade Go change that closes the VULNERABILITY CLASS, not just the one reported ` +
          `payload. Add a regression/security test that FAILS on the vulnerable code and PASSES on your ` +
          `fix; run it (and go test -race for anything concurrency-related) and report the command and ` +
          `result in the tests field. Apply the change and the test in the working tree but DO NOT ` +
          `commit yet — leave the changes uncommitted so the peer verifier can independently re-check ` +
          `before anything is persisted. Report status (fixed | fixed-with-concerns ` +
          `| blocked), summary, files, and concerns.` +
          (priorFeedback ? `\n\nThis is a REWORK. ${priorFeedback}` : '') +
          `\n\n=== FINDING ===\n${J(finding)}\n\n=== PLAN CONTEXT ===\n${J({ scope: plan.scope, threatModel: plan.threatModel })}`,
        { schema: FIX_DONE_SCHEMA, label: `fix:${id}:a${attempt}`, phase: 'Remediation' },
      )
      review = await agent(
        `You are the ${peer.title} acting as PEER VERIFIER (you are NOT the fixer — do NOT trust the ` +
          `fixer's self-report). Verify INDEPENDENTLY by inspecting the actual diff in the working tree ` +
          `and RE-RUNNING the test yourself (plus go vet / go test -race where relevant): (1) the ` +
          `VULNERABILITY CLASS is actually closed — not just the one reported payload masked; confirm ` +
          `the regression test truly FAILS on the pre-fix code and PASSES after; (2) NO REGRESSION and ` +
          `no new vulnerability was introduced. Base your verdict on what YOU observed, not on the ` +
          `fixer's claims. Set vulnerabilityClosed and noRegression and explain in notes.\n\n` +
          `=== FINDING ===\n${J(finding)}\n\n=== FIXER REPORT (claims to check, not trust) ===\n${J(done)}`,
        { schema: FIX_REVIEW_SCHEMA, label: `verify:${id}:a${attempt}`, phase: 'Remediation' },
      )
      if (done.status === 'fixed' && review.vulnerabilityClosed === true && review.noRegression === true) break
      priorFeedback =
        `The previous attempt was REJECTED. Fixer status: ${done.status}; vulnerabilityClosed: ` +
        `${review.vulnerabilityClosed}; noRegression: ${review.noRegression}. Verifier notes: ` +
        `${review.notes}. Address every point and re-run the test.`
    }

    const accepted = !!done && done.status === 'fixed' && !!review && review.vulnerabilityClosed === true && review.noRegression === true

    // Verify-before-persist: commit ONLY an independently-verified fix; otherwise discard the
    // unverified (possibly-insecure) working-tree changes so they never reach the branch.
    const persistence = accepted
      ? await agent(
          `You are the ${owner.title}. The fix for finding ${id} ("${finding.title}") passed independent ` +
            `peer verification. Commit ONLY the files belonging to this fix, with a clear message such as ` +
            `"fix(security): ${finding.title}${finding.cwe ? ` (${finding.cwe})` : ''}". Report the result.`,
          { schema: PERSIST_SCHEMA, label: `commit:${id}`, phase: 'Remediation' },
        )
      : await agent(
          `You are the ${owner.title}. The fix for finding ${id} ("${finding.title}") FAILED verification ` +
            `after ${attempt} attempt(s). Discard its UNCOMMITTED working-tree changes (e.g. ` +
            `git restore / git checkout -- on: ${J((done && done.files) || [])}) so the unverified, ` +
            `possibly-insecure change does NOT persist. Report what you reverted.`,
          { schema: PERSIST_SCHEMA, label: `revert:${id}`, phase: 'Remediation' },
        )

    results.push({ findingId: id, severity: finding.severity, owner: owner.key, peer: peer.key, attempts: attempt, accepted, persistence, done, review })
    log(`Finding ${id} (${finding.severity}): ${accepted ? 'remediated + committed' : 'NOT remediated — reverted'} after ${attempt} attempt(s).`)
  }

  // Final whole-remediation review (senior): do the fixes interact safely; is residual risk acceptable?
  const finalReview = await agent(
    `You are the ${senior.title}. Perform the FINAL whole-remediation review. Looking at all fixes ` +
      `TOGETHER: do they close the findings without interacting badly or opening new gaps? What ` +
      `residual risk remains? Give an overall verdict and list residual risks.\n\n` +
      `=== PLAN ===\n${J({ scope: plan.scope, threatModel: plan.threatModel, findings: plan.findings })}\n\n` +
      `=== FIX RESULTS ===\n${J(results.map((r) => ({ findingId: r.findingId, severity: r.severity, owner: r.owner, accepted: r.accepted, summary: r.done && r.done.summary })))}`,
    { schema: FINAL_REVIEW_SCHEMA, label: 'final-review', phase: 'Remediation' },
  )

  const failed = results.filter((r) => !r.accepted)
  return {
    mode: 'remediate',
    remediatedCount: results.length - failed.length,
    failedCount: failed.length,
    results,
    finalReview,
    ...(failed.length
      ? { warning: `${failed.length} finding(s) were not verified-closed within ${MAX_FIX_ATTEMPTS} attempts — see results[].review and escalate.` }
      : {}),
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
const mode = A.mode || 'review'
if (mode === 'remediate') {
  return await runRemediate(A.plan)
}
if (mode === 'amend') {
  return await runAmend(A.plan, A.changeRequest)
}

// review mode — refuse an empty brief: without it the roles invent scope from the environment.
const briefEmpty =
  A.brief == null ||
  (typeof A.brief === 'string' && !A.brief.trim()) ||
  (typeof A.brief === 'object' && !Array.isArray(A.brief) && Object.keys(A.brief).length === 0)
if (briefEmpty) {
  throw new Error(
    'review mode requires a non-empty args.brief — without it the roles will invent scope from the ' +
      'environment. Pass { mode: "review", brief: <string | { scope, threatModel, assets, compliance, codebase }> }.',
  )
}
return await runReview(A.brief)
