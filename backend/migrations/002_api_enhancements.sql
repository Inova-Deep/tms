-- Add time and duration to training_events (idempotent)
-- SQLite doesn't support IF NOT EXISTS for ADD COLUMN, so we check pragmas

-- Add time column if not exists
INSERT INTO training_events (id, course_name, date) 
SELECT '___temp_check___', 'temp', '2000-01-01' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('training_events') WHERE name='time'
);
DELETE FROM training_events WHERE id = '___temp_check___';
ALTER TABLE training_events ADD COLUMN time TEXT DEFAULT '';

-- Add duration column if not exists  
INSERT INTO training_events (id, course_name, date) 
SELECT '___temp_check___', 'temp', '2000-01-01' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('training_events') WHERE name='duration'
);
DELETE FROM training_events WHERE id = '___temp_check___';
ALTER TABLE training_events ADD COLUMN duration INTEGER DEFAULT 0;

-- Recreate profile_exclusions with an ID primary key for individual CRUD
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

INSERT OR IGNORE INTO profile_exclusions_new (
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
FROM profile_exclusions;

DROP TABLE IF EXISTS profile_exclusions;

ALTER TABLE profile_exclusions_new RENAME TO profile_exclusions;

CREATE INDEX IF NOT EXISTS idx_exclusions_profile ON profile_exclusions (profile_id);
CREATE INDEX IF NOT EXISTS idx_exclusions_employee ON profile_exclusions (employee_id);
