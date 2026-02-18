<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Tabs,
  TabsList,
  TabsTrigger,
  TabsContent,
} from '@/components/ui/tabs'
import DateTimePicker from '@/components/ui/DateTimePicker.vue'
import { useUpdateEvent } from '@/composables/useEvents'
import { useEmployees } from '@/composables/useEmployees'
import { useCourses } from '@/composables/useCourses'
import type { TrainingEvent } from '@/types/api'
import { Users, UserCheck, UserX } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  session: TrainingEvent | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const isOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const { mutate: updateSession, isPending } = useUpdateEvent()
const { data: employees } = useEmployees()
const { data: courses } = useCourses('course')

const durationOptions = [
  { label: '30 minutes', value: 30 },
  { label: '1 hour', value: 60 },
  { label: '1.5 hours', value: 90 },
  { label: '2 hours', value: 120 },
  { label: '3 hours', value: 180 },
  { label: 'Half Day', value: 240 },
  { label: 'Full Day', value: 480 },
]

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Scheduled', value: 'scheduled' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Completed', value: 'completed' },
  { label: 'Cancelled', value: 'cancelled' },
]

const form = ref({
  courseId: props.session?.course_id ?? '',
  date: props.session?.date ?? '',
  time: props.session?.time ?? '09:00',
  duration: props.session?.duration ?? 120,
  location: props.session?.location ?? '',
  instructorId: props.session?.instructor_id ?? '',
  status: (props.session?.status ?? 'scheduled') as TrainingEvent['status'],
  attendees: props.session?.attendees ? [...props.session.attendees] : [] as string[],
})

// Update form when session prop changes
watch(() => props.session, (newSession) => {
  if (newSession) {
    form.value = {
      courseId: newSession.course_id ?? '',
      date: newSession.date ?? '',
      time: newSession.time ?? '09:00',
      duration: newSession.duration ?? 120,
      location: newSession.location ?? '',
      instructorId: newSession.instructor_id ?? '',
      status: (newSession.status ?? 'scheduled') as TrainingEvent['status'],
      attendees: newSession.attendees ? [...newSession.attendees] : [],
    }
  }
}, { immediate: true })

const formattedDate = computed(() => {
  if (!props.session?.date) return ''
  return new Date(props.session.date).toLocaleDateString('en-GB', {
    weekday: 'short', day: 'numeric', month: 'short', year: 'numeric',
  })
})

const attendeeSearch = ref('')

// Filter out instructor from attendee list
const availableForAttendees = computed(() => {
  if (!employees.value) return []
  if (!form.value.instructorId) return employees.value
  return employees.value.filter(emp => emp.id !== form.value.instructorId)
})

const filteredEmployees = computed(() => {
  const list = availableForAttendees.value
  if (!attendeeSearch.value) return list
  const search = attendeeSearch.value.toLowerCase()
  return list.filter(emp =>
    emp.name.toLowerCase().includes(search) ||
    emp.id.toLowerCase().includes(search) ||
    emp.department.toLowerCase().includes(search)
  )
})

function getEmployeeName(id: string): string {
  return employees.value?.find(e => e.id === id)?.name || id
}

function isAttendeeSelected(id: string): boolean {
  return form.value.attendees.includes(id)
}

function toggleAttendee(id: string) {
  const index = form.value.attendees.indexOf(id)
  if (index === -1) {
    form.value.attendees.push(id)
  } else {
    form.value.attendees.splice(index, 1)
  }
}

function selectAll() {
  form.value.attendees = availableForAttendees.value.map(e => e.id)
}

function clearAll() {
  form.value.attendees = []
}

