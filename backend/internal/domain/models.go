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

// MatrixCell represents a single cell in the training matrix
type MatrixCell struct {
	Status     string  `json:"status"`
	ExpiryDate *string `json:"expiryDate,omitempty"`
}

// MatrixRow represents one employee row in the matrix
type MatrixRow struct {
	EmployeeID   string                `json:"employeeId"`
	EmployeeName string                `json:"employeeName"`
	Department   string                `json:"department"`
	Cells        map[string]MatrixCell `json:"cells"`
}

// MatrixResponse is the full matrix response
type MatrixResponse struct {
	Courses []Course    `json:"courses"`
	Rows    []MatrixRow `json:"rows"`
}

// =============================================================================
// DML Competence Assurance Models
// =============================================================================

// Role represents a job function with defined competence requirements
type Role struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Department       string `json:"department"`
	SkillCategory    string `json:"skillCategory"`
	RiskLevel        string `json:"riskLevel"`       // safety_critical, quality_critical, low, medium
	Criticality      string `json:"criticality"`     // safety, quality, compliance
	MandatoryForRole bool   `json:"mandatoryForRole"`
	MatrixApplicable *bool  `json:"matrixApplicable"` // nil = no decision yet
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type CreateRoleRequest struct {
	Name             string `json:"name"`
	Department       string `json:"department"`
	SkillCategory    string `json:"skillCategory"`
	RiskLevel        string `json:"riskLevel"`
	Criticality      string `json:"criticality"`
	MandatoryForRole bool   `json:"mandatoryForRole"`
}

type UpdateRoleRequest struct {
	Name             *string `json:"name,omitempty"`
	Department       *string `json:"department,omitempty"`
	SkillCategory    *string `json:"skillCategory,omitempty"`
	RiskLevel        *string `json:"riskLevel,omitempty"`
	Criticality      *string `json:"criticality,omitempty"`
	MandatoryForRole *bool   `json:"mandatoryForRole,omitempty"`
}

// ApplicabilityDecision records the 5-question gate result for a role (DML-QA-REG-5024-1)
type ApplicabilityDecision struct {
	ID                   string  `json:"id"`
	RoleID               string  `json:"roleId"`
	Q1                   bool    `json:"q1"`
	Q2                   bool    `json:"q2"`
	Q3                   bool    `json:"q3"`
	Q4                   bool    `json:"q4"`
	Q5                   bool    `json:"q5"`
	Result               string  `json:"result"` // matrix or awareness
	DecidedBy            string  `json:"decidedBy"`
	DecidedAt            string  `json:"decidedAt"`
	Escalated            bool    `json:"escalated"`
	EscalationNotes      *string `json:"escalationNotes,omitempty"`
	EscalationResolvedAt *string `json:"escalationResolvedAt,omitempty"`
	EscalationResolvedBy *string `json:"escalationResolvedBy,omitempty"`
}

type CreateApplicabilityDecisionRequest struct {
	Q1              bool   `json:"q1"`
	Q2              bool   `json:"q2"`
	Q3              bool   `json:"q3"`
	Q4              bool   `json:"q4"`
	Q5              bool   `json:"q5"`
	DecidedBy       string `json:"decidedBy"`
	Escalated       bool   `json:"escalated"`
	EscalationNotes string `json:"escalationNotes,omitempty"`
}

type ResolveEscalationRequest struct {
	ResolvedBy string `json:"resolvedBy"`
	Notes      string `json:"notes"`
}

// InductionRecord records HR new starter induction (Section 3.3)
type InductionRecord struct {
	ID                   string   `json:"id"`
	EmployeeID           string   `json:"employeeId"`
	ConductedBy          string   `json:"conductedBy"`
	InductionDate        string   `json:"inductionDate"`
	TopicsCovered        []string `json:"topicsCovered"`
	EmployeeSignatureRef *string  `json:"employeeSignatureRef,omitempty"`
	Completed            bool     `json:"completed"`
	CreatedAt            string   `json:"createdAt"`
}

