-- DML Phase 1: Roles as first-class entity (distinct from profiles)
-- Roles = job functions with applicability decisions
-- Profiles = collections of training requirements
CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    department TEXT,
    skill_category TEXT,
    risk_level TEXT,
    criticality TEXT,
    mandatory_for_role BOOLEAN DEFAULT FALSE,
    matrix_applicable BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_roles_department ON roles(department);
CREATE INDEX IF NOT EXISTS idx_roles_matrix ON roles(matrix_applicable);
