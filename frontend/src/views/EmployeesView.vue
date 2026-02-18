<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useEmployees } from '@/composables/useEmployees'
import { ChevronDown, Search, Filter, MoreHorizontal, Eye } from 'lucide-vue-next'
import type { Employee } from '@/types/api'

const router = useRouter()
const search = ref('')
const departmentFilter = ref<string | null>(null)
const siteFilter = ref<string | null>(null)

const { data: employees, isLoading } = useEmployees()

const departments = computed(() => {
  const set = new Set((employees?.value ?? []).map((e) => e.department))
  return Array.from(set).sort()
})

const sites = computed(() => {
  const set = new Set((employees?.value ?? []).map((e) => e.site))
  return Array.from(set).sort()
})

const filtered = computed(() => {
  const list = employees?.value ?? []
  return list.filter((e: Employee) => {
    const matchesSearch =
      !search.value ||
      e.name.toLowerCase().includes(search.value.toLowerCase()) ||
      e.id.toLowerCase().includes(search.value.toLowerCase())
    const matchesDept = !departmentFilter.value || e.department === departmentFilter.value
    const matchesSite = !siteFilter.value || e.site === siteFilter.value
    return matchesSearch && matchesDept && matchesSite
  })
})

function goToDetail(id: string) {
  router.push(`/employees/${id}`)
}

function clearFilters() {
  departmentFilter.value = null
  siteFilter.value = null
}

const hasActiveFilters = computed(() => departmentFilter.value || siteFilter.value)
</script>

<template>
  <section class="page-section">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Employees</h1>
        <p class="page-description">Search and view employee profiles and training records.</p>
      </div>
      <div class="page-actions">
        <span v-if="hasActiveFilters" class="page-meta">
          {{ filtered.length }} of {{ employees?.length ?? 0 }} employees
        </span>
        <Button v-if="hasActiveFilters" variant="ghost" size="sm" @click="clearFilters">
          Clear filters
        </Button>
      </div>
    </div>

    <Card>
      <!-- Filters -->
      <div class="filter-bar">
        <div class="filter-row">
          <!-- Search -->
          <div class="filter-search-wrapper">
            <Search class="filter-search-icon" />
            <Input v-model="search" placeholder="Search name or ID..." class="filter-input" />
          </div>
          
          <!-- Department Filter -->
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="outline" size="sm">
                <Filter class="icon-standard mr-2" />
                {{ departmentFilter || 'Department' }}
                <ChevronDown class="icon-standard ml-2" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" class="dropdown-width-sm">
              <DropdownMenuItem @click="departmentFilter = null">
                All departments
              </DropdownMenuItem>
              <DropdownMenuItem
                v-for="dept in departments"
                :key="dept"
                @click="departmentFilter = dept"
              >
                {{ dept }}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <!-- Site Filter -->
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="outline" size="sm">
                {{ siteFilter || 'Site' }}
                <ChevronDown class="icon-standard ml-2" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" class="dropdown-width-md">
              <DropdownMenuItem @click="siteFilter = null">
                All sites
              </DropdownMenuItem>
              <DropdownMenuItem
                v-for="site in sites"
                :key="site"
                @click="siteFilter = site"
              >
                {{ site }}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <!-- Table -->
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Department</TableHead>
              <TableHead class="table-cell-hidden-lg">Site</TableHead>
              <TableHead class="table-cell-hidden-md">Worker</TableHead>
              <TableHead class="table-cell-hidden-xl">Cost Ctr</TableHead>
              <TableHead>Status</TableHead>
              <TableHead class="table-action-cell"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <!-- Loading skeletons -->
            <template v-if="isLoading">
              <TableRow v-for="i in 5" :key="i">
                <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                <TableCell><Skeleton class="skeleton-lg" /></TableCell>
                <TableCell><Skeleton class="skeleton-md" /></TableCell>
                <TableCell class="table-cell-hidden-lg"><Skeleton class="skeleton-lg" /></TableCell>
                <TableCell class="table-cell-hidden-md"><Skeleton class="skeleton-md" /></TableCell>
                <TableCell class="table-cell-hidden-xl"><Skeleton class="skeleton-sm" /></TableCell>
                <TableCell><Skeleton class="skeleton-badge" /></TableCell>
                <TableCell><Skeleton class="skeleton-icon" /></TableCell>
              </TableRow>
            </template>
            
            <!-- Empty state -->
            <TableRow v-else-if="filtered.length === 0">
              <TableCell colspan="8" class="table-empty-cell">
                <p class="table-empty">No employees found.</p>
                <Button v-if="hasActiveFilters" variant="link" size="sm" class="mt-2" @click="clearFilters">
                  Clear filters
                </Button>
              </TableCell>
            </TableRow>

            <!-- Data rows -->
            <TableRow v-for="emp in filtered" :key="emp.id">
              <TableCell class="table-cell-mono">{{ emp.id }}</TableCell>
              <TableCell class="table-cell-primary">{{ emp.name }}</TableCell>
              <TableCell class="table-cell-secondary">{{ emp.department }}</TableCell>
              <TableCell class="table-cell-hidden-lg table-cell-secondary">{{ emp.site }}</TableCell>
              <TableCell class="table-cell-hidden-md table-cell-secondary">{{ emp.worker_type }}</TableCell>
              <TableCell class="table-cell-hidden-xl table-cell-secondary">{{ emp.cost_center }}</TableCell>
              <TableCell>
                <Badge :variant="emp.employment_status === 'Active' ? 'success' : 'secondary'" class="status-badge">
                  {{ emp.employment_status }}
                </Badge>
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon" class="btn-icon-sm">
                      <MoreHorizontal class="icon-standard" />
                      <span class="sr-only">Actions</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="goToDetail(emp.id)">
                      <Eye class="icon-standard mr-2" />
                      View details
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </section>
</template>