type CreateInductionRequest struct {
	EmployeeID           string   `json:"employeeId"`
	ConductedBy          string   `json:"conductedBy"`
	InductionDate        string   `json:"inductionDate"`
	TopicsCovered        []string `json:"topicsCovered"`
	EmployeeSignatureRef *string  `json:"employeeSignatureRef,omitempty"`
}

// WorkAuthorization records manager sign-off for independent work (Section 3.2.1)
type WorkAuthorization struct {
	ID                 string  `json:"id"`
	EmployeeID         string  `json:"employeeId"`
	RequirementID      *string `json:"requirementId,omitempty"`
	AuthorizationType  string  `json:"authorizationType"`  // independent, supervised, restricted
	AuthorizationState string  `json:"authorizationState"` // AUTHORIZED, PENDING_AUTHORIZATION, LEGACY_PENDING, etc.
	AuthorizedBy       *string `json:"authorizedBy,omitempty"`
	AuthorizedAt       *string `json:"authorizedAt,omitempty"`
	ExpiryDate         *string `json:"expiryDate,omitempty"`
	RevokedAt          *string `json:"revokedAt,omitempty"`
	RevokeReason       *string `json:"revokeReason,omitempty"`
	IsLegacy           bool    `json:"isLegacy"`
	Notes              *string `json:"notes,omitempty"`
	CreatedAt          string  `json:"createdAt"`
}

type CreateWorkAuthorizationRequest struct {
	EmployeeID        string  `json:"employeeId"`
	RequirementID     *string `json:"requirementId,omitempty"`
	AuthorizationType string  `json:"authorizationType"`
	AuthorizedBy      string  `json:"authorizedBy"`
	ExpiryDate        *string `json:"expiryDate,omitempty"`
	Notes             *string `json:"notes,omitempty"`
}

type RevokeAuthorizationRequest struct {
	RevokedBy    string `json:"revokedBy"`
	RevokeReason string `json:"revokeReason"`
}

// CompetenceAssessment records formal assessment for critical work (Section 3.2.3)
type CompetenceAssessment struct {
	ID                    string  `json:"id"`
	EmployeeID            string  `json:"employeeId"`
	RequirementID         *string `json:"requirementId,omitempty"`
	AssessedBy            string  `json:"assessedBy"`
	AssessmentDate        string  `json:"assessmentDate"`
	WorkActivitiesCovered string  `json:"workActivitiesCovered"`
	Limitations           *string `json:"limitations,omitempty"`
	SupervisionRequired   bool    `json:"supervisionRequired"`
	Outcome               string  `json:"outcome"` // competent, not_yet_competent, supervised
	EvidenceReference     *string `json:"evidenceReference,omitempty"`
	Notes                 *string `json:"notes,omitempty"`
	CreatedAt             string  `json:"createdAt"`
}

type CreateCompetenceAssessmentRequest struct {
	EmployeeID            string  `json:"employeeId"`
	RequirementID         *string `json:"requirementId,omitempty"`
	AssessedBy            string  `json:"assessedBy"`
	AssessmentDate        string  `json:"assessmentDate"`
	WorkActivitiesCovered string  `json:"workActivitiesCovered"`
	Limitations           *string `json:"limitations,omitempty"`
	SupervisionRequired   bool    `json:"supervisionRequired"`
	Outcome               string  `json:"outcome"`
	EvidenceReference     *string `json:"evidenceReference,omitempty"`
	Notes                 *string `json:"notes,omitempty"`
}

// SupervisionPeriod tracks supervised work periods (Section 3.2.4)
type SupervisionPeriod struct {
	ID                     string  `json:"id"`
	EmployeeID             string  `json:"employeeId"`
	SupervisorID           string  `json:"supervisorId"`
	RequirementID          *string `json:"requirementId,omitempty"`
	ScopeLimitations       *string `json:"scopeLimitations,omitempty"`
	StartDate              string  `json:"startDate"`
	EndDate                *string `json:"endDate,omitempty"`
	Status                 string  `json:"status"` // active, completed, abandoned
	CompletionAssessmentID *string `json:"completionAssessmentId,omitempty"`
	CreatedAt              string  `json:"createdAt"`
}

