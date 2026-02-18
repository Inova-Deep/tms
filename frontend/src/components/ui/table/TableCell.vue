<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <td class="px-3 py-2 border-b border-slate-100 text-sm text-slate-800" :class="cn(props.class)">
    <slot />
  </td>
</template>
