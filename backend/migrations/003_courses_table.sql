-- Wave 1: Create courses table and add FK columns
-- This migration creates the courses table and adds course_id FK columns to related tables

-- Courses table
-- Central registry for all courses and certifications
CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,                    -- 'C001', 'C002', etc.
    name TEXT NOT NULL UNIQUE,              -- 'Manual Handling', 'H&S Induction'
    code TEXT,                              -- 'MH-001' (optional, nullable)
    type TEXT NOT NULL DEFAULT 'course',    -- 'course' | 'certification'
    description TEXT,                       -- nullable
    validity_months INTEGER DEFAULT 12,     -- 0 = never expires
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add course_id to requirements table
-- Links requirements to their corresponding course
ALTER TABLE requirements ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to training_events table
-- Links training events to the course being delivered
ALTER TABLE training_events ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to evidence table
-- Links evidence records to the course they satisfy
ALTER TABLE evidence ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Add course_id to certifications table
-- Links certifications to their corresponding course
ALTER TABLE certifications ADD COLUMN course_id TEXT REFERENCES courses(id);

-- Indexes for the new FK columns
CREATE INDEX IF NOT EXISTS idx_requirements_course ON requirements(course_id);
CREATE INDEX IF NOT EXISTS idx_training_events_course ON training_events(course_id);
CREATE INDEX IF NOT EXISTS idx_evidence_course ON evidence(course_id);
CREATE INDEX IF NOT EXISTS idx_certifications_course ON certifications(course_id);

-- Migration tracking
INSERT OR IGNORE INTO schema_migrations (version) VALUES ('003');
