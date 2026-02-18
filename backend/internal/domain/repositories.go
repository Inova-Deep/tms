package domain

import (
	"database/sql"
	"fmt"
	"strings"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func (r *EmployeeRepository) GetAll() ([]Employee, error) {
	rows, err := r.DB.Query("SELECT id, name, department, site, worker_type, cost_center, employment_status FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Department, &e.Site, &e.WorkerType, &e.CostCenter, &e.EmploymentStatus); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) GetByID(id string) (*Employee, error) {
	var e Employee
	err := r.DB.QueryRow("SELECT id, name, department, site, worker_type, cost_center, employment_status FROM employees WHERE id = ?", id).
		Scan(&e.ID, &e.Name, &e.Department, &e.Site, &e.WorkerType, &e.CostCenter, &e.EmploymentStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EmployeeRepository) Search(query string) ([]Employee, error) {
	query = strings.ToLower(query)
	rows, err := r.DB.Query(`
		SELECT id, name, department, site, worker_type, cost_center, employment_status 
		FROM employees 
		WHERE LOWER(id) LIKE ? OR LOWER(name) LIKE ? OR LOWER(department) LIKE ?
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Department, &e.Site, &e.WorkerType, &e.CostCenter, &e.EmploymentStatus); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func (r *EmployeeRepository) GetAssignments(employeeID string) ([]string, error) {
	rows, err := r.DB.Query("SELECT profile_id FROM profile_assignments WHERE employee_id = ?", employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profileIDs := make([]string, 0)
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		profileIDs = append(profileIDs, pid)
	}
	return profileIDs, nil
}

func (r *EmployeeRepository) GetByFilter(department, site string) ([]Employee, error) {
	query := "SELECT id, name, department, site, worker_type, cost_center, employment_status FROM employees WHERE 1=1"
	args := make([]interface{}, 0)

	if department != "" {
		query += " AND department = ?"
		args = append(args, department)
	}
	if site != "" {
		query += " AND site = ?"
		args = append(args, site)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Department, &e.Site, &e.WorkerType, &e.CostCenter, &e.EmploymentStatus); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, nil
}

// ProfileRepository handles profile database operations
type ProfileRepository struct {
	DB *sql.DB
}

func (r *ProfileRepository) GetAll() ([]Profile, error) {
	rows, err := r.DB.Query("SELECT id, name, description FROM profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]Profile, 0)
	for rows.Next() {
		var p Profile
		var desc sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &desc); err != nil {
			return nil, err
		}
		if desc.Valid {
			p.Description = desc.String
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *ProfileRepository) GetByID(id string) (*Profile, error) {
	var p Profile
	var desc sql.NullString
	err := r.DB.QueryRow("SELECT id, name, description FROM profiles WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &desc)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		p.Description = desc.String
	}
	return &p, nil
}

func (r *ProfileRepository) Create(profile *Profile) error {
	_, err := r.DB.Exec(
		"INSERT INTO profiles (id, name, description) VALUES (?, ?, ?)",
		profile.ID, profile.Name, profile.Description,
	)
	return err
}

func (r *ProfileRepository) GenerateID() (string, error) {
	var maxID int
	err := r.DB.QueryRow("SELECT COALESCE(MAX(CAST(SUBSTR(id, 2) AS INTEGER)), 0) FROM profiles WHERE id LIKE 'P%'").Scan(&maxID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("P%d", maxID+1), nil
}

func (r *ProfileRepository) GetRequirements(profileID string) ([]Requirement, error) {
	rows, err := r.DB.Query(`
		SELECT r.id, r.profile_id, r.course_id, c.name, c.type, c.validity_months
		FROM requirements r
		JOIN courses c ON r.course_id = c.id
		WHERE r.profile_id = ?
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requirements := make([]Requirement, 0)
	for rows.Next() {
		var req Requirement
		if err := rows.Scan(&req.ID, &req.ProfileID, &req.CourseID, &req.Name, &req.Type, &req.ValidityMonths); err != nil {
			return nil, err
		}
		requirements = append(requirements, req)
	}
	return requirements, nil
}

func (r *ProfileRepository) GetAllRequirements() ([]Requirement, error) {
	rows, err := r.DB.Query(`
		SELECT r.id, r.profile_id, r.course_id, c.name, c.type, c.validity_months
		FROM requirements r
		JOIN courses c ON r.course_id = c.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requirements := make([]Requirement, 0)
	for rows.Next() {
		var req Requirement
		if err := rows.Scan(&req.ID, &req.ProfileID, &req.CourseID, &req.Name, &req.Type, &req.ValidityMonths); err != nil {
			return nil, err
		}
		requirements = append(requirements, req)
	}
	return requirements, nil
}

func (r *ProfileRepository) AddRequirement(profileID, courseID string) error {
	// Get course details
	var name, courseType string
	var validityMonths int
	err := r.DB.QueryRow("SELECT name, type, validity_months FROM courses WHERE id = ?", courseID).Scan(&name, &courseType, &validityMonths)
	if err != nil {
		return err
	}

	// Generate requirement ID
	var maxID string
	r.DB.QueryRow("SELECT COALESCE(MAX(id), 'R000') FROM requirements WHERE profile_id = ?", profileID).Scan(&maxID)

	num := 0
	if len(maxID) > 1 {
		fmt.Sscanf(maxID, "R%d", &num)
	}
	newID := fmt.Sprintf("R%d", num+1)

	_, err = r.DB.Exec(
		"INSERT INTO requirements (id, profile_id, course_id, name, type, validity_months) VALUES (?, ?, ?, ?, ?, ?)",
		newID, profileID, courseID, name, courseType, validityMonths,
	)
	return err
}

func (r *ProfileRepository) RemoveRequirement(profileID, requirementID string) error {
	_, err := r.DB.Exec("DELETE FROM requirements WHERE id = ? AND profile_id = ?", requirementID, profileID)
	return err
}

func (r *ProfileRepository) AddMember(profileID, employeeID string) error {
	_, err := r.DB.Exec("INSERT OR IGNORE INTO profile_assignments (employee_id, profile_id) VALUES (?, ?)", employeeID, profileID)
	return err
}

func (r *ProfileRepository) RemoveMember(profileID, employeeID string) error {
	_, err := r.DB.Exec("DELETE FROM profile_assignments WHERE employee_id = ? AND profile_id = ?", employeeID, profileID)
	return err
}

func (r *ProfileRepository) GetAllWithCounts() ([]ProfileWithCounts, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.name, COALESCE(p.description, ''),
			(SELECT COUNT(*) FROM profile_assignments pa WHERE pa.profile_id = p.id) AS member_count,
			(SELECT COUNT(*) FROM profile_exclusions pe WHERE pe.profile_id = p.id) AS exclusion_count
		FROM profiles p
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]ProfileWithCounts, 0)
	for rows.Next() {
		var p ProfileWithCounts
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.MemberCount, &p.ExclusionCount); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *ProfileRepository) GetMembers(profileID string) ([]ProfileMember, error) {
	rows, err := r.DB.Query(`
		SELECT e.id, e.name, e.department, e.site, pa.assigned_at
		FROM employees e
		JOIN profile_assignments pa ON e.id = pa.employee_id
		WHERE pa.profile_id = ?
		ORDER BY e.name
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]ProfileMember, 0)
	for rows.Next() {
		var m ProfileMember
		var joinedAt sql.NullString
		if err := rows.Scan(&m.ID, &m.Name, &m.Department, &m.Site, &joinedAt); err != nil {
			return nil, err
		}
		if joinedAt.Valid {
			m.JoinedAt = joinedAt.String
		}
		// Status will be computed by the handler using ComputationService
		m.Status = "not_started"
		members = append(members, m)
	}
	return members, nil
}

func (r *ProfileRepository) GetExclusions(profileID string) ([]ProfileExclusion, error) {
	rows, err := r.DB.Query(`
		SELECT pe.id, pe.profile_id, pe.employee_id, COALESCE(e.name, ''), pe.reason, pe.expiry_date, pe.created_at, COALESCE(pe.created_by, 'admin')
		FROM profile_exclusions pe
		LEFT JOIN employees e ON pe.employee_id = e.id
		WHERE pe.profile_id = ?
		ORDER BY pe.created_at DESC
	`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exclusions := make([]ProfileExclusion, 0)
	for rows.Next() {
		var ex ProfileExclusion
		var createdAt sql.NullString
		if err := rows.Scan(&ex.ID, &ex.ProfileID, &ex.EmployeeID, &ex.EmployeeName, &ex.Reason, &ex.ExpiryDate, &createdAt, &ex.CreatedBy); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			ex.CreatedAt = createdAt.String
		}
		exclusions = append(exclusions, ex)
	}
	return exclusions, nil
}

func (r *ProfileRepository) CreateExclusion(profileID string, req CreateExclusionRequest) (*ProfileExclusion, error) {
	id := fmt.Sprintf("EX-%s-%s", req.EmployeeID, profileID)

	// Remove existing assignment first
	r.DB.Exec("DELETE FROM profile_assignments WHERE employee_id = ? AND profile_id = ?", req.EmployeeID, profileID)

	_, err := r.DB.Exec(`
		INSERT OR REPLACE INTO profile_exclusions (id, employee_id, profile_id, reason, expiry_date, created_by)
		VALUES (?, ?, ?, ?, ?, 'admin')
	`, id, req.EmployeeID, profileID, req.Reason, req.ExpiryDate)
	if err != nil {
		return nil, err
	}

	// Get employee name
	var name string
	r.DB.QueryRow("SELECT name FROM employees WHERE id = ?", req.EmployeeID).Scan(&name)

	return &ProfileExclusion{
		ID:           id,
		ProfileID:    profileID,
		EmployeeID:   req.EmployeeID,
		EmployeeName: name,
		Reason:       req.Reason,
		ExpiryDate:   req.ExpiryDate,
		CreatedBy:    "admin",
	}, nil
}

func (r *ProfileRepository) GetExclusionByID(exclusionID string) (*ProfileExclusion, error) {
	var ex ProfileExclusion
	var createdAt sql.NullString
	err := r.DB.QueryRow(`
		SELECT pe.id, pe.profile_id, pe.employee_id, COALESCE(e.name, ''), pe.reason, pe.expiry_date, pe.created_at, COALESCE(pe.created_by, 'admin')
		FROM profile_exclusions pe
		LEFT JOIN employees e ON pe.employee_id = e.id
		WHERE pe.id = ?
	`, exclusionID).Scan(&ex.ID, &ex.ProfileID, &ex.EmployeeID, &ex.EmployeeName, &ex.Reason, &ex.ExpiryDate, &createdAt, &ex.CreatedBy)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		ex.CreatedAt = createdAt.String
	}
	return &ex, nil
}

func (r *ProfileRepository) UpdateExclusion(exclusionID string, req UpdateExclusionRequest) error {
	_, err := r.DB.Exec("UPDATE profile_exclusions SET reason = ?, expiry_date = ? WHERE id = ?",
		req.Reason, req.ExpiryDate, exclusionID)
	return err
}

func (r *ProfileRepository) DeleteExclusion(exclusionID string) error {
	_, err := r.DB.Exec("DELETE FROM profile_exclusions WHERE id = ?", exclusionID)
	return err
}

// AddExclusion is the legacy method kept for backward compatibility
func (r *ProfileRepository) AddExclusion(profileID, employeeID, reason, expiryDate string) error {
	req := CreateExclusionRequest{
		EmployeeID: employeeID,
		Reason:     reason,
		ExpiryDate: expiryDate,
	}
	_, err := r.CreateExclusion(profileID, req)
	return err
}
