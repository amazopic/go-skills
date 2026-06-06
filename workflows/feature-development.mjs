export const meta = {
  name: 'feature-development',
  description: 'Team-of-three Go feature/task flow: independent analysis → consolidate → reverse review → final plan (approval gate) → implementation by the same trio.',
  whenToUse: 'A feature/task touching logic + data + API, a schema/migration change, or cross-cutting work. Not for one-line fixes.',
  phases: [
    { title: 'Frame' },
    { title: 'Independent analysis' },
    { title: 'Consolidate' },
    { title: 'Reverse review' },
    { title: 'Final plan' },
    { title: 'Implementation' },
  ],
}

// =============================================================================
// Feature Development Flow — runnable orchestration of skills/workflow/feature-development.md
//
// NOTE: this runs under the Claude Code Workflow runtime, which injects the
// globals agent()/parallel()/pipeline()/phase()/log()/args and wraps the body in
// an async function. It is therefore NOT meant to parse via `node --check`
// standalone (top-level await/return are legal only under that wrapper).
//
// A background workflow cannot pause for input, so the human APPROVAL GATE splits
// the run into separate invocations:
//
//   mode: "plan"      — Phases 1–4. Returns the vetted final plan for you to approve.
//        Workflow({ scriptPath: "workflows/feature-development.mjs",
//                   args: { mode: "plan", brief: <string | brief object> } })
//
//   mode: "amend"     — re-vet an edited plan (the real Phase-4 loop-back). Runs a
//        targeted reverse review + re-issues the final plan for re-approval.
//        args: { mode: "amend", plan: <edited plan>, changeRequest: <what changed/why> }
//
//   mode: "implement" — Phase 5. Builds the APPROVED plan with the same trio.
//        args: { mode: "implement", plan: <approved plan object> }
//
// The brief (Phase 0) may be a plain string or an object with any of:
//   { request, why, constraints, acceptance, codebase }
// =============================================================================

const MAX_REVIEW_ROUNDS = 2 // reverse-review re-runs before remaining blocks become open questions
const MAX_IMPL_ATTEMPTS = 2 // implement→peer-review reworks per task before it is flagged failed

// ---- The fixed trio ---------------------------------------------------------
const ROLES = [
  {
    key: 'go-mid',
    title: 'Strong-Middle Go Engineer',
    lens:
      'Implementation lens. How to build this in THIS codebase: concrete package/file layout, ' +
      'which existing types/interfaces to reuse, idiomatic Go (context propagation, %w error ' +
      'wrapping, zero-value usability, goroutine lifecycle/cancellation, no premature abstraction), ' +
      'edge cases, validation, partial-failure behavior, and what the tests must cover.',
  },
  {
    key: 'go-sr',
    title: 'Senior Go Engineer',
    lens:
      'Architecture lens (and plan owner). API/interface design, package boundaries, dependency ' +
      'direction, backward compatibility, trade-offs and the decisions behind them, alignment with ' +
      'the service methodology (layering, DI via constructors, no global state), observability, ' +
      'security, failure modes, rollout/rollback, and decomposition into ordered tasks.',
  },
  {
    key: 'db-arch',
    title: 'Middle-Strong DB Architect',
    lens:
      'Data lens. Everything below the repository interface: data model (tables/columns, ' +
      'normalization vs deliberate denormalization, keys), forward + rollback migrations and ' +
      'online/lock-safe strategy for large tables, indexes matched to real access patterns, query ' +
      'shapes and cost, transactions/isolation, contention/locking, consistency, retention/volume.',
  },
]
const roleByKey = (k) => ROLES.find((r) => r.key === k)
const senior = roleByKey('go-sr')
const validKeys = new Set(ROLES.map((r) => r.key))

// ---- Schemas (structured agent output) --------------------------------------
const TASK_SHAPE = {
  type: 'object',
  required: ['id', 'title', 'owner', 'acceptance'],
  properties: {
    id: { type: 'string' },
    title: { type: 'string' },
    owner: { type: 'string', enum: ['go-mid', 'go-sr', 'db-arch'] },
    dependsOn: { type: 'array', items: { type: 'string' } },
    acceptance: { type: 'string' },
  },
}

