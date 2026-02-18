-- Employees table
CREATE TABLE IF NOT EXISTS employees (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    department TEXT,
    site TEXT,
    worker_type TEXT,
    cost_center TEXT,
    employment_status TEXT
);

-- Profiles table
CREATE TABLE IF NOT EXISTS profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT
);

-- Requirements table
CREATE TABLE IF NOT EXISTS requirements (
    id TEXT PRIMARY KEY,
    profile_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    validity_months INTEGER DEFAULT 0,
    FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

-- Profile assignments
CREATE TABLE IF NOT EXISTS profile_assignments (
    employee_id TEXT NOT NULL,
    profile_id TEXT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (employee_id, profile_id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

-- Profile exclusions
CREATE TABLE IF NOT EXISTS profile_exclusions (
    employee_id TEXT NOT NULL,
    profile_id TEXT NOT NULL,
    reason TEXT NOT NULL,
    expiry_date TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (employee_id, profile_id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    FOREIGN KEY (profile_id) REFERENCES profiles (id)
);

-- Evidence table
CREATE TABLE IF NOT EXISTS evidence (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL,
    requirement_name TEXT NOT NULL,
    evidence_type TEXT NOT NULL,
    completion_date DATE NOT NULL,
    expiry_date DATE,
    metadata TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees (id)
);

-- Training events table
CREATE TABLE IF NOT EXISTS training_events (
    id TEXT PRIMARY KEY,
    course_name TEXT NOT NULL,
    date DATE NOT NULL,
    location TEXT,
    instructor_id TEXT,
    status TEXT DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Event attendance table
CREATE TABLE IF NOT EXISTS event_attendance (
    event_id TEXT NOT NULL,
    employee_id TEXT NOT NULL,
    attended_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, employee_id),
    FOREIGN KEY (event_id) REFERENCES training_events (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id)
);

-- Certifications table
CREATE TABLE IF NOT EXISTS certifications (
    id TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL,
    requirement_name TEXT NOT NULL,
    issuer TEXT NOT NULL,
    issue_date DATE NOT NULL,
    expiry_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees (id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_evidence_employee ON evidence(employee_id);
CREATE INDEX IF NOT EXISTS idx_evidence_requirement ON evidence(requirement_name);
CREATE INDEX IF NOT EXISTS idx_requirements_profile ON requirements(profile_id);
CREATE INDEX IF NOT EXISTS idx_assignments_employee ON profile_assignments(employee_id);
CREATE INDEX IF NOT EXISTS idx_assignments_profile ON profile_assignments(profile_id);
CREATE INDEX IF NOT EXISTS idx_certifications_employee ON certifications(employee_id);
CREATE INDEX IF NOT EXISTS idx_event_attendance_event ON event_attendance(event_id);

-- Migration tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO schema_migrations (version) VALUES ('001');
