<script setup lang="ts">
import { ref, computed } from 'vue'
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
import { useCreateEvent } from '@/composables/useEvents'
import { useEmployees } from '@/composables/useEmployees'
import { useCourses } from '@/composables/useCourses'
import { Users, UserCheck, UserX } from 'lucide-vue-next'

const props = defineProps<{
    open: boolean
}>()

const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
}>()

const isOpen = computed({
    get: () => props.open,
    set: (val) => emit('update:open', val),
})

const { mutate: createSession, isPending } = useCreateEvent()
const { data: employees } = useEmployees()
const { data: courses } = useCourses('course')

function getCurrentHour(): string {
    const now = new Date()
    const hour = Math.ceil(now.getHours())
    return `${hour.toString().padStart(2, '0')}:00`
}

const form = ref({
    courseId: '',
    date: '',
    time: getCurrentHour(),
    duration: 120,
    location: '',
    instructorId: '',
    attendees: [] as string[],
})

const durationOptions = [
    { label: '30 minutes', value: 30 },
    { label: '1 hour', value: 60 },
    { label: '1.5 hours', value: 90 },
    { label: '2 hours', value: 120 },
    { label: '3 hours', value: 180 },
    { label: '4 hours', value: 240 },
    { label: 'Half Day', value: 240 },
    { label: 'Full Day', value: 480 },
]

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

function selectAllEligible() {
    const allIds = availableForAttendees.value.map(emp => emp.id)
    const current = new Set(form.value.attendees)
    allIds.forEach(id => current.add(id))
    form.value.attendees = Array.from(current)
}

function clearSelection() {
    form.value.attendees = []
}

function handleSubmit() {
    if (!form.value.courseId || !form.value.date || !form.value.location || !form.value.instructorId) {
        return
    }

    createSession({
        ...form.value,
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
                <SheetTitle class="sheet-title">Schedule Training Session</SheetTitle>
                <SheetDescription class="sheet-description">
                    Create a new training session and assign attendees.
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
                                                <SelectItem v-for="option in durationOptions" :key="option.value"
                                                    :value="option.value">
                                                    {{ option.label }}
                                                </SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </div>

                                <!-- Location -->
                                <div class="form-row">
                                    <Label class="form-label form-label-right">Location</Label>
                                    <Input v-model="form.location" class="form-field"
                                        placeholder="e.g. Bristol Central Plant" />
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
                            </div>
                        </div>
                    </TabsContent>

                    <!-- Attendees Tab -->
                    <TabsContent value="attendees" class="sheet-tab-content">
                        <div class="tab-panel">
                            <div class="attendee-selection-panel">
                                <div class="attendee-selection-header">
                                    <Input v-model="attendeeSearch" placeholder="Search employees..." class="btn-stretch" />
                                    <Button variant="outline" size="sm" @click="selectAllEligible">
                                        <UserCheck class="icon-with-text" />
                                        All ({{ availableForAttendees.length }})
                                    </Button>
                                    <Button variant="ghost" size="sm" @click="clearSelection">
                                        <UserX class="icon-with-text" />
                                        Clear
                                    </Button>
                                </div>

                                <div class="attendee-selection-list">
                                    <div 
                                        v-for="emp in filteredEmployees" 
                                        :key="emp.id" 
                                        class="attendee-selection-item"
                                        :class="{ 'attendee-selection-item-selected': isAttendeeSelected(emp.id) }"
                                    >
                                        <Checkbox 
                                            :id="'attendee-' + emp.id"
                                            :checked="isAttendeeSelected(emp.id)"
                                            @update:checked="() => toggleAttendee(emp.id)" 
                                        />
                                        <label :for="'attendee-' + emp.id" class="flex-1 min-w-0 cursor-pointer">
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
                        {{ isPending ? 'Creating...' : 'Create Session' }}
                    </Button>
                </div>
            </SheetFooter>
        </SheetContent>
    </Sheet>
</template>
