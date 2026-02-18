<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
    DialogDescription,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import type { TrainingEvent, Employee } from '@/types/api'
import { useConfirmAttendance } from '@/composables/useEvents'
import { useEmployees } from '@/composables/useEmployees'

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

const { mutate: confirmAttendance, isPending } = useConfirmAttendance()
const { data: employees } = useEmployees()

const selectedEmployeeIds = ref<string[]>([])

// When dialog opens with a session, pre-select scheduled attendees
watch(() => props.session, (newSession) => {
    if (newSession && newSession.attendees) {
        selectedEmployeeIds.value = [...newSession.attendees]
    } else {
        selectedEmployeeIds.value = []
    }
}, { immediate: true })

function getEmployeeName(id: string) {
    return employees.value?.find((e: Employee) => e.id === id)?.name || id
}

const formattedDate = computed(() => {
    if (!props.session?.date) return ''
    return new Date(props.session.date).toLocaleDateString('en-GB', {
        weekday: 'short', day: 'numeric', month: 'short', year: 'numeric',
    })
})


function handleSubmit() {
    if (!props.session) return

    confirmAttendance({
        eventId: props.session.id,
        data: {
            attendees: selectedEmployeeIds.value.map(employeeId => ({
                employeeId,
                attended: true,
            })),
        },
    }, {
        onSuccess: () => {
            isOpen.value = false
        }
    })
}

function toggleSelection(id: string, checked: boolean) {
    if (checked) {
        if (!selectedEmployeeIds.value.includes(id)) {
            selectedEmployeeIds.value.push(id)
        }
    } else {
        selectedEmployeeIds.value = selectedEmployeeIds.value.filter(x => x !== id)
    }
}
</script>

<template>
    <Dialog v-model:open="isOpen">
        <DialogContent class="dialog-content">
            <DialogHeader>
                <DialogTitle>Confirm Attendance</DialogTitle>
                <DialogDescription>
                    Verify who attended <strong>{{ session?.course?.name || session?.course_name }}</strong> on {{ formattedDate }}.
                </DialogDescription>
            </DialogHeader>

            <div class="dialog-body">
                <div v-if="!session?.attendees?.length" class="table-empty text-center">
                    No attendees were scheduled.
                </div>
                <div v-else class="dialog-scroll-body">
                    <div v-for="empId in session.attendees" :key="empId" class="btn-group">
                        <Checkbox :id="empId" :checked="selectedEmployeeIds.includes(empId)"
                            @update:checked="(val: boolean) => toggleSelection(empId, val)" />
                        <Label :for="empId">
                            {{ getEmployeeName(empId) }}
                        </Label>
                    </div>
                </div>
            </div>

            <DialogFooter>
                <Button variant="outline" @click="isOpen = false">Cancel</Button>
                <Button @click="handleSubmit" :disabled="isPending">
                    {{ isPending ? 'Confirming...' : 'Confirm & Generate Evidence' }}
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>
