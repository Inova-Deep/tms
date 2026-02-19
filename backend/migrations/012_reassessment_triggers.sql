-- DML Phase 3: Non-calendar reassessment triggers (Section 3.2.5)
-- trigger_type: repeated_errors | process_change | extended_absence | cert_expiry | performance_concern
CREATE TABLE IF NOT EXISTS reassessment_triggers (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL REFERENCES employees(id),
    requirement_id TEXT REFERENCES requirements(id),
    trigger_type TEXT NOT NULL,
    triggered_by TEXT NOT NULL,
    triggered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    resolved_at TIMESTAMP,
    resolution_notes TEXT,
    resolution_assessment_id TEXT REFERENCES competence_assessments(id)
);

CREATE INDEX IF NOT EXISTS idx_reassessment_employee ON reassessment_triggers(employee_id);
CREATE INDEX IF NOT EXISTS idx_reassessment_type ON reassessment_triggers(trigger_type);
CREATE INDEX IF NOT EXISTS idx_reassessment_open ON reassessment_triggers(resolved_at) WHERE resolved_at IS NULL;
