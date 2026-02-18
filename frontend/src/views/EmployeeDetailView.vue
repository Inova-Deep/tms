<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from '@/components/ui/table'
import { useEmployeeDetail, useEmployeeRecords } from '@/composables/useEmployeeDetail'
import { statusToBadge } from '@/lib/status'
import { formatDate } from '@/lib/date'
import { 
  ArrowLeft, 
  Building, 
  MapPin, 
  Briefcase, 
  CreditCard, 
  CheckCircle2, 
  AlertCircle, 
  XCircle, 
  HelpCircle,
  Clock
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const id = computed(() => route.params.id as string)

const detailQuery = useEmployeeDetail(id.value)
const recordsQuery = useEmployeeRecords(id.value)
const detail = computed(() => detailQuery.data.value)
const records = computed(() => recordsQuery.data.value ?? [])

const isLoading = computed(() => detailQuery.isLoading?.value)

// Compliance Stats Calculation
const complianceStats = computed(() => {
  if (!detail.value?.requirements) return { percent: 0, valid: 0, expiring: 0, expired: 0, missing: 0 }
  
  const reqs = detail.value.requirements
  const counts = {
    valid: reqs.filter(r => r.status === 'Valid').length,
    expiring: reqs.filter(r => r.status === 'Expiring').length,
    expired: reqs.filter(r => r.status === 'Expired').length,
    missing: reqs.filter(r => r.status === 'Missing').length
  }
  
  const total = reqs.length
  const percent = total > 0 ? Math.round(((counts.valid + counts.expiring) / total) * 100) : 0
  
  return { ...counts, percent, total }
})

const getComplianceColor = (percent: number) => {
  if (percent >= 90) return 'compliance-high'
  if (percent >= 70) return 'compliance-medium'
  return 'compliance-low'
}
</script>

<template>
  <section class="page-section">
    <!-- Breadcrumbs / Back navigation -->
    <div class="page-header">
      <div class="page-header-left">
        <Button variant="ghost" size="sm" class="back-btn" @click="router.push('/employees')">
          <ArrowLeft class="icon-standard mr-2" />
          Back to employees
        </Button>
        <div class="title-row">
          <h1 class="page-title">{{ detail?.name ?? 'Employee Details' }}</h1>
          <Badge v-if="detail?.employment_status" :variant="detail.employment_status === 'Active' ? 'success' : 'secondary'" class="status-badge">
            {{ detail.employment_status }}
          </Badge>
        </div>
        <p class="page-meta">{{ detail?.id ?? id }} · {{ detail?.worker_type }}</p>
      </div>
    </div>

    <!-- Dashboard Grid Layout -->
    <div class="employee-detail-grid">
      
      <!-- Sidebar Column (Compliance and Info) -->
      <aside class="employee-detail-sidebar">
        
        <!-- Compliance Summary KPI -->
        <Card>
          <div class="compliance-status-card">
            <div class="compliance-kpi-value" :class="getComplianceColor(complianceStats.percent)">
              {{ complianceStats.percent }}%
            </div>
            <p class="compliance-kpi-label">Overall Compliance</p>
          </div>
          <CardContent class="p-0">
            <div class="status-summary-list">
              <div class="status-summary-item px-4">
                <div class="flex items-center gap-2">
                  <CheckCircle2 class="icon-standard text-status-valid" />
                  <span>Valid</span>
                </div>
                <span class="status-summary-count">{{ complianceStats.valid }}</span>
              </div>
              <div class="status-summary-item px-4">
                <div class="flex items-center gap-2">
                  <Clock class="icon-standard text-status-expiring" />
                  <span>Expiring</span>
                </div>
                <span class="status-summary-count text-status-expiring">{{ complianceStats.expiring }}</span>
              </div>
              <div class="status-summary-item px-4">
                <div class="flex items-center gap-2">
                  <AlertCircle class="icon-standard text-status-expired" />
                  <span>Expired</span>
                </div>
                <span class="status-summary-count text-status-expired">{{ complianceStats.expired }}</span>
              </div>
              <div class="status-summary-item px-4 border-none">
                <div class="flex items-center gap-2">
                  <HelpCircle class="icon-standard text-status-missing" />
                  <span>Missing</span>
                </div>
                <span class="status-summary-count">{{ complianceStats.missing }}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Employee Metadata Card -->
        <Card>
          <CardHeader class="card-header-compact">
            <CardTitle class="form-label">Key Information</CardTitle>
          </CardHeader>
          <CardContent class="card-content-spaced">
            <div class="meta-list">
              <div class="meta-list-item">
                <Building class="meta-list-icon" />
                <span>{{ detail?.department }}</span>
              </div>
              <div class="meta-list-item">
                <MapPin class="meta-list-icon" />
                <span>{{ detail?.site }}</span>
              </div>
              <div class="meta-list-item">
                <CreditCard class="meta-list-icon" />
                <span>Cost Center: {{ detail?.cost_center }}</span>
              </div>
            </div>

            <Separator />

            <div>
              <p class="form-label mb-2">Assigned Profiles</p>
              <div class="badge-list">
                <Badge
                  v-for="(status, profileId) in detail?.profile_status"
                  :key="profileId"
                  variant="outline"
                  class="badge-text-sm"
                >
                  {{ profileId }}
                </Badge>
                <p v-if="!detail?.profile_status" class="page-meta">No profiles assigned.</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </aside>

      <!-- Main Content Column -->
      <main class="employee-detail-main">
        <Tabs defaultValue="requirements" class="tabs-full-width">
          <TabsList class="tabs-list-full">
            <TabsTrigger value="requirements">Current Requirements</TabsTrigger>
            <TabsTrigger value="transcript">Training Transcript</TabsTrigger>
          </TabsList>

          <!-- Requirements Tab -->
          <TabsContent value="requirements" class="tabs-content">
            <Card>
              <CardContent class="p-0">
                <Table>
                  <TableHeader class="table-header">
                    <TableRow>
                      <TableHead class="table-header-cell">Requirement</TableHead>
                      <TableHead class="table-header-cell">Profile</TableHead>
                      <TableHead class="table-header-cell">Type</TableHead>
                      <TableHead class="table-header-cell">Status</TableHead>
                      <TableHead class="table-header-cell">Expiry</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-if="isLoading" v-for="i in 3" :key="i">
                      <TableCell><Skeleton class="skeleton-xl" /></TableCell>
                      <TableCell><Skeleton class="skeleton-md" /></TableCell>
                      <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                      <TableCell><Skeleton class="skeleton-badge" /></TableCell>
                      <TableCell><Skeleton class="skeleton-md" /></TableCell>
                    </TableRow>
                    <TableRow v-else-if="!detail?.requirements?.length">
                      <TableCell colspan="5" class="table-empty">No active requirements.</TableCell>
                    </TableRow>
                    <TableRow v-for="req in detail?.requirements" :key="req.id" class="table-row">
                      <TableCell class="table-cell table-cell-primary">{{ req.name }}</TableCell>
                      <TableCell class="table-cell table-cell-muted">{{ req.profile_id }}</TableCell>
                      <TableCell class="table-cell">
                        <Badge variant="outline" class="badge-text-sm">{{ req.type }}</Badge>
                      </TableCell>
                      <TableCell class="table-cell">
                        <Badge :variant="statusToBadge(req.status)" class="status-badge">
                          {{ req.status }}
                        </Badge>
                      </TableCell>
                      <TableCell class="table-cell table-cell-mono">
                        {{ req.expiry_date ? formatDate(req.expiry_date) : '—' }}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>

          <!-- Transcript Tab -->
          <TabsContent value="transcript" class="tabs-content">
            <Card>
              <CardContent class="p-0">
                <Table>
                  <TableHeader class="table-header">
                    <TableRow>
                      <TableHead class="table-header-cell">Requirement</TableHead>
                      <TableHead class="table-header-cell">Evidence</TableHead>
                      <TableHead class="table-header-cell">Date</TableHead>
                      <TableHead class="table-header-cell">Expires</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-if="recordsQuery.isLoading?.value" v-for="i in 3" :key="i">
                      <TableCell><Skeleton class="skeleton-xl" /></TableCell>
                      <TableCell><Skeleton class="skeleton-md" /></TableCell>
                      <TableCell><Skeleton class="skeleton-md" /></TableCell>
                      <TableCell><Skeleton class="skeleton-md" /></TableCell>
                    </TableRow>
                    <TableRow v-else-if="records.length === 0">
                      <TableCell colspan="4" class="table-empty">No training records found.</TableCell>
                    </TableRow>
                    <TableRow v-for="rec in records" :key="rec.id" class="table-row">
                      <TableCell class="table-cell table-cell-primary">{{ rec.requirementName }}</TableCell>
                      <TableCell class="table-cell">
                        <Badge variant="outline" class="badge-text-sm">{{ rec.evidenceType }}</Badge>
                      </TableCell>
                      <TableCell class="table-cell table-cell-secondary">{{ formatDate(rec.completionDate) }}</TableCell>
                      <TableCell class="table-cell table-cell-mono">
                        {{ rec.expiryDate ? formatDate(rec.expiryDate) : 'Never' }}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>
    </div>
  </section>
</template>
