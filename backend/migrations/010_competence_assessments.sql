-- DML Phase 2: Formal competence assessment records (Section 3.2.3)
-- Required for critical work: records person assessed, activities, assessor, limitations
CREATE TABLE IF NOT EXISTS competence_assessments (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL REFERENCES employees(id),
    requirement_id TEXT REFERENCES requirements(id),
    assessed_by TEXT NOT NULL,
    assessment_date TEXT NOT NULL,
    work_activities_covered TEXT NOT NULL,
    limitations TEXT,
    supervision_required BOOLEAN DEFAULT FALSE,
    outcome TEXT NOT NULL,
    evidence_reference TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_assessment_employee ON competence_assessments(employee_id);
CREATE INDEX IF NOT EXISTS idx_assessment_requirement ON competence_assessments(requirement_id);
