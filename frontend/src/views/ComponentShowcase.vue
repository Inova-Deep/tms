<script setup lang="ts">
import { ref } from 'vue'
import { today } from '@internationalized/date'
import { getLocalTimeZone } from '@internationalized/date'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import Textarea from '@/components/ui/textarea/Textarea.vue'
import { Separator } from '@/components/ui/separator'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Skeleton } from '@/components/ui/skeleton'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useErrorStore } from '@/stores/errors'
import { MoreHorizontal, Mail, Calendar as CalendarIcon } from 'lucide-vue-next'

const confirmOpen = ref(false)
const errorStore = useErrorStore()

// Form state
const selectedEmployee = ref('')
const checkboxChecked = ref(false)
const textareaValue = ref('')
const inputValue = ref('')
const activeTab = ref('overview')
const switchEnabled = ref(false)
const calendarDate = ref(today(getLocalTimeZone()) as any)
const popoverOpen = ref(false)

const employees = [
  { id: 'E1001', name: 'Alice Johnson' },
  { id: 'E1005', name: 'Jack Thompson' },
  { id: 'E1008', name: 'Sarah Williams' },
]

const tableData = [
  { id: 'E1001', name: 'Alice Johnson', department: 'Engineering', status: 'Active' },
  { id: 'E1005', name: 'Jack Thompson', department: 'Operations', status: 'Active' },
  { id: 'E1008', name: 'Sarah Williams', department: 'HR', status: 'Pending' },
]

function formatDate(date: any) {
  if (!date) return 'Pick a date'
  return date.toString()
}
</script>

