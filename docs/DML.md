# DML Competence Assurance — TMS Roadmap

## Executive Summary: Why This Matters

Your TMS is currently a training administration tool. The DML Competence Procedure requires a competence assurance system. The difference is not academic — it's the gap between knowing someone attended training and knowing they can safely perform critical work without supervision.

**The Stakes:**

- **Safety-Critical Work:** Wrong person performing work = injury or fatality
- **Quality-Critical Work:** Unassessed competence = product nonconformity, customer rejection, regulatory breach
- **Regulatory Compliance:** API Q1, ISO 9001/14001/45001 all require demonstrable competence, not just training records
- **Legal Liability:** Procedure non-compliance becomes evidence in incident investigations

**The Core Problem with Training-Only Systems:**

Training completion proves exposure to information. It does not prove ability to apply that information correctly under operational conditions. The DML procedure explicitly separates:

- **Awareness** (knowing policies/requirements) → Training + induction
- **Competence** (ability to perform work correctly) → Assessment + authorization

Your current TMS tracks the former. This roadmap extends it to cover the latter.

---

## Strategic Assessment: Current State vs. Required State

### What You Have Today

| Capability | Current Approach | Limitation |
|---|---|---|
| Requirement Definition | Profiles group courses | Assumes all roles need same evidence type |
| Completion Tracking | Event attendance + cert upload | No verification of ability to apply |
| Status Management | Valid/Expiring/Expired based on dates | No concept of authorization to work |
| Exception Handling | Exclusions with reason/expiry | No supervision workflow for partial competence |
| Renewal Triggers | Date-based alerts only | Misses performance-based reassessment |

### What DML Procedure Requires

| Requirement | Source Section | Business Driver |
|---|---|---|
| Risk-based inclusion | 3.2.2 | Don't over-burden low-risk roles; don't under-protect high-risk roles |
| Formal assessment for critical work | 3.2.3 | Regulatory and customer mandate for objective evidence |
| Supervised work periods | 3.2.4 | Safe pathway from training to independence |
| Line manager authorization | 3.2.1 | Accountability for competence decisions |
| Reassessment on triggers | 3.2.5 | Competence degrades — errors, changes, and absence matter |
| Training effectiveness verification | 4.4 | Ensure training investment actually works |
| Awareness vs. competence separation | 5.1–5.3 | Efficient resource allocation — don't assess what doesn't need assessing |
| New starter induction records | 3.3 | HR must record and retain signed induction attendance |

---

## Data Model Clarification: Roles vs. Profiles

> **This distinction must be resolved before implementation begins.**

The current TMS uses **profiles** — collections of training requirements assigned to employees. The DML procedure and Skills Matrix are built around **roles** — job functions with defined Job Descriptions, competence expectations, and applicability decisions.

These are different concepts:

| Concept | Current System | DML Procedure |
|---|---|---|
| Unit of assignment | Profile (groups courses) | Role (job function with JD) |
| Applicability gate | None — all profiles assumed in matrix | 5-question decision per role |
| Evidence type | Uniform per profile | Varies by role — cert, observation, supervisor sign-off |
| Maintained by | System admin | HR (Job Descriptions) |

**Migration path:** Existing profiles become the starting point for defining roles. Each profile maps to one or more roles. Roles then receive an applicability decision before entering the new competence workflow. This mapping must be completed during Phase 1 before the authorization layer can function.

---

## The Six Critical Gaps

### Gap 1: No Applicability Decision

**Why It Matters:**

The procedure's 5-question logic (Section 3.2.2, DML-QA-REG-5024-1 Applicability Decision Matrix) prevents two costly errors:

