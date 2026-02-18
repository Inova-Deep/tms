package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const seedVersion = 2

func Run(db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentVersion int
	err = tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM seed_versions").Scan(&currentVersion)
	if err != nil {
		_, createErr := tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS seed_versions (version INTEGER PRIMARY KEY)")
		if createErr != nil {
			return createErr
		}
		currentVersion = 0
	}

	if currentVersion >= seedVersion {
		log.Println("Database already seeded.")
		return nil
	}

	log.Println("Seeding database...")

	employees := []struct {
		ID, Name, Dept, Site, WorkerType, CostCenter, Status string
	}{
		{"E1001", "Oliver Bennett", "Manufacturing", "Bristol Central Plant", "Employee", "CC100", "Active"},
		{"E1002", "Amelia Clarke", "Manufacturing", "Bristol Central Plant", "Employee", "CC100", "Active"},
		{"E1003", "Harry Patel", "Maintenance", "Bristol Central Plant", "Employee", "CC110", "Active"},
		{"E1004", "Isla Hughes", "Quality", "Bristol Central Plant", "Employee", "CC120", "Active"},
		{"E1005", "Jack Thompson", "Manufacturing", "Avonmouth Logistics & Workshop", "Employee", "CC200", "Active"},
		{"E1006", "Emily Walker", "Maintenance", "Avonmouth Logistics & Workshop", "Employee", "CC210", "Active"},
		{"E1007", "George Evans", "Quality", "Avonmouth Logistics & Workshop", "Employee", "CC220", "Active"},
		{"E1008", "Sophia Khan", "HR", "Bristol Central Plant", "Employee", "CC130", "Active"},
		{"E1009", "Noah Wilson", "Manufacturing", "Bristol Central Plant", "Contractor", "CC100", "Active"},
		{"E1010", "Ava Morgan", "Manufacturing", "Avonmouth Logistics & Workshop", "Employee", "CC200", "Active"},
		{"E1011", "Leo Robinson", "Maintenance", "Bristol Central Plant", "Employee", "CC110", "Active"},
		{"E1012", "Mia Carter", "Quality", "Bristol Central Plant", "Employee", "CC120", "Active"},
		{"E1013", "Ethan James", "Manufacturing", "Avonmouth Logistics & Workshop", "Employee", "CC200", "Active"},
		{"E1014", "Grace Murphy", "HR", "Avonmouth Logistics & Workshop", "Employee", "CC230", "Active"},
		{"E1015", "Freddie Brown", "Manufacturing", "Bristol Central Plant", "Employee", "CC100", "Active"},
		{"E1016", "Lily Ward", "Maintenance", "Avonmouth Logistics & Workshop", "Contractor", "CC210", "Active"},
		{"E1017", "Oscar Lewis", "Quality", "Avonmouth Logistics & Workshop", "Employee", "CC220", "Active"},
		{"E1018", "Chloe Smith", "Manufacturing", "Bristol Central Plant", "Employee", "CC100", "Active"},
		{"E1019", "Charlie Green", "Maintenance", "Bristol Central Plant", "Employee", "CC110", "Active"},
		{"E1020", "Ella Davies", "Manufacturing", "Avonmouth Logistics & Workshop", "Employee", "CC200", "Active"},
	}

	for _, e := range employees {
		_, err := tx.ExecContext(ctx, `INSERT INTO employees VALUES (?, ?, ?, ?, ?, ?, ?)`,
			e.ID, e.Name, e.Dept, e.Site, e.WorkerType, e.CostCenter, e.Status)
		if err != nil {
			return fmt.Errorf("failed to insert employee %s: %w", e.ID, err)
		}
	}

	courses := []struct {
		ID, Name, Type string
		Validity       int
	}{
		{"C001", "H&S Induction", "course", 24},
		{"C002", "Fire Safety Awareness", "course", 12},
		{"C003", "First Aid Awareness", "course", 12},
		{"C004", "Manual Handling", "course", 12},
		{"C005", "COSHH Awareness", "course", 12},
		{"C006", "Welding Safety & PPE", "course", 12},
		{"C007", "Hot Work Permit Awareness", "course", 12},
		{"C008", "Respiratory Fit Test", "course", 24},
		{"C009", "LOTO Awareness", "course", 12},
		{"C010", "Warehouse Pedestrian Safety", "course", 12},
		{"C011", "GDPR & Data Protection", "course", 24},
		{"C012", "Right to Work Checks (UK)", "course", 24},
		{"C013", "Equality, Diversity & Inclusion", "course", 24},
		{"C014", "DSE (Display Screen Equipment)", "course", 24},
		{"C015", "Welding Qualification (ISO 9606-1)", "certification", 24},
		{"C016", "Forklift Theory (RTITB)", "course", 36},
		{"C017", "Forklift Practical Assessment", "course", 36},
		{"C018", "Forklift Licence (RTITB/ITSSAR)", "certification", 36},
	}

	for _, c := range courses {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO courses (id, name, type, validity_months)
			VALUES (?, ?, ?, ?)
		`, c.ID, c.Name, c.Type, c.Validity)
		if err != nil {
			return fmt.Errorf("failed to insert course %s: %w", c.ID, err)
		}
	}

	profiles := []struct {
		ID, Name, Description string
	}{
		{"P1", "Safety Base (All Employees)", "Mandatory safety training for all employees"},
		{"P2", "Welder — Basic", "Basic welding competency profile"},
		{"P3", "Forklift Operator", "Forklift operation competency"},
		{"P4", "HR Management (Compliance)", "HR compliance and regulatory training"},
	}
	for _, p := range profiles {
		_, err := tx.ExecContext(ctx, `INSERT INTO profiles (id, name, description) VALUES (?, ?, ?)`, p.ID, p.Name, p.Description)
		if err != nil {
			return fmt.Errorf("failed to insert profile %s: %w", p.ID, err)
		}
	}

	requirements := []struct {
		ID, ProfileID, CourseID, Name, Type string
		Validity                            int
	}{
		{"R101", "P1", "C001", "H&S Induction", "course", 24},
		{"R102", "P1", "C002", "Fire Safety Awareness", "course", 12},
		{"R103", "P1", "C003", "First Aid Awareness", "course", 12},
		{"R104", "P1", "C004", "Manual Handling", "course", 12},
		{"R105", "P1", "C005", "COSHH Awareness", "course", 12},
		{"R201", "P2", "C006", "Welding Safety & PPE", "course", 12},
		{"R202", "P2", "C007", "Hot Work Permit Awareness", "course", 12},
		{"R203", "P2", "C008", "Respiratory Fit Test", "course", 24},
		{"R204", "P2", "C015", "Welding Qualification (ISO 9606-1)", "certification", 24},
		{"R205", "P2", "C009", "LOTO Awareness", "course", 12},
		{"R301", "P3", "C016", "Forklift Theory (RTITB)", "course", 36},
		{"R302", "P3", "C017", "Forklift Practical Assessment", "course", 36},
		{"R303", "P3", "C018", "Forklift Licence (RTITB/ITSSAR)", "certification", 36},
		{"R304", "P3", "C010", "Warehouse Pedestrian Safety", "course", 12},
		{"R401", "P4", "C011", "GDPR & Data Protection", "course", 24},
		{"R402", "P4", "C012", "Right to Work Checks (UK)", "course", 24},
		{"R403", "P4", "C013", "Equality, Diversity & Inclusion", "course", 24},
		{"R404", "P4", "C014", "DSE (Display Screen Equipment)", "course", 24},
	}
	for _, r := range requirements {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO requirements (id, profile_id, name, type, validity_months, course_id)
			VALUES (?, ?, ?, ?, ?, ?)
		`, r.ID, r.ProfileID, r.Name, r.Type, r.Validity, r.CourseID)
		if err != nil {
			return fmt.Errorf("failed to insert requirement %s: %w", r.ID, err)
		}
	}

	for i := 1001; i <= 1020; i++ {
		empID := fmt.Sprintf("E%d", i)
		_, err := tx.ExecContext(ctx, `INSERT INTO profile_assignments (employee_id, profile_id) VALUES (?, ?)`, empID, "P1")
		if err != nil {
			return fmt.Errorf("failed to assign P1 to %s: %w", empID, err)
		}
	}

	welderEmployees := []string{"E1001", "E1005", "E1013", "E1015", "E1018"}
	for _, empID := range welderEmployees {
		_, err := tx.ExecContext(ctx, `INSERT INTO profile_assignments (employee_id, profile_id) VALUES (?, ?)`, empID, "P2")
		if err != nil {
			return fmt.Errorf("failed to assign P2 to %s: %w", empID, err)
		}
	}

	forkliftEmployees := []string{"E1002", "E1005", "E1010", "E1011", "E1016", "E1020"}
	for _, empID := range forkliftEmployees {
		_, err := tx.ExecContext(ctx, `INSERT INTO profile_assignments (employee_id, profile_id) VALUES (?, ?)`, empID, "P3")
		if err != nil {
			return fmt.Errorf("failed to assign P3 to %s: %w", empID, err)
		}
	}

	hrEmployees := []string{"E1008", "E1014"}
	for _, empID := range hrEmployees {
		_, err := tx.ExecContext(ctx, `INSERT INTO profile_assignments (employee_id, profile_id) VALUES (?, ?)`, empID, "P4")
		if err != nil {
			return fmt.Errorf("failed to assign P4 to %s: %w", empID, err)
		}
	}

	now := time.Now()
	expiringDate := now.AddDate(0, 0, 45).Format("2006-01-02")
	expiredDate := now.AddDate(-1, 0, 0).Format("2006-01-02")

	createEvidence := func(id, empID, courseID, evidenceType, completionDate, expiryDate string) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO evidence (id, employee_id, course_id, requirement_name, evidence_type, completion_date, expiry_date, metadata)
			VALUES (?, ?, ?, '', ?, ?, ?, ?)
		`, id, empID, courseID, evidenceType, completionDate, expiryDate, "")
		return err
	}

	evCounter := 1

	safetyCourseIDs := []string{"C001", "C002", "C003", "C004", "C005"}

	for i := 1001; i <= 1020; i++ {
		empID := fmt.Sprintf("E%d", i)
		for j, courseID := range safetyCourseIDs {
			completionDate := now.AddDate(-1, -j, 0).Format("2006-01-02")

			if evCounter%20 < 16 {
				expiry := now.AddDate(1, j%3, 0).Format("2006-01-02")
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, "attendance", completionDate, expiry)
			} else if evCounter%20 < 18 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, "attendance", completionDate, expiringDate)
			} else if evCounter%20 < 19 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, "attendance", completionDate, expiredDate)
			}
			evCounter++
		}
	}

	welderCourseIDs := []string{"C006", "C007", "C008", "C015", "C009"}
	for _, empID := range welderEmployees {
		for j, courseID := range welderCourseIDs {
			completionDate := now.AddDate(-1, -j, 0).Format("2006-01-02")
			evidenceType := "attendance"
			if courseID == "C015" {
				evidenceType = "certification"
			}

			if evCounter%10 < 7 {
				expiry := now.AddDate(1, j%2, 0).Format("2006-01-02")
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiry)
			} else if evCounter%10 < 8 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiringDate)
			} else if evCounter%10 < 9 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiredDate)
			}
			evCounter++
		}
	}

	forkliftCourseIDs := []string{"C016", "C017", "C018", "C010"}
	for _, empID := range forkliftEmployees {
		for j, courseID := range forkliftCourseIDs {
			completionDate := now.AddDate(-2, -j, 0).Format("2006-01-02")
			evidenceType := "attendance"
			if courseID == "C018" {
				evidenceType = "certification"
			}

			if evCounter%10 < 7 {
				expiry := now.AddDate(2, j%3, 0).Format("2006-01-02")
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiry)
			} else if evCounter%10 < 8 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiringDate)
			} else if evCounter%10 < 9 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, evidenceType, completionDate, expiredDate)
			}
			evCounter++
		}
	}

	hrCourseIDs := []string{"C011", "C012", "C013", "C014"}
	for _, empID := range hrEmployees {
		for j, courseID := range hrCourseIDs {
			completionDate := now.AddDate(-1, -j, 0).Format("2006-01-02")

			if evCounter%8 < 6 {
				expiry := now.AddDate(1, j%2, 0).Format("2006-01-02")
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, "attendance", completionDate, expiry)
			} else if evCounter%8 < 7 {
				createEvidence(fmt.Sprintf("EV%d", evCounter), empID, courseID, "attendance", completionDate, expiringDate)
			}
			evCounter++
		}
	}

	events := []struct {
		ID, CourseID, Date, Time, Location, InstructorID, Status string
		Duration                                                 int
	}{
		{"TE001", "C004", now.AddDate(0, -1, 0).Format("2006-01-02"), "09:00", "Bristol Central Plant", "E1003", "completed", 120},
		{"TE002", "C002", now.AddDate(0, 0, -15).Format("2006-01-02"), "14:00", "Avonmouth Logistics & Workshop", "E1006", "completed", 90},
		{"TE003", "C017", now.AddDate(0, 0, 7).Format("2006-01-02"), "08:30", "Avonmouth Logistics & Workshop", "E1006", "scheduled", 180},
	}

	for _, e := range events {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO training_events (id, course_id, course_name, date, time, duration, location, instructor_id, status)
			VALUES (?, ?, '', ?, ?, ?, ?, ?, ?)
		`, e.ID, e.CourseID, e.Date, e.Time, e.Duration, e.Location, e.InstructorID, e.Status)
		if err != nil {
			return fmt.Errorf("failed to insert event %s: %w", e.ID, err)
		}
	}

	te001Attendees := []string{"E1001", "E1002", "E1003", "E1004", "E1009"}
	for _, empID := range te001Attendees {
		tx.ExecContext(ctx, `INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)`, "TE001", empID)
	}

	te002Attendees := []string{"E1005", "E1006", "E1010", "E1013", "E1014", "E1016", "E1020"}
	for _, empID := range te002Attendees {
		tx.ExecContext(ctx, `INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)`, "TE002", empID)
	}

	exclusions := []struct {
		ID, EmployeeID, ProfileID, Reason, ExpiryDate string
	}{
		{"EX-E1009-P3", "E1009", "P3", "Medical exemption - back injury prevents forklift operation", now.AddDate(0, 6, 0).Format("2006-01-02")},
		{"EX-E1016-P2", "E1016", "P2", "Contractor - not authorized for welding activities", now.AddDate(1, 0, 0).Format("2006-01-02")},
	}

	for _, ex := range exclusions {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO profile_exclusions (id, employee_id, profile_id, reason, expiry_date, created_by)
			VALUES (?, ?, ?, ?, ?, 'admin')
		`, ex.ID, ex.EmployeeID, ex.ProfileID, ex.Reason, ex.ExpiryDate)
		if err != nil {
			return fmt.Errorf("failed to insert exclusion %s: %w", ex.ID, err)
		}
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO seed_versions (version) VALUES (?)", seedVersion)
	if err != nil {
		return fmt.Errorf("failed to update seed version: %w", err)
	}

	log.Printf("Seeded %d employees, %d courses, %d profiles, %d requirements, ~%d evidence records, %d events, %d exclusions",
		len(employees), len(courses), len(profiles), len(requirements), evCounter-1, len(events), len(exclusions))

	return tx.Commit()
}
