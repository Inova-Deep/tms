<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
  class?: HTMLAttributes['class']
}>()
</script>

<template>
  <thead
    class="bg-slate-100 text-slate-700 uppercase text-xs font-semibold"
    :class="cn(props.class)"
  >
    <slot />
  </thead>
</template>
