<script setup lang="ts">
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { 
  HelpCircle, 
  LayoutDashboard, 
  Grid3x3,
  Users, 
  BookOpen, 
  BadgeCheck, 
  CalendarRange, 
  FileBadge2, 
  User,
  Mail
} from 'lucide-vue-next'
</script>

<template>
  <section class="page-section">
    <div class="page-header">
      <div class="page-header-left">
        <div class="header-icon-row">
          <HelpCircle class="icon-standard text-slate-600" />
          <h1 class="page-title">TMS Quick Start Guide</h1>
        </div>
        <p class="page-description">Everything you need to know to get started with the Training Management System</p>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">What is TMS?</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="table-cell-secondary">
          TMS (Training Management System) is a comprehensive platform for managing employee training, certifications, and compliance. 
          It helps organizations track training requirements, schedule sessions, and ensure employees maintain up-to-date certifications.
        </p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">Main Features</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="card-content-spaced">
          <div class="requirement-section">
            <div class="header-info-row">
              <LayoutDashboard class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Dashboard</h3>
            </div>
            <p class="table-cell-secondary">View overall compliance metrics, expiring certifications, and coverage statistics across departments and sites.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <Grid3x3 class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Training Matrix</h3>
            </div>
            <p class="table-cell-secondary">Visual grid showing all employees vs all courses. Click any cell to view training details and history. Color-coded status indicators show compliance at a glance.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <Users class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Employees</h3>
            </div>
            <p class="table-cell-secondary">Manage employee records, view individual training history, compliance status, and certification details.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <BookOpen class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Courses</h3>
            </div>
            <p class="table-cell-secondary">Define and manage training courses including course details, validity periods, and associated profiles.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <BadgeCheck class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Profiles</h3>
            </div>
            <p class="table-cell-secondary">Create training profiles that group required courses for specific roles, departments, or job functions.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <CalendarRange class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Sessions</h3>
            </div>
            <p class="table-cell-secondary">Schedule and manage training sessions, assign attendees, track completion, and record session outcomes.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <FileBadge2 class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">Certifications</h3>
            </div>
            <p class="table-cell-secondary">View and manage employee certifications, track expiry dates, and maintain compliance records.</p>
          </div>

          <div class="requirement-section">
            <div class="header-info-row">
              <User class="icon-standard text-slate-600" />
              <h3 class="text-heading-sm">My Training</h3>
            </div>
            <p class="table-cell-secondary">Personal view for employees to see their assigned training, upcoming sessions, and certification status.</p>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">Training Matrix Guide</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="table-cell-secondary mb-4">
          The Training Matrix provides a comprehensive view of training compliance across your entire organization.
        </p>
        <div class="table-container">
          <table class="w-full">
            <thead class="table-header">
              <tr>
                <th class="table-header-cell">Cell Color</th>
                <th class="table-header-cell">Status</th>
                <th class="table-header-cell">Meaning</th>
              </tr>
            </thead>
            <tbody>
              <tr class="table-row">
                <td class="table-cell">
                  <div class="matrix-legend-cell matrix-legend-valid">
                    <span class="matrix-legend-check">✓</span>
                  </div>
                </td>
                <td class="table-cell table-cell-primary">Valid</td>
                <td class="table-cell table-cell-secondary">Training is current and valid (expiry > 90 days)</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell">
                  <div class="matrix-legend-cell matrix-legend-expiring">
                    <span class="matrix-legend-warning">⚠</span>
                  </div>
                </td>
                <td class="table-cell table-cell-primary">Expiring</td>
                <td class="table-cell table-cell-secondary">Training expires within 90 days</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell">
                  <div class="matrix-legend-cell matrix-legend-expired">
                    <span class="matrix-legend-x">✗</span>
                  </div>
                </td>
                <td class="table-cell table-cell-primary">Expired</td>
                <td class="table-cell table-cell-secondary">Training has passed its expiry date</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell">
                  <div class="matrix-legend-cell matrix-legend-not-taken">
                    <span class="matrix-legend-dash">—</span>
                  </div>
                </td>
                <td class="table-cell table-cell-primary">Not Taken</td>
                <td class="table-cell table-cell-secondary">Required training not yet completed</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell">
                  <div class="matrix-legend-cell matrix-legend-not-required">
                    <span class="matrix-legend-empty"></span>
                  </div>
                </td>
                <td class="table-cell table-cell-primary">Not Required</td>
                <td class="table-cell table-cell-secondary">Course not in employee's profile</td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">Status Color Guide</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="table-container">
          <table class="w-full">
            <thead class="table-header">
              <tr>
                <th class="table-header-cell">Status</th>
                <th class="table-header-cell">Badge</th>
                <th class="table-header-cell">Meaning</th>
              </tr>
            </thead>
            <tbody>
              <tr class="table-row">
                <td class="table-cell table-cell-primary">Valid</td>
                <td class="table-cell"><Badge class="status-badge status-badge-success">Valid</Badge></td>
                <td class="table-cell table-cell-secondary">Certification is current and valid</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell table-cell-primary">Expiring</td>
                <td class="table-cell"><Badge class="status-badge status-badge-warning">Expiring</Badge></td>
                <td class="table-cell table-cell-secondary">Certification expires within 90 days</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell table-cell-primary">Expired</td>
                <td class="table-cell"><Badge class="status-badge status-badge-danger">Expired</Badge></td>
                <td class="table-cell table-cell-secondary">Certification has passed its expiry date</td>
              </tr>
              <tr class="table-row">
                <td class="table-cell table-cell-primary">Missing</td>
                <td class="table-cell"><Badge class="status-badge status-badge-neutral">Missing</Badge></td>
                <td class="table-cell table-cell-secondary">Required certification not yet obtained</td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">Quick Workflow Example</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="card-content-spaced">
          <div class="requirement-section">
            <h3 class="requirement-title">1. Set Up Courses</h3>
            <p class="table-cell-secondary">Navigate to <strong>Courses</strong> and create training courses with validity periods.</p>
          </div>

          <div class="requirement-section">
            <h3 class="requirement-title">2. Create Profiles</h3>
            <p class="table-cell-secondary">Go to <strong>Profiles</strong> and group courses into profiles based on job roles or departments.</p>
          </div>

          <div class="requirement-section">
            <h3 class="requirement-title">3. Assign Profiles to Employees</h3>
            <p class="table-cell-secondary">In <strong>Employees</strong>, assign the appropriate profiles to each employee.</p>
          </div>

          <div class="requirement-section">
            <h3 class="requirement-title">4. Schedule Training Sessions</h3>
            <p class="table-cell-secondary">Use <strong>Sessions</strong> to schedule training and invite attendees.</p>
          </div>

          <div class="requirement-section">
            <h3 class="requirement-title">5. Monitor Compliance</h3>
            <p class="table-cell-secondary">Check the <strong>Dashboard</strong> for overall compliance or use the <strong>Training Matrix</strong> for a detailed employee-by-course view.</p>
          </div>
        </div>
      </CardContent>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle class="sheet-title">Need Help?</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="header-info-row">
          <Mail class="icon-standard text-slate-600" />
          <div>
            <p class="table-cell-primary">Contact Support</p>
            <p class="table-cell-secondary">For technical assistance or questions, please reach out to your system administrator or IT support team.</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
