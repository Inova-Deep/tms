Architecture decisions (locked)
Backend: Go 1.22+, chi v5
DB (demo): SQLite file, seed-on-startup if missing
Runtime: Docker (compose optional)
Auth: fake login toggle with 2 roles
TrainingAdmin (full write)
Employee (read-only, self-only)
Frontend: Vue 3 + Tailwind (your skeleton), shadcn-style components
Repo shape (recommended)
tms/
  backend/
    cmd/server/main.go
    internal/
      config/
      http/            # routers, middleware, handlers
      auth/            # fake auth, roles, user context
      db/              # sqlite open, migrations, seed
      domain/          # shared types (IDs, enums)
      employees/
      profiles/        # profile + version + requirements
      membership/      # manual add/exclude (rules optional in demo)
      catalogue/       # items (course/cert)
      events/          # training events + attendance
      evidence/        # cert uploads + completions
      compute/         # status computation + materialized tables
      dashboards/      # grid queries + drill-down
      audit/           # audit events table + writer
    migrations/
    seeds/
  frontend/            # your Vue skeleton here
  docker-compose.yml
  README.md
This keeps the demo “production-shaped” but not over-engineered.
SQLite seed-on-startup behaviour
On server start:
Ensure DB file exists (e.g., /data/tms.db)
Run migrations
If seed_version table absent/empty → run seed:
20 employees
4 profiles + requirements
memberships
~145 evidence records with mixed Valid/Expiring/Expired/Missing
2–3 events + attendance (optional)
Important: Use deterministic IDs so the demo script always works (e.g., employee IDs E1001..E1020, profile IDs P1..P4).
API surface needed for the demo (minimal)
Auth (fake)
POST /api/auth/login { role: "admin" | "employee", employeeId?: "E1005" }
Returns a simple token (or session cookie) with embedded role + employeeId.
Middleware injects ctx.user.
Employees
GET /api/employees?search=...
GET /api/employees/{id} (admin any, employee self-only)
GET /api/employees/{id}/matrix (derived requirements + statuses)
GET /api/employees/{id}/records (transcript table)
Profiles
GET /api/profiles
GET /api/profiles/{id} (requirements)
POST /api/profiles/{id}/members add employee
POST /api/profiles/{id}/exclusions exclude with reason+expiry
(Optional for demo) no editing profiles via UI; keep seeded
Events (attendance)
POST /api/events create event
POST /api/events/{id}/attendance/confirm confirm roster → generates completion records
Certifications
POST /api/certifications add cert evidence (attachment metadata only)
Dashboards
GET /api/dashboard/kpis?site=&department=&profile=&...
GET /api/dashboard/grid?mode=course|profile&groupBy=department|site|costCenter...
GET /api/dashboard/drilldown?... returns employees + statuses for clicked cell
Business rules the backend must implement (hard requirements)
Requirement status: Missing / Valid / Expiring / Expired
Expiring thresholds: support 90/60/30 (config)
Expiring counts as compliant in headline metric
Profile roll-up:
Eligible: all requirements compliant (Valid or Expiring or Waived)
At Risk: any Expiring OR any non-critical Missing (if you want; simplest: any Expiring)
Not Eligible: any Expired OR any Missing in required set
Only Admin can write.
Employee can only read self.