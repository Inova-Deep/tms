package domain

type EmployeeWithStatus struct {
	Employee
	ProfileStatus map[string]string       `json:"profile_status"` // ProfileID -> Eligible/AtRisk/NotEligible
	Requirements  []RequirementWithStatus `json:"requirements"`
}

type RequirementWithStatus struct {
	Requirement
	Status     string  `json:"status"`
	ExpiryDate *string `json:"expiry_date"`
}
