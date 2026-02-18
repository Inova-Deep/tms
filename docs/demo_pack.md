# TMS Client Demo Pack (Minimal Showcase)

## 1) Demo Objective

Deliver a **profile-driven competency matrix** demo that proves TMS can:

1. Maintain an internal employee table (demo-only)
2. Assign competence **Profiles** (competency matrix) to employees
3. Record evidence via **attendance** and **certification upload**
4. Compute statuses (**Missing / Valid / Expiring / Expired**) and roll-up compliance
5. Visualise gaps via a **coverage grid** (course mode + profile mode) and drill-down

**Constraints:**

* No ERP integration in demo (local employee table)
* **Write access:** Training Admin only
* **Employee access:** read-only, self-only
* Evidence “attachments”: metadata only (no files)

---

## 2) Demo Storyline (2–4 minutes)

### Scene A — Create/select employee

* Open **Employees** → search and open **E1005 Jack Thompson**
* Show basic employee attributes (site/department/cost center)

### Scene B — Assign profiles (competency matrix)

* Assign profiles to Jack:

  * **Safety Base (All Employees)**
  * **Forklift Operator**
  * (Optional: show Welder profile exists but not assigned)
* Show the derived matrix instantly with item statuses.

### Scene C — Record evidence

1. **Attendance completion:**

   * Create event “Manual Handling” (Bristol / date)
   * Add Jack to roster
   * Confirm attendance → “Manual Handling” flips to **Valid**
2. **Certification upload:**

   * Upload “Forklift Licence (RTITB/ITSSAR)” metadata (issuer, issue date, expiry)
   * Status flips to **Valid / Expiring** based on expiry date

### Scene D — Monitor compliance

* Open **Coverage Dashboard**

  * Show KPIs (Overall compliance, Expiring soon, Expired)
  * Switch grid mode: **Course Coverage** → **Profile Coverage**
  * Click a red cell → drill-down list of affected employees

---

## 3) Demo Screens (Minimal Set)

### 3.1 Employees

* List + search
* Employee details page (read-only for employee role)
* Employee shows:

  * Assigned profiles
  * Requirements list (matrix)
  * Transcript (records with completion/expiry/status)

### 3.2 Profiles

* Profiles list
* Profile definition (requirements + validity)
* Membership:

  * Manual add/remove
  * Manual exclude (reason + expiry)

### 3.3 Training Events (Attendance)

* Create event: course, date, location, instructor employee_id
* Add roster (employee_id list)
* Confirm attendance → generates completion records

### 3.4 Certifications (Upload)

* Record: cert item, issuer, issue date, expiry date, notes
* Attachment: **metadata only** for demo

### 3.5 Dashboards

* KPI tiles:

  * Overall compliance % (employee-weighted; **Expiring counts as compliant**)
  * Total employees
  * Expiring soon
  * Expired
* Visual Gap Analysis grid modes:

  1. **Course Coverage:** Course × Department/Site with counts and %
  2. **Profile Coverage:** Profile × Department/Site showing Eligible / At Risk / Not Eligible
* Drill-down from any cell → employee list (export optional)

---

## 4) Demo Configuration

### 4.1 Sites (UK)

* **Bristol Central Plant**
* **Avonmouth Logistics & Workshop**

### 4.2 Departments (4)

* Manufacturing
* Maintenance
* Quality
* HR

### 4.3 Profiles (4)

#### P1 — Safety Base (All Employees)

Requirements (validity):

* H&S Induction (24 months)
* Fire Safety Awareness (12 months)
* First Aid Awareness (12 months)
* Manual Handling (12 months)
* COSHH Awareness (12 months)

#### P2 — Welder — Basic

Requirements (validity):

* Welding Safety & PPE (12 months)
* Hot Work Permit Awareness (12 months)
* Respiratory Fit Test (24 months)
* Welding Qualification (ISO 9606-1) (24 months)
* LOTO Awareness (12 months)

#### P3 — Forklift Operator

Requirements (validity):

