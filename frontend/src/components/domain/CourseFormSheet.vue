<script setup lang="ts">
import { ref, watch, computed } from 'vue'
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
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import type { Course, CreateCourseRequest, UpdateCourseRequest } from '@/types/api'

const props = defineProps<{
  open: boolean
  course?: Course | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', data: CreateCourseRequest | UpdateCourseRequest): void
}>()

const isOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
})

const form = ref({
  name: '',
  code: '',
  type: 'course' as 'course' | 'certification',
  description: '',
  validity_months: 12,
})

const isEditing = computed(() => !!props.course)

watch(
  () => props.course,
  (course) => {
    if (course) {
      form.value = {
        name: course.name,
        code: course.code || '',
        type: course.type,
        description: course.description || '',
        validity_months: course.validity_months,
      }
    } else {
      form.value = {
        name: '',
        code: '',
        type: 'course',
        description: '',
        validity_months: 12,
      }
    }
  },
  { immediate: true }
)

function handleSubmit() {
  const data: CreateCourseRequest | UpdateCourseRequest = {
    name: form.value.name,
    code: form.value.code || undefined,
    type: form.value.type,
    description: form.value.description || undefined,
    validity_months: form.value.validity_months,
  }
  emit('submit', data)
}

function close() {
  isOpen.value = false
}
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetContent class="sheet-content">
      <SheetHeader class="sheet-header">
        <SheetTitle class="sheet-title">
          {{ isEditing ? 'Edit Course' : 'Add Course' }}
        </SheetTitle>
        <SheetDescription class="sheet-description">
          {{ isEditing ? 'Update course details.' : 'Create a new course or certification.' }}
        </SheetDescription>
      </SheetHeader>

      <div class="sheet-body">
        <div class="form-grid">
          <div class="form-row-simple">
            <Label class="form-label form-label-required">Name</Label>
            <Input v-model="form.name" placeholder="Course name" />
          </div>

          <div class="form-row-simple">
            <Label class="form-label">Code</Label>
            <Input v-model="form.code" placeholder="Course code (optional)" />
          </div>

          <div class="form-row-simple">
            <Label class="form-label form-label-required">Type</Label>
            <Select v-model="form.type">
              <SelectTrigger>
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="course">Course</SelectItem>
                <SelectItem value="certification">Certification</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="form-row-simple">
            <Label class="form-label">Description</Label>
            <Textarea
              v-model="form.description"
              placeholder="Course description (optional)"
              rows="3"
            />
          </div>

          <div class="form-row-simple">
            <Label class="form-label form-label-required">Validity (months)</Label>
            <Input
              v-model.number="form.validity_months"
              type="number"
              min="1"
              placeholder="12"
            />
          </div>
        </div>
      </div>

      <SheetFooter class="sheet-footer">
        <div class="sheet-footer-actions">
          <Button variant="outline" class="sheet-footer-btn" @click="close">
            Cancel
          </Button>
          <Button
            class="sheet-footer-btn"
            :disabled="!form.name || form.validity_months < 1"
            @click="handleSubmit"
          >
            {{ isEditing ? 'Update' : 'Create' }}
          </Button>
        </div>
      </SheetFooter>
    </SheetContent>
  </Sheet>
</template>
