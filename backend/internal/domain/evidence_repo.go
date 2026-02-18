package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type EvidenceRepository struct {
	DB *sql.DB
}

type Evidence struct {
	ID              string  `json:"id"`
	EmployeeID      string  `json:"employeeId"`
	CourseID        string  `json:"courseId"`
	Course          *Course `json:"course,omitempty"`
	RequirementName string  `json:"requirementName"` // Deprecated: Use CourseID/Course instead
	EvidenceType    string  `json:"evidenceType"`    // "attendance" or "certification"
	CompletionDate  string  `json:"completionDate"`
	ExpiryDate      *string `json:"expiryDate,omitempty"`
	Metadata        string  `json:"metadata,omitempty"`
}

func (r *EvidenceRepository) GetByEmployee(employeeID string) ([]Evidence, error) {
	query := `
		SELECT e.id, e.employee_id, e.course_id, e.evidence_type, e.completion_date, e.expiry_date, e.metadata,
		       c.id, c.name, c.code, c.type, c.description, c.validity_months
		FROM evidence e
		LEFT JOIN courses c ON e.course_id = c.id
		WHERE e.employee_id = ?
	`
	rows, err := r.DB.Query(query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evidenceList := make([]Evidence, 0)
	for rows.Next() {
		var e Evidence
		var expiry sql.NullString
		var metadata sql.NullString
		var courseID, courseName, courseType sql.NullString
		var courseCode, courseDesc sql.NullString
		var courseValidity sql.NullInt64
		if err := rows.Scan(&e.ID, &e.EmployeeID, &e.CourseID, &e.EvidenceType, &e.CompletionDate, &expiry, &metadata,
			&courseID, &courseName, &courseCode, &courseType, &courseDesc, &courseValidity); err != nil {
			return nil, err
		}
		if expiry.Valid {
			val := expiry.String
			e.ExpiryDate = &val
		}
		if metadata.Valid {
			e.Metadata = metadata.String
		}
		if courseID.Valid && courseName.Valid {
			e.Course = &Course{
				ID:             courseID.String,
				Name:           courseName.String,
				Type:           courseType.String,
				ValidityMonths: int(courseValidity.Int64),
			}
			if courseCode.Valid {
				e.Course.Code = &courseCode.String
			}
			if courseDesc.Valid {
				e.Course.Description = &courseDesc.String
			}
			e.RequirementName = courseName.String
		}
		evidenceList = append(evidenceList, e)
	}
	return evidenceList, nil
}

func (r *EvidenceRepository) Create(e Evidence) error {
	var expiry interface{}
	if e.ExpiryDate != nil {
		expiry = *e.ExpiryDate
	}
	_, err := r.DB.Exec(`
		INSERT INTO evidence (id, employee_id, course_id, evidence_type, completion_date, expiry_date, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, e.ID, e.EmployeeID, e.CourseID, e.EvidenceType, e.CompletionDate, expiry, e.Metadata)
	return err
}

func (r *EvidenceRepository) CreateFromAttendance(employeeID, courseID, completionDate string, validityMonths int) error {
	id := fmt.Sprintf("EV-%s-%d", employeeID, time.Now().UnixNano())

	var expiryDate string
	if validityMonths > 0 {
		t, _ := time.Parse("2006-01-02", completionDate)
		expiryDate = t.AddDate(0, validityMonths, 0).Format("2006-01-02")
	}

	_, err := r.DB.Exec(`
		INSERT INTO evidence (id, employee_id, course_id, evidence_type, completion_date, expiry_date, metadata)
		VALUES (?, ?, ?, 'attendance', ?, ?, ?)
	`, id, employeeID, courseID, completionDate, expiryDate, "")
	return err
}
