-- DML Phase 1: New starter induction records (Section 3.3)
-- HR conducts formal induction; signed attendance record retained
CREATE TABLE IF NOT EXISTS induction_records (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL REFERENCES employees(id),
    conducted_by TEXT NOT NULL,
    induction_date TEXT NOT NULL,
    topics_covered TEXT,
    employee_signature_ref TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_induction_employee ON induction_records(employee_id);
