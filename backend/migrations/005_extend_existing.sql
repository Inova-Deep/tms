-- DML Phase 1: Extend employees and requirements with DML competence fields
ALTER TABLE employees ADD COLUMN IF NOT EXISTS employment_type TEXT DEFAULT 'permanent';

ALTER TABLE requirements ADD COLUMN IF NOT EXISTS skill_category TEXT;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS risk_level TEXT;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS criticality TEXT;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS mandatory_for_role BOOLEAN DEFAULT FALSE;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS training_type TEXT;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS assessment_method TEXT;
ALTER TABLE requirements ADD COLUMN IF NOT EXISTS evidence_type TEXT;
