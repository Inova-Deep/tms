<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Calendar as CalendarIcon } from 'lucide-vue-next'
import {
    CalendarDate,
    DateFormatter,
    type DateValue,
    getLocalTimeZone,
    parseDate,
    today,
} from '@internationalized/date'

import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover'

const props = defineProps<{
    modelValue?: string
    placeholder?: string
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', payload: string | undefined): void
}>()

const df = new DateFormatter('en-GB', {
    dateStyle: 'long',
})

const value = ref<DateValue>()

// Sync from prop to internal value — strip time portion if ISO timestamp
watch(() => props.modelValue, (newVal) => {
    if (newVal) {
        try {
            // Accept both 'yyyy-MM-dd' and 'yyyy-MM-ddTHH:mm:ssZ'
            const datePart = newVal.split('T')[0]!
            value.value = parseDate(datePart)
        } catch (e) {
            console.error('Invalid date format', newVal)
        }
    } else {
        value.value = undefined
    }
}, { immediate: true })

// Sync from internal value to prop
watch(value, (newVal) => {
    if (newVal) {
        emit('update:modelValue', newVal.toString())
    } else {
        emit('update:modelValue', undefined)
    }
})
</script>

<template>
    <Popover>
        <PopoverTrigger as-child>
            <Button variant="outline" :class="cn(
                'w-full justify-start text-left font-normal',
                !value && 'text-muted-foreground',
            )">
                <CalendarIcon class="mr-2 h-4 w-4" />
                <span v-if="value">{{ df.format(value.toDate(getLocalTimeZone())) }}</span>
                <span v-else>{{ placeholder || "Pick a date" }}</span>
            </Button>
        </PopoverTrigger>
        <PopoverContent class="w-auto p-0">
            <Calendar v-model="value" initial-focus />
        </PopoverContent>
    </Popover>
</template>
