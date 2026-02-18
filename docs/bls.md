# TMS v1 Business Logic Specification (Locked)

## 0) Principles

* **TMS is the system of record** for training/certification evidence and computed eligibility.
* **Only Training Admins** can create/modify records. All other users are read-only.
* **Employees can only view their own records**, profiles, requirements, and statuses.
* Profiles are the core requirements model; membership is derived from **directory rules + manual lists + manual exclusions**.

---

## 1) Status Model (Requirement Level)

Each requirement for an employee resolves to exactly one status:

1. **Missing**

* No completion/evidence exists for this requirement (or any valid equivalency item).

2. **Valid**

* Completion/evidence exists and:

  * `never_expires = true`, **OR**
  * `expiry_date > (today + expiring_threshold)`

3. **Expiring**

* Completion/evidence exists and `expiry_date` is within the configured expiring window.
* Default thresholds supported: **90 / 60 / 30** days (optionally 14 / 7 for critical items).

4. **Expired**

* Completion/evidence exists and `expiry_date < today`.

**Note:** “Invalid” is not used.

---

## 2) Requirement Satisfaction Rules

### 2.1 Equivalency (OR groups)

A requirement may be satisfied by one of multiple items (e.g., **C1 OR C3**).

* If **any** equivalency item is **Valid** or **Expiring**, the requirement is considered satisfied.

### 2.2 Course-driven validity policy

* Validity/expiry is defined by the **course/cert item**, not by the profile.
* Each item defines one of:

  * `validity_duration` (e.g., 12 months from completion), **or**
  * `never_expires`, **or**
  * evidence-supplied `expiry_date` (common for third-party certifications).

### 2.3 Conflict rule across profiles

If multiple profiles require the same item:

* Status for that item is computed once per employee and reused.
* **Earliest expiry wins** when determining the risk posture.

### 2.4 Evidence selection when multiple records exist

If multiple evidence records exist for the same item:

* Use the record that yields the **most current validity** (valid beats expired; typically the record with the latest valid expiry).
* Retain all historical records for audit.

---

## 3) Waivers (Exemptions)

* **Allowed:** Yes
* **Who can waive:** **Training Manager only**
* **Waiver rules:**

  * Reason is mandatory
  * Waiver expiry date is mandatory
* Waivers must be explicitly visible in employee views, admin views, reports, and audit packs.

---

## 4) Profile Versioning and Effective Dates

### 4.1 Membership vs requirement timing

* **Membership rules apply continuously** (directory changes can add/remove members anytime).
* Requirements apply based on the **active profile version** (effective-date controlled).

### 4.2 Publish model

Profiles are versioned (v1, v2, …) with lifecycle:

* Draft → Published (**Effective From** date) → Retired

### 4.3 Assignment generation lead time

For a future effective version (vNext):

* Generate **Planned** assignments **30 days before** Effective From.
* If an employee becomes a profile member within the 30-day window, generate assignments **immediately**.

### 4.4 Handling prior completions when requirements change

* Requirement substitutions are supported via **configurable equivalency mapping** to avoid unnecessary retraining.

---

## 5) Profile Membership Logic (Hybrid)

Membership sources:

* **RULE** (directory rule)
* **MANUAL_ADD** (explicit list)
* **MANUAL_EXCLUDE** (override)

### 5.1 Manual add persistence

* Manual adds **persist** and are not auto-removed by directory/rule changes.

### 5.2 Manual exclusion behavior

* Manual exclusions override rule membership.
* Requires:

  * Reason (mandatory)
  * Exclusion expiry date (mandatory)

When exclusion expires:

* Employee is re-evaluated for rule membership and re-added if they match.
* Notify Training Admin/Training Manager that the exclusion expired.

### 5.3 Leaver/mover handling

* If `employment_status` becomes inactive/terminated:

  * Stop notifications
  * Freeze/close planned assignments (implementation choice)
  * Retain records and audit trail

* If org/location/cost center changes:

  * Recompute rule memberships
  * Update derived requirements and gaps immediately

---

## 6) Assignments (v1)

Assignments are planning artifacts (computation is evidence-driven).

### 6.1 When assignments are created

Create a **Planned** assignment when:

* employee becomes member of a profile with a **Missing** requirement
* a required item is **Expiring** and needs renewal
* vNext introduces a new requirement (30-day lead rule)

### 6.2 Assignment states (v1)

* **Planned** only

---

## 7) Training Delivery Types

### 7.1 Instructor-led (attendance-based)

* Training Admin creates a Training Event with:

  * course, date, location
  * instructor recorded by **employee_id** (no instructor accounts)
  * attendance roster (employee list)
* **Attendance confirmation = completion** by default.
* Completions are generated automatically upon attendance confirmation.

### 7.2 Third-party certification (upload-based)

Certification evidence record must include:

* issuer (text)
* issue date
* expiry date (or never)
* **attachment mandatory**

---

## 8) Notifications

Recipients:

* Training Admin + Training Manager (always)
* Employee (if email exists)

Rules:

* Expiring thresholds supported: **90 / 60 / 30** days
* Cadence:

  * scheduled digest (daily/weekly configurable)
  * overdue reminders weekly until resolved

---

## 9) Compliance Calculations & Dashboards

### 9.1 Overall compliance (default)

* Default metric is **employee-weighted**:

  * `valid_requirements / total_required_requirements`

### 9.2 Expiring counting rule

* **Expiring counts as compliant** in the headline metric.
* Expiring items are also surfaced as **At Risk** (separate KPI / grid signal).

### 9.3 Toggles

* **Critical-only** compliance
* **By profile** compliance (Eligible / At Risk / Not Eligible)

### 9.4 Visual gap analysis grid modes

* **Course Coverage Grid:** Course × Department with Valid/Missing/Expiring/Expired counts
* **Profile Coverage Grid:** Profile × Department showing Eligible/At Risk/Not Eligible

---

## 10) Retention

* Evidence records and audit logs are retained **per Deep Manufacturing retention policy**.
