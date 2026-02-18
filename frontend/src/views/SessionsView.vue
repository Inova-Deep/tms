<script setup lang="ts">
import { ref, computed } from 'vue'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus, Loader2 } from 'lucide-vue-next'
import { useEvents } from '@/composables/useEvents'
import { useEmployees } from '@/composables/useEmployees'
import SessionCard from '@/views/components/SessionCard.vue'
import CreateSessionDialog from '@/views/components/CreateSessionDialog.vue'
import EditSessionDialog from '@/views/components/EditSessionDialog.vue'
import ConfirmAttendanceDialog from '@/views/components/ConfirmAttendanceDialog.vue'
import type { TrainingEvent } from '@/types/api'

const { data: sessions, isLoading, isError } = useEvents()
const { data: employees } = useEmployees()

const isCreateOpen = ref(false)
const isConfirmOpen = ref(false)
const isEditOpen = ref(false)
const sessionToConfirm = ref<TrainingEvent | null>(null)
const sessionToEdit = ref<TrainingEvent | null>(null)

// Map employee IDs to names for quick lookup
const instructorMap = computed(() => {
  const map = new Map<string, string>()
  if (employees.value) {
    employees.value.forEach(emp => {
      map.set(emp.id, emp.name)
    })
  }
  return map
})

// Sort sessions: Scheduled/Draft first (by date), then Completed (by date desc)
const sortedSessions = computed(() => {
  if (!sessions.value) return []

  return [...sessions.value].sort((a, b) => {
    // Define priority order
    const statusPriority: Record<string, number> = {
      'draft': 0,
      'scheduled': 1,
      'in_progress': 2,
      'completed': 3,
      'cancelled': 4
    }
    
    const aPriority = statusPriority[a.status] ?? 5
    const bPriority = statusPriority[b.status] ?? 5
    
    if (aPriority !== bPriority) {
      return aPriority - bPriority
    }
    
    // Same status: sort by date
    if (a.status === 'scheduled' || a.status === 'draft') {
      // Upcoming: soonest first
      return new Date(a.date).getTime() - new Date(b.date).getTime()
    }
    // Past: latest first
    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
})

function openConfirmDialog(session: TrainingEvent) {
  sessionToConfirm.value = session
  isConfirmOpen.value = true
}

function openEditDialog(session: TrainingEvent) {
  sessionToEdit.value = session
  isEditOpen.value = true
}
</script>

<template>
  <section class="page-section">
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Training Sessions</h1>
        <p class="page-description">Schedule and manage training sessions.</p>
      </div>
      <div class="page-actions">
        <Button @click="isCreateOpen = true">
          <Plus class="icon-with-text" />
          Schedule Session
        </Button>
      </div>
    </div>

    <div v-if="isLoading" class="loading-container">
      <div class="loading-content">
        <Loader2 class="h-6 w-6 animate-spin text-slate-400" />
        <span class="table-cell-secondary">Loading sessions...</span>
      </div>
    </div>

    <div v-else-if="isError" class="error-container">
      <div class="error-text">Failed to load sessions. Please try again.</div>
    </div>

    <div v-else-if="sortedSessions.length === 0" class="empty-state">
      <div class="empty-state-content">
        <div class="empty-state-icon">
          <Plus class="empty-state-icon-inner" />
        </div>
        <h3 class="empty-state-title">No sessions scheduled</h3>
        <p class="empty-state-description">
          You haven't scheduled any training sessions yet.
        </p>
        <Button @click="isCreateOpen = true">
          Schedule Session
        </Button>
      </div>
    </div>

    <div v-else class="sessions-grid">
      <SessionCard 
        v-for="session in sortedSessions" 
        :key="session.id" 
        :session="session"
        :instructor-name="instructorMap.get(session.instructor_id)" 
        @confirm="openConfirmDialog"
        @edit="openEditDialog" 
      />
    </div>

    <CreateSessionDialog v-if="isCreateOpen" v-model:open="isCreateOpen" />

    <EditSessionDialog v-if="isEditOpen" v-model:open="isEditOpen" :session="sessionToEdit" />

    <ConfirmAttendanceDialog v-model:open="isConfirmOpen" :session="sessionToConfirm" />
  </section>
</template>
