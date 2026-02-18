<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from '@/components/ui/table'
import { Pencil, Trash2 } from 'lucide-vue-next'
import type { ProfileExclusion } from '@/types/api'

const props = defineProps<{
  profileId: string
  exclusions: ProfileExclusion[]
}>()

const emit = defineEmits<{
  'add-exclusion': []
  'remove-exclusion': [exclusionId: string]
  'edit-exclusion': [exclusion: ProfileExclusion]
}>()

function isExpired(expiryDate: string | null): boolean {
  if (!expiryDate) return false
  return new Date(expiryDate) < new Date()
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return 'No expiry'
  return new Date(dateStr).toLocaleDateString()
}
</script>

<template>
  <div class="form-row-simple">
    <div class="btn-group">
      <h3 class="requirement-title">Exclusions</h3>
      <Button variant="outline" size="sm" @click="emit('add-exclusion')">
        Add Exclusion
      </Button>
    </div>

    <Table v-if="exclusions.length > 0">
      <TableHeader>
        <TableRow>
          <TableHead>Employee ID</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>Reason</TableHead>
          <TableHead>Expiry Date</TableHead>
          <TableHead class="table-action-cell-wide"></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="exclusion in exclusions" :key="exclusion.id">
          <TableCell class="table-cell-mono">{{ exclusion.employeeId }}</TableCell>
          <TableCell class="table-cell-primary">{{ exclusion.employeeName }}</TableCell>
          <TableCell class="exclusion-reason">{{ exclusion.reason }}</TableCell>
          <TableCell :class="{ 'exclusion-expired': isExpired(exclusion.expiryDate) }">
            {{ formatDate(exclusion.expiryDate) }}
            <Badge v-if="isExpired(exclusion.expiryDate)" variant="secondary" class="badge-spacing">
              Expired
            </Badge>
          </TableCell>
          <TableCell>
            <div class="btn-group">
              <Button variant="ghost" size="icon" class="btn-icon-sm" @click="emit('edit-exclusion', exclusion)">
                <Pencil class="icon-standard" />
                <span class="sr-only">Edit</span>
              </Button>
              <Button variant="ghost" size="icon" class="btn-icon-sm" @click="emit('remove-exclusion', exclusion.id)">
                <Trash2 class="icon-standard" />
                <span class="sr-only">Remove</span>
              </Button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <div v-else class="table-empty">
      No exclusions defined.
    </div>
  </div>
</template>
