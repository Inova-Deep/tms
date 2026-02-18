<script setup lang="ts">
import { ref, computed } from 'vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'
import { useMatrix } from '@/composables/useMatrix'
import { useEmployeeRecords } from '@/composables/useEmployeeDetail'
import { Check, X, AlertTriangle, Minus, Calendar, Award, ExternalLink } from 'lucide-vue-next'
import type { MatrixCell, Course, EvidenceRecord } from '@/types/api'

const { data: matrix, isLoading } = useMatrix()

// Drill-down sheet state
const sheetOpen = ref(false)
const selectedEmployeeId = ref<string | undefined>(undefined)
const selectedEmployeeName = ref<string>('')
const selectedCourseId = ref<string | undefined>(undefined)
const selectedCourse = ref<Course | null>(null)
const selectedCell = ref<MatrixCell | null>(null)

// Fetch employee records when sheet opens
const { data: employeeRecords } = useEmployeeRecords(selectedEmployeeId)

// Filter evidence for the selected course
const courseEvidence = computed(() => {
  if (!employeeRecords.value || !selectedCourseId.value) return []
  return employeeRecords.value.filter((e: EvidenceRecord) => e.courseId === selectedCourseId.value)
})

function getCellClass(cell: MatrixCell | undefined): string {
  if (!cell) return 'matrix-cell-not-required'
  switch (cell.status) {
    case 'valid':
      return 'matrix-cell-valid'
    case 'expiring':
      return 'matrix-cell-expiring'
    case 'expired':
      return 'matrix-cell-expired'
    case 'not_taken':
      return 'matrix-cell-not-taken'
    case 'not_required':
      return 'matrix-cell-not-required'
    default:
      return ''
  }
}

function handleCellClick(
  employeeId: string,
  employeeName: string,
  courseId: string,
  cell: MatrixCell | undefined
) {
  if (!cell || cell.status === 'not_required') return
  
  selectedEmployeeId.value = employeeId
  selectedEmployeeName.value = employeeName
  selectedCourseId.value = courseId
  selectedCourse.value = matrix.value?.courses.find(c => c.id === courseId) || null
  selectedCell.value = cell
  sheetOpen.value = true
}

function formatDate(dateStr: string | null | undefined): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

function getStatusLabel(status: string): string {
  switch (status) {
    case 'valid': return 'Valid'
    case 'expiring': return 'Expiring Soon'
    case 'expired': return 'Expired'
    case 'not_taken': return 'Not Taken'
    default: return status
  }
}

function goToEmployee() {
  if (selectedEmployeeId.value) {
    window.location.href = `/employees/${selectedEmployeeId.value}`
  }
}
</script>

