-- Backfill courses table and update FK references
-- Idempotent: all operations have ON CONFLICT / WHERE IS NULL guards

INSERT INTO courses (id, name, type, validity_months)
SELECT
    'C' || LPAD(row_number::text, 3, '0') AS id,
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
    GROUP BY name, validity_months
) sub
ON CONFLICT DO NOTHING;

UPDATE requirements
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = requirements.name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = requirements.name);

UPDATE training_events
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = training_events.course_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = training_events.course_name);

UPDATE evidence
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = evidence.requirement_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = evidence.requirement_name);

UPDATE certifications
SET course_id = (
    SELECT c.id
    FROM courses c
    WHERE c.name = certifications.requirement_name
)
WHERE course_id IS NULL
  AND EXISTS (SELECT 1 FROM courses c WHERE c.name = certifications.requirement_name);
