<script setup lang="ts">
import { ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableHeader,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
} from '@/components/ui/table'
import { Search, UserMinus, UserPlus } from 'lucide-vue-next'
import type { ProfileMember } from '@/types/api'

const props = defineProps<{
  profileId: string
  members: ProfileMember[]
}>()

const emit = defineEmits<{
  'add-member': []
  'remove-member': [employeeId: string]
}>()

const searchQuery = ref('')

const filteredMembers = computed(() => {
  if (!searchQuery.value) return props.members
  const query = searchQuery.value.toLowerCase()
  return props.members.filter(
    member =>
      member.name.toLowerCase().includes(query) ||
      member.id.toLowerCase().includes(query) ||
      member.department.toLowerCase().includes(query)
  )
})

function handleRemoveMember(employeeId: string) {
  emit('remove-member', employeeId)
}
</script>

<template>
  <div class="members-tab-content">
    <div class="members-header">
      <div class="filter-search-wrapper">
        <Search class="filter-search-icon" />
        <Input
          v-model="searchQuery"
          placeholder="Search members..."
          class="members-search-input"
        />
      </div>
      <Button variant="outline" size="sm" @click="emit('add-member')">
        <UserPlus class="icon-standard mr-2" />
        Add Members
      </Button>
    </div>

    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Employee ID</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>Department</TableHead>
          <TableHead class="table-action-cell"></TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-if="!filteredMembers.length">
          <TableCell colspan="4" class="table-empty">
            No members found.
          </TableCell>
        </TableRow>
        <TableRow v-for="member in filteredMembers" :key="member.id">
          <TableCell class="table-cell-mono">{{ member.id }}</TableCell>
          <TableCell class="table-cell-primary">{{ member.name }}</TableCell>
          <TableCell class="table-cell-secondary">{{ member.department }}</TableCell>
          <TableCell>
            <Button
              variant="ghost"
              size="icon"
              class="btn-icon-sm btn-danger-ghost"
              @click="handleRemoveMember(member.id)"
            >
              <UserMinus class="icon-standard" />
              <span class="sr-only">Remove member</span>
            </Button>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