<template>
  <section class="page-section">
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Training Matrix</h1>
        <p class="page-description">Overview of employee training status across all courses.</p>
      </div>
    </div>

    <Card>
      <CardHeader class="card-header-compact">
        <CardTitle class="text-heading-sm">Compliance Matrix</CardTitle>
      </CardHeader>
      <CardContent class="p-0">
        <div v-if="isLoading" class="loading-container">
          <div class="loading-content">
            <div class="loading-spinner loading-spinner-md"></div>
            <span>Loading matrix...</span>
          </div>
        </div>

        <div v-else-if="!matrix || matrix.rows.length === 0" class="empty-state">
          <div class="empty-state-content">
            <p class="empty-state-title">No data available</p>
            <p class="empty-state-description">No training matrix data found.</p>
          </div>
        </div>

        <div v-else class="matrix-container">
          <table class="matrix-table">
            <thead>
              <tr>
                <th class="matrix-header-sticky matrix-header-name">Employee</th>
                <th class="matrix-header-sticky matrix-header-dept">Department</th>
                <th
                  v-for="course in matrix.courses"
                  :key="course.id"
                  class="matrix-header-sticky matrix-header-course"
                  :title="course.name"
                >
                  {{ course.code || course.name }}
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in matrix.rows" :key="row.employeeId" class="matrix-row">
                <td class="matrix-cell-name">{{ row.employeeName }}</td>
                <td class="matrix-cell-dept">{{ row.department }}</td>
                <td
                  v-for="course in matrix.courses"
                  :key="course.id"
                  class="matrix-cell matrix-cell-clickable"
                  :class="getCellClass(row.cells[course.id])"
                  @click="handleCellClick(row.employeeId, row.employeeName, course.id, row.cells[course.id])"
                >
                  <span class="matrix-cell-content">
                    <Check v-if="row.cells[course.id]?.status === 'valid'" class="matrix-icon-valid" />
                    <AlertTriangle v-else-if="row.cells[course.id]?.status === 'expiring'" class="matrix-icon-expiring" />
                    <X v-else-if="row.cells[course.id]?.status === 'expired'" class="matrix-icon-expired" />
                    <Minus v-else-if="row.cells[course.id]?.status === 'not_taken'" class="matrix-icon-not-taken" />
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>

    <!-- Drill-down Sheet -->
    <Sheet :open="sheetOpen" @update:open="(open) => !open && (sheetOpen = false)">
      <SheetContent class="sheet-content">
        <SheetHeader class="sheet-header">
          <SheetTitle>{{ selectedEmployeeName }}</SheetTitle>
          <SheetDescription>{{ selectedCourse?.name }}</SheetDescription>
        </SheetHeader>

        <div class="sheet-body">
          <!-- Status Badge -->
          <div class="matrix-detail-status">
            <div 
              class="matrix-status-badge"
              :class="{
                'matrix-status-valid': selectedCell?.status === 'valid',
                'matrix-status-expiring': selectedCell?.status === 'expiring',
                'matrix-status-expired': selectedCell?.status === 'expired',
                'matrix-status-not-taken': selectedCell?.status === 'not_taken'
              }"
            >
              <Check v-if="selectedCell?.status === 'valid'" class="icon-sm" />
              <AlertTriangle v-else-if="selectedCell?.status === 'expiring'" class="icon-sm" />
              <X v-else-if="selectedCell?.status === 'expired'" class="icon-sm" />
              <Minus v-else-if="selectedCell?.status === 'not_taken'" class="icon-sm" />
              <span>{{ getStatusLabel(selectedCell?.status || '') }}</span>
            </div>
          </div>

          <!-- Expiry Date -->
          <div v-if="selectedCell?.expiryDate" class="matrix-detail-row">
            <Calendar class="icon-standard text-muted-foreground" />
            <div>
              <p class="text-label">Expiry Date</p>
              <p class="text-value">{{ formatDate(selectedCell.expiryDate) }}</p>
            </div>
          </div>

          <!-- Evidence Records -->
          <div v-if="courseEvidence.length > 0" class="matrix-detail-section">
            <h4 class="matrix-detail-section-title">
              <Award class="icon-standard" />
              Training Records
            </h4>
            <div class="matrix-evidence-list">
              <div v-for="record in courseEvidence" :key="record.id" class="matrix-evidence-item">
                <div class="matrix-evidence-main">
                  <p class="matrix-evidence-type">{{ record.evidenceType === 'attendance' ? 'Training Session' : 'Certification' }}</p>
                  <p class="matrix-evidence-date">Completed {{ formatDate(record.completionDate) }}</p>
                </div>
                <div v-if="record.expiryDate" class="matrix-evidence-expiry">
                  Expires {{ formatDate(record.expiryDate) }}
                </div>
              </div>
            </div>
          </div>

          <!-- No Records -->
          <div v-else-if="selectedCell?.status !== 'not_required'" class="matrix-detail-empty">
            <p>No training records found for this course.</p>
          </div>
        </div>

        <!-- Footer Actions -->
        <div class="sheet-footer-simple">
          <Button variant="outline" @click="goToEmployee">
            <ExternalLink class="icon-standard" />
            View Employee Profile
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  </section>
</template>