- **Over-engineering:** Requiring full competence tracking for roles that only need awareness (e.g., office staff who need to know evacuation routes but don't touch product)
- **Under-protection:** Missing critical roles that need formal assessment because they don't look "operational" (e.g., engineers approving design changes)

**Business Impact:**
- Wasted training spend on non-critical roles
- Audit findings for missing critical role assessments
- Inconsistent application across departments

**What Changes:**

Every role gets evaluated through the 5-question workflow before entering the matrix. This becomes the gateway to requirement assignment. Borderline cases are escalated to QHSE/HR for resolution before a decision is recorded.

---

### Gap 2: Training Completion ≠ Competence Authorization

**Why It Matters:**

The procedure states (Section 3.2.1): *"Before anyone works independently, their line manager must confirm they're ready. This confirmation is based on evidence like qualifications, experience, or direct observation of their work."*

| | Workflow |
|---|---|
| **Current** | Employee completes course → Status = Valid → Can work |
| **Required** | Employee completes course → Manager confirms readiness → Status = Authorized → Can work independently |

**Important nuance — authorization is not the same as formal assessment:**

- For **standard roles**: manager confirmation is sufficient — based on qualifications, experience, or observation (§3.2.1)
- For **critical work**: a formally recorded competence assessment is required, identifying the person, activities covered, assessor, and any limitations (§3.2.3)

The system must support both pathways. Not every authorization requires a formal assessment record — only critical work does.

**The Danger:**

An employee with a valid training certificate but no demonstrated ability to apply it under operational conditions performs safety-critical work. Incident occurs. Investigation reveals gap between training record and competence verification. Liability attaches to the organization and line manager.

**What Changes:**

An explicit authorization layer sits between training completion and work permission. The authorization type — `independent`, `supervised`, or `restricted` — is recorded by the line manager. Training creates eligibility; authorization creates permission.

---

### Gap 3: No New Starter Induction Record

**Why It Matters:**

Section 3.3 of the procedure requires HR to conduct a formal new starter induction and retain a signed attendance record before the individual begins work. There is currently no system support for this.

**Business Impact:**
- No audit trail of induction completion
- No record of which topics were covered
- HR has no system of record — relies on paper

**What Changes:**

Induction records are created in the system at point of hire — employee, date, HR conductor, topics covered, and employee signature reference. On completion, the employee is handed to their line manager and enters the normal competence workflow (§3.2).

---

### Gap 4: No Supervision Workflow

**Why It Matters:**

Not everyone is immediately competent after training. The procedure (Section 3.2.4) requires a structured path:

> Training complete → Supervised work → Evidence of correct performance → Authorization for independent work

- **Current state:** Binary (trained/untrained or valid/expired)
- **Required state:** Multi-state including supervised periods with defined scope limitations

**Business Impact:**
- New hires placed into independent work prematurely → errors, rework, incidents
- No audit trail of who was supervising whom, for what work, for how long
- Inconsistent "shadowing" practices across departments

**What Changes:**

Formal supervision periods tracked in the system with scope restrictions, supervisor assignment, and completion criteria.

---

### Gap 5: Date-Based Only Renewal

**Why It Matters:**

Certification expiry is one trigger for reassessment. The procedure (Section 3.2.5) lists four others:

1. Repeated errors, defects, or quality issues
2. Process, equipment, or job requirement changes
3. Extended absence (>3 months for safety- or quality-critical activities)
4. Performance or safety concerns identified by manager

The current TMS only tracks the calendar. It misses the operational reality that competence degrades or becomes irrelevant independently of expiry dates.

**Example:**

A welder was certified 6 months ago. Last month, a new welding procedure was introduced. The current TMS shows "Valid" for 18 more months. The DML procedure requires immediate reassessment because of a process change.

**What Changes:**

Multi-trigger reassessment system linking to operational events — quality issues, change management, HR absence tracking.

---

### Gap 6: No Distinction Between Awareness and Competence

**Why It Matters:**

The procedure dedicates Section 5 to Awareness — understanding policies, objectives, and the consequences of nonconformance. This is distinct from Section 3 Competence — the ability to perform work correctly.

| Track | Mechanism | Overhead |
|---|---|---|
| **Awareness** | HR induction, line manager communication, annual refreshers | Low — no formal assessment required |
| **Competence** | Role-specific, evidence-based, formally authorized, continuously monitored | Higher — proportionate to risk |

**Current Risk:**

Existing profiles likely mix both. Result: over-assessment of awareness items (treating policy understanding like safety-critical skill) and under-management of competence (treating critical work like general awareness).

**What Changes:**

The applicability decision routes each role to either:
- **Matrix track** — full competence management (assessment, authorization, supervision)
- **Awareness track** — induction/communication only (outside matrix, lighter touch)

---

## Action Required (AUTO) Logic

The DML Skills Matrix includes a computed `Action Required` field that drives operational decisions. The logic is:

| Status | Condition | Action Required (AUTO) |
|---|---|---|
| `REQUIRED` | Authorization not yet granted, training not complete | Training/authorisation required before independent work |
| `IN_SUPERVISION` | Supervision period active | Working under supervision — assessment pending |
| `PENDING_AUTHORIZATION` | Training complete, awaiting manager sign-off | Manager authorization required |
| `AUTHORIZED` | Fully authorized for independent work | No action required |
| `EXPIRING` | Within 60 days of expiry | Renew/refresher required before expiry |
| `EXPIRED` | Past expiry date | Suspend independent work until competence is revalidated |
| `REASSESSMENT_TRIGGERED` | Non-calendar trigger raised and open | Reassessment required — check trigger log |

This logic must be implemented as a computed field, not stored state — it is derived from the current authorization record, supervision status, and expiry date at query time.

---

## Integration Roadmap

### Phase 1: Foundation — Roles, Applicability & Authorization
**Timeline: Weeks 1–4**

**Why First:** Everything else depends on knowing which roles exist, which need full competence control, and whether a manager has authorized each person to work independently.

| Deliverable | Business Outcome | DML Compliance |
|---|---|---|
| Roles as first-class entity (migrated from profiles) | Clear mapping of job functions to competence requirements | §3.1 |
| 5-Question Applicability Wizard | Consistent, defensible decisions on role inclusion | §3.2.2 |
| Role Classification (Matrix vs. Awareness) | Right level of control for right roles | §3.2.2 |
| New Starter Induction Records | HR retains signed induction attendance | §3.3 |
| Work Authorization Workflow (two tiers) | Manager sign-off for standard; formal record for critical | §3.2.1, §3.2.3 |
| Authorization Status Layer | Separation of training completion from work permission | §3.2.1 |
| Legacy data migration to `legacy_pending_confirmation` state | Existing records acknowledged without bypassing procedure | §3.2.1 |

**Key Endpoints:**
```
POST   /api/roles
GET    /api/roles
POST   /api/roles/{id}/applicability-decision
GET    /api/roles/{id}/applicability-decision
POST   /api/escalations/{id}/resolve
POST   /api/induction-records
POST   /api/work-authorizations
PUT    /api/work-authorizations/{id}/revoke
GET    /api/employees/{id}/competence-status
GET    /api/dashboard/competence-health
```

---

### Phase 2: Operational Safety — Assessment & Supervision
**Timeline: Weeks 5–8**

**Why Second:** Protects against the highest-risk scenario — unassessed personnel performing critical work.

| Deliverable | Business Outcome | DML Compliance |
|---|---|---|
| Formal Competence Assessment Records | Objective evidence for auditors and investigators | §3.2.3 |
| Supervision Period Management | Safe, tracked progression to independence | §3.2.4 |
| Assessment Scheduling Queue | Nothing falls through cracks — systematic coverage | §3.2.1 |
| Scope Limitation Tracking | Clear boundaries on what supervised workers can do | §3.2.4 |

**Key Endpoints:**
```
POST   /api/competence-assessments
GET    /api/competence-assessments/pending
GET    /api/employees/{id}/competence-assessments
POST   /api/supervision-periods
PUT    /api/supervision-periods/{id}/complete
GET    /api/supervision/active
GET    /api/supervision/ready-for-assessment
```

---

### Phase 3: Continuous Assurance — Reassessment & Improvement
**Timeline: Weeks 9–12**

**Why Third:** Prevents competence degradation and ensures training effectiveness over time.

| Deliverable | Business Outcome | DML Compliance |
|---|---|---|
| Multi-Trigger Reassessment System | Proactive response to errors, changes, absence | §3.2.5 |
| Training Effectiveness Verification | Confidence that training investment delivers results | §4.4 |
| Performance-Linked Alerts | Early warning before incidents occur | §3.2.5 |
| Reassessment Queue Management | Systematic handling of non-calendar triggers | §3.2.5 |

**Key Endpoints:**
```
POST   /api/reassessment-triggers
GET    /api/reassessment-triggers/open
PUT    /api/reassessment-triggers/{id}/resolve
POST   /api/training-effectiveness-evaluations
GET    /api/dashboard/reassessment-alerts
GET    /api/reports/competence-degradation-risk
```

---

### Phase 4: Integration & Intelligence
**Timeline: Weeks 13–16**

**Why Last:** External integrations and advanced analytics require stable core processes.

| Deliverable | Business Outcome | DML Compliance |
|---|---|---|
| HRIS Integration | Automatic reassessment triggers on absence/transfer | §3.2.5 |
| QHSE System Integration | Incident/change-driven reassessment | §3.2.5, §6.1 |
| Access Control Integration | Automatic enforcement of authorization status | §3.2.1 |
| Advanced Analytics | Predictive competence risk, training ROI | §8 (Monitoring) |

**Key Endpoints:**
```
POST   /api/webhooks/hris/absence
POST   /api/webhooks/qhse/incident
POST   /api/webhooks/change-management/process-change
POST   /api/access-control/sync
GET    /api/analytics/training-effectiveness
GET    /api/analytics/competence-risk-predictions
```

---

### Phase 5 (Future): Procedure Change Communication
**Out of scope for Phases 1–4 but required by §6**

The procedure requires line managers to communicate significant procedure changes to affected personnel, with QHSE classification and verification of application. This is currently managed entirely outside the system. A future phase should add:

- Procedure change classification (significant vs. minor)
- Targeted communication records to affected roles
- Verification of application tracking

---

## New API Endpoints Summary

### Applicability & Role Management

| Endpoint | Purpose |
|---|---|
| `POST /api/roles` | Create a role (migrated from profile or new) |
| `GET /api/roles` | List all roles with applicability status |
| `POST /api/roles/{id}/applicability-decision` | Execute 5-question logic, record decision |
| `GET /api/roles/{id}/applicability-decision` | Retrieve current applicability status |
| `POST /api/escalations/{id}/resolve` | QHSE/HR resolves borderline cases |
| `GET /api/roles/matrix-included` | List roles in competence matrix |
| `GET /api/roles/awareness-only` | List roles outside matrix |

### New Starter Induction

| Endpoint | Purpose |
|---|---|
| `POST /api/induction-records` | Record HR induction for new starter |
| `GET /api/employees/{id}/induction-record` | Retrieve induction record for individual |
| `GET /api/induction-records/pending` | Employees who have not yet completed induction |

### Competence Assessment

| Endpoint | Purpose |
|---|---|
| `POST /api/competence-assessments` | Record formal assessment for critical work |
| `GET /api/competence-assessments/pending` | Queue of required assessments |
| `GET /api/employees/{id}/competence-assessments` | History for individual |
| `GET /api/employees/{id}/competence-status/{requirement_id}` | Current authorization state |

### Work Authorization

| Endpoint | Purpose |
|---|---|
| `POST /api/work-authorizations` | Grant independent/supervised/restricted authorization |
| `PUT /api/work-authorizations/{id}/revoke` | Remove authorization (incident, expiry, performance) |
| `GET /api/employees/{id}/work-authorizations` | Authorization history and current status |
| `GET /api/work-authorizations/expiring` | Time-limited authorizations needing renewal |
| `GET /api/work-authorizations/legacy-pending` | Records requiring manager confirmation at go-live |

### Supervision Management

| Endpoint | Purpose |
|---|---|
| `POST /api/supervision-periods` | Start supervised work period |
| `PUT /api/supervision-periods/{id}/complete` | End supervision, link to assessment |
| `GET /api/supervision/active` | Currently supervised employees |
| `GET /api/supervision/ready-for-assessment` | Completed supervision, awaiting assessment |

### Reassessment & Triggers

| Endpoint | Purpose |
|---|---|
| `POST /api/reassessment-triggers` | Manual trigger (manager observes issue) |
| `GET /api/reassessment-triggers/open` | Pending reassessments |
| `PUT /api/reassessment-triggers/{id}/resolve` | Close trigger after assessment |
| `POST /api/training-effectiveness-evaluations` | Verify training worked in practice |

### Dashboards & Reporting

| Endpoint | Purpose |
|---|---|
| `GET /api/dashboard/competence-health` | Organization-wide authorization status |
| `GET /api/dashboard/supervision-queue` | Assessment workload for managers |
| `GET /api/dashboard/reassessment-alerts` | Non-calendar triggers requiring action |
| `GET /api/reports/dml-qa-reg-5024` | Audit-ready skills matrix report |
| `GET /api/reports/competence-degradation-risk` | Predictive risk analysis |

---

## Success Metrics

| Metric | Current State | Target State | Measurement |
|---|---|---|---|
| Roles with Applicability Decision | 0% (assumed all in matrix) | 100% | Decision recorded in system |
| Critical Work with Formal Assessment | Unknown | 100% | Assessment record exists |
| Training Completion to Authorization Gap | Invisible | Zero | No "Valid" without "Authorized" or "Legacy Pending" |
| Supervised Work Tracked | Informal | 100% | Supervision period record exists |
| Reassessment Triggers Missed | Unknown | Zero | Open triggers report |
| Legacy Records Confirmed by Manager | N/A | 100% within 90 days of go-live | `legacy_pending_confirmation` queue cleared |
| Audit Finding Response Time | Days to gather evidence | Minutes to generate report | Time to produce DML-QA-REG-5024 |

---

## Risk Mitigation

| Risk | Mitigation |
|---|---|
| Manager Resistance (extra authorization step) | System enforces workflow — cannot skip. Dashboard shows queue, makes it manageable. |
| Legacy Data Migration | Existing "Valid" records migrated as `legacy_pending_confirmation` — acknowledged but not auto-authorized. Managers clear the queue for critical roles within 90 days of go-live. This prevents bypassing the §3.2.1 requirement under audit scrutiny. |
| Roles vs. Profiles Ambiguity | Profile-to-role mapping completed as first task in Phase 1 before any authorization workflow is enabled. |
| Over-Complexity | Awareness-track roles bypass the entire competence workflow — only matrix roles enter the authorization, assessment, and supervision processes. |
| Integration Delays | Phase 4 external integrations are optional — core compliance is achieved in Phases 1–3. |

---

*This roadmap transforms TMS from a training record repository into a comprehensive competence assurance platform that satisfies DML procedure requirements (DML-QA-PRE-1023) while providing operational safety and audit readiness.*
