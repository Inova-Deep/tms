package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type CompetenceAssessmentRepository struct {
	DB *sql.DB
}

func (r *CompetenceAssessmentRepository) Create(req CreateCompetenceAssessmentRequest) (*CompetenceAssessment, error) {
	id := fmt.Sprintf("CA-%d", time.Now().UnixNano())

	var requirementID sql.NullString
	if req.RequirementID != nil {
		requirementID = sql.NullString{String: *req.RequirementID, Valid: true}
	}

	var limitations sql.NullString
	if req.Limitations != nil {
		limitations = sql.NullString{String: *req.Limitations, Valid: true}
	}

	var evidenceReference sql.NullString
	if req.EvidenceReference != nil {
		evidenceReference = sql.NullString{String: *req.EvidenceReference, Valid: true}
	}

	var notes sql.NullString
	if req.Notes != nil {
		notes = sql.NullString{String: *req.Notes, Valid: true}
	}

	_, err := r.DB.Exec(`
		INSERT INTO competence_assessments
		  (id, employee_id, requirement_id, assessed_by, assessment_date,
		   work_activities_covered, limitations, supervision_required,
		   outcome, evidence_reference, notes, created_at)
		VALUES
		  ($1, $2, $3, $4, $5,
		   $6, $7, $8,
		   $9, $10, $11, NOW())
	`,
		id, req.EmployeeID, requirementID, req.AssessedBy, req.AssessmentDate,
		req.WorkActivitiesCovered, limitations, req.SupervisionRequired,
		req.Outcome, evidenceReference, notes,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *CompetenceAssessmentRepository) GetByEmployee(employeeID string) ([]CompetenceAssessment, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, assessed_by, assessment_date,
		       work_activities_covered, limitations, supervision_required,
		       outcome, evidence_reference, notes, created_at
		FROM competence_assessments
		WHERE employee_id = $1
		ORDER BY assessment_date DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCompetenceAssessmentRows(rows)
}

func (r *CompetenceAssessmentRepository) GetByID(id string) (*CompetenceAssessment, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, requirement_id, assessed_by, assessment_date,
		       work_activities_covered, limitations, supervision_required,
		       outcome, evidence_reference, notes, created_at
		FROM competence_assessments
		WHERE id = $1
	`, id)

	var ca CompetenceAssessment
	var requirementID, limitations, evidenceReference, notes sql.NullString
	var createdAt sql.NullString

	err := row.Scan(
		&ca.ID, &ca.EmployeeID, &requirementID, &ca.AssessedBy, &ca.AssessmentDate,
		&ca.WorkActivitiesCovered, &limitations, &ca.SupervisionRequired,
		&ca.Outcome, &evidenceReference, &notes, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	applyCompetenceAssessmentNullables(&ca, requirementID, limitations, evidenceReference, notes, createdAt)
	return &ca, nil
}

// GetPending returns assessments where:
//   - outcome = 'not_yet_competent', OR
//   - the assessment is linked to a supervision period that is completed but has no follow-up assessment.
func (r *CompetenceAssessmentRepository) GetPending() ([]CompetenceAssessment, error) {
	rows, err := r.DB.Query(`
		SELECT DISTINCT ca.id, ca.employee_id, ca.requirement_id, ca.assessed_by, ca.assessment_date,
		       ca.work_activities_covered, ca.limitations, ca.supervision_required,
		       ca.outcome, ca.evidence_reference, ca.notes, ca.created_at
		FROM competence_assessments ca
		WHERE ca.outcome = 'not_yet_competent'
		UNION
		SELECT DISTINCT ca.id, ca.employee_id, ca.requirement_id, ca.assessed_by, ca.assessment_date,
		       ca.work_activities_covered, ca.limitations, ca.supervision_required,
		       ca.outcome, ca.evidence_reference, ca.notes, ca.created_at
		FROM competence_assessments ca
		JOIN supervision_periods sp ON sp.completion_assessment_id = ca.id
		WHERE sp.status = 'completed'
		  AND NOT EXISTS (
		      SELECT 1 FROM competence_assessments ca2
		      WHERE ca2.employee_id = ca.employee_id
		        AND ca2.requirement_id = ca.requirement_id
		        AND ca2.created_at > ca.created_at
		  )
		ORDER BY assessment_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCompetenceAssessmentRows(rows)
}

// scanCompetenceAssessmentRows scans multiple rows into []CompetenceAssessment.
func scanCompetenceAssessmentRows(rows *sql.Rows) ([]CompetenceAssessment, error) {
	results := make([]CompetenceAssessment, 0)
	for rows.Next() {
		var ca CompetenceAssessment
		var requirementID, limitations, evidenceReference, notes sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&ca.ID, &ca.EmployeeID, &requirementID, &ca.AssessedBy, &ca.AssessmentDate,
			&ca.WorkActivitiesCovered, &limitations, &ca.SupervisionRequired,
			&ca.Outcome, &evidenceReference, &notes, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		applyCompetenceAssessmentNullables(&ca, requirementID, limitations, evidenceReference, notes, createdAt)
		results = append(results, ca)
	}
	return results, nil
}

// applyCompetenceAssessmentNullables maps sql.NullString values onto *string fields.
func applyCompetenceAssessmentNullables(
	ca *CompetenceAssessment,
	requirementID, limitations, evidenceReference, notes, createdAt sql.NullString,
) {
	if requirementID.Valid {
		ca.RequirementID = &requirementID.String
	}
	if limitations.Valid {
		ca.Limitations = &limitations.String
	}
	if evidenceReference.Valid {
		ca.EvidenceReference = &evidenceReference.String
	}
	if notes.Valid {
		ca.Notes = &notes.String
	}
	if createdAt.Valid {
		ca.CreatedAt = createdAt.String
	}
}
