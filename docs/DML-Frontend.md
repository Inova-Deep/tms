# DML Competence Assurance — Frontend Specification

> **Status:** Specification only. No frontend code changes are in scope until the backend is complete.
> All new pages and components depend on the API endpoints defined in [DML.md](./DML.md).

---

## Stack Reference

| Layer | Technology |
|---|---|
| Framework | Vue 3 + TypeScript |
| Routing | Vue Router |
| State | Pinia |
| Data fetching | TanStack Vue Query |
| Tables | TanStack Vue Table |
| Components | Reka UI |
| Styling | Tailwind CSS 4 |

---

## New Routes & Pages

| Route | New View File | Purpose |
|---|---|---|
| `/roles` | `RolesView.vue` | List all roles with applicability status badge |
| `/roles/:id` | `RoleDetailView.vue` | Role detail — requirements, applicability decision, employees assigned |
| `/roles/:id/applicability` | `ApplicabilityWizardView.vue` | 5-question wizard to classify role as Matrix or Awareness track |
| `/inductions` | `InductionsView.vue` | New starter induction records — HR view, pending completions |
| `/work-authorizations` | `WorkAuthorizationsView.vue` | Manager view of all authorization statuses, legacy-pending queue |
| `/assessments` | `AssessmentsView.vue` | Competence assessment queue — pending, in-progress, completed |
| `/supervision` | `SupervisionView.vue` | Active supervision periods, ready-for-assessment list |
| `/reassessments` | `ReassessmentsView.vue` | Open reassessment triggers by type — errors, absences, changes, concerns |

---

## Changes to Existing Pages

### `DashboardView.vue`

Add new KPI cards alongside existing ones:

| KPI | Source Endpoint | Description |
|---|---|---|
| Authorization Health | `GET /api/dashboard/competence-health` | % of matrix roles with active authorization |
| Legacy Pending | `GET /api/work-authorizations/legacy-pending` | Count of records awaiting manager confirmation at go-live |
| Active Supervision | `GET /api/supervision/active` | Employees currently working under supervision |
| Open Reassessments | `GET /api/reassessment-triggers/open` | Non-calendar triggers requiring action |
| Pending Assessments | `GET /api/competence-assessments/pending` | Critical work roles awaiting formal assessment |

The existing grid view should gain a third mode alongside Course/Profile: **Authorization Status** — showing each employee's current authorization state per role.

---

### `EmployeeDetailView.vue`

Add new tabs to the existing employee detail layout:

| Tab | Content |
|---|---|
| **Induction** | HR induction record — date, conductor, topics, signature status |
| **Work Authorizations** | Authorization history — type, granted by, date, expiry, revocations |
| **Assessments** | Competence assessment records — activity, assessor, outcome, limitations |
| **Supervision** | Supervision periods — supervisor, scope limits, start/end, status |
| **Reassessments** | Open and resolved reassessment triggers |

---

### `MatrixView.vue`

Add new columns to the skills matrix table:

| Column | Source | Notes |
|---|---|---|
| Risk Level | Role field | Safety-Critical / Quality-Critical / Low |
| Criticality | Role field | Safety / Quality / Compliance |
| Mandatory for Role | Requirement field | Y / N badge |
| Training Type | Requirement field | Certification/Licence, Induction/Orientation, Experience/Supervisor-validated |
| Assessment Method | Requirement field | Record Review, Supervisor Sign-off, Observation |
| Authorization Status | Computed from authorization record | See Authorization Status below |
| Action Required (AUTO) | Computed field | See Action Required logic below |

---

### `ProfilesView.vue`

Profiles remain in the system as the mechanism for grouping requirements. A **Roles** section is added separately (see new `/roles` route above). A banner should appear on the Profiles page pointing users to the new Roles workflow for applicability decisions.

No structural changes to profiles themselves — they continue to define what training is required. Roles define *who* needs it and *at what level*.

---

### `AppSidebar.vue`

Add new navigation group — **Competence** — between the existing training and admin sections:

```
Competence
  ├── Roles                  /roles
  ├── Work Authorizations    /work-authorizations
  ├── Assessments            /assessments
  ├── Supervision            /supervision
  └── Reassessments          /reassessments

Administration
  ├── Inductions             /inductions   (move from HR-only to sidebar)
  └── ... existing items
```

---

## New Components

### `ApplicabilityWizard.vue`

