-- Add time and duration to training_events
ALTER TABLE training_events ADD COLUMN IF NOT EXISTS time TEXT DEFAULT '';
ALTER TABLE training_events ADD COLUMN IF NOT EXISTS duration INTEGER DEFAULT 0;

-- Recreate profile_exclusions with an ID primary key for individual CRUD.
-- PostgreSQL approach: create new table, migrate data, drop old, rename.
CREATE TABLE IF NOT EXISTS profile_exclusions_new (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL,
    profile_id TEXT NOT NULL,
    reason TEXT NOT NULL,
    expiry_date TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT DEFAULT 'admin',
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

INSERT INTO profile_exclusions_new (
    id,
    employee_id,
    profile_id,
    reason,
    expiry_date,
    created_at
)
SELECT
    'EX-' || employee_id || '-' || profile_id,
    employee_id,
    profile_id,
    reason,
    expiry_date,
    created_at
FROM profile_exclusions
ON CONFLICT DO NOTHING;

DROP TABLE IF EXISTS profile_exclusions;

ALTER TABLE profile_exclusions_new RENAME TO profile_exclusions;

CREATE INDEX IF NOT EXISTS idx_exclusions_profile ON profile_exclusions (profile_id);
CREATE INDEX IF NOT EXISTS idx_exclusions_employee ON profile_exclusions (employee_id);
