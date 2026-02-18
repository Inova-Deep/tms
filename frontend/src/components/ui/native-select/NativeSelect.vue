<script setup lang="ts">
import type { AcceptableValue } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit, useVModel } from "@vueuse/core"
import { ChevronDownIcon } from "lucide-vue-next"
import { cn } from "@/lib/utils"

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{ modelValue?: AcceptableValue | AcceptableValue[], class?: HTMLAttributes["class"] }>()

const emit = defineEmits<{
  "update:modelValue": AcceptableValue
}>()

const modelValue = useVModel(props, "modelValue", emit, {
  passive: true,
  defaultValue: "",
})

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <div
    class="group/native-select relative w-fit has-[select:disabled]:opacity-50"
    data-slot="native-select-wrapper"
  >
    <select
      v-bind="{ ...$attrs, ...delegatedProps }"
      v-model="modelValue"
      data-slot="native-select"
      :class="cn(
        'border-slate-200 placeholder:text-slate-400 selection:bg-slate-700 selection:text-white h-9 w-full min-w-0 appearance-none rounded border bg-white px-3 py-2 pr-9 text-sm transition-colors outline-none disabled:pointer-events-none disabled:cursor-not-allowed',
        'focus-visible:border-slate-400 focus-visible:ring-1 focus-visible:ring-slate-400',
        props.class,
      )"
    >
      <slot />
    </select>
    <ChevronDownIcon
      class="text-slate-400 pointer-events-none absolute top-1/2 right-3.5 size-4 -translate-y-1/2 opacity-50 select-none"
      aria-hidden="true"
      data-slot="native-select-icon"
    />
  </div>
</template>
