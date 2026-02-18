<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder?: string
    options: { label: string; value: string }[]
    class?: string
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    placeholder: 'Select…',
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const open = ref(false)

const selectedLabel = computed(() => props.options.find((o) => o.value === props.modelValue)?.label)

function pick(val: string) {
  emit('update:modelValue', val)
  open.value = false
}

function closeOnOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('[data-select-root]')) {
    open.value = false
  }
}

watch(
  () => props.modelValue,
  () => {
    // ensure menu closes when value changes from outside
    open.value = false
  },
)

onMounted(() => {
  document.addEventListener('click', closeOnOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', closeOnOutside)
})
</script>

<template>
  <div class="relative" data-select-root :class="class">
    <button
      type="button"
      class="flex h-9 w-full items-center justify-between rounded-md border border-slate-300 bg-white px-3 text-sm text-slate-800 shadow-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-500 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
      :disabled="disabled"
      @click.stop="open = !open"
    >
      <span class="truncate">
        {{ selectedLabel || placeholder }}
      </span>
      <span class="text-xs text-slate-500">▾</span>
    </button>
    <div
      v-if="open"
      class="absolute z-30 mt-1 w-full rounded-md border border-slate-200 bg-white shadow-lg max-h-64 overflow-y-auto"
    >
      <button
        v-for="opt in options"
        :key="opt.value"
        type="button"
        class="w-full text-left px-3 py-2 text-sm hover:bg-slate-50"
        @click.stop="pick(opt.value)"
      >
        {{ opt.label }}
      </button>
    </div>
  </div>
</template>
