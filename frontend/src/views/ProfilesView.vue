<script setup lang="ts">
import { ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
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
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '@/components/ui/sheet'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import ProfileSheetMembers from './components/ProfileSheetMembers.vue'
import ProfileSheetExclusions from './components/ProfileSheetExclusions.vue'
import { 
  useProfiles, 
  useProfile, 
  useProfileRequirements, 
  useAddProfileMember, 
  useAddProfileExclusion,
  useProfileMembers,
  useProfileExclusions,
  useRemoveProfileMember,
  useRemoveProfileExclusion,
  useCreateProfile,
  useAddProfileRequirement,
  useRemoveProfileRequirement,
} from '@/composables/useProfiles'
import type { ProfileExclusion } from '@/types/api'
import { useEmployees } from '@/composables/useEmployees'
import { useCourses } from '@/composables/useCourses'
import { MoreHorizontal, Eye, UserPlus, UserMinus, Package, Plus, Trash2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const selectedProfileId = ref<string | null>(null)

const { data: profiles, isLoading } = useProfiles()
const { data: profileDetail } = useProfile(
  () => selectedProfileId.value ?? undefined
)
const { data: profileRequirements } = useProfileRequirements(
  () => selectedProfileId.value ?? undefined
)
const { data: employees } = useEmployees()
const addMemberMutation = useAddProfileMember()
const addExclusionMutation = useAddProfileExclusion()
const removeMemberMutation = useRemoveProfileMember()
const removeExclusionMutation = useRemoveProfileExclusion()

const { data: profileMembers } = useProfileMembers(
  computed(() => selectedProfileId.value ?? undefined)
)
const { data: profileExclusions } = useProfileExclusions(
  computed(() => selectedProfileId.value ?? undefined)
)

const sheetOpen = ref(false)
const confirmDialogOpen = ref(false)
const confirmAction = ref<'member' | 'exclusion' | null>(null)
const createSheetOpen = ref(false)
const newProfileName = ref('')
const newProfileDescription = ref('')
const createProfileMutation = useCreateProfile()

const selectedEmployee = ref('')
const exclusionReason = ref('')
const exclusionDate = ref('')

const addCourseSheetOpen = ref(false)
const selectedCourseId = ref('')
const addRequirementMutation = useAddProfileRequirement()
const removeRequirementMutation = useRemoveProfileRequirement()
const { data: allCourses } = useCourses()

const isAdmin = computed(() => auth.is_admin)

function viewProfile(id: string) {
  selectedProfileId.value = id
  sheetOpen.value = true
}

function closeSheet() {
  sheetOpen.value = false
  selectedProfileId.value = null
}

function openCreateProfile() {
  newProfileName.value = ''
  newProfileDescription.value = ''
  createSheetOpen.value = true
}

async function handleCreateProfile() {
  if (!newProfileName.value.trim()) return
  await createProfileMutation.mutateAsync({
    name: newProfileName.value,
    description: newProfileDescription.value || undefined,
  })
  createSheetOpen.value = false
}

function openAddMember() {
  confirmAction.value = 'member'
  confirmDialogOpen.value = true
}

function openExclude() {
  confirmAction.value = 'exclusion'
  confirmDialogOpen.value = true
}

async function handleConfirm() {
  if (!selectedProfileId.value || !selectedEmployee.value) return

  if (confirmAction.value === 'member') {
    await addMemberMutation.mutateAsync({
      profileId: selectedProfileId.value,
      employeeId: selectedEmployee.value,
    })
  } else if (confirmAction.value === 'exclusion') {
    await addExclusionMutation.mutateAsync({
      profileId: selectedProfileId.value,
      employeeId: selectedEmployee.value,
      reason: exclusionReason.value,
      expiryDate: exclusionDate.value,
    })
  }

  confirmDialogOpen.value = false
  selectedEmployee.value = ''
  exclusionReason.value = ''
  exclusionDate.value = ''
  confirmAction.value = null
}

function cancelConfirm() {
  confirmDialogOpen.value = false
  confirmAction.value = null
  selectedEmployee.value = ''
  exclusionReason.value = ''
  exclusionDate.value = ''
}

async function handleRemoveMember(employeeId: string) {
  if (!selectedProfileId.value) return
  await removeMemberMutation.mutateAsync({
    profileId: selectedProfileId.value,
    employeeId
  })
}

async function handleRemoveExclusion(exclusionId: string) {
  if (!selectedProfileId.value) return
  await removeExclusionMutation.mutateAsync({
    profileId: selectedProfileId.value,
    exclusionId
  })
}

function handleEditExclusion(exclusion: ProfileExclusion) {
  console.log('Edit exclusion:', exclusion)
}

const availableEmployees = computed(() => {
  if (!employees?.value) return []
  return employees.value.map(e => ({ id: e.id, label: `${e.id} — ${e.name} (${e.department})` }))
})

const availableCourses = computed(() => {
  if (!allCourses?.value) return []
  const existingIds = new Set(profileRequirements.value?.map(r => r.course_id) || [])
  return allCourses.value.filter(c => !existingIds.has(c.id))
})

function openAddCourse() {
  selectedCourseId.value = ''
  addCourseSheetOpen.value = true
}

async function handleAddCourse() {
  if (!selectedProfileId.value || !selectedCourseId.value) return
  await addRequirementMutation.mutateAsync({
    profileId: selectedProfileId.value,
    courseId: selectedCourseId.value,
  })
  addCourseSheetOpen.value = false
}

async function handleRemoveRequirement(reqId: string) {
  if (!selectedProfileId.value) return
  await removeRequirementMutation.mutateAsync({
    profileId: selectedProfileId.value,
    requirementId: reqId,
  })
}
</script>

<template>
  <section class="page-section">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Profiles</h1>
        <p class="page-description">Competency profiles and their requirements.</p>
      </div>
      <Button v-if="isAdmin" @click="openCreateProfile">
        <Plus class="icon-standard mr-2" />
        Add Profile
      </Button>
    </div>

    <!-- Profiles List -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">All Profiles</CardTitle>
        <CardDescription>Available competency profiles in the system</CardDescription>
      </CardHeader>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead class="table-action-cell"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="isLoading" v-for="i in 3" :key="i">
              <TableCell><Skeleton class="h-4 w-12" /></TableCell>
              <TableCell><Skeleton class="h-4 w-32" /></TableCell>
              <TableCell><Skeleton class="h-4 w-48" /></TableCell>
              <TableCell><Skeleton class="h-8 w-8" /></TableCell>
            </TableRow>
            <TableRow v-else-if="!profiles?.length">
              <TableCell colspan="4" class="table-empty">No profiles found.</TableCell>
            </TableRow>
            <TableRow v-for="profile in profiles" :key="profile.id">
              <TableCell class="table-cell-mono">{{ profile.id }}</TableCell>
              <TableCell>
                <div class="btn-group">
                  <span class="table-cell-primary">{{ profile.name }}</span>
                  <Badge v-if="profile.memberCount" variant="secondary" class="badge-spacing">
                    {{ profile.memberCount }} members
                  </Badge>
                </div>
              </TableCell>
              <TableCell class="table-cell-secondary">{{ profile.description }}</TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon" class="btn-icon-sm">
                      <MoreHorizontal class="icon-standard" />
                      <span class="sr-only">Actions</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="viewProfile(profile.id)">
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

    <!-- Profile Detail Sheet -->
    <Sheet :open="sheetOpen" @update:open="(open) => !open && closeSheet()">
      <SheetContent class="sheet-content-wide">

        <!-- Sheet Header: Identity Block -->
        <SheetHeader class="sheet-header-identity">
          <div class="sheet-identity-row">
            <div class="sheet-identity-icon">
              <Package class="icon-standard" />
            </div>
            <SheetTitle class="sheet-identity-title">{{ profileDetail?.name ?? selectedProfileId }}</SheetTitle>
          </div>
          <SheetDescription class="sheet-identity-desc">{{ profileDetail?.description }}</SheetDescription>
          <Badge variant="outline" class="sheet-identity-badge">{{ selectedProfileId }}</Badge>
        </SheetHeader>

        <!-- Sheet Body: Scrollable Tab Content -->
        <div class="sheet-body-scroll">
          <Tabs defaultValue="requirements" class="sheet-tabs">
            <div class="sheet-tabs-header">
              <TabsList class="tabs-list-full">
                <TabsTrigger value="requirements">Requirements</TabsTrigger>
                <TabsTrigger value="members">Members</TabsTrigger>
                <TabsTrigger value="exclusions">Exclusions</TabsTrigger>
              </TabsList>
            </div>

            <!-- Requirements Tab -->
            <TabsContent value="requirements" class="sheet-tab-content">
              <div class="tab-panel">
                <div class="tab-panel-header">
                  <h3 class="tab-panel-title">Requirements</h3>
                  <Button v-if="isAdmin" size="sm" @click="openAddCourse">
                    <Plus class="icon-standard mr-2" />
                    Add Course
                  </Button>
                </div>
                <Table>
                  <TableHeader class="table-header">
                    <TableRow>
                      <TableHead class="table-header-cell">ID</TableHead>
                      <TableHead class="table-header-cell">Requirement</TableHead>
                      <TableHead class="table-header-cell">Type</TableHead>
                      <TableHead class="table-header-cell">Validity</TableHead>
                      <TableHead class="table-header-cell table-action-cell"></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-if="!profileRequirements?.length">
                      <TableCell colspan="5" class="table-empty">No requirements defined for this profile.</TableCell>
                    </TableRow>
                    <TableRow v-for="req in profileRequirements" :key="req.id" class="table-row">
                      <TableCell class="table-cell table-cell-mono">{{ req.id }}</TableCell>
                      <TableCell class="table-cell table-cell-primary">{{ req.name }}</TableCell>
                      <TableCell class="table-cell">
                        <Badge :variant="req.type === 'course' ? 'default' : 'secondary'" class="badge-text-sm">
                          {{ req.type }}
                        </Badge>
                      </TableCell>
                      <TableCell class="table-cell table-cell-secondary">{{ req.validity_months }} months</TableCell>
                      <TableCell class="table-cell">
                        <Button v-if="isAdmin" variant="ghost" size="icon" class="btn-icon-sm" @click="handleRemoveRequirement(req.id)">
                          <Trash2 class="icon-standard" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </TabsContent>

            <!-- Members Tab -->
            <TabsContent value="members" class="sheet-tab-content">
              <div class="tab-panel">
                <ProfileSheetMembers
                  :profile-id="selectedProfileId ?? ''"
                  :members="profileMembers?.members ?? []"
                  @add-member="openAddMember"
                  @remove-member="handleRemoveMember"
                />
              </div>
            </TabsContent>

            <!-- Exclusions Tab -->
            <TabsContent value="exclusions" class="sheet-tab-content">
              <div class="tab-panel">
                <ProfileSheetExclusions
                  :profile-id="selectedProfileId ?? ''"
                  :exclusions="profileExclusions ?? []"
                  @add-exclusion="openExclude"
                  @remove-exclusion="handleRemoveExclusion"
                  @edit-exclusion="handleEditExclusion"
                />
              </div>
            </TabsContent>
          </Tabs>
        </div>

        <!-- Sheet Footer: Fixed Actions -->
        <SheetFooter v-if="isAdmin" class="sheet-footer-fixed">
          <div class="sheet-footer-actions">
            <Button variant="default" class="sheet-footer-btn" @click="openAddMember">
              <UserPlus class="icon-standard mr-2" />
              Add Member
            </Button>
            <Button variant="outline" class="sheet-footer-btn" @click="openExclude">
              <UserMinus class="icon-standard mr-2" />
              Exclude Employee
            </Button>
          </div>
        </SheetFooter>

      </SheetContent>
    </Sheet>

    <!-- Create Profile Sheet -->
    <Sheet :open="createSheetOpen" @update:open="(open) => !open && (createSheetOpen = false)">
      <SheetContent class="sheet-content">
        <SheetHeader class="sheet-header">
          <SheetTitle>Create Profile</SheetTitle>
          <SheetDescription>Add a new competency profile.</SheetDescription>
        </SheetHeader>
        <div class="sheet-body">
          <div class="form-grid">
            <div class="form-row-simple">
              <Label>Name *</Label>
              <Input v-model="newProfileName" placeholder="e.g., Safety Advanced" />
            </div>
            <div class="form-row-simple">
              <Label>Description</Label>
              <Textarea v-model="newProfileDescription" placeholder="Profile description..." />
            </div>
          </div>
        </div>
        <SheetFooter>
          <div class="sheet-footer-actions">
            <Button variant="outline" class="sheet-footer-btn" @click="createSheetOpen = false">Cancel</Button>
            <Button class="sheet-footer-btn" @click="handleCreateProfile" :disabled="!newProfileName.trim() || createProfileMutation.isPending">
              Create
            </Button>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>

    <!-- Add Member / Exclude Dialog -->
    <ConfirmDialog
      v-model:open="confirmDialogOpen"
      :title="confirmAction === 'member' ? 'Add Employee to Profile' : 'Exclude Employee from Profile'"
      :description="confirmAction === 'member' 
        ? 'This will assign the employee to this profile.' 
        : 'This will exclude the employee from this profile with a reason.'"
      :confirm-label="confirmAction === 'member' ? 'Add' : 'Exclude'"
      variant="destructive"
      @confirm="handleConfirm"
      @cancel="cancelConfirm"
    >
      <div class="form-grid">
        <div class="form-row-simple">
          <Label>Select Employee</Label>
          <Select v-model="selectedEmployee">
            <SelectTrigger>
              <SelectValue placeholder="Choose an employee..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="emp in availableEmployees" :key="emp.id" :value="emp.id">
                {{ emp.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <template v-if="confirmAction === 'exclusion'">
          <div class="form-row-simple">
            <Label>Reason for Exclusion</Label>
            <Input v-model="exclusionReason" placeholder="e.g., Medical exemption" />
          </div>
          <div class="form-row-simple">
            <Label>Expiry Date</Label>
            <Input v-model="exclusionDate" type="date" />
          </div>
        </template>
      </div>
    </ConfirmDialog>

    <!-- Add Course to Profile Sheet -->
    <Sheet :open="addCourseSheetOpen" @update:open="(open) => !open && (addCourseSheetOpen = false)">
      <SheetContent class="sheet-content">
        <SheetHeader class="sheet-header">
          <SheetTitle>Add Course to Profile</SheetTitle>
          <SheetDescription>Select a course to add as a requirement.</SheetDescription>
        </SheetHeader>
        <div class="sheet-body">
          <div class="form-row-simple">
            <Label>Select Course</Label>
            <Select v-model="selectedCourseId">
              <SelectTrigger>
                <SelectValue placeholder="Choose a course..." />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="course in availableCourses" :key="course.id" :value="course.id">
                  {{ course.name }} ({{ course.validity_months }}mo)
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <p v-if="!availableCourses.length" class="page-meta">All courses are already added to this profile.</p>
        </div>
        <SheetFooter>
          <div class="sheet-footer-actions">
            <Button variant="outline" class="sheet-footer-btn" @click="addCourseSheetOpen = false">Cancel</Button>
            <Button class="sheet-footer-btn" @click="handleAddCourse" :disabled="!selectedCourseId || addRequirementMutation.isPending">
              Add
            </Button>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  </section>
</template>
