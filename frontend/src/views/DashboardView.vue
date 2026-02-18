<script setup lang="ts">
import { ref, computed } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import Select from '@/components/ui/Select.vue'
import { useDashboardKpis, useDashboardGrid, useDashboardDrilldown, type GridRow } from '@/composables/useDashboard'
import { statusToBadge } from '@/lib/status'

// Filter state
const mode = ref<'course' | 'profile'>('course')
const groupBy = ref<'department' | 'site' | 'costCenter'>('department')
const filterDept = ref('')
const filterSite = ref('')

// Drill-down state
const drilldownOpen = ref(false)
const selectedGroupValue = ref<string | null>(null)
const selectedItemId = ref<string | null>(null)

// Fetch data from API
const kpisQuery = useDashboardKpis(filterDept, filterSite)
const gridQuery = useDashboardGrid(mode, groupBy, filterDept, filterSite)
const drilldownQuery = useDashboardDrilldown(
    mode,
    groupBy,
    selectedGroupValue,
    selectedItemId,
    filterDept,
    filterSite,
)

// Computed data
const kpis = computed(() => kpisQuery.data.value)
const gridRows = computed(() => gridQuery.data.value ?? [])
const drilldownRows = computed(() => drilldownQuery.data.value ?? [])

// Progress bar width class
function getProgressClass(percent: number): string {
  if (percent === 100) return 'progress-w-full'
  if (percent >= 90) return 'progress-w-11-12'
  if (percent >= 80) return 'progress-w-5-6'
  if (percent >= 70) return 'progress-w-3-4'
  if (percent >= 60) return 'progress-w-2-3'
  if (percent >= 50) return 'progress-w-1-2'
  if (percent >= 30) return 'progress-w-1-3'
  return 'progress-w-1-4'
}

function openDrill(row: any) {
    selectedGroupValue.value = row.groupValue
    selectedItemId.value = row.itemId
    drilldownOpen.value = true
}
</script>

