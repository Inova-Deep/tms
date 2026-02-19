package domain

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type RoleRepository struct {
	DB *sql.DB
}

func (r *RoleRepository) GetAll() ([]Role, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, department, skill_category, risk_level, criticality,
		       mandatory_for_role, matrix_applicable, created_at, updated_at
		FROM roles
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}
	return roles, nil
}

func (r *RoleRepository) GetByID(id string) (*Role, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, department, skill_category, risk_level, criticality,
		       mandatory_for_role, matrix_applicable, created_at, updated_at
		FROM roles
		WHERE id = $1
	`, id)

	var ro Role
	var matrixApplicable sql.NullBool
	var createdAt, updatedAt sql.NullString

	err := row.Scan(
		&ro.ID, &ro.Name, &ro.Department, &ro.SkillCategory,
		&ro.RiskLevel, &ro.Criticality, &ro.MandatoryForRole,
		&matrixApplicable, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if matrixApplicable.Valid {
		ro.MatrixApplicable = &matrixApplicable.Bool
	}
	if createdAt.Valid {
		ro.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		ro.UpdatedAt = updatedAt.String
	}
	return &ro, nil
}

func (r *RoleRepository) Create(req CreateRoleRequest) (*Role, error) {
	id := fmt.Sprintf("ROLE-%d", time.Now().UnixNano())

	_, err := r.DB.Exec(`
		INSERT INTO roles
		  (id, name, department, skill_category, risk_level, criticality,
		   mandatory_for_role, matrix_applicable, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, NULL, NOW(), NOW())
	`,
		id, req.Name, req.Department, req.SkillCategory,
		req.RiskLevel, req.Criticality, req.MandatoryForRole,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *RoleRepository) Update(id string, req UpdateRoleRequest) (*Role, error) {
	setClauses := make([]string, 0)
	args := make([]interface{}, 0)
	paramIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIdx))
		args = append(args, *req.Name)
		paramIdx++
	}
	if req.Department != nil {
		setClauses = append(setClauses, fmt.Sprintf("department = $%d", paramIdx))
		args = append(args, *req.Department)
		paramIdx++
	}
	if req.SkillCategory != nil {
		setClauses = append(setClauses, fmt.Sprintf("skill_category = $%d", paramIdx))
		args = append(args, *req.SkillCategory)
		paramIdx++
	}
	if req.RiskLevel != nil {
		setClauses = append(setClauses, fmt.Sprintf("risk_level = $%d", paramIdx))
		args = append(args, *req.RiskLevel)
		paramIdx++
	}
	if req.Criticality != nil {
		setClauses = append(setClauses, fmt.Sprintf("criticality = $%d", paramIdx))
		args = append(args, *req.Criticality)
		paramIdx++
	}
	if req.MandatoryForRole != nil {
		setClauses = append(setClauses, fmt.Sprintf("mandatory_for_role = $%d", paramIdx))
		args = append(args, *req.MandatoryForRole)
		paramIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE roles SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		paramIdx,
	)

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *RoleRepository) GetMatrixRoles() ([]Role, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, department, skill_category, risk_level, criticality,
		       mandatory_for_role, matrix_applicable, created_at, updated_at
		FROM roles
		WHERE matrix_applicable = true
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}
	return roles, nil
}

func (r *RoleRepository) GetAwarenessRoles() ([]Role, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, department, skill_category, risk_level, criticality,
		       mandatory_for_role, matrix_applicable, created_at, updated_at
		FROM roles
		WHERE matrix_applicable = false
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}
	return roles, nil
}

func (r *RoleRepository) SetMatrixApplicable(id string, applicable bool) error {
	_, err := r.DB.Exec(
		"UPDATE roles SET matrix_applicable = $1, updated_at = NOW() WHERE id = $2",
		applicable, id,
	)
	return err
}

// scanRole is a shared helper that scans a roles row.
// It accepts *sql.Rows so it works inside iteration loops.
func scanRole(rows *sql.Rows) (*Role, error) {
	var ro Role
	var matrixApplicable sql.NullBool
	var createdAt, updatedAt sql.NullString

	err := rows.Scan(
		&ro.ID, &ro.Name, &ro.Department, &ro.SkillCategory,
		&ro.RiskLevel, &ro.Criticality, &ro.MandatoryForRole,
		&matrixApplicable, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	if matrixApplicable.Valid {
		ro.MatrixApplicable = &matrixApplicable.Bool
	}
	if createdAt.Valid {
		ro.CreatedAt = createdAt.String
	}
	if updatedAt.Valid {
		ro.UpdatedAt = updatedAt.String
	}
	return &ro, nil
}
