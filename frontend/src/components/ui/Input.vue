<script setup lang="ts">
const modelValue = defineModel<string | number | null>({ default: '' })
defineProps<{
  type?: string
  placeholder?: string
  disabled?: boolean
}>()
</script>

<template>
  <input
    v-model="modelValue"
    :type="type || 'text'"
    :placeholder="placeholder"
    :disabled="disabled"
    class="flex h-9 w-full rounded border border-slate-300 bg-white px-3 text-sm text-slate-900 placeholder:text-slate-400 transition focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-slate-400 disabled:cursor-not-allowed disabled:opacity-50"
  />
</template>
