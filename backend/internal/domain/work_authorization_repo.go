package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type WorkAuthorizationRepository struct {
	DB *sql.DB
}

func (r *WorkAuthorizationRepository) Create(req CreateWorkAuthorizationRequest) (*WorkAuthorization, error) {
	id := fmt.Sprintf("WAUTH-%d", time.Now().UnixNano())

	// Determine state and authorized_at based on whether authorized_by is provided.
	authState := "PENDING_AUTHORIZATION"
	if req.AuthorizedBy != "" {
		authState = "AUTHORIZED"
	}

	var requirementID sql.NullString
	if req.RequirementID != nil {
		requirementID = sql.NullString{String: *req.RequirementID, Valid: true}
	}

	var authorizedBy sql.NullString
	if req.AuthorizedBy != "" {
		authorizedBy = sql.NullString{String: req.AuthorizedBy, Valid: true}
	}

	var expiryDate sql.NullString
	if req.ExpiryDate != nil {
		expiryDate = sql.NullString{String: *req.ExpiryDate, Valid: true}
	}

	var notes sql.NullString
	if req.Notes != nil {
		notes = sql.NullString{String: *req.Notes, Valid: true}
	}

	// authorized_at = NOW() when authorizing immediately, else NULL.
	if req.AuthorizedBy != "" {
		_, err := r.DB.Exec(`
			INSERT INTO work_authorizations
			  (id, employee_id, requirement_id, authorization_type, authorization_state,
			   authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
			   is_legacy, notes, created_at)
			VALUES
			  ($1, $2, $3, $4, $5,
			   $6, NOW(), $7, NULL, NULL,
			   false, $8, NOW())
		`,
			id, req.EmployeeID, requirementID, req.AuthorizationType, authState,
			authorizedBy, expiryDate, notes,
		)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := r.DB.Exec(`
			INSERT INTO work_authorizations
			  (id, employee_id, requirement_id, authorization_type, authorization_state,
			   authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
			   is_legacy, notes, created_at)
			VALUES
			  ($1, $2, $3, $4, $5,
			   $6, NULL, $7, NULL, NULL,
			   false, $8, NOW())
		`,
			id, req.EmployeeID, requirementID, req.AuthorizationType, authState,
			authorizedBy, expiryDate, notes,
		)
		if err != nil {
			return nil, err
		}
	}

	return r.getByID(id)
}

func (r *WorkAuthorizationRepository) GetByEmployee(employeeID string) ([]WorkAuthorization, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, authorization_type, authorization_state,
		       authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
		       is_legacy, notes, created_at
		FROM work_authorizations
		WHERE employee_id = $1
		ORDER BY created_at DESC
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWorkAuthorizationRows(rows)
}

func (r *WorkAuthorizationRepository) Revoke(id string, req RevokeAuthorizationRequest) error {
	_, err := r.DB.Exec(`
		UPDATE work_authorizations
		SET revoked_at = NOW(), revoke_reason = $1
		WHERE id = $2
	`, req.RevokeReason, id)
	return err
}

func (r *WorkAuthorizationRepository) GetLegacyPending() ([]WorkAuthorization, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, authorization_type, authorization_state,
		       authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
		       is_legacy, notes, created_at
		FROM work_authorizations
		WHERE is_legacy = true AND revoked_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWorkAuthorizationRows(rows)
}

func (r *WorkAuthorizationRepository) GetExpiring(daysAhead int) ([]WorkAuthorization, error) {
	rows, err := r.DB.Query(`
		SELECT id, employee_id, requirement_id, authorization_type, authorization_state,
		       authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
		       is_legacy, notes, created_at
		FROM work_authorizations
		WHERE expiry_date <= (NOW() + ($1 || ' days')::interval)
		  AND revoked_at IS NULL
		ORDER BY expiry_date ASC
	`, fmt.Sprintf("%d", daysAhead))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWorkAuthorizationRows(rows)
}