Step-by-step form implementing the 5-question decision logic from DML-QA-REG-5024-1.

| Step | Question | On YES | On NO |
|---|---|---|---|
| Q1 | Does this role perform hands-on operational work? | Continue | Likely outside matrix (unless Q2 = YES) |
| Q2 | Does this role approve/release/inspect/certify work affecting conformity? | Continue | Go to Q3 |
| Q3 | Can an error in this role cause safety, quality, environmental, or compliance impact? | Continue | NOT in matrix |
| Q4 | Is specific competence required (not just general awareness)? | Continue | Manage via Awareness/induction |
| Q5 | Is objective evidence required to prove competence? | INCLUDE in matrix | NOT in matrix |

Result displayed as: **Applicable — Include in Skills Matrix** or **Not Applicable — Awareness Track Only**

Borderline cases (manual override) route to escalation form — QHSE/HR resolution.

---

### `AuthorizationStatusBadge.vue`

Replaces and extends the existing status badge. New states:

| State | Colour | Meaning |
|---|---|---|
| `AUTHORIZED` | Green | Manager confirmed — independent work permitted |
| `IN_SUPERVISION` | Blue | Working under supervision |
| `PENDING_AUTHORIZATION` | Amber | Training complete — awaiting manager sign-off |
| `LEGACY_PENDING` | Grey | Migrated from pre-DML system — confirmation required |
| `REQUIRED` | Red | Training and/or authorization not yet started |
| `EXPIRING` | Amber | Within 60 days of expiry |
| `EXPIRED` | Red | Past expiry — independent work suspended |
| `REASSESSMENT_TRIGGERED` | Orange | Open reassessment trigger — action required |

---

### `ActionRequiredChip.vue`

Computed display component for the `Action Required (AUTO)` column in the matrix. Renders the action text with appropriate severity colour. Not a stored value — derived at render time from authorization status, expiry date, and open triggers.

---

### `WorkAuthorizationForm.vue`

Manager-facing form to grant or revoke authorization. Two variants:

- **Standard (§3.2.1):** Manager confirmation based on qualifications, experience, or observation — lightweight sign-off
- **Critical (§3.2.3):** Formal assessment record — person assessed, activities covered, assessor, limitations/supervision requirements

The form selects the correct variant automatically based on the role's risk level.

---

### `SupervisionPeriodForm.vue`

Form to open and close supervision periods:
- Employee, supervisor (competent person), start date
- Scope limitations (what the employee can/cannot do independently)
- Completion: link to competence assessment, evidence notes

---

### `ReassessmentTriggerForm.vue`

Manager-facing form to raise a non-calendar reassessment trigger:

| Trigger Type | Procedure Reference |
|---|---|
| Repeated errors / defects / quality issues | §3.2.5 |
| Process / equipment / job requirement change | §3.2.5 |
| Extended absence (>3 months, critical roles) | §3.2.5 |
| Performance or safety concern | §3.2.5 |

Cert/licence expiry is auto-triggered by the system — not raised manually.

---

### `LegacyPendingBanner.vue`

Go-live migration component. Shown on the Dashboard and Work Authorizations page until the `legacy_pending_confirmation` queue is cleared.

Displays: count of records pending manager confirmation, link to the queue, 90-day target date.

---

### `InductionRecordForm.vue`

HR-facing form for new starter induction:
- Employee, induction date, HR conductor
- Topics covered (checklist from §3.3: quality/safety framework, mandatory policies, site access, emergency arrangements, admin requirements)
- Employee signature reference
- Submission records to HR

---

## New Composables

| File | Endpoints Used |
|---|---|
| `useRoles.ts` | `GET/POST /api/roles`, `GET /api/roles/matrix-included`, `GET /api/roles/awareness-only` |
| `useApplicabilityDecisions.ts` | `POST/GET /api/roles/{id}/applicability-decision`, `POST /api/escalations/{id}/resolve` |
| `useWorkAuthorizations.ts` | `POST /api/work-authorizations`, `PUT /api/work-authorizations/{id}/revoke`, `GET /api/work-authorizations/legacy-pending` |
| `useCompetenceAssessments.ts` | `POST /api/competence-assessments`, `GET /api/competence-assessments/pending`, `GET /api/employees/{id}/competence-assessments` |
| `useSupervisionPeriods.ts` | `POST /api/supervision-periods`, `PUT /api/supervision-periods/{id}/complete`, `GET /api/supervision/active` |
| `useReassessmentTriggers.ts` | `POST /api/reassessment-triggers`, `GET /api/reassessment-triggers/open`, `PUT /api/reassessment-triggers/{id}/resolve` |
| `useInductionRecords.ts` | `POST /api/induction-records`, `GET /api/employees/{id}/induction-record`, `GET /api/induction-records/pending` |
| `useTrainingEffectiveness.ts` | `POST /api/training-effectiveness-evaluations` |

