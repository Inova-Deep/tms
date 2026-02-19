package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type TrainingEffectivenessRepository struct {
	DB *sql.DB
}

func (r *TrainingEffectivenessRepository) Create(req CreateEffectivenessEvaluationRequest) (*TrainingEffectivenessEvaluation, error) {
	id := fmt.Sprintf("TEE-%d", time.Now().UnixNano())

	var evidenceID sql.NullString
	if req.EvidenceID != nil {
		evidenceID = sql.NullString{String: *req.EvidenceID, Valid: true}
	}

	var followUpAction sql.NullString
	if req.FollowUpAction != nil {
		followUpAction = sql.NullString{String: *req.FollowUpAction, Valid: true}
	}

	_, err := r.DB.Exec(`
		INSERT INTO training_effectiveness_evaluations
		  (id, evidence_id, evaluated_by, evaluation_date, method,
		   outcome, follow_up_action, created_at)
		VALUES
		  ($1, $2, $3, $4, $5,
		   $6, $7, NOW())
	`,
		id, evidenceID, req.EvaluatedBy, req.EvaluationDate, req.Method,
		req.Outcome, followUpAction,
	)
	if err != nil {
		return nil, err
	}

	// Re-fetch the just-inserted row.
	row := r.DB.QueryRow(`
		SELECT id, evidence_id, evaluated_by, evaluation_date, method,
		       outcome, follow_up_action, created_at
		FROM training_effectiveness_evaluations
		WHERE id = $1
	`, id)

	return scanEffectivenessRow(row)
}

func (r *TrainingEffectivenessRepository) GetByEvidence(evidenceID string) ([]TrainingEffectivenessEvaluation, error) {
	rows, err := r.DB.Query(`
		SELECT id, evidence_id, evaluated_by, evaluation_date, method,
		       outcome, follow_up_action, created_at
		FROM training_effectiveness_evaluations
		WHERE evidence_id = $1
		ORDER BY evaluation_date DESC
	`, evidenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]TrainingEffectivenessEvaluation, 0)
	for rows.Next() {
		var eval TrainingEffectivenessEvaluation
		var evID, followUpAction sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&eval.ID, &evID, &eval.EvaluatedBy, &eval.EvaluationDate, &eval.Method,
			&eval.Outcome, &followUpAction, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		if evID.Valid {
			eval.EvidenceID = &evID.String
		}
		if followUpAction.Valid {
			eval.FollowUpAction = &followUpAction.String
		}
		if createdAt.Valid {
			eval.CreatedAt = createdAt.String
		}

		results = append(results, eval)
	}
	return results, nil
}

// scanEffectivenessRow scans a single *sql.Row into a TrainingEffectivenessEvaluation.
func scanEffectivenessRow(row *sql.Row) (*TrainingEffectivenessEvaluation, error) {
	var eval TrainingEffectivenessEvaluation
	var evidenceID, followUpAction sql.NullString
	var createdAt sql.NullString

	err := row.Scan(
		&eval.ID, &evidenceID, &eval.EvaluatedBy, &eval.EvaluationDate, &eval.Method,
		&eval.Outcome, &followUpAction, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if evidenceID.Valid {
		eval.EvidenceID = &evidenceID.String
	}
	if followUpAction.Valid {
		eval.FollowUpAction = &followUpAction.String
	}
	if createdAt.Valid {
		eval.CreatedAt = createdAt.String
	}
	return &eval, nil
}
