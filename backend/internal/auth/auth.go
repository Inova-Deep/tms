package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

// User represents the authenticated user context
type User struct {
	Role       string `json:"role"`
	EmployeeID string `json:"employee_id,omitempty"`
}

// Claims represents the token payload
type Claims struct {
	Role       string `json:"role"`
	EmployeeID string `json:"employeeId,omitempty"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Role       string `json:"role"`
	EmployeeID string `json:"employeeId,omitempty"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token      string `json:"token"`
	Role       string `json:"role"`
	EmployeeID string `json:"employeeId,omitempty"`
}

// GenerateToken creates a simple base64-encoded token
func GenerateToken(role, employeeID string) string {
	claims := Claims{
		Role:       role,
		EmployeeID: employeeID,
	}
	data, _ := json.Marshal(claims)
	return base64.StdEncoding.EncodeToString(data)
}

// ParseToken decodes the token and returns claims
func ParseToken(token string) (*Claims, error) {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// ExtractTokenFromRequest gets token from Authorization header
func ExtractTokenFromRequest(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	// Expect "Bearer <token>"
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

// GetUserFromRequest extracts user from request context or token
func GetUserFromRequest(r *http.Request) *User {
	token := ExtractTokenFromRequest(r)
	if token == "" {
		return nil
	}
	claims, err := ParseToken(token)
	if err != nil {
		return nil
	}
	return &User{
		Role:       claims.Role,
		EmployeeID: claims.EmployeeID,
	}
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u != nil && u.Role == "admin"
}

// CanAccessEmployee checks if user can access employee data
func (u *User) CanAccessEmployee(employeeID string) bool {
	if u == nil {
		return false
	}
	if u.IsAdmin() {
		return true
	}
	return u.EmployeeID == employeeID
}
