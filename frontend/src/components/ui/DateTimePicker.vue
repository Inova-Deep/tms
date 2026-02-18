<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { CalendarIcon, Clock } from 'lucide-vue-next'
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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'

/**
 * DateTimePicker — shadcn-compliant, no native HTML inputs.
 * Props:
 *   modelValue (date):  string in 'yyyy-MM-dd' format (emits same)
 *   modelTime  (time):  string in 'HH:mm' format (emits same)
 *   placeholder: string
 */
const props = defineProps<{
    modelValue?: string   // date: yyyy-MM-dd
    modelTime?: string    // time: HH:mm
    placeholder?: string
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', val: string | undefined): void
    (e: 'update:modelTime', val: string): void
}>()

// ── Date ──────────────────────────────────────────────────────────────────────
const df = new DateFormatter('en-GB', { dateStyle: 'medium' }) // e.g. "25 Feb 2026"

const dateValue = ref<DateValue>()

watch(() => props.modelValue, (newVal) => {
    if (newVal) {
        try {
            const datePart = newVal.split('T')[0]!
            dateValue.value = parseDate(datePart)
        } catch {
            console.error('DateTimePicker: invalid date', newVal)
        }
    } else {
        dateValue.value = undefined
    }
}, { immediate: true })

watch(dateValue, (newVal) => {
    emit('update:modelValue', newVal ? newVal.toString() : undefined)
})

// ── Time ──────────────────────────────────────────────────────────────────────
// Generate hours 00–23 and minutes 00, 15, 30, 45
const hours = Array.from({ length: 24 }, (_, i) => String(i).padStart(2, '0'))
const minutes = ['00', '15', '30', '45']

const selectedHour = ref('09')
const selectedMinute = ref('00')

// Parse initial time prop
watch(() => props.modelTime, (newVal) => {
    if (newVal && /^\d{2}:\d{2}$/.test(newVal)) {
        const [h, m] = newVal.split(':')
        selectedHour.value = h!
        // Snap to nearest quarter
        const minNum = parseInt(m!)
        selectedMinute.value = minNum < 8 ? '00' : minNum < 23 ? '15' : minNum < 38 ? '30' : minNum < 53 ? '45' : '00'
    }
}, { immediate: true })

// Emit combined time whenever hour or minute changes
watch([selectedHour, selectedMinute], ([h, m]) => {
    emit('update:modelTime', `${h}:${m}`)
})

// ── Display ───────────────────────────────────────────────────────────────────
const displayValue = computed(() => {
    if (!dateValue.value) return null
    const dateStr = df.format(dateValue.value.toDate(getLocalTimeZone()))
    return `${dateStr}  ${selectedHour.value}:${selectedMinute.value}`
})
</script>

<template>
    <Popover>
        <PopoverTrigger as-child>
            <Button
                variant="outline"
                :class="cn('datetime-picker-trigger', !dateValue && 'datetime-picker-placeholder')"
            >
                <CalendarIcon class="datetime-picker-icon" />
                <span>{{ displayValue ?? (placeholder || 'Pick date & time') }}</span>
            </Button>
        </PopoverTrigger>
        <PopoverContent class="datetime-picker-popover">
            <!-- Calendar -->
            <Calendar v-model="dateValue" initial-focus />

            <!-- Time row -->
            <div class="datetime-picker-time-row">
                <Clock class="datetime-picker-time-icon" />
                <Select v-model="selectedHour">
                    <SelectTrigger class="datetime-picker-time-select">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="h in hours" :key="h" :value="h">{{ h }}</SelectItem>
                    </SelectContent>
                </Select>
                <span class="datetime-picker-time-sep">:</span>
                <Select v-model="selectedMinute">
                    <SelectTrigger class="datetime-picker-time-select">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem v-for="m in minutes" :key="m" :value="m">{{ m }}</SelectItem>
                    </SelectContent>
                </Select>
            </div>
        </PopoverContent>
    </Popover>
</template>
