-- DML Phase 1: Work authorization layer (Section 3.2.1)
-- Manager sign-off required before independent work
-- is_legacy=true marks records migrated from pre-DML system (legacy_pending_confirmation state)
CREATE TABLE IF NOT EXISTS work_authorizations (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL REFERENCES employees(id),
    requirement_id TEXT REFERENCES requirements(id),
    authorization_type TEXT NOT NULL,
    authorization_state TEXT NOT NULL,
    authorized_by TEXT,
    authorized_at TIMESTAMP,
    expiry_date TEXT,
    revoked_at TIMESTAMP,
    revoke_reason TEXT,
    is_legacy BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_work_auth_employee ON work_authorizations(employee_id);
CREATE INDEX IF NOT EXISTS idx_work_auth_requirement ON work_authorizations(requirement_id);
CREATE INDEX IF NOT EXISTS idx_work_auth_state ON work_authorizations(authorization_state);