function handleSubmit() {
  if (!props.session) return
  if (!form.value.courseId || !form.value.date || !form.value.location || !form.value.instructorId) return

  updateSession({
    id: props.session.id,
    data: { ...form.value },
  }, {
    onSuccess: () => {
      isOpen.value = false
    }
  })
}
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetContent class="sheet-content-wide">
      <SheetHeader class="sheet-header">
        <SheetTitle class="sheet-title">Edit Session</SheetTitle>
        <SheetDescription class="sheet-description">
          Update the details for <strong>{{ session?.course?.name || session?.course_name }}</strong> — {{ formattedDate }}
        </SheetDescription>
      </SheetHeader>

      <div class="sheet-body">
        <Tabs defaultValue="details" class="sheet-tabs">
          <TabsList class="sheet-tabs-header">
            <TabsTrigger value="details">Details</TabsTrigger>
            <TabsTrigger value="attendees">
              Attendees
              <Badge variant="secondary" class="ml-2">{{ form.attendees.length }}</Badge>
            </TabsTrigger>
          </TabsList>

          <!-- Details Tab -->
          <TabsContent value="details" class="sheet-tab-content">
            <div class="tab-panel">
              <div class="form-grid">
                <!-- Course -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Course</Label>
                  <div class="form-field">
                    <Select v-model="form.courseId">
                      <SelectTrigger>
                        <SelectValue placeholder="Select course" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="course in courses" :key="course.id" :value="course.id">
                          {{ course.name }} ({{ course.validity_months }}mo)
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <!-- Date & Time -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Date &amp; Time</Label>
                  <div class="form-field">
                    <DateTimePicker v-model="form.date" v-model:model-time="form.time"
                      placeholder="Select date &amp; time" />
                  </div>
                </div>

                <!-- Duration -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Duration</Label>
                  <div class="form-field">
                    <Select v-model="form.duration">
                      <SelectTrigger>
                        <SelectValue placeholder="Select duration" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="option in durationOptions" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <!-- Location -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Location</Label>
                  <Input v-model="form.location" class="form-field" placeholder="e.g. Bristol Central Plant" />
                </div>

                <!-- Instructor -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Instructor</Label>
                  <div class="form-field">
                    <Select v-model="form.instructorId">
                      <SelectTrigger>
                        <SelectValue placeholder="Select instructor" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="emp in employees" :key="emp.id" :value="emp.id">
                          {{ emp.name }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <!-- Status -->
                <div class="form-row">
                  <Label class="form-label form-label-right">Status</Label>
                  <div class="form-field">
                    <Select v-model="form.status">
                      <SelectTrigger>
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="option in statusOptions" :key="option.value" :value="option.value">
                          {{ option.label }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              </div>
            </div>
          </TabsContent>

          <!-- Attendees Tab -->
          <TabsContent value="attendees" class="sheet-tab-content">
            <div class="tab-panel">
              <div class="attendee-selection-panel">
                <div class="attendee-selection-header">
                  <Input v-model="attendeeSearch" placeholder="Search employees..." class="btn-stretch" />
                  <Button variant="outline" size="sm" @click="selectAll">
                    <UserCheck class="icon-with-text" />
                    All ({{ availableForAttendees.length }})
                  </Button>
                  <Button variant="ghost" size="sm" @click="clearAll">
                    <UserX class="icon-with-text" />
                    Clear
                  </Button>
                </div>

                <div class="attendee-selection-list">
                  <div 
                    v-for="emp in filteredEmployees" 
                    :key="emp.id" 
                    class="attendee-selection-item"
                    :class="{ 'attendee-selection-item-selected': form.attendees.includes(emp.id) }"
                    @click="toggleAttendee(emp.id)"
                  >
                    <Checkbox 
                      :id="'edit-attendee-' + emp.id"
                      :model-value="form.attendees.includes(emp.id)"
                      @update:model-value="toggleAttendee(emp.id)"
                      @click.stop
                    />
                    <label :for="'edit-attendee-' + emp.id" class="flex-1 min-w-0 cursor-pointer" @click.stop>
                      <span class="table-cell-primary">{{ emp.name }}</span>
                      <span class="table-cell-mono ml-2">{{ emp.id }}</span>
                    </label>
                    <span class="table-cell-secondary">{{ emp.department }}</span>
                  </div>
                  <div v-if="!filteredEmployees.length" class="table-empty py-4">
                    <template v-if="form.instructorId">
                      No other employees available.
                    </template>
                    <template v-else>
                      No employees found.
                    </template>
                  </div>
                </div>

                <div class="attendee-selection-count">
                  <Users class="icon-standard icon-muted" />
                  <span>{{ form.attendees.length }} selected</span>
                  <span v-if="form.attendees.length > 0" class="table-cell-secondary">
                    — {{ form.attendees.slice(0, 3).map(getEmployeeName).join(', ') }}
                    {{ form.attendees.length > 3 ? `+${form.attendees.length - 3} more` : '' }}
                  </span>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>

      <SheetFooter class="sheet-footer">
        <div class="sheet-footer-actions">
          <Button variant="outline" class="sheet-footer-btn" @click="isOpen = false">Cancel</Button>
          <Button class="sheet-footer-btn" @click="handleSubmit" :disabled="isPending">
            {{ isPending ? 'Saving...' : 'Save Changes' }}
          </Button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