<template>
    <section class="page-section">
        <!-- Page Header with Filters -->
        <Card>
            <CardHeader>
                <div class="card-row-compact">
                    <div>
                        <CardTitle class="page-title">Coverage Dashboard</CardTitle>
                        <p class="page-description">Expiring counts as compliant · UK dates</p>
                    </div>

                    <div class="btn-group">
                        <Select v-model="mode" :options="[
                            { label: 'Course mode', value: 'course' },
                            { label: 'Profile mode', value: 'profile' },
                        ]" class="w-40" />
                        <Select v-model="groupBy" :options="[
                            { label: 'Department', value: 'department' },
                            { label: 'Site', value: 'site' },
                            { label: 'Cost Center', value: 'costCenter' },
                        ]" class="w-40" />
                        <Select v-model="filterDept" :options="[
                            { label: 'All Departments', value: '' },
                            { label: 'Manufacturing', value: 'Manufacturing' },
                            { label: 'Maintenance', value: 'Maintenance' },
                            { label: 'Quality', value: 'Quality' },
                        ]" class="w-44" />
                        <Select v-model="filterSite" :options="[
                            { label: 'All Sites', value: '' },
                            { label: 'Bristol Central', value: 'Bristol Central Plant' },
                            { label: 'Avonmouth', value: 'Avonmouth Logistics & Workshop' },
                        ]" class="w-44" />
                    </div>
                </div>
            </CardHeader>
        </Card>

        <!-- KPI Tiles -->
        <div class="kpi-grid">
            <Card class="kpi-card">
                <p class="kpi-label">Overall compliance</p>
                <p class="kpi-value">
                    {{ kpis?.overallCompliance?.toFixed(1) ?? '—' }}%
                </p>
            </Card>
            <Card class="kpi-card">
                <p class="kpi-label">Expiring soon</p>
                <p class="kpi-value-warning">{{ kpis?.expiringSoon ?? '—' }}</p>
            </Card>
            <Card class="kpi-card">
                <p class="kpi-label">Expired</p>
                <p class="kpi-value-danger">{{ kpis?.expired ?? '—' }}</p>
            </Card>
            <Card class="kpi-card">
                <p class="kpi-label">Employees</p>
                <p class="kpi-value">{{ kpis?.totalEmployees ?? '—' }}</p>
            </Card>
        </div>

        <!-- Coverage Grid -->
        <Card>
            <CardContent class="p-0">
                <div class="table-container">
                    <Table>
                        <TableHeader class="table-header">
                            <TableRow>
                                <TableHead class="table-header-cell">
                                    {{ mode === 'course' ? 'Course' : 'Profile' }}
                                </TableHead>
                                <TableHead class="table-header-cell">
                                    {{ groupBy === 'department' ? 'Department' : groupBy === 'site' ? 'Site' : 'Cost Center' }}
                                </TableHead>
                                <TableHead class="table-header-cell">% Compliance</TableHead>
                                <TableHead class="table-header-cell">Valid</TableHead>
                                <TableHead class="table-header-cell">Expiring</TableHead>
                                <TableHead class="table-header-cell">Expired</TableHead>
                                <TableHead class="table-header-cell">Missing</TableHead>
                                <TableHead class="table-header-cell"></TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            <TableRow v-for="(row, idx) in gridRows" :key="idx" class="table-row">
                                <TableCell class="table-cell">{{ row.itemName }}</TableCell>
                                <TableCell class="table-cell">
                                    {{ row.groupValue }} 
                                    <span class="page-meta">({{ row.itemId }})</span>
                                </TableCell>
                                <TableCell class="table-cell">
                                    <div class="progress-bar-container">
                                        <div class="progress-bar">
                                            <div class="progress-bar-fill" :class="getProgressClass(row.percent)" />
                                        </div>
                                        <span class="progress-bar-label">{{ row.percent.toFixed(2) }}%</span>
                                    </div>
                                </TableCell>
                                <TableCell class="table-cell">{{ row.valid }}</TableCell>
                                <TableCell class="table-cell table-cell-warning">{{ row.expiring }}</TableCell>
                                <TableCell class="table-cell table-cell-danger">{{ row.expired }}</TableCell>
                                <TableCell class="table-cell">{{ row.missing }}</TableCell>
                                <TableCell class="table-cell text-right">
                                    <Button variant="ghost" size="sm" @click="openDrill(row)">Drill down</Button>
                                </TableCell>
                            </TableRow>
                        </TableBody>
                    </Table>
                </div>
            </CardContent>
        </Card>

        <!-- Drill-down Sheet -->
        <Sheet v-model:open="drilldownOpen">
            <SheetContent class="sheet-content-wide">
                <SheetHeader class="sheet-header-section">
                    <SheetTitle>Employee Details</SheetTitle>
                    <SheetDescription v-if="selectedGroupValue && selectedItemId">
                        {{(gridRows as GridRow[]).find((r: GridRow) => r.groupValue === selectedGroupValue && r.itemId === selectedItemId)?.itemName ?? '' }} · {{ selectedGroupValue }}
                    </SheetDescription>
                </SheetHeader>

                <div class="sheet-body">
                    <div v-if="drilldownQuery.isFetching.value" class="loading-container">
                        <span class="table-cell-muted">Loading employee details...</span>
                    </div>
                    <div v-else-if="drilldownRows.length === 0" class="loading-container">
                        <div class="text-center">
                            <span class="table-cell-muted">No employees found for this selection.</span>
                            <span class="page-meta block mt-1">Debug: Group='{{ selectedGroupValue }}', Item='{{ selectedItemId }}'</span>
                        </div>
                    </div>
                    <div v-else class="table-container">
                        <Table>
                            <TableHeader class="table-header-alt">
                                <TableRow>
                                    <TableHead class="table-header-cell-alt">Employee</TableHead>
                                    <TableHead class="table-header-cell-alt">Department</TableHead>
                                    <TableHead class="table-header-cell-alt">Site</TableHead>
                                    <TableHead class="table-header-cell-alt">Status</TableHead>
                                    <TableHead class="table-header-cell-alt">Expiry Date</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                <TableRow v-for="emp in drilldownRows" :key="emp.id" class="table-row-alt">
                                    <TableCell class="table-cell-alt">
                                        <div class="flex flex-col">
                                            <span class="table-cell-primary">{{ emp.name }}</span>
                                            <span class="table-cell-mono">{{ emp.id }}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell class="table-cell-alt table-cell-secondary">{{ emp.department }}</TableCell>
                                    <TableCell class="table-cell-alt table-cell-secondary">{{ emp.site }}</TableCell>
                                    <TableCell class="table-cell-alt">
                                        <Badge :variant="statusToBadge(emp.status)">{{ emp.status }}</Badge>
                                    </TableCell>
                                    <TableCell class="table-cell-alt table-cell-secondary">
                                        {{ emp.expiryDate && emp.expiryDate !== '0001-01-01' ? emp.expiryDate : '—' }}
                                    </TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </div>
                </div>
            </SheetContent>
        </Sheet>

    </section>
</template>
