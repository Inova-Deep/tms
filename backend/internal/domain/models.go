package domain

import (
	"time"
)

// Employee represents a user in the system (read-only from directory)
type Employee struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Department       string `json:"department"`
	Site             string `json:"site"`
	WorkerType       string `json:"worker_type"`
	CostCenter       string `json:"cost_center"`
	EmploymentStatus string `json:"employment_status"`
}

// Profile represents a competency matrix (assigned to employees)
type Profile struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Requirements []Requirement `json:"requirements,omitempty"`
}

// ProfileWithCounts extends Profile with member and exclusion counts
type ProfileWithCounts struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	MemberCount    int    `json:"memberCount"`
	ExclusionCount int    `json:"exclusionCount"`
}

// ProfileMember represents a member of a profile with compliance status
type ProfileMember struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Site       string `json:"site"`
	Status     string `json:"status"` // compliant, expiring, expired, not_started
	JoinedAt   string `json:"joinedAt"`
}

// ProfileMembersResponse wraps the members list with total count
type ProfileMembersResponse struct {
	Members []ProfileMember `json:"members"`
	Total   int             `json:"total"`
}

// ProfileExclusion represents a profile exclusion record
type ProfileExclusion struct {
	ID           string `json:"id"`
	ProfileID    string `json:"profileId"`
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	Reason       string `json:"reason"`
	ExpiryDate   string `json:"expiryDate"`
	CreatedAt    string `json:"createdAt"`
	CreatedBy    string `json:"createdBy"`
}

// CreateExclusionRequest represents the request to create an exclusion
type CreateExclusionRequest struct {
	EmployeeID string `json:"employeeId"`
	Reason     string `json:"reason"`
	ExpiryDate string `json:"expiryDate"`
}

// UpdateExclusionRequest represents the request to update an exclusion
type UpdateExclusionRequest struct {
	Reason     string `json:"reason"`
	ExpiryDate string `json:"expiryDate"`
}

type CreateCourseRequest struct {
	Name           string  `json:"name"`
	Code           *string `json:"code,omitempty"`
	Type           string  `json:"type"`
	Description    *string `json:"description,omitempty"`
	ValidityMonths int     `json:"validity_months"`
}

type UpdateCourseRequest struct {
	Name           *string `json:"name,omitempty"`
	Code           *string `json:"code,omitempty"`
	Type           *string `json:"type,omitempty"`
	Description    *string `json:"description,omitempty"`
	ValidityMonths *int    `json:"validity_months,omitempty"`
}

type Requirement struct {
	ID             string `json:"id"`
	ProfileID      string `json:"profile_id"`
	CourseID       string `json:"course_id"`
	Name           string `json:"name"`
	Type           string `json:"type"` // "course" or "certification"
	ValidityMonths int    `json:"validity_months"`
}

// RequirementStatus represents the computed status for a requirement
type RequirementStatus struct {
	RequirementID string     `json:"requirement_id"`
	Status        string     `json:"status"` // Missing, Valid, Expiring, Expired
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
}

type TrainingEvent struct {
	ID           string   `json:"id"`
	CourseID     string   `json:"course_id"`
	Course       *Course  `json:"course,omitempty"`
	CourseName   string   `json:"course_name"` // Deprecated: Use CourseID/Course instead
	Date         string   `json:"date"`
	Time         string   `json:"time"`     // HH:MM format (e.g., "09:00")
	Duration     int      `json:"duration"` // Duration in minutes
	Location     string   `json:"location"`
	InstructorID string   `json:"instructor_id"`
	Attendees    []string `json:"attendees,omitempty"`
	Status       string   `json:"status"` // draft, scheduled, in_progress, completed, cancelled
}

type CreateEventRequest struct {
	CourseID     string   `json:"courseId"`
	Date         string   `json:"date"`
	Time         string   `json:"time"`
	Duration     int      `json:"duration"`
	Location     string   `json:"location"`
	InstructorID string   `json:"instructorId"`
	Attendees    []string `json:"attendees"`
}

type UpdateEventRequest struct {
	CourseID     *string  `json:"courseId,omitempty"`
	Date         *string  `json:"date,omitempty"`
	Time         *string  `json:"time,omitempty"`
	Duration     *int     `json:"duration,omitempty"`
	Location     *string  `json:"location,omitempty"`
	InstructorID *string  `json:"instructorId,omitempty"`
	Status       *string  `json:"status,omitempty"`
	Attendees    []string `json:"attendees,omitempty"`
}

// AttendeeConfirmation represents a single attendee's confirmation status
type AttendeeConfirmation struct {
	EmployeeID    string `json:"employeeId"`
	Attended      bool   `json:"attended"`
	AbsenceReason string `json:"absenceReason,omitempty"` // sick_leave, no_show, emergency, other
	AbsenceNotes  string `json:"absenceNotes,omitempty"`
}

// WalkIn represents a walk-in attendee
type WalkIn struct {
	EmployeeID string `json:"employeeId"`
}

// ConfirmAttendanceRequest represents the enhanced attendance confirmation
type ConfirmAttendanceRequest struct {
	Attendees []AttendeeConfirmation `json:"attendees"`
	WalkIns   []WalkIn               `json:"walkIns,omitempty"`
}

type Certification struct {
	ID              string `json:"id"`
	EmployeeID      string `json:"employeeId"`
	CourseID        string `json:"course_id"`
	RequirementName string `json:"requirementName"` // Deprecated: Use CourseID instead
	Issuer          string `json:"issuer"`
	IssueDate       string `json:"issueDate"`
	ExpiryDate      string `json:"expiryDate,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

type CreateCertificationRequest struct {
	EmployeeID string  `json:"employeeId"`
	CourseID   string  `json:"courseId"`
	Issuer     string  `json:"issuer"`
	IssueDate  string  `json:"issueDate"`
	ExpiryDate *string `json:"expiryDate,omitempty"`
	Notes      *string `json:"notes,omitempty"`
}

// DashboardKPIs represents the KPI tile data
type DashboardKPIs struct {
	TotalEmployees    int     `json:"totalEmployees"`
	OverallCompliance float64 `json:"overallCompliance"`
	Valid             int     `json:"valid"`
	ExpiringSoon      int     `json:"expiringSoon"`
	Expired           int     `json:"expired"`
	Missing           int     `json:"missing"`
}

// GridCell represents a cell in the coverage grid
type GridCell struct {
	GroupValue string  `json:"groupValue"`
	ItemID     string  `json:"itemId,omitempty"`
	ItemName   string  `json:"itemName,omitempty"`
	Total      int     `json:"total"`
	Valid      int     `json:"valid"`
	Expiring   int     `json:"expiring"`
	Expired    int     `json:"expired"`
	Missing    int     `json:"missing"`
	Percent    float64 `json:"percent"`
}

// DrilldownEmployee represents an employee in drilldown results
type DrilldownEmployee struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Site       string `json:"site"`
	Status     string `json:"status"`
	ExpiryDate string `json:"expiryDate,omitempty"`
}
