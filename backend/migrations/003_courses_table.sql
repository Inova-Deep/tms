-- Wave 1: Create courses table and add FK columns
-- Idempotent migration using SQLite pragma checks

CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    code TEXT,
    type TEXT NOT NULL DEFAULT 'course',
    description TEXT,
    validity_months INTEGER DEFAULT 12,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add course_id to requirements table if not exists
INSERT INTO requirements (id, profile_id, name, type) 
SELECT '___temp_check___', 'temp', 'temp', 'temp' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('requirements') WHERE name='course_id'
);
DELETE FROM requirements WHERE id = '___temp_check___';
ALTER TABLE requirements ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to training_events table if not exists
INSERT INTO training_events (id, course_name, date) 
SELECT '___temp_check___', 'temp', '2000-01-01' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('training_events') WHERE name='course_id'
);
DELETE FROM training_events WHERE id = '___temp_check___';
ALTER TABLE training_events ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to evidence table if not exists
INSERT INTO evidence (id, employee_id, requirement_name, evidence_type, completion_date) 
SELECT '___temp_check___', 'temp', 'temp', 'temp', '2000-01-01' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('evidence') WHERE name='course_id'
);
DELETE FROM evidence WHERE id = '___temp_check___';
ALTER TABLE evidence ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to certifications table if not exists
INSERT INTO certifications (id, employee_id, requirement_name, issuer, issue_date) 
SELECT '___temp_check___', 'temp', 'temp', 'temp', '2000-01-01' 
WHERE NOT EXISTS (
    SELECT 1 FROM pragma_table_info('certifications') WHERE name='course_id'
);
DELETE FROM certifications WHERE id = '___temp_check___';
ALTER TABLE certifications ADD COLUMN course_id TEXT REFERENCES courses(id);

CREATE INDEX IF NOT EXISTS idx_requirements_course ON requirements(course_id);
CREATE INDEX IF NOT EXISTS idx_training_events_course ON training_events(course_id);
CREATE INDEX IF NOT EXISTS idx_evidence_course ON evidence(course_id);
CREATE INDEX IF NOT EXISTS idx_certifications_course ON certifications(course_id);
