-- Create courses table and add FK columns

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

-- Add course_id FK columns using PostgreSQL IF NOT EXISTS syntax
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS course_id TEXT REFERENCES courses(id);
ALTER TABLE training_events ADD COLUMN IF NOT EXISTS course_id TEXT REFERENCES courses(id);
ALTER TABLE evidence ADD COLUMN IF NOT EXISTS course_id TEXT REFERENCES courses(id);
ALTER TABLE certifications ADD COLUMN IF NOT EXISTS course_id TEXT REFERENCES courses(id);

CREATE INDEX IF NOT EXISTS idx_requirements_course ON requirements(course_id);
CREATE INDEX IF NOT EXISTS idx_training_events_course ON training_events(course_id);
CREATE INDEX IF NOT EXISTS idx_evidence_course ON evidence(course_id);
CREATE INDEX IF NOT EXISTS idx_certifications_course ON certifications(course_id);
