<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    variant?: 'neutral' | 'success' | 'warning' | 'danger' | 'info'
  }>(),
  { variant: 'neutral' },
)

const variants: Record<string, string> = {
  neutral: 'bg-slate-50 text-slate-600 border border-slate-200',
  success: 'bg-emerald-50 text-emerald-600 border border-emerald-200',
  warning: 'bg-amber-50 text-amber-700 border border-amber-200',
  danger: 'bg-rose-50 text-rose-600 border border-rose-200',
  info: 'bg-slate-50 text-slate-600 border border-slate-200',
}
</script>

<template>
  <span
    class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
    :class="variants[props.variant]"
  >
    <slot />
  </span>
</template>
