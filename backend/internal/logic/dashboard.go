package logic

import (
	"tms/internal/domain"
)

// DashboardService handles dashboard computations
type DashboardService struct {
	EmployeeRepo      *domain.EmployeeRepository
	ProfileRepo       *domain.ProfileRepository
	EvidenceRepo      *domain.EvidenceRepository
	ComputationService *ComputationService
}

// GetKPIs calculates dashboard KPIs
func (s *DashboardService) GetKPIs(department, site, profileID string) (*domain.DashboardKPIs, error) {
	employees, err := s.EmployeeRepo.GetByFilter(department, site)
	if err != nil {
		return nil, err
	}

	kpis := &domain.DashboardKPIs{
		TotalEmployees: len(employees),
	}

	totalRequirements := 0
	compliantRequirements := 0

	for _, emp := range employees {
		profileIDs, err := s.EmployeeRepo.GetAssignments(emp.ID)
		if err != nil {
			continue
		}

		// Filter by profile if specified
		if profileID != "" {
			found := false
			for _, pid := range profileIDs {
				if pid == profileID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		evidenceList, err := s.EvidenceRepo.GetByEmployee(emp.ID)
		if err != nil {
			continue
		}

		for _, pid := range profileIDs {
			if profileID != "" && pid != profileID {
				continue
			}

			reqs, err := s.ProfileRepo.GetRequirements(pid)
			if err != nil {
				continue
			}

			for _, req := range reqs {
				totalRequirements++
				status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)
				
				switch status.Status {
				case "Valid":
					kpis.Valid++
					compliantRequirements++
				case "Expiring":
					kpis.ExpiringSoon++
					compliantRequirements++
				case "Expired":
					kpis.Expired++
				case "Missing":
					kpis.Missing++
				}
			}
		}
	}

	if totalRequirements > 0 {
		kpis.OverallCompliance = float64(compliantRequirements) / float64(totalRequirements) * 100
	}

	return kpis, nil
}

// GetCoverageGrid generates coverage grid data
func (s *DashboardService) GetCoverageGrid(mode, groupBy, department, site string) ([]domain.GridCell, error) {
	employees, err := s.EmployeeRepo.GetByFilter(department, site)
	if err != nil {
		return nil, err
	}

	// Group employees
	groupMap := make(map[string][]domain.Employee)
	for _, emp := range employees {
		var key string
		switch groupBy {
		case "department":
			key = emp.Department
		case "site":
			key = emp.Site
		case "costCenter":
			key = emp.CostCenter
		default:
			key = emp.Department
		}
		groupMap[key] = append(groupMap[key], emp)
	}

	cells := make([]domain.GridCell, 0)

	if mode == "profile" {
		// Profile coverage mode
		profiles, _ := s.ProfileRepo.GetAll()
		
		for groupValue, groupEmployees := range groupMap {
			for _, profile := range profiles {
				cell := domain.GridCell{
					GroupValue: groupValue,
					ItemID:     profile.ID,
					ItemName:   profile.Name,
				}

				for _, emp := range groupEmployees {
					profileIDs, _ := s.EmployeeRepo.GetAssignments(emp.ID)
					hasProfile := false
					for _, pid := range profileIDs {
						if pid == profile.ID {
							hasProfile = true
							break
						}
					}
					if !hasProfile {
						continue
					}

					cell.Total++
					evidenceList, _ := s.EvidenceRepo.GetByEmployee(emp.ID)
					
					reqs, _ := s.ProfileRepo.GetRequirements(profile.ID)
					allCompliant := true
					anyExpiring := false

					for _, req := range reqs {
						status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)
						if status.Status == "Missing" || status.Status == "Expired" {
							allCompliant = false
						}
						if status.Status == "Expiring" {
							anyExpiring = true
						}
					}

					if allCompliant {
						if anyExpiring {
							cell.Expiring++
						} else {
							cell.Valid++
						}
					} else {
						cell.Missing++
					}
				}

				if cell.Total > 0 {
					cell.Percent = float64(cell.Valid+cell.Expiring) / float64(cell.Total) * 100
				}
				cells = append(cells, cell)
			}
		}
	} else {
		// Course coverage mode
		requirements, _ := s.ProfileRepo.GetAllRequirements()
		reqMap := make(map[string]domain.Requirement)
		for _, req := range requirements {
			reqMap[req.Name] = req
		}

		for groupValue, groupEmployees := range groupMap {
			for reqName, req := range reqMap {
				cell := domain.GridCell{
					GroupValue: groupValue,
					ItemID:     req.ID,
					ItemName:   reqName,
				}

				for _, emp := range groupEmployees {
					profileIDs, _ := s.EmployeeRepo.GetAssignments(emp.ID)
					
					// Check if employee needs this requirement
					needsRequirement := false
					for _, pid := range profileIDs {
						reqs, _ := s.ProfileRepo.GetRequirements(pid)
						for _, r := range reqs {
							if r.Name == reqName {
								needsRequirement = true
								break
							}
						}
						if needsRequirement {
							break
						}
					}
					if !needsRequirement {
						continue
					}

					cell.Total++
					evidenceList, _ := s.EvidenceRepo.GetByEmployee(emp.ID)
					status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)

					switch status.Status {
					case "Valid":
						cell.Valid++
					case "Expiring":
						cell.Expiring++
					case "Expired":
						cell.Expired++
					case "Missing":
						cell.Missing++
					}
				}

				if cell.Total > 0 {
					cell.Percent = float64(cell.Valid+cell.Expiring) / float64(cell.Total) * 100
				}
				cells = append(cells, cell)
			}
		}
	}

	return cells, nil
}

