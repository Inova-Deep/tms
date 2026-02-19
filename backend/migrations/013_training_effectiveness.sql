-- DML Phase 3: Training effectiveness evaluations (Section 4.4)
-- Manager confirms training outcomes applied in practice
-- method: observation | supervision | error_reduction | output_review | feedback
-- outcome: effective | partially_effective | ineffective
CREATE TABLE IF NOT EXISTS training_effectiveness_evaluations (
    id TEXT PRIMARY KEY,
    evidence_id TEXT REFERENCES evidence(id),
    evaluated_by TEXT NOT NULL,
    evaluation_date TEXT NOT NULL,
    method TEXT NOT NULL,
    outcome TEXT NOT NULL,
    follow_up_action TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_effectiveness_evidence ON training_effectiveness_evaluations(evidence_id);
