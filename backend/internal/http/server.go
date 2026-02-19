package http

import (
	"encoding/json"
	"net/http"
	"tms/internal/auth"
	"tms/internal/domain"
	"tms/internal/logic"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	EmployeeRepo            *domain.EmployeeRepository
	ProfileRepo             *domain.ProfileRepository
	EvidenceRepo            *domain.EvidenceRepository
	EventRepo               *domain.EventRepository
	CertificationRepo       *domain.CertificationRepository
	CourseRepo              *domain.CourseRepository
	ComputationService      *logic.ComputationService
	DashboardService        *logic.DashboardService
	MatrixRepo              *domain.MatrixRepository
	RoleRepo                *domain.RoleRepository
	ApplicabilityRepo       *domain.ApplicabilityRepository
	InductionRepo           *domain.InductionRepository
	WorkAuthorizationRepo   *domain.WorkAuthorizationRepository
	CompetenceAssessmentRepo *domain.CompetenceAssessmentRepository
	SupervisionRepo         *domain.SupervisionRepository
	ReassessmentRepo        *domain.ReassessmentRepository
	EffectivenessRepo       *domain.TrainingEffectivenessRepository
}

// Routes returns the API router
func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()

	// Auth endpoints (public)
	r.Post("/auth/login", s.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware, RequireAuth)

		// Employee routes
		r.Route("/employees", func(r chi.Router) {
			r.Get("/", s.GetEmployees)
			r.Get("/search", s.SearchEmployees)
			r.Get("/{id}", s.GetEmployee)
			r.Get("/{id}/records", s.GetEmployeeRecords)
			// DML sub-routes
			r.Get("/{id}/induction-record", s.GetEmployeeInduction)
			r.Get("/{id}/work-authorizations", s.GetEmployeeAuthorizations)
			r.Get("/{id}/competence-assessments", s.GetEmployeeAssessments)
			r.Get("/{id}/competence-status/{reqId}", s.GetEmployeeCompetenceStatus)
		})

		// Profile routes
		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", s.GetProfiles)
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateProfile)).ServeHTTP)
			r.Get("/{id}", s.GetProfile)
			r.Get("/{id}/requirements", s.GetProfileRequirements)

			// Members
			r.Get("/{id}/members", s.GetProfileMembers)
			r.Post("/{id}/members", s.AddProfileMember)
			r.Delete("/{id}/members/{employeeId}", s.RemoveProfileMember)

			// Exclusions
			r.Get("/{id}/exclusions", s.GetProfileExclusions)
			r.Post("/{id}/exclusions", s.AddProfileExclusion)
			r.Put("/{id}/exclusions/{exclusionId}", s.UpdateProfileExclusion)
			r.Delete("/{id}/exclusions/{exclusionId}", s.DeleteProfileExclusion)

			// Requirements
			r.Post("/{id}/requirements", RequireAdmin(http.HandlerFunc(s.AddProfileRequirement)).ServeHTTP)
			r.Delete("/{id}/requirements/{reqId}", RequireAdmin(http.HandlerFunc(s.RemoveProfileRequirement)).ServeHTTP)
		})

		// Event routes (admin only for writes)
		r.Route("/events", func(r chi.Router) {
			r.Get("/", s.GetEvents)
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateEvent)).ServeHTTP)
			r.Put("/{id}", RequireAdmin(http.HandlerFunc(s.UpdateEvent)).ServeHTTP)
			r.Post("/{id}/attendance/confirm", RequireAdmin(http.HandlerFunc(s.ConfirmAttendance)).ServeHTTP)
		})

		// Certification routes (admin only for writes)
		r.Route("/certifications", func(r chi.Router) {
			r.Get("/", s.GetCertifications)
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateCertification)).ServeHTTP)
		})

		r.Route("/courses", func(r chi.Router) {
			r.Get("/", s.GetCourses)
			r.Group(func(r chi.Router) {
				r.Use(RequireAdmin)
				r.Post("/", s.CreateCourse)
				r.Put("/{id}", s.UpdateCourse)
				r.Delete("/{id}", s.DeleteCourse)
			})
		})

		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/kpis", s.GetDashboardKPIs)
			r.Get("/grid", s.GetDashboardGrid)
			r.Get("/drilldown", s.GetDashboardDrilldown)
		})

		r.Get("/matrix", s.handleGetMatrix)

		// ── DML: Roles ──────────────────────────────────────────────
		r.Route("/roles", func(r chi.Router) {
			r.Get("/", s.GetRoles)
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateRole)).ServeHTTP)
			r.Get("/matrix-included", s.GetMatrixRoles)
			r.Get("/awareness-only", s.GetAwarenessRoles)
			r.Get("/{id}", s.GetRole)
			r.Put("/{id}", RequireAdmin(http.HandlerFunc(s.UpdateRole)).ServeHTTP)
			r.Get("/{id}/applicability-decision", s.GetApplicabilityDecision)
			r.Post("/{id}/applicability-decision", RequireAdmin(http.HandlerFunc(s.CreateApplicabilityDecision)).ServeHTTP)
		})

		// ── DML: Escalations ────────────────────────────────────────
		r.Post("/escalations/{id}/resolve", RequireAdmin(http.HandlerFunc(s.ResolveEscalation)).ServeHTTP)

		// ── DML: Inductions ─────────────────────────────────────────
		r.Route("/induction-records", func(r chi.Router) {
			r.Get("/", s.GetPendingInductions)
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateInductionRecord)).ServeHTTP)
		})

		// ── DML: Work Authorizations ────────────────────────────────
		r.Route("/work-authorizations", func(r chi.Router) {
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateWorkAuthorization)).ServeHTTP)
			r.Get("/expiring", s.GetExpiringAuthorizations)
			r.Get("/legacy-pending", s.GetLegacyPendingAuthorizations)
			r.Put("/{id}/revoke", RequireAdmin(http.HandlerFunc(s.RevokeWorkAuthorization)).ServeHTTP)
		})

		// ── DML: Competence Assessments ─────────────────────────────
		r.Route("/competence-assessments", func(r chi.Router) {
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateCompetenceAssessment)).ServeHTTP)
			r.Get("/pending", s.GetPendingAssessments)
		})

		// ── DML: Supervision ────────────────────────────────────────
		r.Route("/supervision-periods", func(r chi.Router) {
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateSupervisionPeriod)).ServeHTTP)
			r.Put("/{id}/complete", RequireAdmin(http.HandlerFunc(s.CompleteSupervision)).ServeHTTP)
		})
		r.Get("/supervision/active", s.GetActiveSupervision)
		r.Get("/supervision/ready-for-assessment", s.GetReadyForAssessment)

		// ── DML: Reassessment Triggers ──────────────────────────────
		r.Route("/reassessment-triggers", func(r chi.Router) {
			r.Post("/", RequireAdmin(http.HandlerFunc(s.CreateReassessmentTrigger)).ServeHTTP)
			r.Get("/open", s.GetOpenReassessments)
			r.Put("/{id}/resolve", RequireAdmin(http.HandlerFunc(s.ResolveReassessmentTrigger)).ServeHTTP)
		})

		// ── DML: Training Effectiveness ─────────────────────────────
		r.Post("/training-effectiveness-evaluations", RequireAdmin(http.HandlerFunc(s.CreateEffectivenessEvaluation)).ServeHTTP)

		// ── DML: Dashboard KPIs ──────────────────────────────────────
		r.Get("/dashboard/competence-health", s.GetCompetenceHealth)
		r.Get("/dashboard/supervision-queue", s.GetSupervisionQueue)
		r.Get("/dashboard/reassessment-alerts", s.GetReassessmentAlerts)

		// ── DML: Reports ─────────────────────────────────────────────
		r.Get("/reports/dml-qa-reg-5024", s.GetDMLQAReport)
		r.Get("/reports/competence-degradation-risk", s.GetCompetenceDegradationRisk)
	})

	return r
}

