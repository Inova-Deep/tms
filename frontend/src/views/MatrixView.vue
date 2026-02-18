<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useMatrix } from '@/composables/useMatrix'
import { Check, X, AlertTriangle, Minus } from 'lucide-vue-next'
import type { MatrixCell } from '@/types/api'

const { data: matrix, isLoading } = useMatrix()

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

function handleCellClick(employeeId: string, courseId: string, cell: MatrixCell | undefined) {
  if (!cell) return
  console.log('Cell clicked:', { employeeId, courseId, status: cell.status, expiryDate: cell.expiryDate })
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
                  class="matrix-cell"
                  :class="getCellClass(row.cells[course.id])"
                  @click="handleCellClick(row.employeeId, course.id, row.cells[course.id])"
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
  </section>
</template>
