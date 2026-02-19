package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type ApplicabilityRepository struct {
	DB *sql.DB
}

func (r *ApplicabilityRepository) Create(roleID string, req CreateApplicabilityDecisionRequest) (*ApplicabilityDecision, error) {
	id := fmt.Sprintf("APPL-%d", time.Now().UnixNano())

	// Compute result from Q answers (DML-QA-REG-5024-1 5-question gate):
	// result = "matrix" if Q5 AND (Q1 OR Q2) AND Q3 AND Q4
	// result = "awareness" otherwise
	result := "awareness"
	if req.Q5 && (req.Q1 || req.Q2) && req.Q3 && req.Q4 {
		result = "matrix"
	}

	var escalationNotes sql.NullString
	if req.EscalationNotes != "" {
		escalationNotes = sql.NullString{String: req.EscalationNotes, Valid: true}
	}

	_, err := r.DB.Exec(`
		INSERT INTO applicability_decisions
		  (id, role_id, q1, q2, q3, q4, q5, result, decided_by, decided_at,
		   escalated, escalation_notes, escalation_resolved_at, escalation_resolved_by)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(),
		   $10, $11, NULL, NULL)
	`,
		id, roleID, req.Q1, req.Q2, req.Q3, req.Q4, req.Q5,
		result, req.DecidedBy, req.Escalated, escalationNotes,
	)
	if err != nil {
		return nil, err
	}

	// Update the role's matrix_applicable flag to match the computed result.
	_, err = r.DB.Exec(
		"UPDATE roles SET matrix_applicable = $1, updated_at = NOW() WHERE id = $2",
		result == "matrix",
		roleID,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByRoleID(roleID)
}

func (r *ApplicabilityRepository) GetByRoleID(roleID string) (*ApplicabilityDecision, error) {
	row := r.DB.QueryRow(`
		SELECT id, role_id, q1, q2, q3, q4, q5, result, decided_by, decided_at,
		       escalated, escalation_notes, escalation_resolved_at, escalation_resolved_by
		FROM applicability_decisions
		WHERE role_id = $1
		ORDER BY decided_at DESC
		LIMIT 1
	`, roleID)

	var d ApplicabilityDecision
	var decidedAt sql.NullString
	var escalationNotes, escalationResolvedAt, escalationResolvedBy sql.NullString

	err := row.Scan(
		&d.ID, &d.RoleID, &d.Q1, &d.Q2, &d.Q3, &d.Q4, &d.Q5,
		&d.Result, &d.DecidedBy, &decidedAt,
		&d.Escalated, &escalationNotes, &escalationResolvedAt, &escalationResolvedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if decidedAt.Valid {
		d.DecidedAt = decidedAt.String
	}
	if escalationNotes.Valid {
		d.EscalationNotes = &escalationNotes.String
	}
	if escalationResolvedAt.Valid {
		d.EscalationResolvedAt = &escalationResolvedAt.String
	}
	if escalationResolvedBy.Valid {
		d.EscalationResolvedBy = &escalationResolvedBy.String
	}
	return &d, nil
}

func (r *ApplicabilityRepository) ResolveEscalation(id string, req ResolveEscalationRequest) (*ApplicabilityDecision, error) {
	_, err := r.DB.Exec(`
		UPDATE applicability_decisions
		SET escalation_resolved_at = NOW(),
		    escalation_resolved_by = $1,
		    escalation_notes = COALESCE(escalation_notes || E'\n' || $2, $2)
		WHERE id = $3
	`, req.ResolvedBy, req.Notes, id)
	if err != nil {
		return nil, err
	}

	// Re-fetch by id using a direct query.
	row := r.DB.QueryRow(`
		SELECT id, role_id, q1, q2, q3, q4, q5, result, decided_by, decided_at,
		       escalated, escalation_notes, escalation_resolved_at, escalation_resolved_by
		FROM applicability_decisions
		WHERE id = $1
	`, id)

	var d ApplicabilityDecision
	var decidedAt sql.NullString
	var escalationNotes, escalationResolvedAt, escalationResolvedBy sql.NullString

	err = row.Scan(
		&d.ID, &d.RoleID, &d.Q1, &d.Q2, &d.Q3, &d.Q4, &d.Q5,
		&d.Result, &d.DecidedBy, &decidedAt,
		&d.Escalated, &escalationNotes, &escalationResolvedAt, &escalationResolvedBy,
	)
	if err != nil {
		return nil, err
	}

	if decidedAt.Valid {
		d.DecidedAt = decidedAt.String
	}
	if escalationNotes.Valid {
		d.EscalationNotes = &escalationNotes.String
	}
	if escalationResolvedAt.Valid {
		d.EscalationResolvedAt = &escalationResolvedAt.String
	}
	if escalationResolvedBy.Valid {
		d.EscalationResolvedBy = &escalationResolvedBy.String
	}
	return &d, nil
}
