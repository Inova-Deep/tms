# TMS API Documentation

**Version:** 2.0  
**Base URL:** `http://localhost:8080/api`  
**Content-Type:** `application/json`

---

## Table of Contents

1. [Authentication](#1-authentication)
2. [Employees](#2-employees)
3. [Profiles](#3-profiles)
4. [Courses](#4-courses)
5. [Events](#5-events)
6. [Certifications](#6-certifications)
7. [Dashboard](#7-dashboard)
8. [Error Responses](#8-error-responses)
9. [Data Types](#9-data-types)

---

## Authentication

All protected endpoints require an `Authorization` header with a Bearer token.

```
Authorization: Bearer <token>
```

### Roles

| Role | Permissions |
|------|-------------|
| `admin` | Full read/write access to all resources |
| `employee` | Read-only access to own data only |

---

## 1. Authentication

### 1.1 Login

Authenticate and receive a token.

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "role": "admin"
}
```

For employee login:
```json
{
  "role": "employee",
  "employeeId": "E1005"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `role` | string | Yes | Either `"admin"` or `"employee"` |
| `employeeId` | string | Conditional | Required when `role` is `"employee"` |

**Response:** `200 OK`
```json
{
  "token": "eyJyb2xlIjoiYWRtaW4iLCJlbXBsb3llZUlkIjoiIn0=",
  "role": "admin",
  "employeeId": ""
}
```

**Error Responses:**

| Status | Description |
|--------|-------------|
| 400 | Invalid request body or role |
| 404 | Employee not found (for employee role) |

**Example:**
```bash
# Admin login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'

# Employee login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"role": "employee", "employeeId": "E1005"}'
```

---

## 2. Employees

### 2.1 Get All Employees

Retrieve all employees (admin) or self (employee).

**Endpoint:** `GET /api/employees`

**Auth:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "E1001",
    "name": "Oliver Bennett",
    "department": "Manufacturing",
    "site": "Bristol Central Plant",
    "worker_type": "Employee",
    "cost_center": "CC100",
    "employment_status": "Active"
  },
  ...
]
```

**Example:**
```bash
curl http://localhost:8080/api/employees \
  -H "Authorization: Bearer <token>"
```

---

### 2.2 Search Employees

Search employees by ID, name, or department.

**Endpoint:** `GET /api/employees/search?q={query}`

**Auth:** Required (Admin only)

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query (matches ID, name, or department) |

**Response:** `200 OK`
```json
[
  {
    "id": "E1005",
    "name": "Jack Thompson",
    "department": "Manufacturing",
    "site": "Avonmouth Logistics & Workshop",
    "worker_type": "Employee",
    "cost_center": "CC200",
    "employment_status": "Active"
  }
]
```

**Example:**
```bash
curl "http://localhost:8080/api/employees/search?q=jack" \
  -H "Authorization: Bearer <admin-token>"
```

---

### 2.3 Get Employee by ID

Retrieve detailed employee information with profile status and requirements.

**Endpoint:** `GET /api/employees/{id}`

**Auth:** Required (Admin: any employee, Employee: self only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Employee ID (e.g., `E1005`) |

**Response:** `200 OK`
```json
{
  "id": "E1005",
  "name": "Jack Thompson",
  "department": "Manufacturing",
  "site": "Avonmouth Logistics & Workshop",
  "worker_type": "Employee",
  "cost_center": "CC200",
  "employment_status": "Active",
  "profile_status": {
    "P1": "Eligible",
    "P3": "At Risk"
  },
  "requirements": [
    {
      "id": "R101",
      "profile_id": "P1",
      "name": "H&S Induction",
      "type": "course",
      "validity_months": 24,
      "status": "Valid",
      "expiry_date": "2026-02-18"
    },
    {
      "id": "R102",
      "profile_id": "P1",
      "name": "Fire Safety Awareness",
      "type": "course",
      "validity_months": 12,
      "status": "Expiring",
      "expiry_date": "2026-04-03"
    },
    {
      "id": "R303",
      "profile_id": "P3",
      "name": "Forklift Licence (RTITB/ITSSAR)",
      "type": "certification",
      "validity_months": 36,
      "status": "Missing",
      "expiry_date": null
    }
  ]
}
```

**Profile Status Values:**

| Status | Description |
|--------|-------------|
| `Eligible` | All requirements are Valid |
| `At Risk` | All requirements compliant, but at least one Expiring |
| `Not Eligible` | At least one Missing or Expired requirement |

**Requirement Status Values:**

| Status | Description |
|--------|-------------|
| `Valid` | Evidence exists and expiry > 90 days |
| `Expiring` | Evidence exists and expiry within 90 days |
| `Expired` | Evidence exists and expiry date has passed |
| `Missing` | No evidence record exists |

**Error Responses:**

| Status | Description |
|--------|-------------|
| 403 | Access denied (employee trying to view another) |
| 404 | Employee not found |

**Example:**
```bash
curl http://localhost:8080/api/employees/E1005 \
  -H "Authorization: Bearer <token>"
```

---

### 2.4 Get Employee Records (Transcript)

Retrieve all evidence records for an employee.

**Endpoint:** `GET /api/employees/{id}/records`

**Auth:** Required (Admin: any employee, Employee: self only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Employee ID |

**Response:** `200 OK`
```json
[
  {
    "id": "EV1",
    "employeeId": "E1005",
    "requirementName": "H&S Induction",
    "evidenceType": "attendance",
    "completionDate": "2024-02-18",
    "expiryDate": "2026-02-18",
    "metadata": ""
  },
  {
    "id": "EV2",
    "employeeId": "E1005",
    "requirementName": "Forklift Licence (RTITB/ITSSAR)",
    "evidenceType": "certification",
    "completionDate": "2024-06-15",
    "expiryDate": "2027-06-15",
    "metadata": "{\"issuer\":\"RTITB\"}"
  }
]
```

**Evidence Types:**

| Type | Description |
|------|-------------|
| `attendance` | Generated from training event attendance |
| `certification` | Uploaded third-party certification |

**Example:**
```bash
curl http://localhost:8080/api/employees/E1005/records \
  -H "Authorization: Bearer <token>"
```

---

## 3. Profiles

### 3.1 Get All Profiles

Retrieve all competency profiles with member and exclusion counts.

**Endpoint:** `GET /api/profiles`

**Auth:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "P1",
    "name": "Safety Base (All Employees)",
    "description": "Mandatory safety training for all employees",
    "memberCount": 20,
    "exclusionCount": 0
  },
  {
    "id": "P2",
    "name": "Welder — Basic",
    "description": "Basic welding competency profile",
    "memberCount": 5,
    "exclusionCount": 1
  },
  {
    "id": "P3",
    "name": "Forklift Operator",
    "description": "Forklift operation competency",
    "memberCount": 6,
    "exclusionCount": 1
  },
  {
    "id": "P4",
    "name": "HR Management (Compliance)",
    "description": "HR compliance and regulatory training",
    "memberCount": 2,
    "exclusionCount": 0
  }
]
```

**Example:**
```bash
curl http://localhost:8080/api/profiles \
  -H "Authorization: Bearer <token>"
```

---

### 3.2 Get Profile by ID

Retrieve a single profile.

**Endpoint:** `GET /api/profiles/{id}`

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID (e.g., `P1`) |

**Response:** `200 OK`
```json
{
  "id": "P1",
  "name": "Safety Base (All Employees)",
  "description": "Mandatory safety training for all employees"
}
```

**Example:**
```bash
curl http://localhost:8080/api/profiles/P1 \
  -H "Authorization: Bearer <token>"
```

---

### 3.3 Get Profile Requirements

Retrieve all requirements for a profile.

**Endpoint:** `GET /api/profiles/{id}/requirements`

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |

**Response:** `200 OK`
```json
[
  {
    "id": "R101",
    "profile_id": "P1",
    "name": "H&S Induction",
    "type": "course",
    "validity_months": 24
  },
  {
    "id": "R102",
    "profile_id": "P1",
    "name": "Fire Safety Awareness",
    "type": "course",
    "validity_months": 12
  },
  {
    "id": "R103",
    "profile_id": "P1",
    "name": "First Aid Awareness",
    "type": "course",
    "validity_months": 12
  },
  {
    "id": "R104",
    "profile_id": "P1",
    "name": "Manual Handling",
    "type": "course",
    "validity_months": 12
  },
  {
    "id": "R105",
    "profile_id": "P1",
    "name": "COSHH Awareness",
    "type": "course",
    "validity_months": 12
  }
]
```

**Requirement Types:**

| Type | Description |
|------|-------------|
| `course` | Internal training course (attendance-based) |
| `certification` | External certification (upload-based) |

**Example:**
```bash
curl http://localhost:8080/api/profiles/P1/requirements \
  -H "Authorization: Bearer <token>"
```

---

### 3.4 Get Profile Members

Retrieve all members assigned to a profile with their compliance status.

**Endpoint:** `GET /api/profiles/{id}/members`

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |

**Response:** `200 OK`
```json
{
  "members": [
    {
      "id": "E1001",
      "name": "Oliver Bennett",
      "department": "Manufacturing",
      "site": "Bristol Central Plant",
      "status": "compliant",
      "joinedAt": "2026-02-18T09:42:33Z"
    },
    {
      "id": "E1020",
      "name": "Ella Davies",
      "department": "Manufacturing",
      "site": "Avonmouth Logistics & Workshop",
      "status": "expired",
      "joinedAt": "2026-02-18T09:42:33Z"
    }
  ],
  "total": 20
}
```

**Member Status Values:**

| Status | Description |
|--------|-------------|
| `compliant` | All requirements are Valid |
| `expiring` | All requirements compliant, but at least one Expiring |
| `expired` | At least one Expired requirement |
| `not_started` | No evidence records exist |

**Example:**
```bash
curl http://localhost:8080/api/profiles/P1/members \
  -H "Authorization: Bearer <token>"
```

---

### 3.5 Add Profile Member

Assign a profile to an employee.

**Endpoint:** `POST /api/profiles/{id}/members`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |

**Request Body:**
```json
{
  "employeeId": "E1005"
}
```

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/profiles/P3/members \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"employeeId": "E1005"}'
```

---

### 3.6 Remove Profile Member

Remove an employee from a profile.

**Endpoint:** `DELETE /api/profiles/{id}/members/{employeeId}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |
| `employeeId` | string | Employee ID to remove |

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/profiles/P3/members/E1020 \
  -H "Authorization: Bearer <admin-token>"
```

---

### 3.7 Get Profile Exclusions

Retrieve all exclusions for a profile.

**Endpoint:** `GET /api/profiles/{id}/exclusions`

**Auth:** Required

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |

**Response:** `200 OK`
```json
[
  {
    "id": "EX-E1009-P3",
    "profileId": "P3",
    "employeeId": "E1009",
    "employeeName": "Noah Wilson",
    "reason": "Medical exemption - back injury prevents forklift operation",
    "expiryDate": "2026-08-18",
    "createdAt": "2026-02-18T09:42:33Z",
    "createdBy": "admin"
  }
]
```

**Example:**
```bash
curl http://localhost:8080/api/profiles/P3/exclusions \
  -H "Authorization: Bearer <token>"
```

---

### 3.8 Add Profile Exclusion

Exclude an employee from a profile with a reason. This also removes the employee's profile assignment.

**Endpoint:** `POST /api/profiles/{id}/exclusions`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |

**Request Body:**
```json
{
  "employeeId": "E1005",
  "reason": "Medical exemption from forklift operation",
  "expiryDate": "2025-12-31"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employeeId` | string | Yes | Employee to exclude |
| `reason` | string | Yes | Reason for exclusion |
| `expiryDate` | string | Yes | Date when exclusion expires (YYYY-MM-DD) |

**Response:** `200 OK`
```json
{
  "id": "EX-E1005-P3",
  "profileId": "P3",
  "employeeId": "E1005",
  "employeeName": "Jack Thompson",
  "reason": "Medical exemption from forklift operation",
  "expiryDate": "2025-12-31",
  "createdAt": "",
  "createdBy": "admin"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/profiles/P3/exclusions \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeId": "E1005",
    "reason": "Medical exemption",
    "expiryDate": "2025-12-31"
  }'
```

---

### 3.9 Update Profile Exclusion

Update an existing exclusion's reason or expiry date.

**Endpoint:** `PUT /api/profiles/{id}/exclusions/{exclusionId}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |
| `exclusionId` | string | Exclusion ID |

**Request Body:**
```json
{
  "reason": "Updated reason",
  "expiryDate": "2026-06-30"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `reason` | string | Yes | Updated reason for exclusion |
| `expiryDate` | string | Yes | Updated expiry date (YYYY-MM-DD) |

**Response:** `200 OK`
Returns the updated exclusion object.

**Error Responses:**

| Status | Description |
|--------|-------------|
| 404 | Exclusion not found |

**Example:**
```bash
curl -X PUT http://localhost:8080/api/profiles/P3/exclusions/EX-E1009-P3 \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Extended medical leave", "expiryDate": "2026-12-31"}'
```

---

### 3.10 Delete Profile Exclusion

Remove a profile exclusion.

**Endpoint:** `DELETE /api/profiles/{id}/exclusions/{exclusionId}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Profile ID |
| `exclusionId` | string | Exclusion ID |

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

**Error Responses:**

| Status | Description |
|--------|-------------|
| 404 | Exclusion not found |

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/profiles/P3/exclusions/EX-E1009-P3 \
  -H "Authorization: Bearer <admin-token>"
```

---

## 4. Courses

### 4.1 Get All Courses

Retrieve all courses and certifications.

**Endpoint:** `GET /api/courses`

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | string | No | Filter by type: `course` or `certification` |

**Response:** `200 OK`
```json
[
  {
    "id": "C001",
    "name": "H&S Induction",
    "code": "HS-001",
    "type": "course",
    "description": "Health and Safety induction for all employees",
    "validity_months": 24,
    "created_at": "2026-01-15T10:00:00Z",
    "updated_at": "2026-01-15T10:00:00Z"
  },
  {
    "id": "C018",
    "name": "Forklift Licence (RTITB/ITSSAR)",
    "code": null,
    "type": "certification",
    "description": "External forklift operator certification",
    "validity_months": 36,
    "created_at": "2026-01-15T10:00:00Z",
    "updated_at": "2026-01-15T10:00:00Z"
  }
]
```

**Example:**
```bash
curl http://localhost:8080/api/courses \
  -H "Authorization: Bearer <token>"
```

---

### 4.2 Create Course (Admin Only)

Create a new course or certification.

**Endpoint:** `POST /api/courses`

**Auth:** Required (Admin only)

**Request Body:**
```json
{
  "name": "New Safety Course",
  "code": "NSC-001",
  "type": "course",
  "description": "Description of the course",
  "validity_months": 12
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Course/certification name |
| `code` | string | No | Internal course code |
| `type` | string | Yes | `course` or `certification` |
| `description` | string | No | Description |
| `validity_months` | number | No | Validity period in months (0 = never expires) |

**Response:** `201 Created`
```json
{
  "id": "C025",
  "name": "New Safety Course",
  "code": "NSC-001",
  "type": "course",
  "description": "Description of the course",
  "validity_months": 12,
  "created_at": "2026-02-18T10:00:00Z",
  "updated_at": "2026-02-18T10:00:00Z"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/courses \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Safety Course",
    "code": "NSC-001",
    "type": "course",
    "description": "Description of the course",
    "validity_months": 12
  }'
```

---

### 4.3 Update Course (Admin Only)

Update an existing course.

**Endpoint:** `PUT /api/courses/{id}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Course ID |

**Request Body:** (all fields optional)
```json
{
  "name": "Updated Course Name",
  "validity_months": 24
}
```

**Response:** `200 OK`
Returns the updated course object.

**Example:**
```bash
curl -X PUT http://localhost:8080/api/courses/C001 \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Course Name", "validity_months": 24}'
```

---

### 4.4 Delete Course (Admin Only)

Delete a course. Courses in use by profiles cannot be deleted.

**Endpoint:** `DELETE /api/courses/{id}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Course ID |

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

**Error Responses:**

| Status | Description |
|--------|-------------|
| 400 | Course is in use by one or more profiles |
| 404 | Course not found |

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/courses/C025 \
  -H "Authorization: Bearer <admin-token>"
```

---

## 5. Events (Sessions)

### 5.1 Get All Events

Retrieve all training events/sessions.

**Endpoint:** `GET /api/events`

**Auth:** Required

**Response:** `200 OK`
```json
[
  {
    "id": "TE001",
    "course_id": "C004",
    "course": {
      "id": "C004",
      "name": "Manual Handling",
      "type": "course"
    },
    "date": "2026-01-18",
    "time": "09:00",
    "duration": 120,
    "location": "Bristol Central Plant",
    "instructor_id": "E1003",
    "status": "completed",
    "attendees": ["E1001", "E1002", "E1003", "E1004", "E1009"]
  },
  {
    "id": "TE003",
    "course_id": "C012",
    "course": {
      "id": "C012",
      "name": "Forklift Practical Assessment",
      "type": "course"
    },
    "date": "2026-02-25",
    "time": "08:30",
    "duration": 180,
    "location": "Avonmouth Logistics & Workshop",
    "instructor_id": "E1006",
    "status": "scheduled",
    "attendees": []
  }
]
```

**Event Status:**

| Status | Description |
|--------|-------------|
| `draft` | Event is being planned, not yet finalized |
| `scheduled` | Event is confirmed and planned |
| `in_progress` | Event is currently running |
| `completed` | Attendance confirmed, evidence records generated |
| `cancelled` | Event has been cancelled |

**Example:**
```bash
curl http://localhost:8080/api/events \
  -H "Authorization: Bearer <token>"
```

---

### 5.2 Create Event

Create a new training event/session.

**Endpoint:** `POST /api/events`

**Auth:** Required (Admin only)

**Request Body:**
```json
{
  "courseId": "C004",
  "date": "2025-03-15",
  "time": "09:00",
  "duration": 120,
  "location": "Avonmouth Logistics & Workshop",
  "instructorId": "E1003",
  "attendees": ["E1005", "E1010", "E1013"]
}```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `courseId` | string | Yes | Course ID from the courses table |
| `date` | string | Yes | Event date (YYYY-MM-DD) |
| `time` | string | Yes | Start time in HH:MM format (e.g., "09:00") |
| `duration` | number | Yes | Duration in minutes (e.g., 120 for 2 hours) |
| `location` | string | Yes | Event location |
| `instructorId` | string | Yes | Employee ID of the instructor |
| `attendees` | array | No | List of employee IDs attending |

**Response:** `200 OK`
```json
{
  "id": "TE004",
  "course_id": "C004",
  "course": {
    "id": "C004",
    "name": "Manual Handling",
    "type": "course"
  },
  "date": "2025-03-15",
  "time": "09:00",
  "duration": 120,
  "location": "Avonmouth Logistics & Workshop",
  "instructor_id": "E1003",
  "status": "scheduled",
  "attendees": ["E1005", "E1010", "E1013"]
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "courseId": "C004",
    "date": "2025-03-15",
    "time": "09:00",
    "duration": 120,
    "location": "Avonmouth Logistics & Workshop",
    "instructorId": "E1003",
    "attendees": ["E1005", "E1010", "E1013"]
  }'
```

---

### 5.3 Update Event

Update an existing event's details. Only provided fields are updated.

**Endpoint:** `PUT /api/events/{id}`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Event ID |

**Request Body:**
```json
{
  "courseId": "C004",
  "date": "2025-04-01",
  "time": "10:00",
  "duration": 90,
  "location": "Bristol Central Plant",
  "instructorId": "E1003",
  "status": "scheduled",
  "attendees": ["E1005", "E1010"]
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `courseId` | string | No | Course ID |
| `date` | string | No | Event date (YYYY-MM-DD) |
| `time` | string | No | Start time (HH:MM) |
| `duration` | number | No | Duration in minutes |
| `location` | string | No | Event location |
| `instructorId` | string | No | Instructor employee ID |
| `status` | string | No | Event status |
| `attendees` | array | No | Replaces attendee list if provided |

**Response:** `200 OK`
Returns the updated event object.

**Example:**
```bash
curl -X PUT http://localhost:8080/api/events/TE003 \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"time": "10:00", "duration": 90, "status": "draft"}'
```

---

### 5.4 Confirm Attendance (Enhanced)

Confirm attendance for an event with per-attendee status and walk-in support. This generates evidence records for attendees who were present.

**Endpoint:** `POST /api/events/{id}/attendance/confirm`

**Auth:** Required (Admin only)

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | Event ID |

**Request Body:**
```json
{
  "attendees": [
    {
      "employeeId": "E1005",
      "attended": true
    },
    {
      "employeeId": "E1010",
      "attended": false,
      "absenceReason": "sick_leave",
      "absenceNotes": "Called in sick"
    }
  ],
  "walkIns": [
    {
      "employeeId": "E1020"
    }
  ]
}
```

**Attendee Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employeeId` | string | Yes | Employee ID |
| `attended` | boolean | Yes | Whether the employee attended |
| `absenceReason` | string | No | Reason for absence (if `attended` is false) |
| `absenceNotes` | string | No | Additional notes about absence |

**Absence Reason Values:**

| Value | Description |
|-------|-------------|
| `sick_leave` | Employee was on sick leave |
| `no_show` | Employee did not appear without notice |
| `emergency` | Employee had an emergency |
| `other` | Other reason (see absenceNotes) |

**Walk-In Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employeeId` | string | Yes | Employee ID of the walk-in |

**Response:** `200 OK`
```json
{
  "status": "ok"
}
```

**Notes:**
- Evidence records are created only for attendees with `attended: true` and all walk-ins
- The evidence `completionDate` is set to the event date
- The evidence `expiryDate` is calculated based on the course validity
- Event status is updated to `completed`
- Non-attending employees are tracked but do not receive evidence records

**Example:**
```bash
curl -X POST http://localhost:8080/api/events/TE003/attendance/confirm \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "attendees": [
      {"employeeId": "E1005", "attended": true},
      {"employeeId": "E1010", "attended": false, "absenceReason": "sick_leave"}
    ],
    "walkIns": [{"employeeId": "E1020"}]
  }'
```

---

## 6. Certifications

### 6.1 Get All Certifications

Retrieve all certification records, optionally filtered by employee.

**Endpoint:** `GET /api/certifications`

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `employeeId` | string | No | Filter by employee ID |

**Response:** `200 OK`
```json
[
  {
    "id": "CERT-1",
    "employeeId": "E1005",
    "requirementName": "Forklift Licence (RTITB/ITSSAR)",
    "issuer": "RTITB",
    "issueDate": "2024-03-15",
    "expiryDate": "2027-03-15",
    "notes": "Renewal due Q1 2027"
  }
]
```

**Example:**
```bash
# Get all certifications
curl http://localhost:8080/api/certifications \
  -H "Authorization: Bearer <admin-token>"

# Get certifications for specific employee
curl "http://localhost:8080/api/certifications?employeeId=E1005" \
  -H "Authorization: Bearer <token>"
```

---

### 6.2 Create Certification

Upload a third-party certification record.

**Endpoint:** `POST /api/certifications`

**Auth:** Required (Admin only)

**Request Body:**
```json
{
  "employeeId": "E1005",
  "courseId": "C018",
  "issuer": "RTITB",
  "issueDate": "2024-03-15",
  "expiryDate": "2027-03-15",
  "notes": "Original certification"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `employeeId` | string | Yes | Employee ID |
| `courseId` | string | Yes | Course ID from the courses table |
| `issuer` | string | Yes | Certifying body or issuer |
| `issueDate` | string | Yes | Date issued (YYYY-MM-DD) |
| `expiryDate` | string | No | Expiry date (YYYY-MM-DD), omit if never expires |
| `notes` | string | No | Additional notes |

**Response:** `201 Created`
```json
{
  "id": "CERT-1234567890",
  "employeeId": "E1005",
  "courseId": "C018",
  "course": {
    "id": "C018",
    "name": "Forklift Licence (RTITB/ITSSAR)",
    "type": "certification"
  },
  "issuer": "RTITB",
  "issueDate": "2024-03-15",
  "expiryDate": "2027-03-15",
  "notes": "Original certification"
}
```

**Notes:**
- This endpoint also creates an evidence record automatically
- The evidence type will be set to `certification`

**Example:**
```bash
curl -X POST http://localhost:8080/api/certifications \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeId": "E1005",
    "courseId": "C018",
    "issuer": "RTITB",
    "issueDate": "2024-03-15",
    "expiryDate": "2027-03-15",
    "notes": "Original certification"
  }'
```

---

## 7. Dashboard

### 7.1 Get Dashboard KPIs

Retrieve key performance indicators for the dashboard.

**Endpoint:** `GET /api/dashboard/kpis`

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `department` | string | No | Filter by department |
| `site` | string | No | Filter by site |
| `profile` | string | No | Filter by profile ID |

**Response:** `200 OK`
```json
{
  "totalEmployees": 20,
  "overallCompliance": 86.0,
  "valid": 122,
  "expiringSoon": 13,
  "expired": 10,
  "missing": 12
}
```

**Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| `totalEmployees` | number | Total employees matching filters |
| `overallCompliance` | number | Percentage of compliant requirements (Valid + Expiring) |
| `valid` | number | Requirements with Valid status |
| `expiringSoon` | number | Requirements with Expiring status (within 90 days) |
| `expired` | number | Requirements with Expired status |
| `missing` | number | Requirements with Missing status |

**Notes:**
- `overallCompliance` counts Expiring as compliant (per BLS spec)
- Filters are applied additively (AND logic)

**Example:**
```bash
# Get all KPIs
curl http://localhost:8080/api/dashboard/kpis \
  -H "Authorization: Bearer <admin-token>"

# Get KPIs for Manufacturing department
curl "http://localhost:8080/api/dashboard/kpis?department=Manufacturing" \
  -H "Authorization: Bearer <admin-token>"

# Get KPIs for specific profile
curl "http://localhost:8080/api/dashboard/kpis?profile=P3" \
  -H "Authorization: Bearer <admin-token>"
```

---

### 7.2 Get Coverage Grid

Retrieve coverage grid data for visual gap analysis.

**Endpoint:** `GET /api/dashboard/grid`

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mode` | string | No | Grid mode: `course` or `profile` (default: `course`) |
| `groupBy` | string | No | Grouping: `department`, `site`, or `costCenter` (default: `department`) |
| `department` | string | No | Filter by department |
| `site` | string | No | Filter by site |

**Response:** `200 OK`

**Course Mode Response:**
```json
[
  {
    "groupValue": "Manufacturing",
    "itemId": "R101",
    "itemName": "H&S Induction",
    "total": 8,
    "valid": 6,
    "expiring": 1,
    "expired": 0,
    "missing": 1,
    "percent": 87.5
  },
  {
    "groupValue": "Manufacturing",
    "itemId": "R102",
    "itemName": "Fire Safety Awareness",
    "total": 8,
    "valid": 5,
    "expiring": 2,
    "expired": 1,
    "missing": 0,
    "percent": 87.5
  }
]
```

**Profile Mode Response:**
```json
[
  {
    "groupValue": "Manufacturing",
    "itemId": "P1",
    "itemName": "Safety Base (All Employees)",
    "total": 8,
    "valid": 5,
    "expiring": 2,
    "expired": 0,
    "missing": 1,
    "percent": 87.5
  },
  {
    "groupValue": "Manufacturing",
    "itemId": "P3",
    "itemName": "Forklift Operator",
    "total": 3,
    "valid": 1,
    "expiring": 1,
    "expired": 0,
    "missing": 1,
    "percent": 66.67
  }
]
```

**Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| `groupValue` | string | The group (department/site/costCenter) |
| `itemId` | string | Course requirement ID or profile ID |
| `itemName` | string | Course name or profile name |
| `total` | number | Total employees in this group requiring this item |
| `valid` | number | Employees with Valid status |
| `expiring` | number | Employees with Expiring status |
| `expired` | number | Employees with Expired status |
| `missing` | number | Employees with Missing status |
| `percent` | number | Compliance percentage ((valid + expiring) / total * 100) |

**Example:**
```bash
# Course coverage by department
curl "http://localhost:8080/api/dashboard/grid?mode=course&groupBy=department" \
  -H "Authorization: Bearer <admin-token>"

# Profile coverage by site
curl "http://localhost:8080/api/dashboard/grid?mode=profile&groupBy=site" \
  -H "Authorization: Bearer <admin-token>"

# Filtered grid
curl "http://localhost:8080/api/dashboard/grid?mode=course&department=Manufacturing" \
  -H "Authorization: Bearer <admin-token>"
```

---

### 7.3 Get Drilldown

Retrieve employee-level details for a specific grid cell.

**Endpoint:** `GET /api/dashboard/drilldown`

**Auth:** Required

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mode` | string | Yes | Grid mode: `course` or `profile` |
| `groupBy` | string | Yes | Grouping used in grid |
| `groupValue` | string | Yes | The group value to drill into |
| `itemId` | string | Yes | Course ID or profile ID |
| `department` | string | No | Additional department filter |
| `site` | string | No | Additional site filter |

**Response:** `200 OK`

**Course Mode Response:**
```json
[
  {
    "id": "E1001",
    "name": "Oliver Bennett",
    "department": "Manufacturing",
    "site": "Bristol Central Plant",
    "status": "Valid",
    "expiryDate": "2026-02-18"
  },
  {
    "id": "E1005",
    "name": "Jack Thompson",
    "department": "Manufacturing",
    "site": "Avonmouth Logistics & Workshop",
    "status": "Expiring",
    "expiryDate": "2025-04-03"
  },
  {
    "id": "E1009",
    "name": "Noah Wilson",
    "department": "Manufacturing",
    "site": "Bristol Central Plant",
    "status": "Missing",
    "expiryDate": ""
  }
]
```

**Profile Mode Response:**
```json
[
  {
    "id": "E1002",
    "name": "Amelia Clarke",
    "department": "Manufacturing",
    "site": "Bristol Central Plant",
    "status": "Eligible",
    "expiryDate": ""
  },
  {
    "id": "E1005",
    "name": "Jack Thompson",
    "department": "Manufacturing",
    "site": "Avonmouth Logistics & Workshop",
    "status": "At Risk",
    "expiryDate": ""
  },
  {
    "id": "E1010",
    "name": "Ava Morgan",
    "department": "Manufacturing",
    "site": "Avonmouth Logistics & Workshop",
    "status": "Not Eligible",
    "expiryDate": ""
  }
]
```

**Example:**
```bash
# Drill down into Manufacturing - H&S Induction (course mode)
curl "http://localhost:8080/api/dashboard/drilldown?mode=course&groupBy=department&groupValue=Manufacturing&itemId=R101" \
  -H "Authorization: Bearer <admin-token>"

# Drill down into Manufacturing - Forklift Operator (profile mode)
curl "http://localhost:8080/api/dashboard/drilldown?mode=profile&groupBy=department&groupValue=Manufacturing&itemId=P3" \
  -H "Authorization: Bearer <admin-token>"
```

---

## 8. Error Responses

All error responses follow a consistent format:

```json
{
  "error": "Error message describing what went wrong",
  "code": 400
}
```

### Common Error Codes

| Status | Description |
|--------|-------------|
| 400 | Bad Request - Invalid request body or parameters |
| 401 | Unauthorized - Missing or invalid authentication token |
| 403 | Forbidden - Authenticated but not authorized for this resource |
| 404 | Not Found - Resource does not exist |
| 500 | Internal Server Error - Something went wrong on the server |

### Error Examples

**401 Unauthorized:**
```json
{
  "error": "Unauthorized",
  "code": 401
}
```

**403 Forbidden:**
```json
{
  "error": "Admin access required",
  "code": 403
}
```

**404 Not Found:**
```json
{
  "error": "Employee not found",
  "code": 404
}
```

---

## 9. Data Types

### Employee

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique employee identifier (e.g., `E1001`) |
| `name` | string | Full name |
| `department` | string | Department name |
| `site` | string | Work site location |
| `worker_type` | string | `Employee` or `Contractor` |
| `cost_center` | string | Cost center code |
| `employment_status` | string | Employment status (e.g., `Active`) |

### Profile (List Response)

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique profile identifier (e.g., `P1`) |
| `name` | string | Profile name |
| `description` | string | Profile description |
| `memberCount` | number | Number of employees assigned to this profile |
| `exclusionCount` | number | Number of active exclusions for this profile |

### Course

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique course identifier (e.g., `C001`) |
| `name` | string | Course/certification name |
| `code` | string | Internal course code (nullable) |
| `type` | string | `course` or `certification` |
| `description` | string | Course description |
| `validity_months` | number | Validity period in months (0 = never expires) |
| `created_at` | string | Creation timestamp (ISO 8601) |
| `updated_at` | string | Last update timestamp (ISO 8601) |

### Requirement

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique requirement identifier (e.g., `R101`) |
| `profile_id` | string | Parent profile ID |
| `name` | string | Requirement name |
| `type` | string | `course` or `certification` |
| `validity_months` | number | Validity period in months (0 = never expires) |

### Evidence

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique evidence identifier |
| `employeeId` | string | Employee ID |
| `course_id` | string | Course ID this satisfies |
| `requirementName` | string | Requirement name this satisfies (legacy) |
| `evidenceType` | string | `attendance` or `certification` |
| `completionDate` | string | Date completed (YYYY-MM-DD) |
| `expiryDate` | string | Expiry date (YYYY-MM-DD), null if never expires |
| `metadata` | string | JSON metadata (e.g., issuer info) |

### TrainingEvent

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique event identifier |
| `course_id` | string | Course ID |
| `course` | object | Embedded course object with `id`, `name`, `type` |
| `date` | string | Event date (YYYY-MM-DD) |
| `time` | string | Start time (HH:MM format, e.g., `09:00`) |
| `duration` | number | Duration in minutes |
| `location` | string | Event location |
| `instructor_id` | string | Instructor employee ID |
| `status` | string | `draft`, `scheduled`, `in_progress`, `completed`, or `cancelled` |
| `attendees` | array | List of employee IDs |

### ProfileMember

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Employee ID |
| `name` | string | Employee full name |
| `department` | string | Department |
| `site` | string | Site location |
| `status` | string | `compliant`, `expiring`, `expired`, or `not_started` |
| `joinedAt` | string | Date/time assigned to profile |

### ProfileExclusion

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique exclusion identifier |
| `profileId` | string | Profile ID |
| `employeeId` | string | Employee ID |
| `employeeName` | string | Employee full name |
| `reason` | string | Reason for exclusion |
| `expiryDate` | string | Expiry date (YYYY-MM-DD) |
| `createdAt` | string | Creation timestamp |
| `createdBy` | string | Creator (e.g., `admin`) |

### Certification

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique certification identifier |
| `employeeId` | string | Employee ID |
| `courseId` | string | Course ID from the courses table |
| `course` | object | Embedded course object with `id`, `name`, `type` |
| `issuer` | string | Certifying body |
| `issueDate` | string | Date issued (YYYY-MM-DD) |
| `expiryDate` | string | Expiry date (YYYY-MM-DD), empty if never expires |
| `notes` | string | Additional notes |

---

## Appendix: Sample Demo Flow

### 1. Login as Admin
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}' | jq -r '.token')
```

### 2. Get Employee Details
```bash
curl http://localhost:8080/api/employees/E1005 \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 3. View Dashboard KPIs
```bash
curl http://localhost:8080/api/dashboard/kpis \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 4. Create Training Event
```bash
EVENT_ID=$(curl -s -X POST http://localhost:8080/api/events \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "courseId": "C004",
    "date": "2025-03-15",
    "location": "Bristol Central Plant",
    "instructorId": "E1003",
    "attendees": ["E1005"]
  }' | jq -r '.id')
```

### 5. Confirm Attendance
```bash
curl -X POST http://localhost:8080/api/events/$EVENT_ID/attendance/confirm \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"employeeIds": ["E1005"]}'
```

### 6. Verify Updated Status
```bash
curl http://localhost:8080/api/employees/E1005 \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## Appendix: Docker Deployment

### Quick Start with Docker Compose

From the project root:

```bash
# Build and run all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Backend Only

```bash
# From backend directory
cd backend
docker-compose up -d
```

### Manual Docker Build

```bash
# Build the image
cd backend
docker build -t tms-backend:latest .

# Run the container
docker run -d \
  --name tms-backend \
  -p 8080:8080 \
  -v tms-data:/app/data \
  -e PORT=8080 \
  -e DB_PATH=/app/data/tms.db \
  tms-backend:latest
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_PATH` | `./tms.db` | SQLite database path |

### Health Check

```bash
curl http://localhost:8080/health
# Response: {"status":"healthy"}
```

### Docker Compose Services

| Service | Port | Description |
|---------|------|-------------|
| `backend` | 8080 | Go API server |
| `frontend` | 80 | Vue frontend (nginx) |

### Volume Mounts

| Volume | Purpose |
|--------|---------|
| `tms-data` | Persists SQLite database |

---

## Appendix: Development Setup

### Prerequisites

- Go 1.22+
- Node.js 20+ (for frontend)
- Docker & Docker Compose (optional)

### Local Development (Backend)

```bash
cd backend

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go

# Server will start on http://localhost:8080
```

### Local Development (Frontend)

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Server will start on http://localhost:5173
```

### Testing the API

```bash
# Health check
curl http://localhost:8080/health

# Login as admin
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'

# Get employees (with token)
curl http://localhost:8080/api/employees \
  -H "Authorization: Bearer YOUR_TOKEN"
```