<template>
  <section class="page-section">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">UI Kit</h1>
        <p class="page-description">Component showcase following Industrial Density design rules.</p>
      </div>
    </div>

    <!-- Header Card -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">UI Kit Showcase</CardTitle>
        <CardDescription>
          Components following <span class="form-label">Industrial Density</span> rules:
          subtle borders, minimal shadows, compact spacing.
        </CardDescription>
      </CardHeader>
      <CardFooter class="btn-group">
        <Button variant="secondary" size="sm" @click="errorStore.push('Example warning', 'Check styling rules')">
          Trigger Error Banner
        </Button>
      </CardFooter>
    </Card>

    <!-- Buttons -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Buttons</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="btn-group">
          <Button>Primary</Button>
          <Button variant="secondary">Secondary</Button>
          <Button variant="ghost">Ghost</Button>
          <Button variant="outline">Outline</Button>
          <Button variant="destructive">Destructive</Button>
          <Button variant="link">Link</Button>
        </div>
        <div class="btn-group">
          <Button size="sm">Small</Button>
          <Button size="default">Default</Button>
          <Button size="lg">Large</Button>
          <Button size="icon" aria-label="Icon button">
            <Mail class="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Badges -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Badges</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="btn-group">
          <Badge>Neutral</Badge>
          <Badge variant="success">Valid</Badge>
          <Badge variant="warning">Expiring</Badge>
          <Badge variant="destructive">Expired</Badge>
          <Badge variant="outline">Outline</Badge>
        </div>
      </CardContent>
    </Card>

    <!-- Form Elements -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Form Elements</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 md:grid-cols-2">
          <div class="form-row-simple">
            <Label for="input-demo">Text Input</Label>
            <Input id="input-demo" v-model="inputValue" placeholder="Type here" />
          </div>
          <div class="form-row-simple">
            <Label>Select</Label>
            <Select v-model="selectedEmployee">
              <SelectTrigger class="w-full">
                <SelectValue placeholder="Select employee" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="emp in employees" :key="emp.id" :value="emp.id">
                  {{ emp.id }} — {{ emp.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="form-row-simple">
            <Label for="textarea-demo">Textarea</Label>
            <Textarea id="textarea-demo" v-model="textareaValue" placeholder="Enter notes..." />
          </div>
          <div class="form-row-simple">
            <Label>Checkbox</Label>
            <div class="btn-group pt-2">
              <Checkbox id="checkbox-demo" v-model:checked="checkboxChecked" />
              <label for="checkbox-demo" class="form-label cursor-pointer">
                I agree to the terms
              </label>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Calendar, Popover & Switch -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Calendar, Popover & Switch</CardTitle>
      </CardHeader>
      <CardContent class="space-y-6">
        <div class="grid gap-6 md:grid-cols-2">
          <!-- Calendar -->
          <div class="form-row-simple">
            <Label>Calendar</Label>
            <Popover v-model:open="popoverOpen">
              <PopoverTrigger as-child>
                <Button variant="outline" class="w-full justify-start text-left font-normal">
                  <CalendarIcon class="mr-2 h-4 w-4" />
                  {{ formatDate(calendarDate) }}
                </Button>
              </PopoverTrigger>
              <PopoverContent class="w-auto p-0" align="start">
                <Calendar v-model="calendarDate" />
              </PopoverContent>
            </Popover>
          </div>
          <!-- Switch -->
          <div class="form-row-simple">
            <Label>Switch</Label>
            <div class="btn-group pt-2">
              <Switch v-model:checked="switchEnabled" />
              <span class="form-label">
                {{ switchEnabled ? 'Enabled' : 'Disabled' }}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Tabs -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Tabs</CardTitle>
      </CardHeader>
      <CardContent>
        <Tabs v-model="activeTab">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="details">Details</TabsTrigger>
            <TabsTrigger value="settings">Settings</TabsTrigger>
          </TabsList>
          <TabsContent value="overview" class="p-4 table-cell-secondary">
            Overview content goes here. This tab shows a summary view.
          </TabsContent>
          <TabsContent value="details" class="p-4 table-cell-secondary">
            Details content with more specific information.
          </TabsContent>
          <TabsContent value="settings" class="p-4 table-cell-secondary">
            Settings and configuration options.
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Data Table</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Department</TableHead>
              <TableHead>Status</TableHead>
              <TableHead class="w-12"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="row in tableData" :key="row.id">
              <TableCell class="table-cell-mono">{{ row.id }}</TableCell>
              <TableCell>{{ row.name }}</TableCell>
              <TableCell>{{ row.department }}</TableCell>
              <TableCell>
                <Badge :variant="row.status === 'Active' ? 'success' : 'warning'">
                  {{ row.status }}
                </Badge>
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon" class="btn-icon-sm">
                      <MoreHorizontal class="h-4 w-4" />
                      <span class="sr-only">Open menu</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>Edit</DropdownMenuItem>
                    <DropdownMenuItem>Duplicate</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem class="text-rose-600">Delete</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Avatar & Skeleton -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Avatar & Skeleton</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="btn-group">
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
          <Avatar>
            <AvatarFallback class="sidebar-avatar-employee">JD</AvatarFallback>
          </Avatar>
          <Avatar>
            <AvatarFallback class="sidebar-avatar-admin">AB</AvatarFallback>
          </Avatar>
        </div>
        <Separator />
        <div class="space-y-2">
          <Skeleton class="h-4 w-64" />
          <Skeleton class="h-4 w-48" />
          <Skeleton class="h-4 w-56" />
        </div>
      </CardContent>
    </Card>

    <!-- Confirm Dialog -->
    <Card>
      <CardHeader>
        <CardTitle class="form-label">Confirm Dialog</CardTitle>
      </CardHeader>
      <CardContent class="space-y-3">
        <p class="table-cell-secondary">
          All destructive actions must present this confirmation (per rules: "No delete without prompt").
        </p>
        <Button variant="destructive" size="sm" @click="confirmOpen = true">
          Show Confirm
        </Button>
        <ConfirmDialog v-model:open="confirmOpen" title="Delete record?"
          description="This will remove the item permanently." confirm-label="Delete" cancel-label="Cancel"
          variant="destructive" />
      </CardContent>
    </Card>
  </section>
</template>