func (r *WorkAuthorizationRepository) GetByEmployeeAndRequirement(employeeID, requirementID string) (*WorkAuthorization, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, requirement_id, authorization_type, authorization_state,
		       authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
		       is_legacy, notes, created_at
		FROM work_authorizations
		WHERE employee_id = $1 AND requirement_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, employeeID, requirementID)

	return scanWorkAuthorizationRow(row)
}

// getByID is an internal helper to fetch a single work_authorization by id.
func (r *WorkAuthorizationRepository) getByID(id string) (*WorkAuthorization, error) {
	row := r.DB.QueryRow(`
		SELECT id, employee_id, requirement_id, authorization_type, authorization_state,
		       authorized_by, authorized_at, expiry_date, revoked_at, revoke_reason,
		       is_legacy, notes, created_at
		FROM work_authorizations
		WHERE id = $1
	`, id)

	return scanWorkAuthorizationRow(row)
}

// scanWorkAuthorizationRow scans a *sql.Row into a WorkAuthorization.
func scanWorkAuthorizationRow(row *sql.Row) (*WorkAuthorization, error) {
	var wa WorkAuthorization
	var requirementID, authorizedBy, authorizedAt, expiryDate sql.NullString
	var revokedAt, revokeReason, notes sql.NullString
	var createdAt sql.NullString

	err := row.Scan(
		&wa.ID, &wa.EmployeeID, &requirementID, &wa.AuthorizationType, &wa.AuthorizationState,
		&authorizedBy, &authorizedAt, &expiryDate, &revokedAt, &revokeReason,
		&wa.IsLegacy, &notes, &createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	applyWorkAuthNullables(&wa, requirementID, authorizedBy, authorizedAt, expiryDate, revokedAt, revokeReason, notes, createdAt)
	return &wa, nil
}

// scanWorkAuthorizationRows scans multiple rows into []WorkAuthorization.
func scanWorkAuthorizationRows(rows *sql.Rows) ([]WorkAuthorization, error) {
	results := make([]WorkAuthorization, 0)
	for rows.Next() {
		var wa WorkAuthorization
		var requirementID, authorizedBy, authorizedAt, expiryDate sql.NullString
		var revokedAt, revokeReason, notes sql.NullString
		var createdAt sql.NullString

		err := rows.Scan(
			&wa.ID, &wa.EmployeeID, &requirementID, &wa.AuthorizationType, &wa.AuthorizationState,
			&authorizedBy, &authorizedAt, &expiryDate, &revokedAt, &revokeReason,
			&wa.IsLegacy, &notes, &createdAt,
		)
		if err != nil {
			return nil, err
		}

		applyWorkAuthNullables(&wa, requirementID, authorizedBy, authorizedAt, expiryDate, revokedAt, revokeReason, notes, createdAt)
		results = append(results, wa)
	}
	return results, nil
}

// applyWorkAuthNullables maps sql.NullString values onto *string fields of WorkAuthorization.
func applyWorkAuthNullables(
	wa *WorkAuthorization,
	requirementID, authorizedBy, authorizedAt, expiryDate,
	revokedAt, revokeReason, notes, createdAt sql.NullString,
) {
	if requirementID.Valid {
		wa.RequirementID = &requirementID.String
	}
	if authorizedBy.Valid {
		wa.AuthorizedBy = &authorizedBy.String
	}
	if authorizedAt.Valid {
		wa.AuthorizedAt = &authorizedAt.String
	}
	if expiryDate.Valid {
		wa.ExpiryDate = &expiryDate.String
	}
	if revokedAt.Valid {
		wa.RevokedAt = &revokedAt.String
	}
	if revokeReason.Valid {
		wa.RevokeReason = &revokeReason.String
	}
	if notes.Valid {
		wa.Notes = &notes.String
	}
	if createdAt.Valid {
		wa.CreatedAt = createdAt.String
	}
}