// GetDrilldown returns employee drilldown for a grid cell
func (s *DashboardService) GetDrilldown(mode, groupBy, groupValue, itemID, department, site string) ([]domain.DrilldownEmployee, error) {
	employees, err := s.EmployeeRepo.GetByFilter(department, site)
	if err != nil {
		return nil, err
	}

	results := make([]domain.DrilldownEmployee, 0)

	for _, emp := range employees {
		// Filter by group
		var empGroup string
		switch groupBy {
		case "department":
			empGroup = emp.Department
		case "site":
			empGroup = emp.Site
		case "costCenter":
			empGroup = emp.CostCenter
		}
		if empGroup != groupValue {
			continue
		}

		profileIDs, _ := s.EmployeeRepo.GetAssignments(emp.ID)
		evidenceList, _ := s.EvidenceRepo.GetByEmployee(emp.ID)

		if mode == "profile" {
			// Check if employee has this profile
			hasProfile := false
			for _, pid := range profileIDs {
				if pid == itemID {
					hasProfile = true
					break
				}
			}
			if !hasProfile {
				continue
			}

			reqs, _ := s.ProfileRepo.GetRequirements(itemID)
			allCompliant := true
			anyExpiring := false

			for _, req := range reqs {
				status := s.ComputationService.ComputeRequirementStatus(req, evidenceList)
				if status.Status == "Missing" || status.Status == "Expired" {
					allCompliant = false
				}
				if status.Status == "Expiring" {
					anyExpiring = true
				}
			}

			statusStr := "Eligible"
			if !allCompliant {
				statusStr = "Not Eligible"
			} else if anyExpiring {
				statusStr = "At Risk"
			}

			results = append(results, domain.DrilldownEmployee{
				ID:         emp.ID,
				Name:       emp.Name,
				Department: emp.Department,
				Site:       emp.Site,
				Status:     statusStr,
			})
		} else {
			// Course mode - find the requirement
			requirements, _ := s.ProfileRepo.GetAllRequirements()
			var targetReq *domain.Requirement
			for _, req := range requirements {
				if req.ID == itemID {
					targetReq = &req
					break
				}
			}
			if targetReq == nil {
				continue
			}

			// Check if employee needs this requirement
			needsRequirement := false
			for _, pid := range profileIDs {
				reqs, _ := s.ProfileRepo.GetRequirements(pid)
				for _, r := range reqs {
					if r.ID == itemID {
						needsRequirement = true
						break
					}
				}
				if needsRequirement {
					break
				}
			}
			if !needsRequirement {
				continue
			}

			status := s.ComputationService.ComputeRequirementStatus(*targetReq, evidenceList)
			expiryStr := ""
			if status.ExpiryDate != nil {
				expiryStr = status.ExpiryDate.Format("2006-01-02")
			}

			results = append(results, domain.DrilldownEmployee{
				ID:         emp.ID,
				Name:       emp.Name,
				Department: emp.Department,
				Site:       emp.Site,
				Status:     status.Status,
				ExpiryDate: expiryStr,
			})
		}
	}

	return results, nil
}
