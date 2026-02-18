export type Role = 'admin' | 'employee'

export interface Course {
  id: string
  name: string
  code?: string
  type: 'course' | 'certification'
  description?: string
  validity_months: number
  created_at: string
  updated_at: string
}

export interface CreateCourseRequest {
  name: string
  code?: string
  type: 'course' | 'certification'
  description?: string
  validity_months: number
}

export interface UpdateCourseRequest {
  name?: string
  code?: string
  type?: 'course' | 'certification'
  description?: string
  validity_months?: number
}

export interface Employee {
  id: string
  name: string
  department: string
  site: string
  worker_type: string
  cost_center: string
  employment_status: string
}

export type RequirementStatusValue = 'Valid' | 'Expiring' | 'Expired' | 'Missing'
export type ProfileStatusValue = 'Eligible' | 'At Risk' | 'Not Eligible'

export interface RequirementWithStatus {
  id: string
  profile_id: string
  name: string
  type: 'course' | 'certification'
  validity_months: number
  status: RequirementStatusValue
  expiry_date: string | null
}

export interface EmployeeDetail extends Employee {
  profile_status: Record<string, ProfileStatusValue>
  requirements: RequirementWithStatus[]
}

export interface EvidenceRecord {
  id: string
  employeeId: string
  courseId: string
  course?: Course
  requirementName: string
  evidenceType: 'attendance' | 'certification'
  completionDate: string
  expiryDate: string | null
  metadata: string
}

export interface DashboardKpi {
  totalEmployees: number
  overallCompliance: number
  valid: number
  expiringSoon: number
  expired: number
  missing: number
}

export interface Profile {
  id: string
  name: string
  description: string
  memberCount?: number
  exclusionCount?: number
}

export interface Requirement {
  course_id: string
  id: string
  profile_id: string
  name: string
  type: 'course' | 'certification'
  validity_months: number
}

export interface ProfileMember {
  id: string
  name: string
  department: string
  site: string
  status: 'compliant' | 'expiring' | 'expired' | 'not_started'
  joinedAt: string
}

export interface ProfileExclusion {
  id: string
  profileId: string
  employeeId: string
  employeeName: string
  reason: string
  expiryDate: string | null
  createdAt: string
  createdBy: string
}

export interface TrainingEvent {
  id: string
  course_id: string
  course?: Course
  course_name?: string
  date: string
  time: string
  duration: number
  location: string
  instructor_id: string
  status: 'draft' | 'scheduled' | 'in_progress' | 'completed' | 'cancelled'
  attendees: string[]
}

export interface AttendanceAttendee {
  employeeId: string
  attended: boolean
  absenceReason?: 'sick_leave' | 'no_show' | 'emergency' | 'other'
  absenceNotes?: string
}

export interface AttendanceWalkIn {
  employeeId: string
}

export interface ConfirmAttendanceRequest {
  attendees: AttendanceAttendee[]
  walkIns?: AttendanceWalkIn[]
}

export interface CreateEventRequest {
  courseId: string
  date: string
  time: string
  duration: number
  location: string
  instructorId: string
  attendees?: string[]
}

export interface UpdateEventRequest {
  courseId?: string
  date?: string
  time?: string
  duration?: number
  location?: string
  instructorId?: string
  status?: 'draft' | 'scheduled' | 'in_progress' | 'completed' | 'cancelled'
  attendees?: string[]
}

export interface AddExclusionRequest {
  employeeId: string
  reason: string
  expiryDate: string
}

export interface UpdateExclusionRequest {
  reason: string
  expiryDate: string
}

export interface Certification {
  id: string
  employee_id: string
  course_id: string
  course?: Course
  issued_date: string
  expiry_date: string | null
  certificate_number?: string
  created_at: string
  updated_at: string
}

export interface CreateCertificationRequest {
  employeeId: string
  courseId: string
  issuedDate: string
  expiryDate?: string
  certificateNumber?: string
}
