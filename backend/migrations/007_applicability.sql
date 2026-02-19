-- DML Phase 1: 5-question applicability decisions (DML-QA-REG-5024-1, Section 3.2.2)
-- Determines whether a role enters the Skills Matrix or Awareness track only
CREATE TABLE IF NOT EXISTS applicability_decisions (
    id TEXT PRIMARY KEY,
    role_id TEXT NOT NULL REFERENCES roles(id),
    q1 BOOLEAN NOT NULL,
    q2 BOOLEAN NOT NULL,
    q3 BOOLEAN NOT NULL,
    q4 BOOLEAN NOT NULL,
    q5 BOOLEAN NOT NULL,
    result TEXT NOT NULL,
    decided_by TEXT NOT NULL,
    decided_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    escalated BOOLEAN DEFAULT FALSE,
    escalation_notes TEXT,
    escalation_resolved_at TIMESTAMP,
    escalation_resolved_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_applicability_role ON applicability_decisions(role_id);