---

## New Types (`types/api.ts`)

```typescript
// Roles
type Role = {
  id: string
  name: string
  department: string
  skill_category: 'role_specific_technical' | 'management_system' | 'governance'
  risk_level: 'safety_critical' | 'quality_critical' | 'low' | 'medium'
  criticality: 'safety' | 'quality' | 'compliance'
  mandatory_for_role: boolean
  matrix_applicable: boolean | null  // null = no decision yet
}

// Applicability
type ApplicabilityDecision = {
  id: string
  role_id: string
  q1: boolean; q2: boolean; q3: boolean; q4: boolean; q5: boolean
  result: 'matrix' | 'awareness'
  decided_by: string
  decided_at: string
  escalated: boolean
  escalation_notes?: string
}

// Authorization
type AuthorizationType = 'independent' | 'supervised' | 'restricted'
type AuthorizationState =
  | 'AUTHORIZED'
  | 'IN_SUPERVISION'
  | 'PENDING_AUTHORIZATION'
  | 'LEGACY_PENDING'
  | 'REQUIRED'
  | 'EXPIRING'
  | 'EXPIRED'
  | 'REASSESSMENT_TRIGGERED'

type WorkAuthorization = {
  id: string
  employee_id: string
  requirement_id: string
  authorization_type: AuthorizationType
  authorized_by: string
  authorized_at: string
  expiry_date?: string
  revoked_at?: string
  revoke_reason?: string
  is_legacy: boolean
}

// Competence Assessments
type CompetenceAssessment = {
  id: string
  employee_id: string
  requirement_id: string
  assessed_by: string
  assessment_date: string
  work_activities_covered: string
  limitations?: string
  supervision_required: boolean
  outcome: 'competent' | 'not_yet_competent' | 'supervised'
  evidence_reference?: string
}

// Supervision
type SupervisionPeriod = {
  id: string
  employee_id: string
  supervisor_id: string
  requirement_id: string
  scope_limitations: string
  start_date: string
  end_date?: string
  status: 'active' | 'completed' | 'abandoned'
  completion_assessment_id?: string
}

// Reassessment
type TriggerType =
  | 'repeated_errors'
  | 'process_change'
  | 'extended_absence'
  | 'cert_expiry'
  | 'performance_concern'

type ReassessmentTrigger = {
  id: string
  employee_id: string
  requirement_id: string
  trigger_type: TriggerType
  triggered_by: string
  triggered_at: string
  notes: string
  resolved_at?: string
  resolution_notes?: string
  resolution_assessment_id?: string
}

// Induction
type InductionRecord = {
  id: string
  employee_id: string
  conducted_by: string
  induction_date: string
  topics_covered: string[]
  employee_signature_ref?: string
  completed: boolean
}

// Training Effectiveness
type TrainingEffectivenessEvaluation = {
  id: string
  evidence_id: string
  evaluated_by: string
  evaluation_date: string
  method: 'observation' | 'supervision' | 'error_reduction' | 'output_review' | 'feedback'
  outcome: 'effective' | 'partially_effective' | 'ineffective'
  follow_up_action?: string
}
```

---

## Implementation Notes

- The `Action Required (AUTO)` field is **never stored** — always computed on the frontend from the combination of `AuthorizationState`, expiry date, and open trigger count. The backend provides the raw data; the frontend derives the display value.
- The `AuthorizationStatusBadge` replaces the current `Valid/Expiring/Expired` badge system in the matrix. The existing three states map to: `AUTHORIZED` (Valid), `EXPIRING` (Expiring), `EXPIRED` (Expired).
- At go-live, the `LegacyPendingBanner` is the primary call-to-action for managers. It should be the first thing visible on the Dashboard until the queue is clear.
- Awareness-track roles (applicability result = `awareness`) do not appear in the matrix, the authorization workflow, or the assessment queue. They appear only in the Roles list with an "Awareness Only" badge.
