# TMS PRD (v1) — Deep Manufacturing

## 1. Product Overview

TMS (Training Management System) is Deep Manufacturing’s centrally governed system of record for training, certifications, competence profiles, evidence retention, expiry management, visual gap analysis, and audit-ready reporting. TMS replaces spreadsheets and fragmented records with a single, deterministic ruleset that computes employee eligibility against profile requirements and provides operational planning visibility through dashboards and notifications.

TMS is designed for high auditability and low operational friction:

* Only Training Admins can create/modify records (single-writer governance).
* Employees have read-only, self-only access to their own profiles/requirements/records.
* Profiles are dynamic and owned by Training Management (not hardcoded role taxonomies).
* Eligibility is computed from evidence + expiry + equivalencies.

---

## 2. Goals and Non‑Goals

### 2.1 Goals

* Replace Excel-based tracking with a centralized, auditable system.
* Provide a profile-driven competence model that Training Managers can change without engineering.
* Support two primary training modes:

  1. Instructor-led (attendance-based)
  2. Third-party certification (evidence upload, attachment mandatory)
* Provide clear employee-facing “My Requirements / My Records” visibility.
* Provide Training Admin/Manager dashboards:

  * Overall compliance
  * Expiring soon / expired
  * Visual gap analysis (Course × Org slice and Profile × Org slice)
* Provide audit packs and evidence completeness reporting.
* Integrate with Deep’s Employee Directory via API for membership rules and scoping.

### 2.2 Non‑Goals (v1)

* Authorisation to work / permit-to-work workflows (out of scope; TMS publishes status only).
* Self-service record creation by employees or instructors.
* Full LMS features (content delivery, quizzes, SCORM).
* Session scheduling complexity (capacity planning, waitlists) beyond simple events + attendance.

---

## 3. Users and Personas

### 3.1 Training Admin (Primary Writer)

* Creates/edits profiles, requirements, items.
* Creates training events, captures instructor and attendance.
* Uploads certification evidence and validates metadata.
* Runs reports, manages expiries, executes training plans.

### 3.2 Training Manager (Governance + Oversight)

* Cross-cost-center visibility.
* Defines/approves profile standards, waivers, and policy rules.
* Reviews compliance posture, risk exposure, and upcoming workload.

### 3.3 Employee (Self-Service Read Only)

* Views assigned profiles, requirements, and statuses.
* Views personal training/certification transcript and expiry timeline.
* Receives expiry notifications (if email exists).

(Optionally later) Auditor / QHSE / Department Management (read-only reporting roles).

---

## 4. Scope: Capabilities

### 4.1 Profile Management

* Create profile (name, description, owner, tags, domain).
* Add requirements to profile:

  * Course completion requirements
  * Certification/licence requirements
  * Optional: acknowledgement requirement (future)
* Equivalency groups (OR logic): requirement satisfied by any valid item in group.
* Profile lifecycle: Draft → Published (future Effective From date) → Retired.
* Profile versioning:

  * Create vNext while vCurrent remains active.
  * vNext becomes active at Effective From date.
  * Preview impact (who becomes non-compliant; assignment volume).

### 4.2 Membership and Auto‑Assignment

* Membership sources:

  * RULE (directory rule)
  * MANUAL_ADD (explicit employee IDs)
  * MANUAL_EXCLUDE (override)
* Directory-driven rules must support the guaranteed fields:

  * employee_id, employment_status, worker_type, org_unit/department, cost_center, location
* Manual exclusions require reason + expiry.
* Membership rules apply continuously; recompute on directory changes.

### 4.3 Course & Certification Catalogue

* Catalogue items:

  * Course item (internal)
  * Certification/licence item (external)
* Each item defines validity policy:

  * validity duration (months) OR never expires OR expiry provided by evidence
* Domains (e.g., QHSE, Quality, Operations) and optional course type/category.

### 4.4 Instructor‑Led Training Events (Attendance)

* Create event: course, date, location, instructor (employee_id), notes.
* Attendance roster: list of employee_id.
* Confirm attendance (Training Admin only).
* Completion generation:

  * Attendance confirmation = completion (default).
* Event history retained for audit.

### 4.5 Third‑Party Certifications (Evidence Upload)

* Create certification record: item, issuer, issue date, expiry date (or never).
* Attachment mandatory.
* Store evidence in TMS-managed storage.

### 4.6 Evidence & Records

* Store completion records with metadata:

  * item, completion date, expiry date, evidence type, attachments, admin user, timestamps.
* Maintain historical records; select active evidence for computation.
* Retention per Deep Manufacturing policy.

### 4.7 Eligibility & Gap Computation

* Requirement statuses (locked): Missing, Valid, Expiring, Expired.
* Profile status per employee: Eligible / At Risk / Not Eligible.
* Expiring counts as compliant in headline compliance metric but is highlighted as At Risk.
* Materialize computed statuses for performance and consistent reporting.

### 4.8 Dashboards and Reporting

* KPI tiles: overall compliance, total employees in scope, expiring soon, expired.
* Filters: org unit/department, cost center, location, domain, course type, search.
* Visual gap analysis grids:

  1. Course Coverage: Course × Org slice with counts and %
  2. Profile Coverage: Profile × Org slice with Eligible/At Risk/Not Eligible
* Drill-down:

  * click any cell → employee list with status breakdown; export.
* Reports:

  * audit pack per profile / org slice (evidence completeness)
  * expiring/expired lists by horizon (30/60/90)
  * compliance trend (optional)

### 4.9 Notifications