// ==================== Auth Handlers ====================

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != "admin" && req.Role != "employee" {
		respondWithError(w, "Invalid role. Must be 'admin' or 'employee'", http.StatusBadRequest)
		return
	}

	// For employee role, employeeId is required
	if req.Role == "employee" && req.EmployeeID == "" {
		respondWithError(w, "employeeId is required for employee role", http.StatusBadRequest)
		return
	}

	// For employee role, verify employee exists
	if req.Role == "employee" {
		emp, err := s.EmployeeRepo.GetByID(req.EmployeeID)
		if err != nil || emp == nil {
			respondWithError(w, "Employee not found", http.StatusNotFound)
			return
		}
	}

	token := auth.GenerateToken(req.Role, req.EmployeeID)

	respondWithJSON(w, auth.LoginResponse{
		Token:      token,
		Role:       req.Role,
		EmployeeID: req.EmployeeID,
	})
}

// ==================== Employee Handlers ====================

func (s *Server) GetEmployees(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)

	var employees []domain.Employee
	var err error

	if user.IsAdmin() {
		employees, err = s.EmployeeRepo.GetAll()
	} else {
		// Employee can only see themselves
		emp, err := s.EmployeeRepo.GetByID(user.EmployeeID)
		if err != nil {
			respondWithError(w, "Failed to get employee", http.StatusInternalServerError)
			return
		}
		if emp != nil {
			employees = []domain.Employee{*emp}
		}
	}

	if err != nil {
		respondWithError(w, "Failed to get employees", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, employees)
}

