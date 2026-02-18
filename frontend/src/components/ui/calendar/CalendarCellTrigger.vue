<script lang="ts" setup>
import type { CalendarCellTriggerProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { CalendarCellTrigger, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"
import { buttonVariants } from '@/components/ui/button'

const props = withDefaults(defineProps<CalendarCellTriggerProps & { class?: HTMLAttributes["class"] }>(), {
  as: "button",
})

const delegatedProps = reactiveOmit(props, "class")

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <CalendarCellTrigger
    data-slot="calendar-cell-trigger"
    :class="cn(
      buttonVariants({ variant: 'ghost' }),
      'size-8 p-0 font-normal text-slate-700',
      '[&[data-today]:not([data-selected])]:bg-slate-100 [&[data-today]:not([data-selected])]:text-slate-700',
      'data-[selected]:bg-slate-700 data-[selected]:text-white data-[selected]:opacity-100 data-[selected]:hover:bg-slate-600 data-[selected]:hover:text-white data-[selected]:focus:bg-slate-700 data-[selected]:focus:text-white',
      'data-[disabled]:text-slate-300 data-[disabled]:opacity-50',
      'data-[unavailable]:text-rose-600 data-[unavailable]:line-through',
      'data-[outside-view]:text-slate-300',
      props.class,
    )"
    v-bind="forwardedProps"
  >
    <slot />
  </CalendarCellTrigger>
</template>
