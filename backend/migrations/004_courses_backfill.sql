-- Wave 1: Backfill courses table and update FK references
-- This migration extracts unique courses from requirements and populates the courses table

-- Step 1: Insert unique courses from requirements table
-- Generate IDs as C001, C002, etc.
-- Set type to 'certification' if name contains 'Licence' or 'Certification'
INSERT INTO courses (id, name, type, validity_months)
SELECT
    'C' || printf('%03d', row_number) AS id,
    name,
    CASE
        WHEN LOWER(name) LIKE '%licence%' OR LOWER(name) LIKE '%certification%'
        THEN 'certification'
        ELSE 'course'
    END AS type,
    validity_months
FROM (
    SELECT
        name,
        validity_months,
        ROW_NUMBER() OVER (ORDER BY name) AS row_number
    FROM requirements
    GROUP BY name
);

-- Step 2: Backfill course_id in requirements table
-- Match by name
UPDATE requirements
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = requirements.name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = requirements.name);

-- Step 3: Backfill course_id in training_events table
-- Match by course_name
UPDATE training_events
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = training_events.course_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = training_events.course_name);

-- Step 4: Backfill course_id in evidence table
-- Match by requirement_name
UPDATE evidence
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = evidence.requirement_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = evidence.requirement_name);

-- Step 5: Backfill course_id in certifications table
-- Match by requirement_name
UPDATE certifications
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = certifications.requirement_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = certifications.requirement_name);

-- Note: Old columns (requirements.name, training_events.course_name, etc.)
-- are NOT dropped yet. They will be removed in a later migration after verification.

-- Migration tracking
INSERT OR IGNORE INTO schema_migrations (version) VALUES ('004');