func (s *Server) SearchEmployees(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	query := r.URL.Query().Get("q")
	employees, err := s.EmployeeRepo.Search(query)
	if err != nil {
		respondWithError(w, "Failed to search employees", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, employees)
}

func (s *Server) GetEmployee(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	id := chi.URLParam(r, "id")

	// Validate access
	if !user.CanAccessEmployee(id) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}

	employee, err := s.EmployeeRepo.GetByID(id)
	if err != nil {
		respondWithError(w, "Failed to get employee", http.StatusInternalServerError)
		return
	}
	if employee == nil {
		respondWithError(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Fetch assignments
	profileIDs, err := s.EmployeeRepo.GetAssignments(id)
	if err != nil {
		respondWithError(w, "Failed to get assignments", http.StatusInternalServerError)
		return
	}

	// Fetch evidence
	evidenceList, err := s.EvidenceRepo.GetByEmployee(id)
	if err != nil {
		respondWithError(w, "Failed to get evidence", http.StatusInternalServerError)
		return
	}

	for i := range evidenceList {
		if evidenceList[i].CourseID != "" {
			course, err := s.CourseRepo.GetByID(evidenceList[i].CourseID)
			if err == nil {
				evidenceList[i].Course = course
				evidenceList[i].RequirementName = course.Name
			}
		}
	}

	// Compute status
	extendedEmployee := domain.EmployeeWithStatus{
		Employee:      *employee,
		ProfileStatus: make(map[string]string),
		Requirements:  []domain.RequirementWithStatus{},
	}

	for _, pid := range profileIDs {
		reqs, err := s.ProfileRepo.GetRequirements(pid)
		if err != nil {
			continue
		}

		profileEligible := true
		profileAtRisk := false

		for _, req := range reqs {
			status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)
			extendedEmployee.Requirements = append(extendedEmployee.Requirements, domain.RequirementWithStatus{
				Requirement: req,
				Status:      status.Status,
				ExpiryDate: func() *string {
					if status.ExpiryDate == nil {
						return nil
					}
					s := status.ExpiryDate.Format("2006-01-02")
					return &s
				}(),
			})

			if status.Status == "Missing" || status.Status == "Expired" {
				profileEligible = false
			} else if status.Status == "Expiring" {
				profileAtRisk = true
			}
		}

		if profileEligible {
			if profileAtRisk {
				extendedEmployee.ProfileStatus[pid] = "At Risk"
			} else {
				extendedEmployee.ProfileStatus[pid] = "Eligible"
			}
		} else {
			extendedEmployee.ProfileStatus[pid] = "Not Eligible"
		}
	}

	respondWithJSON(w, extendedEmployee)
}

func (s *Server) GetEmployeeRecords(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	id := chi.URLParam(r, "id")

	// Validate access
	if !user.CanAccessEmployee(id) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}

	records, err := s.EvidenceRepo.GetByEmployee(id)
	if err != nil {
		respondWithError(w, "Failed to get records", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, records)
}

// ==================== Profile Handlers ====================

func (s *Server) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := s.ProfileRepo.GetAllWithCounts()
	if err != nil {
		respondWithError(w, "Failed to get profiles", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, profiles)
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	profile, err := s.ProfileRepo.GetByID(id)
	if err != nil {
		respondWithError(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}
	if profile == nil {
		respondWithError(w, "Profile not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, profile)
}

func (s *Server) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		respondWithError(w, "Name is required", http.StatusBadRequest)
		return
	}

	id, err := s.ProfileRepo.GenerateID()
	if err != nil {
		respondWithError(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}

	profile := &domain.Profile{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.ProfileRepo.Create(profile); err != nil {
		respondWithError(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (s *Server) GetProfileRequirements(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requirements, err := s.ProfileRepo.GetRequirements(id)
	if err != nil {
		respondWithError(w, "Failed to get requirements", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, requirements)
}

func (s *Server) GetProfileMembers(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")

	members, err := s.ProfileRepo.GetMembers(profileID)
	if err != nil {
		respondWithError(w, "Failed to get members", http.StatusInternalServerError)
		return
	}

	// Compute compliance status for each member
	reqs, _ := s.ProfileRepo.GetRequirements(profileID)
	for i, m := range members {
		evidenceList, _ := s.EvidenceRepo.GetByEmployee(m.ID)
		if len(evidenceList) == 0 {
			members[i].Status = "not_started"
			continue
		}

		allValid := true
		anyExpiring := false
		anyExpired := false

		for _, req := range reqs {
			status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)
			switch status.Status {
			case "Missing":
				allValid = false
			case "Expired":
				allValid = false
				anyExpired = true
			case "Expiring":
				anyExpiring = true
			}
		}

		if !allValid {
			if anyExpired {
				members[i].Status = "expired"
			} else {
				members[i].Status = "not_started"
			}
		} else if anyExpiring {
			members[i].Status = "expiring"
		} else {
			members[i].Status = "compliant"
		}
	}

	respondWithJSON(w, domain.ProfileMembersResponse{
		Members: members,
		Total:   len(members),
	})
}

func (s *Server) AddProfileMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	profileID := chi.URLParam(r, "id")
	var req struct {
		EmployeeID string `json:"employeeId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.ProfileRepo.AddMember(profileID, req.EmployeeID); err != nil {
		respondWithError(w, "Failed to add member", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) RemoveProfileMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	profileID := chi.URLParam(r, "id")
	employeeID := chi.URLParam(r, "employeeId")

	if err := s.ProfileRepo.RemoveMember(profileID, employeeID); err != nil {
		respondWithError(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) AddProfileRequirement(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")

	var req struct {
		CourseID string `json:"courseId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.CourseID == "" {
		respondWithError(w, "courseId is required", http.StatusBadRequest)
		return
	}

	if _, err := s.CourseRepo.GetByID(req.CourseID); err != nil {
		respondWithError(w, "Course not found", http.StatusNotFound)
		return
	}

	if err := s.ProfileRepo.AddRequirement(profileID, req.CourseID); err != nil {
		respondWithError(w, "Failed to add requirement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) RemoveProfileRequirement(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")
	reqID := chi.URLParam(r, "reqId")

	if err := s.ProfileRepo.RemoveRequirement(profileID, reqID); err != nil {
		respondWithError(w, "Failed to remove requirement", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ==================== Profile Exclusion Handlers ====================

func (s *Server) GetProfileExclusions(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "id")

	exclusions, err := s.ProfileRepo.GetExclusions(profileID)
	if err != nil {
		respondWithError(w, "Failed to get exclusions", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, exclusions)
}

func (s *Server) AddProfileExclusion(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	profileID := chi.URLParam(r, "id")
	var req domain.CreateExclusionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Reason == "" || req.ExpiryDate == "" {
		respondWithError(w, "Reason and expiryDate are required", http.StatusBadRequest)
		return
	}

	exclusion, err := s.ProfileRepo.CreateExclusion(profileID, req)
	if err != nil {
		respondWithError(w, "Failed to add exclusion", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, exclusion)
}

func (s *Server) UpdateProfileExclusion(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	exclusionID := chi.URLParam(r, "exclusionId")
	var req domain.UpdateExclusionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Reason == "" || req.ExpiryDate == "" {
		respondWithError(w, "Reason and expiryDate are required", http.StatusBadRequest)
		return
	}

	// Verify exclusion exists
	existing, err := s.ProfileRepo.GetExclusionByID(exclusionID)
	if err != nil {
		respondWithError(w, "Failed to get exclusion", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		respondWithError(w, "Exclusion not found", http.StatusNotFound)
		return
	}

	if err := s.ProfileRepo.UpdateExclusion(exclusionID, req); err != nil {
		respondWithError(w, "Failed to update exclusion", http.StatusInternalServerError)
		return
	}

	// Return updated exclusion
	updated, _ := s.ProfileRepo.GetExclusionByID(exclusionID)
	respondWithJSON(w, updated)
}

func (s *Server) DeleteProfileExclusion(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if !user.IsAdmin() {
		respondWithError(w, "Admin access required", http.StatusForbidden)
		return
	}

	exclusionID := chi.URLParam(r, "exclusionId")

	// Verify exclusion exists
	existing, err := s.ProfileRepo.GetExclusionByID(exclusionID)
	if err != nil {
		respondWithError(w, "Failed to get exclusion", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		respondWithError(w, "Exclusion not found", http.StatusNotFound)
		return
	}

	if err := s.ProfileRepo.DeleteExclusion(exclusionID); err != nil {
		respondWithError(w, "Failed to delete exclusion", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, map[string]string{"status": "ok"})
}

// ==================== Event Handlers ====================

func (s *Server) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.EventRepo.GetAll()
	if err != nil {
		respondWithError(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	for i := range events {
		if events[i].CourseID != "" {
			course, err := s.CourseRepo.GetByID(events[i].CourseID)
			if err == nil {
				events[i].Course = course
				events[i].CourseName = course.Name
			}
		}
	}

	respondWithJSON(w, events)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CourseID == "" {
		respondWithError(w, "courseId is required", http.StatusBadRequest)
		return
	}

	_, err := s.CourseRepo.GetByID(req.CourseID)
	if err == domain.ErrNotFound {
		respondWithError(w, "Course not found", http.StatusBadRequest)
		return
	}
	if err != nil {
		respondWithError(w, "Failed to validate course", http.StatusInternalServerError)
		return
	}

	event, err := s.EventRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create event", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, event)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	var req domain.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CourseID != nil {
		_, err := s.CourseRepo.GetByID(*req.CourseID)
		if err == domain.ErrNotFound {
			respondWithError(w, "Course not found", http.StatusBadRequest)
			return
		}
		if err != nil {
			respondWithError(w, "Failed to validate course", http.StatusInternalServerError)
			return
		}
	}

	event, err := s.EventRepo.Update(eventID, req)
	if err != nil {
		respondWithError(w, "Failed to update event", http.StatusInternalServerError)
		return
	}
	if event == nil {
		respondWithError(w, "Event not found", http.StatusNotFound)
		return
	}

	if event.CourseID != "" {
		course, err := s.CourseRepo.GetByID(event.CourseID)
		if err == nil {
			event.Course = course
			event.CourseName = course.Name
		}
	}

	respondWithJSON(w, event)
}

func (s *Server) ConfirmAttendance(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	// Try to parse as enhanced format first
	var enhancedReq domain.ConfirmAttendanceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&enhancedReq); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if this is the enhanced format (has attendees array with objects)
	if len(enhancedReq.Attendees) > 0 {
		if err := s.EventRepo.EnhancedConfirmAttendance(eventID, enhancedReq, s.EvidenceRepo); err != nil {
			respondWithError(w, "Failed to confirm attendance", http.StatusInternalServerError)
			return
		}
	} else {
		status := "completed"
		_, err := s.EventRepo.Update(eventID, domain.UpdateEventRequest{Status: &status})
		if err != nil {
			respondWithError(w, "Failed to confirm attendance", http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, map[string]string{"status": "ok"})
}

// ==================== Certification Handlers ====================

func (s *Server) GetCertifications(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("employeeId")

	var certs []domain.Certification
	var err error

	if employeeID != "" {
		certs, err = s.CertificationRepo.GetByEmployee(employeeID)
	} else {
		certs, err = s.CertificationRepo.GetAll()
	}

	if err != nil {
		respondWithError(w, "Failed to get certifications", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, certs)
}

func (s *Server) CreateCertification(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCertificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CourseID == "" {
		respondWithError(w, "courseId is required", http.StatusBadRequest)
		return
	}

	_, err := s.CourseRepo.GetByID(req.CourseID)
	if err == domain.ErrNotFound {
		respondWithError(w, "Course not found", http.StatusBadRequest)
		return
	}
	if err != nil {
		respondWithError(w, "Failed to validate course", http.StatusInternalServerError)
		return
	}

	cert, err := s.CertificationRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create certification", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, cert)
}

func (s *Server) GetCourses(w http.ResponseWriter, r *http.Request) {
	courseType := r.URL.Query().Get("type")

	var courses []domain.Course
	var err error

	if courseType != "" {
		if courseType != "course" && courseType != "certification" {
			respondWithError(w, "Invalid type. Must be 'course' or 'certification'", http.StatusBadRequest)
			return
		}
		courses, err = s.CourseRepo.GetByType(courseType)
	} else {
		courses, err = s.CourseRepo.GetAll()
	}

	if err != nil {
		respondWithError(w, "Failed to get courses", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, courses)
}

func (s *Server) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		respondWithError(w, "name is required", http.StatusBadRequest)
		return
	}

	if req.Type != "course" && req.Type != "certification" {
		respondWithError(w, "type must be 'course' or 'certification'", http.StatusBadRequest)
		return
	}

	course := &domain.Course{
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		Description:    req.Description,
		ValidityMonths: req.ValidityMonths,
	}

	if err := s.CourseRepo.Create(course); err != nil {
		respondWithError(w, "Failed to create course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, course)
}

func (s *Server) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	existing, err := s.CourseRepo.GetByID(id)
	if err == domain.ErrNotFound {
		respondWithError(w, "Course not found", http.StatusNotFound)
		return
	}
	if err != nil {
		respondWithError(w, "Failed to get course", http.StatusInternalServerError)
		return
	}

	var req domain.UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Code != nil {
		existing.Code = req.Code
	}
	if req.Type != nil {
		if *req.Type != "course" && *req.Type != "certification" {
			respondWithError(w, "type must be 'course' or 'certification'", http.StatusBadRequest)
			return
		}
		existing.Type = *req.Type
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.ValidityMonths != nil {
		existing.ValidityMonths = *req.ValidityMonths
	}

	if err := s.CourseRepo.Update(existing); err != nil {
		respondWithError(w, "Failed to update course", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, existing)
}

func (s *Server) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := s.CourseRepo.GetByID(id); err == domain.ErrNotFound {
		respondWithError(w, "Course not found", http.StatusNotFound)
		return
	}

	if err := s.CourseRepo.Delete(id); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, map[string]string{"status": "ok"})
}

// ==================== Dashboard Handlers ====================

func (s *Server) GetDashboardKPIs(w http.ResponseWriter, r *http.Request) {
	// Get filter params
	department := r.URL.Query().Get("department")
	site := r.URL.Query().Get("site")
	profileID := r.URL.Query().Get("profile")

	kpis, err := s.DashboardService.GetKPIs(department, site, profileID)
	if err != nil {
		respondWithError(w, "Failed to get KPIs", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, kpis)
}

func (s *Server) GetDashboardGrid(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")       // "course" or "profile"
	groupBy := r.URL.Query().Get("groupBy") // "department", "site", "costCenter"
	department := r.URL.Query().Get("department")
	site := r.URL.Query().Get("site")

	if mode == "" {
		mode = "course"
	}
	if groupBy == "" {
		groupBy = "department"
	}

	grid, err := s.DashboardService.GetCoverageGrid(mode, groupBy, department, site)
	if err != nil {
		respondWithError(w, "Failed to get grid data", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, grid)
}

func (s *Server) GetDashboardDrilldown(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	groupBy := r.URL.Query().Get("groupBy")
	groupValue := r.URL.Query().Get("groupValue")
	itemID := r.URL.Query().Get("itemId")
	department := r.URL.Query().Get("department")
	site := r.URL.Query().Get("site")

	employees, err := s.DashboardService.GetDrilldown(mode, groupBy, groupValue, itemID, department, site)
	if err != nil {
		respondWithError(w, "Failed to get drilldown data", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, employees)
}

func (s *Server) handleGetMatrix(w http.ResponseWriter, r *http.Request) {
	matrix, err := s.MatrixRepo.GetMatrix()
	if err != nil {
		respondWithError(w, "Failed to get matrix data", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, matrix)
}

// ==================== DML: Role Handlers ====================

func (s *Server) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.RoleRepo.GetAll()
	if err != nil {
		respondWithError(w, "Failed to get roles", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, roles)
}

func (s *Server) GetRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, err := s.RoleRepo.GetByID(id)
	if err != nil {
		respondWithError(w, "Failed to get role", http.StatusInternalServerError)
		return
	}
	if role == nil {
		respondWithError(w, "Role not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, role)
}

func (s *Server) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		respondWithError(w, "name is required", http.StatusBadRequest)
		return
	}
	role, err := s.RoleRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, role)
}

func (s *Server) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	role, err := s.RoleRepo.Update(id, req)
	if err != nil {
		respondWithError(w, "Failed to update role", http.StatusInternalServerError)
		return
	}
	if role == nil {
		respondWithError(w, "Role not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, role)
}

func (s *Server) GetMatrixRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.RoleRepo.GetMatrixRoles()
	if err != nil {
		respondWithError(w, "Failed to get matrix roles", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, roles)
}

func (s *Server) GetAwarenessRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.RoleRepo.GetAwarenessRoles()
	if err != nil {
		respondWithError(w, "Failed to get awareness roles", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, roles)
}

// ==================== DML: Applicability Decision Handlers ====================

func (s *Server) GetApplicabilityDecision(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	decision, err := s.ApplicabilityRepo.GetByRoleID(roleID)
	if err != nil {
		respondWithError(w, "Failed to get applicability decision", http.StatusInternalServerError)
		return
	}
	if decision == nil {
		respondWithError(w, "No applicability decision found for this role", http.StatusNotFound)
		return
	}
	respondWithJSON(w, decision)
}

func (s *Server) CreateApplicabilityDecision(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	var req domain.CreateApplicabilityDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.DecidedBy == "" {
		respondWithError(w, "decidedBy is required", http.StatusBadRequest)
		return
	}
	decision, err := s.ApplicabilityRepo.Create(roleID, req)
	if err != nil {
		respondWithError(w, "Failed to create applicability decision", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, decision)
}

func (s *Server) ResolveEscalation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.ResolveEscalationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.ResolvedBy == "" {
		respondWithError(w, "resolvedBy is required", http.StatusBadRequest)
		return
	}
	decision, err := s.ApplicabilityRepo.ResolveEscalation(id, req)
	if err != nil {
		respondWithError(w, "Failed to resolve escalation", http.StatusInternalServerError)
		return
	}
	if decision == nil {
		respondWithError(w, "Escalation not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, decision)
}

// ==================== DML: Induction Handlers ====================

func (s *Server) GetPendingInductions(w http.ResponseWriter, r *http.Request) {
	records, err := s.InductionRepo.GetPending()
	if err != nil {
		respondWithError(w, "Failed to get pending inductions", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, records)
}

func (s *Server) CreateInductionRecord(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateInductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EmployeeID == "" || req.ConductedBy == "" || req.InductionDate == "" {
		respondWithError(w, "employeeId, conductedBy, and inductionDate are required", http.StatusBadRequest)
		return
	}
	record, err := s.InductionRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create induction record", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, record)
}

func (s *Server) GetEmployeeInduction(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "id")
	user := GetUserFromContext(r)
	if !user.CanAccessEmployee(employeeID) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}
	records, err := s.InductionRepo.GetByEmployee(employeeID)
	if err != nil {
		respondWithError(w, "Failed to get induction records", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, records)
}

// ==================== DML: Work Authorization Handlers ====================

func (s *Server) CreateWorkAuthorization(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateWorkAuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EmployeeID == "" || req.AuthorizationType == "" || req.AuthorizedBy == "" {
		respondWithError(w, "employeeId, authorizationType, and authorizedBy are required", http.StatusBadRequest)
		return
	}
	auth, err := s.WorkAuthorizationRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create work authorization", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, auth)
}

func (s *Server) GetExpiringAuthorizations(w http.ResponseWriter, r *http.Request) {
	daysAhead := 30
	auths, err := s.WorkAuthorizationRepo.GetExpiring(daysAhead)
	if err != nil {
		respondWithError(w, "Failed to get expiring authorizations", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, auths)
}

func (s *Server) GetLegacyPendingAuthorizations(w http.ResponseWriter, r *http.Request) {
	auths, err := s.WorkAuthorizationRepo.GetLegacyPending()
	if err != nil {
		respondWithError(w, "Failed to get legacy pending authorizations", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, auths)
}

func (s *Server) RevokeWorkAuthorization(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.RevokeAuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.RevokeReason == "" {
		respondWithError(w, "revokeReason is required", http.StatusBadRequest)
		return
	}
	if err := s.WorkAuthorizationRepo.Revoke(id, req); err != nil {
		respondWithError(w, "Failed to revoke authorization", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, map[string]string{"status": "ok"})
}

func (s *Server) GetEmployeeAuthorizations(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "id")
	user := GetUserFromContext(r)
	if !user.CanAccessEmployee(employeeID) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}
	auths, err := s.WorkAuthorizationRepo.GetByEmployee(employeeID)
	if err != nil {
		respondWithError(w, "Failed to get work authorizations", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, auths)
}

func (s *Server) GetEmployeeCompetenceStatus(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "id")
	reqID := chi.URLParam(r, "reqId")
	user := GetUserFromContext(r)
	if !user.CanAccessEmployee(employeeID) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}
	auth, err := s.WorkAuthorizationRepo.GetByEmployeeAndRequirement(employeeID, reqID)
	if err != nil {
		respondWithError(w, "Failed to get competence status", http.StatusInternalServerError)
		return
	}
	if auth == nil {
		respondWithJSON(w, map[string]string{"status": "REQUIRED"})
		return
	}
	respondWithJSON(w, auth)
}

// ==================== DML: Competence Assessment Handlers ====================

func (s *Server) CreateCompetenceAssessment(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCompetenceAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EmployeeID == "" || req.AssessedBy == "" || req.AssessmentDate == "" || req.WorkActivitiesCovered == "" || req.Outcome == "" {
		respondWithError(w, "employeeId, assessedBy, assessmentDate, workActivitiesCovered, and outcome are required", http.StatusBadRequest)
		return
	}
	assessment, err := s.CompetenceAssessmentRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create competence assessment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, assessment)
}

func (s *Server) GetPendingAssessments(w http.ResponseWriter, r *http.Request) {
	assessments, err := s.CompetenceAssessmentRepo.GetPending()
	if err != nil {
		respondWithError(w, "Failed to get pending assessments", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, assessments)
}

func (s *Server) GetEmployeeAssessments(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "id")
	user := GetUserFromContext(r)
	if !user.CanAccessEmployee(employeeID) {
		respondWithError(w, "Access denied", http.StatusForbidden)
		return
	}
	assessments, err := s.CompetenceAssessmentRepo.GetByEmployee(employeeID)
	if err != nil {
		respondWithError(w, "Failed to get assessments", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, assessments)
}

// ==================== DML: Supervision Handlers ====================

func (s *Server) CreateSupervisionPeriod(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSupervisionPeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EmployeeID == "" || req.SupervisorID == "" || req.StartDate == "" {
		respondWithError(w, "employeeId, supervisorId, and startDate are required", http.StatusBadRequest)
		return
	}
	period, err := s.SupervisionRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create supervision period", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, period)
}

func (s *Server) CompleteSupervision(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.CompleteSupervisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EndDate == "" {
		respondWithError(w, "endDate is required", http.StatusBadRequest)
		return
	}
	period, err := s.SupervisionRepo.Complete(id, req)
	if err != nil {
		respondWithError(w, "Failed to complete supervision period", http.StatusInternalServerError)
		return
	}
	if period == nil {
		respondWithError(w, "Supervision period not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, period)
}

func (s *Server) GetActiveSupervision(w http.ResponseWriter, r *http.Request) {
	periods, err := s.SupervisionRepo.GetActive()
	if err != nil {
		respondWithError(w, "Failed to get active supervision periods", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, periods)
}

func (s *Server) GetReadyForAssessment(w http.ResponseWriter, r *http.Request) {
	periods, err := s.SupervisionRepo.GetReadyForAssessment()
	if err != nil {
		respondWithError(w, "Failed to get supervision periods ready for assessment", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, periods)
}

// ==================== DML: Reassessment Trigger Handlers ====================

func (s *Server) CreateReassessmentTrigger(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateReassessmentTriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EmployeeID == "" || req.TriggerType == "" || req.TriggeredBy == "" {
		respondWithError(w, "employeeId, triggerType, and triggeredBy are required", http.StatusBadRequest)
		return
	}
	trigger, err := s.ReassessmentRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create reassessment trigger", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, trigger)
}

func (s *Server) GetOpenReassessments(w http.ResponseWriter, r *http.Request) {
	triggers, err := s.ReassessmentRepo.GetOpen()
	if err != nil {
		respondWithError(w, "Failed to get open reassessments", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, triggers)
}

func (s *Server) ResolveReassessmentTrigger(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.ResolveReassessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.ResolutionNotes == "" {
		respondWithError(w, "resolutionNotes is required", http.StatusBadRequest)
		return
	}
	trigger, err := s.ReassessmentRepo.Resolve(id, req)
	if err != nil {
		respondWithError(w, "Failed to resolve reassessment trigger", http.StatusInternalServerError)
		return
	}
	if trigger == nil {
		respondWithError(w, "Reassessment trigger not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, trigger)
}

// ==================== DML: Training Effectiveness Handlers ====================

func (s *Server) CreateEffectivenessEvaluation(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateEffectivenessEvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.EvaluatedBy == "" || req.EvaluationDate == "" || req.Method == "" || req.Outcome == "" {
		respondWithError(w, "evaluatedBy, evaluationDate, method, and outcome are required", http.StatusBadRequest)
		return
	}
	eval, err := s.EffectivenessRepo.Create(req)
	if err != nil {
		respondWithError(w, "Failed to create effectiveness evaluation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respondWithJSON(w, eval)
}

// ==================== DML: Competence Dashboard Handlers ====================

func (s *Server) GetCompetenceHealth(w http.ResponseWriter, r *http.Request) {
	matrixRoles, err := s.RoleRepo.GetMatrixRoles()
	if err != nil {
		respondWithError(w, "Failed to get roles", http.StatusInternalServerError)
		return
	}

	legacyPending, err := s.WorkAuthorizationRepo.GetLegacyPending()
	if err != nil {
		respondWithError(w, "Failed to get legacy pending", http.StatusInternalServerError)
		return
	}

	activeSupervision, err := s.SupervisionRepo.GetActive()
	if err != nil {
		respondWithError(w, "Failed to get supervision", http.StatusInternalServerError)
		return
	}

	openReassessments, err := s.ReassessmentRepo.GetOpen()
	if err != nil {
		respondWithError(w, "Failed to get reassessments", http.StatusInternalServerError)
		return
	}

	authorized := 0
	pendingAuth := 0
	for _, role := range matrixRoles {
		auths, _ := s.WorkAuthorizationRepo.GetByEmployee(role.ID)
		for _, a := range auths {
			switch a.AuthorizationState {
			case "AUTHORIZED":
				authorized++
			case "PENDING_AUTHORIZATION":
				pendingAuth++
			}
		}
	}

	total := authorized + pendingAuth
	healthPct := 0.0
	if total > 0 {
		healthPct = float64(authorized) / float64(total) * 100
	}

	resp := domain.CompetenceHealthResponse{
		TotalMatrixRoles:          len(matrixRoles),
		AuthorizedCount:           authorized,
		PendingAuthorizationCount: pendingAuth,
		LegacyPendingCount:        len(legacyPending),
		ActiveSupervisionCount:    len(activeSupervision),
		OpenReassessmentsCount:    len(openReassessments),
		AuthorizationHealthPct:    healthPct,
	}
	respondWithJSON(w, resp)
}

func (s *Server) GetSupervisionQueue(w http.ResponseWriter, r *http.Request) {
	active, err := s.SupervisionRepo.GetActive()
	if err != nil {
		respondWithError(w, "Failed to get active supervision", http.StatusInternalServerError)
		return
	}

	ready, err := s.SupervisionRepo.GetReadyForAssessment()
	if err != nil {
		respondWithError(w, "Failed to get ready for assessment", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, domain.SupervisionQueueResponse{
		Active:             active,
		ReadyForAssessment: ready,
	})
}

func (s *Server) GetReassessmentAlerts(w http.ResponseWriter, r *http.Request) {
	open, err := s.ReassessmentRepo.GetOpen()
	if err != nil {
		respondWithError(w, "Failed to get open reassessments", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, domain.ReassessmentAlertsResponse{
		Open:  open,
		Total: len(open),
	})
}

// ==================== DML: Report Handlers ====================

func (s *Server) GetDMLQAReport(w http.ResponseWriter, r *http.Request) {
	legacyPending, _ := s.WorkAuthorizationRepo.GetLegacyPending()
	expiring, _ := s.WorkAuthorizationRepo.GetExpiring(30)
	openReassessments, _ := s.ReassessmentRepo.GetOpen()
	activeSupervision, _ := s.SupervisionRepo.GetActive()
	pendingInductions, _ := s.InductionRepo.GetPending()

	report := map[string]interface{}{
		"reportTitle":           "DML Competence Assurance Summary (DML-QA-REG-5024)",
		"generatedAt":           "now",
		"legacyPendingCount":    len(legacyPending),
		"expiringIn30DaysCount": len(expiring),
		"openReassessmentsCount": len(openReassessments),
		"activeSupervisionCount": len(activeSupervision),
		"pendingInductionsCount": len(pendingInductions),
		"legacyPending":          legacyPending,
		"expiring":               expiring,
		"openReassessments":      openReassessments,
	}
	respondWithJSON(w, report)
}

func (s *Server) GetCompetenceDegradationRisk(w http.ResponseWriter, r *http.Request) {
	expiring, _ := s.WorkAuthorizationRepo.GetExpiring(30)
	openReassessments, _ := s.ReassessmentRepo.GetOpen()
	pendingAssessments, _ := s.CompetenceAssessmentRepo.GetPending()

	report := map[string]interface{}{
		"reportTitle":              "Competence Degradation Risk Assessment",
		"generatedAt":              "now",
		"expiringAuthorizationsCount": len(expiring),
		"openReassessmentsCount":   len(openReassessments),
		"pendingAssessmentsCount":  len(pendingAssessments),
		"expiringAuthorizations":   expiring,
		"openReassessments":        openReassessments,
		"pendingAssessments":       pendingAssessments,
	}
	respondWithJSON(w, report)
}

// ==================== Helpers ====================

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}
