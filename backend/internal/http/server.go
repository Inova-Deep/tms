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
	EmployeeRepo       *domain.EmployeeRepository
	ProfileRepo        *domain.ProfileRepository
	EvidenceRepo       *domain.EvidenceRepository
	EventRepo          *domain.EventRepository
	CertificationRepo  *domain.CertificationRepository
	CourseRepo         *domain.CourseRepository
	ComputationService *logic.ComputationService
	DashboardService   *logic.DashboardService
	MatrixRepo         *domain.MatrixRepository
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