* Forklift Theory (RTITB) (36 months)
* Forklift Practical Assessment (36 months)
* Forklift Licence (RTITB/ITSSAR) (36 months)
* Warehouse Pedestrian Safety (12 months)

#### P4 — HR Management (Compliance)

Requirements (validity):

* GDPR & Data Protection (24 months)
* Right to Work Checks (UK) (24 months)
* Equality, Diversity & Inclusion (24 months)
* DSE (Display Screen Equipment) (24 months)

---

## 5) Demo Employees (20, UK-style)

| ID    | Name           | Dept          | Site                           | Worker Type | Cost Center | Profiles                      |
| ----- | -------------- | ------------- | ------------------------------ | ----------- | ----------- | ----------------------------- |
| E1001 | Oliver Bennett | Manufacturing | Bristol Central Plant          | Employee    | CC100       | Safety Base; Welder           |
| E1002 | Amelia Clarke  | Manufacturing | Bristol Central Plant          | Employee    | CC100       | Safety Base; Forklift         |
| E1003 | Harry Patel    | Maintenance   | Bristol Central Plant          | Employee    | CC110       | Safety Base                   |
| E1004 | Isla Hughes    | Quality       | Bristol Central Plant          | Employee    | CC120       | Safety Base                   |
| E1005 | Jack Thompson  | Manufacturing | Avonmouth Logistics & Workshop | Employee    | CC200       | Safety Base; Welder; Forklift |
| E1006 | Emily Walker   | Maintenance   | Avonmouth Logistics & Workshop | Employee    | CC210       | Safety Base                   |
| E1007 | George Evans   | Quality       | Avonmouth Logistics & Workshop | Employee    | CC220       | Safety Base                   |
| E1008 | Sophia Khan    | HR            | Bristol Central Plant          | Employee    | CC130       | Safety Base; HR Mgmt          |
| E1009 | Noah Wilson    | Manufacturing | Bristol Central Plant          | Contractor  | CC100       | Safety Base                   |
| E1010 | Ava Morgan     | Manufacturing | Avonmouth Logistics & Workshop | Employee    | CC200       | Safety Base; Forklift         |
| E1011 | Leo Robinson   | Maintenance   | Bristol Central Plant          | Employee    | CC110       | Safety Base; Forklift         |
| E1012 | Mia Carter     | Quality       | Bristol Central Plant          | Employee    | CC120       | Safety Base                   |
| E1013 | Ethan James    | Manufacturing | Avonmouth Logistics & Workshop | Employee    | CC200       | Safety Base; Welder           |
| E1014 | Grace Murphy   | HR            | Avonmouth Logistics & Workshop | Employee    | CC230       | Safety Base; HR Mgmt          |
| E1015 | Freddie Brown  | Manufacturing | Bristol Central Plant          | Employee    | CC100       | Safety Base; Welder           |
| E1016 | Lily Ward      | Maintenance   | Avonmouth Logistics & Workshop | Contractor  | CC210       | Safety Base; Forklift         |
| E1017 | Oscar Lewis    | Quality       | Avonmouth Logistics & Workshop | Employee    | CC220       | Safety Base                   |
| E1018 | Chloe Smith    | Manufacturing | Bristol Central Plant          | Employee    | CC100       | Safety Base; Welder           |
| E1019 | Charlie Green  | Maintenance   | Bristol Central Plant          | Employee    | CC110       | Safety Base                   |
| E1020 | Ella Davies    | Manufacturing | Avonmouth Logistics & Workshop | Employee    | CC200       | Safety Base; Forklift         |

---

## 6) Pre-Seeded Records (to make dashboards “alive”)

### 6.1 Volume and distribution

* Total requirement instances (all employees, all assigned profiles): **157**
* Evidence/records seeded (admin-created): **145** (100+ as requested)
* Status distribution (requirement-level):

  * **Valid:** 122
  * **Expiring:** 13
  * **Expired:** 10
  * **Missing:** 12