type CreateSupervisionPeriodRequest struct {
	EmployeeID       string  `json:"employeeId"`
	SupervisorID     string  `json:"supervisorId"`
	RequirementID    *string `json:"requirementId,omitempty"`
	ScopeLimitations *string `json:"scopeLimitations,omitempty"`
	StartDate        string  `json:"startDate"`
}

type CompleteSupervisionRequest struct {
	EndDate      string  `json:"endDate"`
	AssessmentID *string `json:"assessmentId,omitempty"`
}

// ReassessmentTrigger records non-calendar reassessment triggers (Section 3.2.5)
type ReassessmentTrigger struct {
	ID                     string  `json:"id"`
	EmployeeID             string  `json:"employeeId"`
	RequirementID          *string `json:"requirementId,omitempty"`
	TriggerType            string  `json:"triggerType"` // repeated_errors, process_change, extended_absence, cert_expiry, performance_concern
	TriggeredBy            string  `json:"triggeredBy"`
	TriggeredAt            string  `json:"triggeredAt"`
	Notes                  string  `json:"notes"`
	ResolvedAt             *string `json:"resolvedAt,omitempty"`
	ResolutionNotes        *string `json:"resolutionNotes,omitempty"`
	ResolutionAssessmentID *string `json:"resolutionAssessmentId,omitempty"`
}

type CreateReassessmentTriggerRequest struct {
	EmployeeID    string  `json:"employeeId"`
	RequirementID *string `json:"requirementId,omitempty"`
	TriggerType   string  `json:"triggerType"`
	TriggeredBy   string  `json:"triggeredBy"`
	Notes         string  `json:"notes"`
}

type ResolveReassessmentRequest struct {
	ResolutionNotes string  `json:"resolutionNotes"`
	AssessmentID    *string `json:"assessmentId,omitempty"`
}

// TrainingEffectivenessEvaluation records whether training was applied in practice (Section 4.4)
type TrainingEffectivenessEvaluation struct {
	ID             string  `json:"id"`
	EvidenceID     *string `json:"evidenceId,omitempty"`
	EvaluatedBy    string  `json:"evaluatedBy"`
	EvaluationDate string  `json:"evaluationDate"`
	Method         string  `json:"method"`  // observation, supervision, error_reduction, output_review, feedback
	Outcome        string  `json:"outcome"` // effective, partially_effective, ineffective
	FollowUpAction *string `json:"followUpAction,omitempty"`
	CreatedAt      string  `json:"createdAt"`
}

type CreateEffectivenessEvaluationRequest struct {
	EvidenceID     *string `json:"evidenceId,omitempty"`
	EvaluatedBy    string  `json:"evaluatedBy"`
	EvaluationDate string  `json:"evaluationDate"`
	Method         string  `json:"method"`
	Outcome        string  `json:"outcome"`
	FollowUpAction *string `json:"followUpAction,omitempty"`
}

// =============================================================================
// DML Dashboard & Reporting Models
// =============================================================================

type CompetenceHealthResponse struct {
	TotalMatrixRoles          int     `json:"totalMatrixRoles"`
	AuthorizedCount           int     `json:"authorizedCount"`
	PendingAuthorizationCount int     `json:"pendingAuthorizationCount"`
	LegacyPendingCount        int     `json:"legacyPendingCount"`
	ActiveSupervisionCount    int     `json:"activeSupervisionCount"`
	OpenReassessmentsCount    int     `json:"openReassessmentsCount"`
	AuthorizationHealthPct    float64 `json:"authorizationHealthPct"`
}

type SupervisionQueueResponse struct {
	Active             []SupervisionPeriod `json:"active"`
	ReadyForAssessment []SupervisionPeriod `json:"readyForAssessment"`
}

type ReassessmentAlertsResponse struct {
	Open  []ReassessmentTrigger `json:"open"`
	Total int                   `json:"total"`
}