const PROPOSAL_SCHEMA = {
  type: 'object',
  required: ['understanding', 'plan', 'risks', 'openQuestions', 'tasks'],
  properties: {
    // `role` is stamped deterministically by the script after generation (do not rely on the model).
    role: { type: 'string', description: 'set by the orchestrator; leave as your role key if asked' },
    understanding: { type: 'string' },
    plan: { type: 'string', description: 'the role-domain plan' },
    risks: { type: 'array', items: { type: 'string' } },
    openQuestions: { type: 'array', items: { type: 'string' } },
    tasks: {
      type: 'array',
      items: {
        type: 'object',
        required: ['title', 'effort'],
        properties: {
          title: { type: 'string' },
          effort: { type: 'string', enum: ['S', 'M', 'L'] },
          note: { type: 'string' },
        },
      },
    },
  },
}

const DRAFT_SCHEMA = {
  type: 'object',
  required: ['summary', 'conflictsResolved', 'tasks'],
  properties: {
    summary: { type: 'string' },
    convergenceCheck: {
      type: 'string',
      description: 'do the three understandings agree on WHAT is being built? if not, say so',
    },
    conflictsResolved: {
      type: 'array',
      items: {
        type: 'object',
        required: ['conflict', 'decision', 'rationale'],
        properties: {
          conflict: { type: 'string' },
          decision: { type: 'string' },
          rationale: { type: 'string' },
        },
      },
    },
    tasks: { type: 'array', items: TASK_SHAPE },
  },
}

