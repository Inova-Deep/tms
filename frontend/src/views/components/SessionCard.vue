<script setup lang="ts">
import { computed } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import type { TrainingEvent } from '@/types/api'
import { CalendarDays, MapPin, User, Users, Clock, MoreHorizontal, Pencil, CheckSquare } from 'lucide-vue-next'

const props = defineProps<{
    session: TrainingEvent
    instructorName?: string
}>()

const emit = defineEmits<{
    (e: 'confirm', session: TrainingEvent): void
    (e: 'edit', session: TrainingEvent): void
}>()

const statusClass = computed(() => {
    switch (props.session.status) {
        case 'completed': return 'session-status-completed'
        case 'scheduled': return 'session-status-scheduled'
        case 'in_progress': return 'session-status-in-progress'
        case 'cancelled': return 'session-status-cancelled'
        case 'draft': return 'session-status-draft'
        default: return 'session-status-scheduled'
    }
})

const sessionAny = props.session as TrainingEvent & { time?: string; duration?: number }

const formattedDate = computed(() => {
    return new Date(props.session.date).toLocaleDateString('en-GB', {
        weekday: 'short',
        day: 'numeric',
        month: 'short',
        year: 'numeric',
    })
})

function formatDuration(minutes: number | undefined): string {
    if (!minutes) return ''
    if (minutes < 60) return `${minutes} min`
    if (minutes === 60) return '1 hour'
    if (minutes === 480) return 'Full Day'
    if (minutes === 240) return 'Half Day'
    const hours = minutes / 60
    return hours % 1 === 0 ? `${hours} hours` : `${hours} hours`
}
</script>

<template>
    <Card>
        <CardHeader class="card-header-compact">
            <div class="card-row">
                <CardTitle class="form-label">{{ session.course?.name || session.course_name }}</CardTitle>
                <div class="btn-group">
                    <span :class="['session-status-badge', statusClass]">{{ session.status }}</span>
                    <DropdownMenu>
                        <DropdownMenuTrigger as-child>
                            <Button variant="ghost" size="icon" class="btn-icon-sm">
                                <MoreHorizontal class="icon-standard" />
                                <span class="sr-only">Actions</span>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuItem @click="emit('edit', session)">
                                <Pencil class="icon-with-text" />
                                Edit Session
                            </DropdownMenuItem>
                            <DropdownMenuSeparator v-if="session.status !== 'cancelled'" />
                            <DropdownMenuItem
                                v-if="session.status === 'scheduled' || session.status === 'draft' || session.status === 'in_progress'"
                                @click="emit('confirm', session)"
                            >
                                <CheckSquare class="icon-with-text" />
                                Confirm Attendance
                            </DropdownMenuItem>
                            <DropdownMenuItem
                                v-if="session.status === 'completed'"
                                @click="emit('confirm', session)"
                            >
                                <Users class="icon-with-text" />
                                View Attendees
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>
        </CardHeader>
        <CardContent class="card-content-compact">
            <div class="employee-meta-item">
                <CalendarDays class="employee-meta-icon" />
                {{ formattedDate }}<span v-if="sessionAny.time"> • {{ sessionAny.time }}</span>
            </div>
            <div v-if="sessionAny.duration" class="employee-meta-item">
                <Clock class="employee-meta-icon" />
                {{ formatDuration(sessionAny.duration) }}
            </div>
            <div class="employee-meta-item">
                <MapPin class="employee-meta-icon" />
                {{ session.location }}
            </div>
            <div class="employee-meta-item">
                <User class="employee-meta-icon" />
                {{ instructorName || session.instructor_id }}
            </div>
            <div class="employee-meta-item">
                <Users class="employee-meta-icon" />
                {{ session.attendees?.length || 0 }} attendees
            </div>
        </CardContent>
    </Card>
</template>
