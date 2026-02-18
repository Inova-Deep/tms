<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <th
    class="text-left px-3 py-2 border-b border-slate-200 font-semibold text-xs uppercase tracking-wide"
    :class="cn(props.class)"
  >
    <slot />
  </th>
</template>
