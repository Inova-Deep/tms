package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type SupervisionRepository struct {
	DB *sql.DB
}

func (r *SupervisionRepository) Create(req CreateSupervisionPeriodRequest) (*SupervisionPeriod, error) {
	id := fmt.Sprintf("SUP-%d", time.Now().UnixNano())

	var requirementID sql.NullString
	if req.RequirementID != nil {
		requirementID = sql.NullString{String: *req.RequirementID, Valid: true}
	}

	var scopeLimitations sql.NullString
	if req.ScopeLimitations != nil {
		scopeLimitations = sql.NullString{String: *req.ScopeLimitations, Valid: true}
	}

	_, err := r.DB.Exec(`
		INSERT INTO supervision_periods
		  (id, employee_id, supervisor_id, requirement_id, scope_limitations,
		   start_date, end_date, status, completion_assessment_id, created_at)
		VALUES
		  ($1, $2, $3, $4, $5,
		   $6, NULL, 'active', NULL, NOW())
	`,
		id, req.EmployeeID, req.SupervisorID, requirementID, scopeLimitations,
		req.StartDate,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *SupervisionRepository) Complete(id string, req CompleteSupervisionRequest) (*SupervisionPeriod, error) {
	var assessmentID sql.NullString
	if req.AssessmentID != nil {
		assessmentID = sql.NullString{String: *req.AssessmentID, Valid: true}
	}

	_, err := r.DB.Exec(`
		UPDATE supervision_periods
		SET status = 'completed',
		    end_date = $1,
		    completion_assessment_id = $2
		WHERE id = $3
	`, req.EndDate, assessmentID, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *SupervisionRepository) GetActive() ([]SupervisionPeriod, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, supervisor_id, requirement_id, scope_limitations,
		       start_date, end_date, status, completion_assessment_id, created_at
		FROM supervision_periods
		WHERE status = 'active'
		ORDER BY start_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSupervisionRows(rows)
}

func (r *SupervisionRepository) GetReadyForAssessment() ([]SupervisionPeriod, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, supervisor_id, requirement_id, scope_limitations,
		       start_date, end_date, status, completion_assessment_id, created_at
		FROM supervision_periods
		WHERE status = 'completed' AND completion_assessment_id IS NULL
		ORDER BY end_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSupervisionRows(rows)
}

func (r *SupervisionRepository) GetByEmployee(employeeID string) ([]SupervisionPeriod, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, supervisor_id, requirement_id, scope_limitations,
		       start_date, end_date, status, completion_assessment_id, created_at
		FROM supervision_periods
		WHERE employee_id = $1
		ORDER BY start_date DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSupervisionRows(rows)
}

func (r *SupervisionRepository) GetByID(id string) (*SupervisionPeriod, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, supervisor_id, requirement_id, scope_limitations,
		       start_date, end_date, status, completion_assessment_id, created_at
		FROM supervision_periods
		WHERE id = $1
	`, id)

	var sp SupervisionPeriod
	var requirementID, scopeLimitations, endDate, completionAssessmentID sql.NullString
	var createdAt sql.NullString

	err := row.Scan(
		&sp.ID, &sp.EmployeeID, &sp.SupervisorID, &requirementID, &scopeLimitations,
		&sp.StartDate, &endDate, &sp.Status, &completionAssessmentID, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	applySupervisionNullables(&sp, requirementID, scopeLimitations, endDate, completionAssessmentID, createdAt)
	return &sp, nil
}

// scanSupervisionRows scans multiple rows into []SupervisionPeriod.
func scanSupervisionRows(rows *sql.Rows) ([]SupervisionPeriod, error) {
	results := make([]SupervisionPeriod, 0)
	for rows.Next() {
		var sp SupervisionPeriod
		var requirementID, scopeLimitations, endDate, completionAssessmentID sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&sp.ID, &sp.EmployeeID, &sp.SupervisorID, &requirementID, &scopeLimitations,
			&sp.StartDate, &endDate, &sp.Status, &completionAssessmentID, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		applySupervisionNullables(&sp, requirementID, scopeLimitations, endDate, completionAssessmentID, createdAt)
		results = append(results, sp)
	}
	return results, nil
}

// applySupervisionNullables maps sql.NullString values onto *string fields of SupervisionPeriod.
func applySupervisionNullables(
	sp *SupervisionPeriod,
	requirementID, scopeLimitations, endDate, completionAssessmentID, createdAt sql.NullString,
) {
	if requirementID.Valid {
		sp.RequirementID = &requirementID.String
	}
	if scopeLimitations.Valid {
		sp.ScopeLimitations = &scopeLimitations.String
	}
	if endDate.Valid {
		sp.EndDate = &endDate.String
	}
	if completionAssessmentID.Valid {
		sp.CompletionAssessmentID = &completionAssessmentID.String
	}
	if createdAt.Valid {
		sp.CreatedAt = createdAt.String
	}
}