### 6.2 KPI snapshot (demo landing)

* Total employees: **20**
* Overall compliance (employee-weighted; Expiring counts compliant): **86.0%**
* Expiring soon: **13**
* Expired: **10**

### 6.3 Breakdown examples (for filters)

By Profile (Valid / Expiring / Expired / Missing):

* Safety Base: 80 / 8 / 5 / 7
* Welder: 17 / 1 / 3 / 4
* Forklift: 18 / 3 / 2 / 1
* HR Mgmt: 7 / 1 / 0 / 0

By Department (Valid / Expiring / Expired / Missing):

* Manufacturing: 64 / 6 / 7 / 9
* Maintenance: 26 / 4 / 1 / 2
* Quality: 17 / 1 / 1 / 1
* HR: 15 / 2 / 1 / 0

---

## 7) Demo “Wow” Moments (scripted)

1. Assign **Welder — Basic** to an employee → matrix populates instantly.
2. Upload an external cert with expiry next month → status becomes **Expiring** immediately.
3. Confirm attendance on a training event → completion record generated and grid updates.
4. Switch the dashboard from **Course Coverage** to **Profile Coverage** → shows eligibility by org slice.

---

## 8) Minimal Viable Implementation Notes (for the demo build)

* Keep computation deterministic and visible:

  * status rules: Missing/Valid/Expiring/Expired
  * profile roll-up: Eligible/At Risk/Not Eligible
* Keep write surfaces limited to Training Admin only.
* Evidence files are mocked as metadata.
* Drill-down from grid cells is essential (even if export is not).

---

## 9) Optional Enhancements (only if time allows)

* “vNext effective date preview” widget (impact: new gaps and planned assignments 30 days ahead)
* Simple “Create Planned Assignment” view (v1 planning)
* Audit pack export (CSV/PDF) placeholder

---

## 10) Demo Talk-Track (Presenter Script)

### Timing

* **2 minutes** (tight) or **4 minutes** (with drill-down + filters).

### Opening (10–15s)

* “This is **TMS**, Deep Manufacturing’s training and competency matrix system. It’s profile-driven: assign a profile like ‘Forklift Operator’ and TMS instantly shows what’s required, what’s missing, and what’s expiring — with full evidence records and auditability.”

### Step 1 — Employee (20–30s)

1. Navigate: **Employees → Search: Jack Thompson (E1005)**
2. Say:

   * “For the demo we’re using an internal employee table, but in production TMS can connect to your employee directory.”
3. Show fields:

   * Site, Department, Cost Center, Worker Type
4. Point out:

   * “Employees have read-only access to their own records; only Training Admin can modify records.”

### Step 2 — Assign Profiles (30–60s)

1. In Jack’s page: **Assigned Profiles → Add Profile**
2. Add:

   * Safety Base (All Employees)
   * Forklift Operator
3. Say:

   * “The competency matrix is derived from the profile definitions — no spreadsheets.”
4. Show the Requirements Matrix and call out a mix of statuses:

   * “Here we have **Valid**, **Expiring**, and **Missing** requirements.”
5. Optional (if asked):

   * “Profiles are versioned and can be published with a future effective date for controlled change.”

### Step 3 — Record Evidence (60–120s)

#### Option A: Attendance completion (fast)

1. Navigate: **Training Events → Create Event**
2. Choose:

   * Course: Manual Handling
   * Site: Avonmouth Logistics & Workshop
   * Date: today
   * Instructor: choose an employee_id
3. Add roster:

   * Add Jack (E1005)
4. Click: **Confirm Attendance**
5. Say:

   * “Once attendance is confirmed, TMS generates the completion record and recalculates status immediately.”

#### Option B: Certification upload (high impact)

1. Navigate: **Certifications → Add Record** (or on Jack’s page)
2. Choose:

   * Cert item: Forklift Licence (RTITB/ITSSAR)
   * Issuer: RTITB
   * Issue date: e.g., 2024-03-01
   * Expiry date: e.g., 2027-03-01
