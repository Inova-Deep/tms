-- DML Phase 2: Supervision periods (Section 3.2.4)
-- Tracks supervised work with scope limits, supervisor assignment, completion evidence
CREATE TABLE IF NOT EXISTS supervision_periods (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL REFERENCES employees(id),
    supervisor_id TEXT NOT NULL REFERENCES employees(id),
    requirement_id TEXT REFERENCES requirements(id),
    scope_limitations TEXT,
    start_date TEXT NOT NULL,
    end_date TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    completion_assessment_id TEXT REFERENCES competence_assessments(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_supervision_employee ON supervision_periods(employee_id);
CREATE INDEX IF NOT EXISTS idx_supervision_supervisor ON supervision_periods(supervisor_id);
CREATE INDEX IF NOT EXISTS idx_supervision_status ON supervision_periods(status);
