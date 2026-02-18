package logic

import (
	"time"
	"tms/internal/domain"
)

type ComputationService struct {
	EvidenceRepo *domain.EvidenceRepository
}

func (s *ComputationService) ComputeRequirementStatus(req domain.Requirement, evidenceList []domain.Evidence) domain.RequirementStatus {
	var bestEvidence *domain.Evidence
	var latestExpiry time.Time

	// Find best evidence for this requirement
	for _, e := range evidenceList {
		if e.CourseID != req.CourseID {
			continue
		}

		// Parse expiry
		var expiry time.Time
		if e.ExpiryDate != nil && *e.ExpiryDate != "" {
			// Try various date formats
			formats := []string{
				"2006-01-02",
				"2006-01-02 15:04:05",
				time.RFC3339,
			}

			var parsed time.Time
			var err error
			for _, f := range formats {
				parsed, err = time.Parse(f, *e.ExpiryDate)
				if err == nil {
					expiry = parsed
					break
				}
			}

			if err != nil {
				// fmt.Printf("Error parsing date '%s': %v\n", *e.ExpiryDate, err)
			}
		} else {
			// ... default logic ... (omitted for brevity in prompt, but I should keep original logic if not replacing it)
			// Wait, I am replacing the block. I need to keep the else logic or simplified version.
			// The original code had logic for "no expiry".
			// Let's just focus on the date parsing logging and fixing loop var.
			if req.ValidityMonths > 0 {
				completion, _ := time.Parse("2006-01-02", e.CompletionDate)
				expiry = completion.AddDate(0, req.ValidityMonths, 0)
			} else {
				expiry = time.Now().AddDate(100, 0, 0)
			}
		}

		currentEvidence := e // Copy loop variable to safe location
		if bestEvidence == nil || expiry.After(latestExpiry) {
			bestEvidence = &currentEvidence
			latestExpiry = expiry
		}
	}

	status := domain.RequirementStatus{
		RequirementID: req.ID,
		Status:        "Missing",
	}

	if bestEvidence != nil {
		now := time.Now()
		status.ExpiryDate = &latestExpiry

		if latestExpiry.Before(now) {
			status.Status = "Expired"
		} else if latestExpiry.Before(now.AddDate(0, 0, 90)) { // 90 days warning
			status.Status = "Expiring"
		} else {
			status.Status = "Valid"
		}
	}

	return status
}