* Recipients:

  * Training Admin + Training Manager (always)
  * Employee (if email exists)
* Thresholds: 90/60/30 days (optionally 14/7 for critical).
* Cadence:

  * digest (daily/weekly configurable)
  * weekly overdue reminders until resolved.

### 4.10 Waivers

* Allowed (Training Manager only).
* Requires reason + expiry.
* Waivers are visible and reportable (never hidden).

---

## 5. Business Logic (Normative)

This PRD references the locked business logic document:

* **TMS v1 Business Logic Specification (Locked)** (Canvas)

Key highlights:

* Evidence-driven computation; assignments are planning artifacts (v1: Planned only).
* Earliest-expiry rule governs risk posture when duplicates exist.
* Course validity is defined by the course/cert item.

---

## 6. User Experience Requirements

### 6.1 Admin: Coverage Dashboard (Client-Liked Pattern)

* KPI tiles across the top.
* Filters bar: Department/Org slice, Domain, Type, Search.
* Grid view with card-style cells showing:

  * % coverage
  * Valid / Expiring / Expired / Missing counts
  * total (x/y)
* Cell click opens drill-down drawer:

  * employee list + status breakdown + export.

### 6.2 Admin: Training Events

* Create event; choose course; capture instructor by employee_id; date/location.
* Add roster (search by employee).
* Confirm attendance; generate completions.

### 6.3 Admin: Certification Upload

* Choose employee; choose cert item; capture issuer/issue/expiry.
* Upload attachment (mandatory).

### 6.4 Employee: My Profile & Records

* My Profiles (read-only): list of profiles that apply.
* My Requirements: grouped by profile and domain; statuses visible.
* My Records (transcript): course/cert list with completion/expiry/status and evidence view.

---

## 7. Data and Domain Model (High Level)

### 7.1 Core Entities

* Employee (cached from Directory)
* Profile
* ProfileVersion
* Requirement (within ProfileVersion)
* Item (Course / Certification)
* Membership (RULE / MANUAL_ADD / MANUAL_EXCLUDE)
* TrainingEvent
* Attendance
* EvidenceRecord
* Waiver
* ComputedStatus (RequirementStatus, ProfileStatus)
* AuditEvent

### 7.2 Required Directory Fields

* employee_id (immutable key)
* employment_status (active/inactive/terminated)
* worker_type (employee/contractor)
* org_unit/department
* cost_center
* location

(Recommended additional fields if available: email, manager_id, job_title.)

---

## 8. Integrations

### 8.1 Inbound

* Employee Directory API

  * sync schedule + delta handling
  * supports mover/leaver updates

### 8.2 Outbound (Optional / Future)

* Eligibility status API for other systems to query:

  * by employee_id, profile_id
  * returns Eligible/At Risk/Not Eligible + evidence summary

---

## 9. Security, Access Control, and Audit

### 9.1 RBAC

* Training Admin: full write
* Training Manager: cross-cost-center read; (optional write for governance actions)
* Employee: self-only read

### 9.2 Org Scoping

* Training Manager has cross-cost-center access by policy.
* Admin access can be global or scoped (configurable).

### 9.3 Auditability

* Immutable audit log for:

  * profile/version changes
  * membership rules and overrides
  * exclusions and waiver actions
  * evidence record creation/modification
  * attendance confirmation

### 9.4 Retention

* Evidence and audit logs retained per **Deep Manufacturing retention policy**.

---

## 10. Non‑Functional Requirements

* Performance: dashboards must remain responsive with materialized status tables.
* Availability: business-hours critical; define RTO/RPO per IT standards.
* Data integrity: single-writer governance + audit trail.
* Privacy: employee self-only views; strict access checks.
* Export: CSV/PDF/Excel export (format choice can be phased; core is data availability).

---

## 11. Acceptance Criteria (MVP-Level)

* Training Admin can:

  * create profile + publish with future effective date
  * attach requirements and equivalencies
  * define rule membership and manual add/exclude with reason+expiry
  * create training event, capture instructor, add attendance, confirm attendance
  * upload certification evidence with mandatory attachment
  * see dashboards and drill-down lists
  * generate audit pack for a profile

* Employee can:

  * view own profiles and requirements
  * view own transcript and evidence
  * receive expiry notification if email exists

* System computes:

  * requirement statuses (Missing/Valid/Expiring/Expired)
  * profile statuses (Eligible/At Risk/Not Eligible)
  * overall compliance (employee-weighted; expiring counts as compliant)

---

## 12. Delivery Plan (Agile, 3 Teams)

This is a planning outline suitable for stakeholder alignment.

### Team Streams

* Stream A: Platform/Directory/RBAC/Audit/Notifications
* Stream B: Profiles/Membership/Rules/Versioning
* Stream C: Events/Evidence/Computation/Dashboards/Reports

### Iterative Delivery (Recommended Sequence)

1. Foundations: Directory sync, RBAC, audit backbone
2. Profiles + membership rules + versioning + employee self-view
3. Training events + attendance completion + certification upload
4. Dashboards + drill-down + audit packs + exports
5. Hardening: performance, retention controls, UAT, rollout

(Exact sprint counts can be added once you confirm delivery expectations and environment constraints.)

---

## 13. Open Items (For Final Sign‑Off)

* Export formats required for go-live (CSV only vs PDF vs Excel).
* Notification cadence defaults (daily vs weekly digest).
* Whether “criticality” is introduced in v1 to drive prioritization and escalation.
* Directory field list finalization beyond the guaranteed 6 (email recommended).
