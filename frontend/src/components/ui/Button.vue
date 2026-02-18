<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import { buttonVariants } from '@/components/ui/button'

const props = withDefaults(
  defineProps<{
    variant?: 'default' | 'secondary' | 'ghost' | 'outline' | 'destructive' | 'link'
    size?: 'sm' | 'default' | 'lg' | 'icon'
    as?: 'button' | 'a'
    href?: string
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'default',
    size: 'default',
    as: 'button',
    type: 'button',
  },
)

const classes = computed(() => cn(buttonVariants({ variant: props.variant, size: props.size })))
</script>

<template>
  <component :is="as" :href="href" :type="as === 'button' ? type : undefined"
    class="focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-slate-400" :class="classes">
    <slot />
  </component>
</template>
