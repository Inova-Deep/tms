<script setup lang="ts">
import { ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
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
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import CourseFormSheet from '@/components/domain/CourseFormSheet.vue'
import {
  useCourses,
  useCreateCourse,
  useUpdateCourse,
  useDeleteCourse,
} from '@/composables/useCourses'
import type { Course, CreateCourseRequest, UpdateCourseRequest } from '@/types/api'
import { Plus, Search, MoreHorizontal, Pencil, Trash2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const activeTab = ref<'course' | 'certification'>('course')
const search = ref('')
const formSheetOpen = ref(false)
const deleteDialogOpen = ref(false)
const selectedCourse = ref<Course | null>(null)

const courseType = computed(() => activeTab.value)

const { data: courses, isLoading } = useCourses(courseType)
const createMutation = useCreateCourse()
const updateMutation = useUpdateCourse()
const deleteMutation = useDeleteCourse()

const filtered = computed(() => {
  const list = courses?.value ?? []
  return list.filter((c: Course) => {
    const matchesSearch =
      !search.value ||
      c.name.toLowerCase().includes(search.value.toLowerCase()) ||
      (c.code?.toLowerCase().includes(search.value.toLowerCase()) ?? false) ||
      c.id.toLowerCase().includes(search.value.toLowerCase())
    return matchesSearch
  })
})

const isAdmin = computed(() => auth.is_admin)

function openCreate() {
  selectedCourse.value = null
  formSheetOpen.value = true
}

function openEdit(course: Course) {
  selectedCourse.value = course
  formSheetOpen.value = true
}

function openDelete(course: Course) {
  selectedCourse.value = course
  deleteDialogOpen.value = true
}

async function handleFormSubmit(data: CreateCourseRequest | UpdateCourseRequest) {
  if (selectedCourse.value) {
    await updateMutation.mutateAsync({
      id: selectedCourse.value.id,
      data: data as UpdateCourseRequest,
    })
  } else {
    await createMutation.mutateAsync(data as CreateCourseRequest)
  }
  formSheetOpen.value = false
  selectedCourse.value = null
}

async function handleDelete() {
  if (!selectedCourse.value) return
  await deleteMutation.mutateAsync(selectedCourse.value.id)
  deleteDialogOpen.value = false
  selectedCourse.value = null
}
</script>

<template>
  <section class="page-section">
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Courses</h1>
        <p class="page-description">Manage courses and certifications for training requirements.</p>
      </div>
      <div class="page-actions" v-if="isAdmin">
        <Button @click="openCreate">
          <Plus class="icon-standard mr-2" />
          Add Course
        </Button>
      </div>
    </div>

    <Card>
      <Tabs v-model="activeTab" class="tabs-full-width">
        <CardHeader class="card-header-compact">
          <div class="tabs-header-row">
            <CardTitle class="form-label">Course Library</CardTitle>
            <TabsList class="tabs-list-full">
              <TabsTrigger value="course">Courses</TabsTrigger>
              <TabsTrigger value="certification">Certifications</TabsTrigger>
            </TabsList>
          </div>
          <CardDescription>
            {{ activeTab === 'course' ? 'Training courses with validity periods' : 'Certification requirements' }}
          </CardDescription>
        </CardHeader>

        <TabsContent value="course" class="tab-content-no-top">
          <div class="filter-bar">
            <div class="filter-row">
              <div class="filter-search-wrapper">
                <Search class="filter-search-icon" />
                <Input v-model="search" placeholder="Search courses..." class="filter-input" />
              </div>
            </div>
          </div>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Code</TableHead>
                  <TableHead>Validity</TableHead>
                  <TableHead class="table-action-cell"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <template v-if="isLoading">
                  <TableRow v-for="i in 5" :key="i">
                    <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                    <TableCell><Skeleton class="skeleton-lg" /></TableCell>
                    <TableCell><Skeleton class="skeleton-md" /></TableCell>
                    <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                    <TableCell><Skeleton class="skeleton-icon" /></TableCell>
                  </TableRow>
                </template>
                <TableRow v-else-if="filtered.length === 0">
                  <TableCell colspan="5" class="table-empty-cell">
                    <p class="table-empty">No courses found.</p>
                  </TableCell>
                </TableRow>
                <TableRow v-for="course in filtered" :key="course.id">
                  <TableCell class="table-cell-mono">{{ course.id }}</TableCell>
                  <TableCell>
                    <div class="btn-group">
                      <span class="table-cell-primary">{{ course.name }}</span>
                    </div>
                  </TableCell>
                  <TableCell class="table-cell-secondary">{{ course.code || '—' }}</TableCell>
                  <TableCell class="table-cell-secondary">{{ course.validity_months }} months</TableCell>
                  <TableCell>
                    <DropdownMenu v-if="isAdmin">
                      <DropdownMenuTrigger as-child>
                        <Button variant="ghost" size="icon" class="btn-icon-sm">
                          <MoreHorizontal class="icon-standard" />
                          <span class="sr-only">Actions</span>
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem @click="openEdit(course)">
                          <Pencil class="icon-standard mr-2" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem class="btn-danger-ghost" @click="openDelete(course)">
                          <Trash2 class="icon-standard mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </TabsContent>

        <TabsContent value="certification" class="tab-content-no-top">
          <div class="filter-bar">
            <div class="filter-row">
              <div class="filter-search-wrapper">
                <Search class="filter-search-icon" />
                <Input v-model="search" placeholder="Search certifications..." class="filter-input" />
              </div>
            </div>
          </div>
          <CardContent class="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Code</TableHead>
                  <TableHead>Validity</TableHead>
                  <TableHead class="table-action-cell"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <template v-if="isLoading">
                  <TableRow v-for="i in 5" :key="i">
                    <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                    <TableCell><Skeleton class="skeleton-lg" /></TableCell>
                    <TableCell><Skeleton class="skeleton-md" /></TableCell>
                    <TableCell><Skeleton class="skeleton-sm" /></TableCell>
                    <TableCell><Skeleton class="skeleton-icon" /></TableCell>
                  </TableRow>
                </template>
                <TableRow v-else-if="filtered.length === 0">
                  <TableCell colspan="5" class="table-empty-cell">
                    <p class="table-empty">No certifications found.</p>
                  </TableCell>
                </TableRow>
                <TableRow v-for="course in filtered" :key="course.id">
                  <TableCell class="table-cell-mono">{{ course.id }}</TableCell>
                  <TableCell>
                    <div class="btn-group">
                      <span class="table-cell-primary">{{ course.name }}</span>
                    </div>
                  </TableCell>
                  <TableCell class="table-cell-secondary">{{ course.code || '—' }}</TableCell>
                  <TableCell class="table-cell-secondary">{{ course.validity_months }} months</TableCell>
                  <TableCell>
                    <DropdownMenu v-if="isAdmin">
                      <DropdownMenuTrigger as-child>
                        <Button variant="ghost" size="icon" class="btn-icon-sm">
                          <MoreHorizontal class="icon-standard" />
                          <span class="sr-only">Actions</span>
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem @click="openEdit(course)">
                          <Pencil class="icon-standard mr-2" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem class="btn-danger-ghost" @click="openDelete(course)">
                          <Trash2 class="icon-standard mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </TabsContent>
      </Tabs>
    </Card>

    <CourseFormSheet
      v-model:open="formSheetOpen"
      :course="selectedCourse"
      @submit="handleFormSubmit"
    />

    <ConfirmDialog
      v-model:open="deleteDialogOpen"
      title="Delete Course"
      description="Are you sure you want to delete this course? This action cannot be undone."
      confirm-label="Delete"
      variant="destructive"
      @confirm="handleDelete"
    />
  </section>
</template>
