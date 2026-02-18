package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type CertificationRepository struct {
	DB *sql.DB
}

func (r *CertificationRepository) GetAll() ([]Certification, error) {
	rows, err := r.DB.Query("SELECT id, employee_id, course_id, issuer, issue_date, expiry_date, notes FROM certifications ORDER BY issue_date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	certs := make([]Certification, 0)
	for rows.Next() {
		var c Certification
		var expiry sql.NullString
		var notes sql.NullString
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.CourseID, &c.Issuer, &c.IssueDate, &expiry, &notes); err != nil {
			return nil, err
		}
		if expiry.Valid {
			c.ExpiryDate = expiry.String
		}
		if notes.Valid {
			c.Notes = notes.String
		}
		certs = append(certs, c)
	}
	return certs, nil
}

func (r *CertificationRepository) GetByEmployee(employeeID string) ([]Certification, error) {
	rows, err := r.DB.Query("SELECT id, employee_id, course_id, issuer, issue_date, expiry_date, notes FROM certifications WHERE employee_id = ? ORDER BY issue_date DESC", employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	certs := make([]Certification, 0)
	for rows.Next() {
		var c Certification
		var expiry sql.NullString
		var notes sql.NullString
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.CourseID, &c.Issuer, &c.IssueDate, &expiry, &notes); err != nil {
			return nil, err
		}
		if expiry.Valid {
			c.ExpiryDate = expiry.String
		}
		if notes.Valid {
			c.Notes = notes.String
		}
		certs = append(certs, c)
	}
	return certs, nil
}

func (r *CertificationRepository) Create(req CreateCertificationRequest) (*Certification, error) {
	id := fmt.Sprintf("CERT-%d", time.Now().UnixNano())

	var expiry interface{}
	if req.ExpiryDate != nil {
		expiry = *req.ExpiryDate
	}
	var notes interface{}
	if req.Notes != nil {
		notes = *req.Notes
	}

	_, err := r.DB.Exec(`
		INSERT INTO certifications (id, employee_id, course_id, issuer, issue_date, expiry_date, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, req.EmployeeID, req.CourseID, req.Issuer, req.IssueDate, expiry, notes)
	if err != nil {
		return nil, err
	}

	evidenceID := fmt.Sprintf("EV-%d", time.Now().UnixNano())
	r.DB.Exec(`
		INSERT INTO evidence (id, employee_id, course_id, evidence_type, completion_date, expiry_date, metadata)
		VALUES (?, ?, ?, 'certification', ?, ?, ?)
	`, evidenceID, req.EmployeeID, req.CourseID, req.IssueDate, expiry, fmt.Sprintf(`{"issuer":"%s"}`, req.Issuer))

	var expiryStr, notesStr string
	if req.ExpiryDate != nil {
		expiryStr = *req.ExpiryDate
	}
	if req.Notes != nil {
		notesStr = *req.Notes
	}

	return &Certification{
		ID:         id,
		EmployeeID: req.EmployeeID,
		CourseID:   req.CourseID,
		Issuer:     req.Issuer,
		IssueDate:  req.IssueDate,
		ExpiryDate: expiryStr,
		Notes:      notesStr,
	}, nil
}
