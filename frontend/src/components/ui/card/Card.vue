<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<{
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <div
    data-slot="card"
    :class="
      cn(
        'bg-white text-slate-900 flex flex-col gap-4 rounded border border-slate-200 py-4',
        props.class,
      )
    "
  >
    <slot />
  </div>
</template>
