package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var ErrNotFound = errors.New("not found")

type Course struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Code           *string `json:"code,omitempty"`
	Type           string  `json:"type"`
	Description    *string `json:"description,omitempty"`
	ValidityMonths int     `json:"validity_months"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type CourseRepository struct {
	DB *sql.DB
}

func (r *CourseRepository) GetAll() ([]Course, error) {
	rows, err := r.DB.Query("SELECT id, name, code, type, description, validity_months, created_at, updated_at FROM courses ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]Course, 0)
	for rows.Next() {
		var c Course
		var code, description sql.NullString
		var createdAt, updatedAt sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &code, &c.Type, &description, &c.ValidityMonths, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		if code.Valid {
			c.Code = &code.String
		}
		if description.Valid {
			c.Description = &description.String
		}
		if createdAt.Valid {
			c.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			c.UpdatedAt = updatedAt.String
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *CourseRepository) GetByID(id string) (*Course, error) {
	var c Course
	var code, description sql.NullString
	var createdAt, updatedAt sql.NullString
	err := r.DB.QueryRow("SELECT id, name, code, type, description, validity_months, created_at, updated_at FROM courses WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &code, &c.Type, &description, &c.ValidityMonths, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if code.Valid {
		c.Code = &code.String
	}
	if description.Valid {
		c.Description = &description.String
	}
	if createdAt.Valid {
		c.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		c.UpdatedAt = updatedAt.String
	}
	return &c, nil
}

func (r *CourseRepository) GetByType(courseType string) ([]Course, error) {
	rows, err := r.DB.Query("SELECT id, name, code, type, description, validity_months, created_at, updated_at FROM courses WHERE type = $1 ORDER BY name", courseType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]Course, 0)
	for rows.Next() {
		var c Course
		var code, description sql.NullString
		var createdAt, updatedAt sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &code, &c.Type, &description, &c.ValidityMonths, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		if code.Valid {
			c.Code = &code.String
		}
		if description.Valid {
			c.Description = &description.String
		}
		if createdAt.Valid {
			c.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			c.UpdatedAt = updatedAt.String
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *CourseRepository) Create(course *Course) error {
	if course.ID == "" {
		id, err := r.generateCourseID()
		if err != nil {
			return err
		}
		course.ID = id
	}

	now := time.Now().Format(time.RFC3339)
	_, err := r.DB.Exec(`
		INSERT INTO courses (id, name, code, type, description, validity_months, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, course.ID, course.Name, course.Code, course.Type, course.Description, course.ValidityMonths, now, now)
	if err != nil {
		return err
	}
	course.CreatedAt = now
	course.UpdatedAt = now
	return nil
}

func (r *CourseRepository) Update(course *Course) error {
	now := time.Now().Format(time.RFC3339)
	_, err := r.DB.Exec(`
		UPDATE courses SET name = $1, code = $2, type = $3, description = $4, validity_months = $5, updated_at = $6 WHERE id = $7
	`, course.Name, course.Code, course.Type, course.Description, course.ValidityMonths, now, course.ID)
	if err != nil {
		return err
	}
	course.UpdatedAt = now
	return nil
}

func (r *CourseRepository) Delete(id string) error {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM requirements WHERE course_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete course: %d requirements reference this course", count)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM training_events WHERE course_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete course: %d training events reference this course", count)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM evidence WHERE course_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete course: %d evidence records reference this course", count)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM certifications WHERE course_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete course: %d certifications reference this course", count)
	}

	_, err = r.DB.Exec("DELETE FROM courses WHERE id = $1", id)
	return err
}

func (r *CourseRepository) generateCourseID() (string, error) {
	var maxID sql.NullString
	err := r.DB.QueryRow("SELECT MAX(id) FROM courses WHERE id LIKE 'C%'").Scan(&maxID)
	if err != nil {
		return "", err
	}

	if !maxID.Valid || maxID.String == "" {
		return "C001", nil
	}

	numStr := maxID.String[1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse course ID: %s", maxID.String)
	}

	return fmt.Sprintf("C%03d", num+1), nil
}