3. Save
4. Say:

   * “External certifications are stored in TMS with expiry tracking — attachment can be mandatory in production.”

### Step 4 — Dashboards & Gap Analysis (120–240s)

1. Navigate: **Coverage Dashboard**
2. Say:

   * “This dashboard is designed for planning — you can see overall compliance, what’s expiring, and where the gaps are by site and department.”
3. Show KPI tiles:

   * Overall compliance, Expiring soon, Expired
4. Show grid mode:

   * Start with **Course Coverage** (Course × Department)
   * Switch to **Profile Coverage** (Profile × Department)
5. Drill-down:

   * Click a red cell and say:
   * “Drill-down shows exactly which employees are missing or expired for that requirement, so you can schedule training.”

### Closing (10–15s)

* “That’s the core: profiles define the competency matrix, admin records training and evidence, and TMS gives you real-time compliance and gap visibility with audit-ready records.”

---

## 11) Seed-Data Checklist (Admin Setup in < 30 minutes)

### A) Create reference data

* Create 4 Departments: Manufacturing, Maintenance, Quality, HR
* Create 2 Sites: Bristol Central Plant; Avonmouth Logistics & Workshop
* Create Cost Centers: CC100, CC110, CC120, CC130, CC200, CC210, CC220, CC230

### B) Create catalogue items (Courses + Certifications)

Create the following items with validity:

* H&S Induction (24 months)
* Fire Safety Awareness (12 months)
* First Aid Awareness (12 months)
* Manual Handling (12 months)
* COSHH Awareness (12 months)
* Welding Safety & PPE (12 months)
* Hot Work Permit Awareness (12 months)
* Respiratory Fit Test (24 months)
* Welding Qualification (ISO 9606-1) (24 months)
* LOTO Awareness (12 months)
* Forklift Theory (RTITB) (36 months)
* Forklift Practical Assessment (36 months)
* Forklift Licence (RTITB/ITSSAR) (36 months)
* Warehouse Pedestrian Safety (12 months)
* GDPR & Data Protection (24 months)
* Right to Work Checks (UK) (24 months)
* Equality, Diversity & Inclusion (24 months)
* DSE (Display Screen Equipment) (24 months)

### C) Create 4 Profiles and attach requirements

* P1 Safety Base (All Employees): 5 items
* P2 Welder — Basic: 5 items
* P3 Forklift Operator: 4 items
* P4 HR Management (Compliance): 4 items

### D) Create 20 employees

* Add the 20 employee records from the employee table (Section 5).

### E) Assign profiles (bulk)

* Assign Safety Base to all employees
* Assign Welder to E1001, E1005, E1013, E1015, E1018
* Assign Forklift to E1002, E1005, E1010, E1011, E1016, E1020
* Assign HR Mgmt to E1008, E1014

### F) Seed evidence (145 records)

Seed completions/certs to achieve a mixed dashboard:

* For most Safety Base items: set **Valid** for all, but:

  * 7 Missing
  * 5 Expired
  * 8 Expiring
* For Welder items:

  * 4 Missing
  * 3 Expired
  * 1 Expiring
* For Forklift items:

  * 1 Missing
  * 2 Expired
  * 3 Expiring
* For HR Mgmt items:

  * mostly Valid; 1 Expiring

**Tip:** Seed “Expiring” by setting expiry dates within the next 30–60 days.

### G) Create 2–3 training events

* Manual Handling (Avonmouth)
* Fire Safety Awareness (Bristol)
* Forklift Practical Assessment (Avonmouth)

Add rosters and confirm attendance for 5–10 people to show the event flow.

---

## 12) Demo Delivery Notes

* Keep the UI uncluttered: the grid + drill-down is the headline.
* Avoid deep configuration screens during the presentation; pre-seed everything.
* Use one employee as the hero (E1005 Jack Thompson) for repeatable flow.
* If asked about integrations, say:

  * “For the demo we keep it standalone; production uses your directory for auto-assignment rules.”