const REVIEW_SCHEMA = {
  type: 'object',
  required: ['role', 'verdict', 'findings', 'crossDomainCheck'],
  properties: {
    role: { type: 'string' },
    verdict: { type: 'string', enum: ['approve', 'approve-with-changes', 'block'] },
    consolidationLoss: { type: 'array', items: { type: 'string' } },
    mergeDefects: { type: 'array', items: { type: 'string' } },
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
  required: ['title', 'goal', 'architecture', 'dataModel', 'tasks', 'risks', 'rollout', 'testStrategy', 'openQuestions'],
  properties: {
    title: { type: 'string' },
    goal: { type: 'string' },
    architecture: { type: 'string' },
    dataModel: { type: 'string' },
    tasks: { type: 'array', items: TASK_SHAPE },
    risks: {
      type: 'array',
      items: {
        type: 'object',
        required: ['risk', 'mitigation'],
        properties: { risk: { type: 'string' }, mitigation: { type: 'string' } },
      },
    },
    rollout: { type: 'string' },
    testStrategy: { type: 'string' },
    openQuestions: { type: 'array', items: { type: 'string' } },
  },
}

const TASK_DONE_SCHEMA = {
  type: 'object',
  required: ['taskId', 'status', 'summary'],
  properties: {
    taskId: { type: 'string' },
    status: { type: 'string', enum: ['done', 'done-with-concerns', 'blocked'] },
    summary: { type: 'string' },
    files: { type: 'array', items: { type: 'string' } },
    tests: {
      type: 'object',
      description: 'captured test evidence — fill when status is done / done-with-concerns',
      properties: {
        command: { type: 'string', description: 'e.g. go test -race ./...' },
        passed: { type: 'boolean' },
        output: { type: 'string', description: 'short excerpt of the actual run output' },
      },
    },
    concerns: { type: 'array', items: { type: 'string' } },
  },
}

const TASK_REVIEW_SCHEMA = {
  type: 'object',
  required: ['taskId', 'specCompliant', 'qualityApproved', 'notes'],
  properties: {
    taskId: { type: 'string' },
    specCompliant: { type: 'boolean' },
    qualityApproved: { type: 'boolean' },
    notes: { type: 'string' },
  },
}

const FINAL_REVIEW_SCHEMA = {
  type: 'object',
  required: ['overallVerdict', 'gaps', 'notes'],
  properties: {
    overallVerdict: { type: 'string', enum: ['ship', 'ship-with-followups', 'block'] },
    gaps: { type: 'array', items: { type: 'string' }, description: 'integration gaps / acceptance criteria missed across the whole' },
    notes: { type: 'string' },
  },
}

// ---- Helpers ----------------------------------------------------------------
function briefText(b) {
  if (b == null) return '(no brief provided)'
  if (typeof b === 'string') return b
  const parts = []
  if (b.request) parts.push(`Request: ${b.request}`)
  if (b.why) parts.push(`Why: ${b.why}`)
  if (b.constraints) parts.push(`Constraints: ${JSON.stringify(b.constraints)}`)
  if (b.acceptance) parts.push(`Acceptance criteria: ${JSON.stringify(b.acceptance)}`)
  if (b.codebase) parts.push(`Codebase context: ${JSON.stringify(b.codebase)}`)
  return parts.length ? parts.join('\n') : JSON.stringify(b)
}

const J = (x) => JSON.stringify(x, null, 2)

function requirePlan(plan, mode, { checkOwners } = {}) {
  if (!plan || !Array.isArray(plan.tasks) || !plan.tasks.length) {
    throw new Error(`${mode} mode requires args.plan with a non-empty tasks[] array.`)
  }
  if (checkOwners) {
    const ids = new Set(plan.tasks.map((t) => t.id))
    for (const t of plan.tasks) {
      if (!validKeys.has(t.owner)) {
        throw new Error(
          `${mode} mode: task ${t.id || '(no id)'} has invalid owner "${t.owner}" — expected one of ${[...validKeys].join(', ')}.`,
        )
      }
      for (const dep of t.dependsOn || []) {
        if (!ids.has(dep)) {
          throw new Error(`${mode} mode: task ${t.id || '(no id)'} dependsOn unknown task "${dep}".`)
        }
      }
    }
  }
}

// =============================================================================
// MODE: plan  — Phases 1–4, returns the final plan for the approval gate
// =============================================================================
async function runPlan(brief) {
  phase('Frame')
  const BRIEF = briefText(brief)
  log('Framed brief; handing the identical brief to all three roles.')

  // --- Phase 1: independent analysis (fan-out; no role sees another's output)
  phase('Independent analysis')
  const proposals = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are a ${r.title} on a Go service team.\n` +
            `${r.lens}\n\n` +
            `Analyze the task below INDEPENDENTLY — you cannot see your teammates' work; ` +
            `divergence is wanted. If the brief references real codebase paths/files, INSPECT them ` +
            `(read-only) so your analysis is grounded in what actually exists — do not invent APIs. ` +
            `Do NOT write code or edit files. Produce a structured proposal: your understanding, your ` +
            `domain plan, risks, open questions, and a domain task breakdown with effort (S/M/L).\n\n` +
            `=== TASK BRIEF ===\n${BRIEF}`,
          { schema: PROPOSAL_SCHEMA, label: `analyze:${r.key}`, phase: 'Independent analysis' },
        ),
      ),
    )
  )
    // parallel() preserves input order — stamp the role key by index BEFORE filtering,
    // so re-association in Phase 3 never depends on model-authored strings.
    .map((p, i) => (p ? { ...p, role: ROLES[i].key } : null))
    .filter(Boolean)

  if (proposals.length < 2) {
    throw new Error('Independent analysis produced too few proposals to consolidate.')
  }

  // --- Phase 2: consolidation (senior owns the merge)
  phase('Consolidate')
  let draft = await agent(
    `You are the ${senior.title} and you OWN this plan. First run a CONVERGENCE CHECK: do the three ` +
      `proposals' understandings agree on WHAT is being built? If they materially disagree, say so in ` +
      `convergenceCheck (the brief may be ambiguous and need fixing). Then consolidate the three ` +
      `proposals into ONE coherent draft: reconcile every conflict EXPLICITLY (record the conflict, ` +
      `your decision, and the rationale — never drop an idea silently), de-duplicate, fill the gaps ` +
      `the individual lenses missed at their seams, and output a single ORDERED task list where every ` +
      `task has an owner role (go-mid | go-sr | db-arch), dependencies, and acceptance criteria.\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== PROPOSALS ===\n${J(proposals)}`,
    { schema: DRAFT_SCHEMA, label: 'consolidate', phase: 'Consolidate' },
  )

  // --- Phase 3: reverse review (same trio attacks the draft) — up to MAX_REVIEW_ROUNDS
  let reviews = []
  let round = 0
  while (true) {
    round++
    phase('Reverse review')
    reviews = (
      await parallel(
        ROLES.map((r) => () => {
          const own = proposals.find((p) => p.role === r.key) || null
          return agent(
            `You are the ${r.title}. Run a REVERSE REVIEW of the consolidated draft: instead of ` +
              `producing, ATTACK it. Hunt specifically for (1) CONSOLIDATION LOSS — anything from ` +
              `your own Phase-1 proposal that the merge dropped or distorted; (2) MERGE-INTRODUCED ` +
              `DEFECTS — problems that exist only because the parts were combined; and (3) a ` +
              `CROSS-DOMAIN CHECK of a neighbor's area (Go⇒is the data plan ergonomic from the ` +
              `caller; DB⇒are the access patterns serviceable; Senior⇒does it still meet the ` +
              `methodology + acceptance criteria). Be genuinely skeptical: attempt to refute the ` +
              `draft, and if you approve, cite a specific risk you considered and dismissed. Give a ` +
              `verdict (approve | approve-with-changes | block) and concrete, actionable findings.\n\n` +
              `=== BRIEF ===\n${BRIEF}\n\n=== YOUR PHASE-1 PROPOSAL ===\n${own ? J(own) : '(unavailable)'}\n\n` +
              `=== CONSOLIDATED DRAFT ===\n${J(draft)}`,
            { schema: REVIEW_SCHEMA, label: `reverse:${r.key}:r${round}`, phase: 'Reverse review' },
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

    log(`Round ${round}: ${blocking.length} block verdict(s) — senior re-consolidates and the trio re-reviews.`)
    phase('Consolidate')
    draft = await agent(
      `You are the ${senior.title}. The reverse review returned blocking findings on your draft. ` +
        `Produce a revised draft that resolves every blocker and folds in the actionable ` +
        `approve-with-changes findings. Keep the explicit conflict log up to date.\n\n` +
        `=== CURRENT DRAFT ===\n${J(draft)}\n\n=== REVIEW FINDINGS ===\n${J(reviews)}`,
      { schema: DRAFT_SCHEMA, label: `reconsolidate:r${round}`, phase: 'Consolidate' },
    )
  }

  // --- Phase 4: final plan (senior folds review findings) → approval gate
  phase('Final plan')
  const plan = await agent(
    `You are the ${senior.title}. Fold the reverse-review findings into the FINAL implementation ` +
      `plan — the canonical artifact. Be concrete and decision-complete. If the change has no ` +
      `data-layer/storage work, set dataModel to "N/A — no data-layer changes" rather than ` +
      `inventing one. Any remaining blocking concern that needs a human decision goes into ` +
      `openQuestions (surfaced, not hidden).\n\n` +
      `=== BRIEF ===\n${BRIEF}\n\n=== DRAFT ===\n${J(draft)}\n\n=== REVIEWS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'final-plan', phase: 'Final plan' },
  )

  // Deterministically surface any unresolved blocking findings as open questions, so the
  // "remaining blocks become open questions" guarantee holds regardless of the agent's wording.
  const unresolvedBlocks = reviews
    .filter((rv) => rv.verdict === 'block')
    .flatMap((rv) => (rv.findings || []).filter((f) => f.severity === 'blocker').map((f) => `[${rv.role}] ${f.issue}`))
  if (unresolvedBlocks.length) {
    plan.openQuestions = [...(plan.openQuestions || []), ...unresolvedBlocks]
  }

  return {
    mode: 'plan',
    plan,
    draft,
    reviews,
    proposals,
    approvalGate:
      'Present `plan` to the user for approval. Nothing is implemented yet. After sign-off, re-run ' +
      'with { mode: "implement", plan }. To re-vet user edits first, run { mode: "amend", plan, changeRequest }.',
  }
}

// =============================================================================
// MODE: amend — the real Phase-4 loop-back: re-vet an edited plan, re-issue it
// =============================================================================
async function runAmend(plan, changeRequest) {
  requirePlan(plan, 'amend')
  const cr = changeRequest || '(no change request text supplied — review the plan as given)'

  phase('Reverse review')
  const reviews = (
    await parallel(
      ROLES.map((r) => () =>
        agent(
          `You are the ${r.title}. The approved plan below has been AMENDED. Change request: ${cr}\n\n` +
            `Reverse-review the amended plan from your lens: does the change introduce consolidation ` +
            `loss or new defects? Cross-check a neighbor's domain. Be skeptical; verdict + findings.\n\n` +
            `=== AMENDED PLAN ===\n${J(plan)}`,
          { schema: REVIEW_SCHEMA, label: `amend-review:${r.key}`, phase: 'Reverse review' },
        ),
      ),
    )
  ).filter(Boolean)

  phase('Final plan')
  const revised = await agent(
    `You are the ${senior.title}. Fold the amend reviews into a REVISED final plan that honors the ` +
      `change request while keeping the plan coherent. Surface any unresolved concern in openQuestions.\n\n` +
      `Change request: ${cr}\n\n=== CURRENT PLAN ===\n${J(plan)}\n\n=== REVIEWS ===\n${J(reviews)}`,
    { schema: PLAN_SCHEMA, label: 'amend-final-plan', phase: 'Final plan' },
  )

  return {
    mode: 'amend',
    plan: revised,
    reviews,
    approvalGate: 'Present the revised `plan` to the user. After sign-off, run { mode: "implement", plan }.',
  }
}

// =============================================================================
// MODE: implement — Phase 5, the same trio builds the approved plan
// Tasks run in plan order (the senior already ordered them to honor dependsOn).
// Each task: owner implements → a PEER role reviews (spec then quality), gated,
// with bounded rework. A task that never passes is flagged, not silently shipped.
// =============================================================================
async function runImplement(plan) {
  requirePlan(plan, 'implement', { checkOwners: true })
  phase('Implementation')
  log(`Implementing ${plan.tasks.length} approved task(s) in order, peer-reviewed (spec then quality), up to ${MAX_IMPL_ATTEMPTS} attempts each.`)

  const results = []
  for (let i = 0; i < plan.tasks.length; i++) {
    const task = plan.tasks[i]
    const owner = roleByKey(task.owner) // guaranteed present by requirePlan(checkOwners)
    // Rotate the peer reviewer among the non-owner roles by task index, so review-lens coverage
    // is even across a multi-task plan (every role both implements and reviews) rather than one
    // role never reviewing and another always reviewing.
    const candidates = ROLES.filter((r) => r.key !== owner.key)
    const peer = candidates[i % candidates.length]
    const id = task.id || `t${i + 1}`

    let attempt = 0
    let done = null
    let review = null
    let priorFeedback = ''
    while (attempt < MAX_IMPL_ATTEMPTS) {
      attempt++
      done = await agent(
        `You are the ${owner.title}. Implement this approved task with idiomatic, production-grade ` +
          `Go. Write race-safe code and tests, then RUN them and report the exact command, whether ` +
          `they passed, and a short output excerpt in the tests field. Make a focused commit. Report ` +
          `status (done | done-with-concerns | blocked), a summary, files touched, and concerns.` +
          (priorFeedback ? `\n\nThis is a REWORK. ${priorFeedback}` : '') +
          `\n\n=== TASK ===\n${J(task)}\n\n=== PLAN CONTEXT ===\n${J({ goal: plan.goal, architecture: plan.architecture, dataModel: plan.dataModel })}`,
        { schema: TASK_DONE_SCHEMA, label: `impl:${id}:a${attempt}`, phase: 'Implementation' },
      )
      review = await agent(
        `You are the ${peer.title} acting as PEER REVIEWER (you are NOT the implementer). Review in ` +
          `two stages: (1) SPEC COMPLIANCE against the task's acceptance criteria — nothing missing, ` +
          `nothing extra; (2) CODE QUALITY — idiomatic Go, race-safety, and the captured test ` +
          `evidence (verify tests were actually run and passed; do not take the prose summary on ` +
          `faith). Set specCompliant and qualityApproved and explain in notes.\n\n` +
          `=== TASK ===\n${J(task)}\n\n=== IMPLEMENTER REPORT ===\n${J(done)}`,
        { schema: TASK_REVIEW_SCHEMA, label: `review:${id}:a${attempt}`, phase: 'Implementation' },
      )
      if (done.status === 'done' && review.specCompliant === true && review.qualityApproved === true) break
      priorFeedback =
        `The previous attempt was REJECTED. Implementer status: ${done.status}; ` +
        `specCompliant: ${review.specCompliant}; qualityApproved: ${review.qualityApproved}. ` +
        `Reviewer notes: ${review.notes}. Address every point and re-run the tests.`
    }

    const accepted = !!done && done.status === 'done' && !!review && review.specCompliant === true && review.qualityApproved === true
    results.push({ taskId: id, owner: owner.key, peer: peer.key, attempts: attempt, accepted, done, review })
    log(`Task ${id}: ${accepted ? 'accepted' : 'NOT accepted'} after ${attempt} attempt(s).`)
  }

  // Final whole-implementation review (senior): do the delivered tasks TOGETHER satisfy the
  // plan's goal and acceptance criteria? Catches integration gaps a per-task review can't see.
  const finalReview = await agent(
    `You are the ${senior.title}. Perform the FINAL whole-implementation review. Looking at all ` +
      `delivered tasks TOGETHER, do they satisfy the plan's goal and every acceptance criterion? ` +
      `Find integration gaps, inconsistencies between tasks, and any acceptance criterion that no ` +
      `single task fully owns. Give an overall verdict and list concrete gaps.\n\n` +
      `=== PLAN ===\n${J({ goal: plan.goal, architecture: plan.architecture, tasks: plan.tasks })}\n\n` +
      `=== TASK RESULTS ===\n${J(results.map((r) => ({ taskId: r.taskId, owner: r.owner, accepted: r.accepted, summary: r.done && r.done.summary })))}`,
    { schema: FINAL_REVIEW_SCHEMA, label: 'final-review', phase: 'Implementation' },
  )

  const failed = results.filter((r) => !r.accepted)
  return {
    mode: 'implement',
    acceptedCount: results.length - failed.length,
    failedCount: failed.length,
    results,
    finalReview,
    ...(failed.length
      ? { warning: `${failed.length} task(s) did not pass peer review within ${MAX_IMPL_ATTEMPTS} attempts — see results[].review and escalate.` }
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
const mode = A.mode || 'plan'
if (mode === 'implement') {
  return await runImplement(A.plan)
}
if (mode === 'amend') {
  return await runAmend(A.plan, A.changeRequest)
}

// plan mode — refuse an empty brief: without it the roles invent scope from the environment
// (observed failure mode). Fail loudly with usage rather than fabricate a task.
const briefEmpty =
  A.brief == null ||
  (typeof A.brief === 'string' && !A.brief.trim()) ||
  (typeof A.brief === 'object' && !Array.isArray(A.brief) && Object.keys(A.brief).length === 0)
if (briefEmpty) {
  throw new Error(
    'plan mode requires a non-empty args.brief — without it the roles will invent scope from the ' +
      'environment. Pass { mode: "plan", brief: <string | { request, why, constraints, acceptance, codebase }> }.',
  )
}
return await runPlan(A.brief)
