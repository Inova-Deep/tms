package domain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type InductionRepository struct {
	DB *sql.DB
}

func (r *InductionRepository) Create(req CreateInductionRequest) (*InductionRecord, error) {
	id := fmt.Sprintf("IND-%d", time.Now().UnixNano())

	topicsJSON, err := json.Marshal(req.TopicsCovered)
	if err != nil {
		return nil, err
	}

	var sigRef sql.NullString
	if req.EmployeeSignatureRef != nil {
		sigRef = sql.NullString{String: *req.EmployeeSignatureRef, Valid: true}
	}

	_, err = r.DB.Exec(`
		INSERT INTO induction_records
		  (id, employee_id, conducted_by, induction_date, topics_covered,
		   employee_signature_ref, completed, created_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, false, NOW())
	`,
		id, req.EmployeeID, req.ConductedBy, req.InductionDate,
		string(topicsJSON), sigRef,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByEmployee(req.EmployeeID)
}

func (r *InductionRepository) GetByEmployee(employeeID string) (*InductionRecord, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, conducted_by, induction_date, topics_covered,
		       employee_signature_ref, completed, created_at
		FROM induction_records
		WHERE employee_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, employeeID)

	return scanInductionRow(row)
}

func (r *InductionRepository) GetPending() ([]InductionRecord, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, conducted_by, induction_date, topics_covered,
		       employee_signature_ref, completed, created_at
		FROM induction_records
		WHERE completed = false
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]InductionRecord, 0)
	for rows.Next() {
		var rec InductionRecord
		var topicsCoveredRaw sql.NullString
		var sigRef sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&rec.ID, &rec.EmployeeID, &rec.ConductedBy, &rec.InductionDate,
			&topicsCoveredRaw, &sigRef, &rec.Completed, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		rec.TopicsCovered = unmarshalTopics(topicsCoveredRaw)

		if sigRef.Valid {
			rec.EmployeeSignatureRef = &sigRef.String
		}
		if createdAt.Valid {
			rec.CreatedAt = createdAt.String
		}

		records = append(records, rec)
	}
	return records, nil
}

func (r *InductionRepository) MarkComplete(id string) error {
	_, err := r.DB.Exec(
		"UPDATE induction_records SET completed = true WHERE id = $1",
		id,
	)
	return err
}

// scanInductionRow scans a single *sql.Row returned by QueryRow.
func scanInductionRow(row *sql.Row) (*InductionRecord, error) {
	var rec InductionRecord
	var topicsCoveredRaw sql.NullString
	var sigRef sql.NullString
	var createdAt sql.NullString

	err := row.Scan(
		&rec.ID, &rec.EmployeeID, &rec.ConductedBy, &rec.InductionDate,
		&topicsCoveredRaw, &sigRef, &rec.Completed, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rec.TopicsCovered = unmarshalTopics(topicsCoveredRaw)

	if sigRef.Valid {
		rec.EmployeeSignatureRef = &sigRef.String
	}
	if createdAt.Valid {
		rec.CreatedAt = createdAt.String
	}
	return &rec, nil
}

// unmarshalTopics safely unmarshal a nullable JSON text column into []string.
// Returns an empty (non-nil) slice on NULL or invalid JSON.
func unmarshalTopics(raw sql.NullString) []string {
	topics := make([]string, 0)
	if !raw.Valid || raw.String == "" {
		return topics
	}
	_ = json.Unmarshal([]byte(raw.String), &topics)
	return topics
}
