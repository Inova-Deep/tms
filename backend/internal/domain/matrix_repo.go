package domain

import (
	"database/sql"
	"time"
)

type MatrixRepository struct {
	DB              *sql.DB
	EmployeeRepo    *EmployeeRepository
	CourseRepo      *CourseRepository
	EvidenceRepo    *EvidenceRepository
	ProfileRepo     *ProfileRepository
	ComputationRepo *ComputationRepository
}

type ComputationRepository struct {
	DB *sql.DB
}

func NewMatrixRepository(db *sql.DB, empRepo *EmployeeRepository, courseRepo *CourseRepository, evidenceRepo *EvidenceRepository, profileRepo *ProfileRepository) *MatrixRepository {
	return &MatrixRepository{
		DB:           db,
		EmployeeRepo: empRepo,
		CourseRepo:   courseRepo,
		EvidenceRepo: evidenceRepo,
		ProfileRepo:  profileRepo,
	}
}

func (r *MatrixRepository) GetMatrix() (*MatrixResponse, error) {
	courses, err := r.CourseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	employees, err := r.EmployeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	rows := make([]MatrixRow, 0, len(employees))

	for _, emp := range employees {
		row := MatrixRow{
			EmployeeID:   emp.ID,
			EmployeeName: emp.Name,
			Department:   emp.Department,
			Cells:        make(map[string]MatrixCell),
		}

		profileIDs, err := r.EmployeeRepo.GetAssignments(emp.ID)
		if err != nil {
			return nil, err
		}

		requiredCourseIDs := make(map[string]bool)
		for _, pid := range profileIDs {
			reqs, err := r.ProfileRepo.GetRequirements(pid)
			if err != nil {
				continue
			}
			for _, req := range reqs {
				requiredCourseIDs[req.CourseID] = true
			}
		}

		evidenceList, err := r.EvidenceRepo.GetByEmployee(emp.ID)
		if err != nil {
			return nil, err
		}

		evidenceByCourse := make(map[string]Evidence)
		for _, e := range evidenceList {
			if existing, ok := evidenceByCourse[e.CourseID]; !ok || e.CompletionDate > existing.CompletionDate {
				evidenceByCourse[e.CourseID] = e
			}
		}

		for _, course := range courses {
			cell := MatrixCell{}

			if !requiredCourseIDs[course.ID] {
				cell.Status = "not_required"
			} else {
				evidence, hasEvidence := evidenceByCourse[course.ID]
				if !hasEvidence {
					cell.Status = "not_taken"
				} else {
					expiryDate := r.computeExpiryDate(evidence, course.ValidityMonths)
					if expiryDate != nil {
						cell.ExpiryDate = expiryDate
						status := r.computeStatus(*expiryDate)
						cell.Status = status
					} else {
						cell.Status = "valid"
					}
				}
			}

			row.Cells[course.ID] = cell
		}

		rows = append(rows, row)
	}

	return &MatrixResponse{
		Courses: courses,
		Rows:    rows,
	}, nil
}

func (r *MatrixRepository) computeExpiryDate(evidence Evidence, validityMonths int) *string {
	if evidence.ExpiryDate != nil && *evidence.ExpiryDate != "" {
		return evidence.ExpiryDate
	}

	if validityMonths > 0 {
		completion, err := time.Parse("2006-01-02", evidence.CompletionDate)
		if err != nil {
			return nil
		}
		expiry := completion.AddDate(0, validityMonths, 0).Format("2006-01-02")
		return &expiry
	}

	return nil
}

func (r *MatrixRepository) computeStatus(expiryDateStr string) string {
	now := time.Now()

	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	var expiry time.Time
	var err error
	for _, f := range formats {
		expiry, err = time.Parse(f, expiryDateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return "valid"
	}

	if expiry.Before(now) {
		return "expired"
	} else if expiry.Before(now.AddDate(0, 0, 90)) {
		return "expiring"
	}

	return "valid"
}
