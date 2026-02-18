package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type EventRepository struct {
	DB *sql.DB
}

func (r *EventRepository) GetAll() ([]TrainingEvent, error) {
	rows, err := r.DB.Query("SELECT id, course_id, date, COALESCE(time, ''), COALESCE(duration, 0), location, instructor_id, status FROM training_events ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]TrainingEvent, 0)
	for rows.Next() {
		var e TrainingEvent
		if err := rows.Scan(&e.ID, &e.CourseID, &e.Date, &e.Time, &e.Duration, &e.Location, &e.InstructorID, &e.Status); err != nil {
			return nil, err
		}
		attendees, _ := r.getAttendees(e.ID)
		e.Attendees = attendees
		events = append(events, e)
	}
	return events, nil
}

func (r *EventRepository) GetByID(id string) (*TrainingEvent, error) {
	var e TrainingEvent
	err := r.DB.QueryRow("SELECT id, course_id, date, COALESCE(time, ''), COALESCE(duration, 0), location, instructor_id, status FROM training_events WHERE id = ?", id).
		Scan(&e.ID, &e.CourseID, &e.Date, &e.Time, &e.Duration, &e.Location, &e.InstructorID, &e.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	attendees, _ := r.getAttendees(e.ID)
	e.Attendees = attendees
	return &e, nil
}

func (r *EventRepository) getAttendees(eventID string) ([]string, error) {
	rows, err := r.DB.Query("SELECT employee_id FROM event_attendance WHERE event_id = ?", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendees := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		attendees = append(attendees, id)
	}
	return attendees, nil
}

func (r *EventRepository) Create(req CreateEventRequest) (*TrainingEvent, error) {
	id := fmt.Sprintf("EV-%d", time.Now().UnixNano())

	_, err := r.DB.Exec(`
		INSERT INTO training_events (id, course_id, date, time, duration, location, instructor_id, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'scheduled')
	`, id, req.CourseID, req.Date, req.Time, req.Duration, req.Location, req.InstructorID)
	if err != nil {
		return nil, err
	}

	for _, attendeeID := range req.Attendees {
		r.DB.Exec("INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)", id, attendeeID)
	}

	return r.GetByID(id)
}

func (r *EventRepository) Update(id string, req UpdateEventRequest) (*TrainingEvent, error) {
	query := "UPDATE training_events SET "
	args := make([]interface{}, 0)
	setClauses := make([]string, 0)

	if req.CourseID != nil {
		setClauses = append(setClauses, "course_id = ?")
		args = append(args, *req.CourseID)
	}
	if req.Date != nil {
		setClauses = append(setClauses, "date = ?")
		args = append(args, *req.Date)
	}
	if req.Time != nil {
		setClauses = append(setClauses, "time = ?")
		args = append(args, *req.Time)
	}
	if req.Duration != nil {
		setClauses = append(setClauses, "duration = ?")
		args = append(args, *req.Duration)
	}
	if req.Location != nil {
		setClauses = append(setClauses, "location = ?")
		args = append(args, *req.Location)
	}
	if req.InstructorID != nil {
		setClauses = append(setClauses, "instructor_id = ?")
		args = append(args, *req.InstructorID)
	}
	if req.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *req.Status)
	}

	if len(setClauses) == 0 && len(req.Attendees) == 0 {
		return r.GetByID(id)
	}

	if len(setClauses) > 0 {
		query += setClauses[0]
		for i := 1; i < len(setClauses); i++ {
			query += ", " + setClauses[i]
		}
		query += " WHERE id = ?"
		args = append(args, id)

		_, err := r.DB.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	if req.Attendees != nil {
		r.DB.Exec("DELETE FROM event_attendance WHERE event_id = ?", id)
		for _, attendeeID := range req.Attendees {
			r.DB.Exec("INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)", id, attendeeID)
		}
	}

	return r.GetByID(id)
}

func (r *EventRepository) ConfirmAttendance(eventID string, employeeIDs []string, evidenceRepo *EvidenceRepository) error {
	event, err := r.GetByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	var validityMonths int
	err = r.DB.QueryRow("SELECT validity_months FROM courses WHERE id = ?", event.CourseID).Scan(&validityMonths)
	if err != nil {
		validityMonths = 12
	}

	r.DB.Exec("DELETE FROM event_attendance WHERE event_id = ?", eventID)

	for _, empID := range employeeIDs {
		r.DB.Exec("INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)", eventID, empID)
		evidenceRepo.CreateFromAttendance(empID, event.CourseID, event.Date, validityMonths)
	}

	_, err = r.DB.Exec("UPDATE training_events SET status = 'completed' WHERE id = ?", eventID)
	return err
}

func (r *EventRepository) EnhancedConfirmAttendance(eventID string, req ConfirmAttendanceRequest, evidenceRepo *EvidenceRepository) error {
	event, err := r.GetByID(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event not found")
	}

	var validityMonths int
	err = r.DB.QueryRow("SELECT validity_months FROM courses WHERE id = ?", event.CourseID).Scan(&validityMonths)
	if err != nil {
		validityMonths = 12
	}

	r.DB.Exec("DELETE FROM event_attendance WHERE event_id = ?", eventID)

	for _, a := range req.Attendees {
		if a.Attended {
			r.DB.Exec("INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)", eventID, a.EmployeeID)
			evidenceRepo.CreateFromAttendance(a.EmployeeID, event.CourseID, event.Date, validityMonths)
		}
	}

	for _, w := range req.WalkIns {
		r.DB.Exec("INSERT INTO event_attendance (event_id, employee_id) VALUES (?, ?)", eventID, w.EmployeeID)
		evidenceRepo.CreateFromAttendance(w.EmployeeID, event.CourseID, event.Date, validityMonths)
	}

	_, err = r.DB.Exec("UPDATE training_events SET status = 'completed' WHERE id = ?", eventID)
	return err
}
