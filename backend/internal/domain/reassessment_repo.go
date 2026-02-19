package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type ReassessmentRepository struct {
	DB *sql.DB
}

func (r *ReassessmentRepository) Create(req CreateReassessmentTriggerRequest) (*ReassessmentTrigger, error) {
	id := fmt.Sprintf("REAST-%d", time.Now().UnixNano())

	var requirementID sql.NullString
	if req.RequirementID != nil {
		requirementID = sql.NullString{String: *req.RequirementID, Valid: true}
	}

	_, err := r.DB.Exec(`
		INSERT INTO reassessment_triggers
		  (id, employee_id, requirement_id, trigger_type, triggered_by, triggered_at,
		   notes, resolved_at, resolution_notes, resolution_assessment_id)
		VALUES
		  ($1, $2, $3, $4, $5, NOW(),
		   $6, NULL, NULL, NULL)
	`,
		id, req.EmployeeID, requirementID, req.TriggerType, req.TriggeredBy,
		req.Notes,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *ReassessmentRepository) GetOpen() ([]ReassessmentTrigger, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, trigger_type, triggered_by, triggered_at,
		       notes, resolved_at, resolution_notes, resolution_assessment_id
		FROM reassessment_triggers
		WHERE resolved_at IS NULL
		ORDER BY triggered_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReassessmentRows(rows)
}

func (r *ReassessmentRepository) GetByEmployee(employeeID string) ([]ReassessmentTrigger, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, trigger_type, triggered_by, triggered_at,
		       notes, resolved_at, resolution_notes, resolution_assessment_id
		FROM reassessment_triggers
		WHERE employee_id = $1
		ORDER BY triggered_at DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReassessmentRows(rows)
}

func (r *ReassessmentRepository) Resolve(id string, req ResolveReassessmentRequest) (*ReassessmentTrigger, error) {
	var assessmentID sql.NullString
	if req.AssessmentID != nil {
		assessmentID = sql.NullString{String: *req.AssessmentID, Valid: true}
	}

	_, err := r.DB.Exec(`
		UPDATE reassessment_triggers
		SET resolved_at = NOW(),
		    resolution_notes = $1,
		    resolution_assessment_id = $2
		WHERE id = $3
	`, req.ResolutionNotes, assessmentID, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *ReassessmentRepository) GetByID(id string) (*ReassessmentTrigger, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, requirement_id, trigger_type, triggered_by, triggered_at,
		       notes, resolved_at, resolution_notes, resolution_assessment_id
		FROM reassessment_triggers
		WHERE id = $1
	`, id)

	var rt ReassessmentTrigger
	var requirementID, triggeredAt sql.NullString
	var resolvedAt, resolutionNotes, resolutionAssessmentID sql.NullString

	err := row.Scan(
		&rt.ID, &rt.EmployeeID, &requirementID, &rt.TriggerType, &rt.TriggeredBy, &triggeredAt,
		&rt.Notes, &resolvedAt, &resolutionNotes, &resolutionAssessmentID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	applyReassessmentNullables(&rt, requirementID, triggeredAt, resolvedAt, resolutionNotes, resolutionAssessmentID)
	return &rt, nil
}

// scanReassessmentRows scans multiple rows into []ReassessmentTrigger.
func scanReassessmentRows(rows *sql.Rows) ([]ReassessmentTrigger, error) {
	results := make([]ReassessmentTrigger, 0)
	for rows.Next() {
		var rt ReassessmentTrigger
		var requirementID, triggeredAt sql.NullString
		var resolvedAt, resolutionNotes, resolutionAssessmentID sql.NullString

		err := rows.Scan(
			&rt.ID, &rt.EmployeeID, &requirementID, &rt.TriggerType, &rt.TriggeredBy, &triggeredAt,
			&rt.Notes, &resolvedAt, &resolutionNotes, &resolutionAssessmentID,
		)
		if err != nil {
			return nil, err
		}

		applyReassessmentNullables(&rt, requirementID, triggeredAt, resolvedAt, resolutionNotes, resolutionAssessmentID)
		results = append(results, rt)
	}
	return results, nil
}

// applyReassessmentNullables maps sql.NullString values onto *string fields of ReassessmentTrigger.
func applyReassessmentNullables(
	rt *ReassessmentTrigger,
	requirementID, triggeredAt, resolvedAt, resolutionNotes, resolutionAssessmentID sql.NullString,
) {
	if requirementID.Valid {
		rt.RequirementID = &requirementID.String
	}
	if triggeredAt.Valid {
		rt.TriggeredAt = triggeredAt.String
	}
	if resolvedAt.Valid {
		rt.ResolvedAt = &resolvedAt.String
	}
	if resolutionNotes.Valid {
		rt.ResolutionNotes = &resolutionNotes.String
	}
	if resolutionAssessmentID.Valid {
		rt.ResolutionAssessmentID = &resolutionAssessmentID.String
	}
}
