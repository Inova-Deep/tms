<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { LayoutDashboard, Users, BookOpen, BadgeCheck, CalendarRange, FileBadge2, Menu, User, ChevronUp, RefreshCw, LogOut, Loader2, HelpCircle, Grid3x3, GraduationCap } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from '@/components/ui/dropdown-menu'
import { useAuthStore } from '@/stores/auth'
import { useEmployees } from '@/composables/useEmployees'

type NavItem = {
  label: string
  to: string
  icon: any
}

type NavGroup = {
  header: string
  items: NavItem[]
}

const route = useRoute()
const auth = useAuthStore()
const isCollapsed = ref(false)

const emit = defineEmits<{
  (e: 'collapse-change', value: boolean): void
}>()

watch(isCollapsed, (val) => {
  emit('collapse-change', val)
})

const navGroups: NavGroup[] = [
  {
    header: 'Overview',
    items: [
      { label: 'Dashboard', to: '/dashboard', icon: LayoutDashboard },
      { label: 'Training Matrix', to: '/matrix', icon: Grid3x3 },
    ],
  },
  {
    header: 'Resources',
    items: [
      { label: 'Employees', to: '/employees', icon: Users },
      { label: 'Courses', to: '/courses', icon: BookOpen },
    ],
  },
  {
    header: 'Training',
    items: [
      { label: 'Profiles', to: '/profiles', icon: BadgeCheck },
      { label: 'Sessions', to: '/sessions', icon: CalendarRange },
      { label: 'Certifications', to: '/certifications', icon: FileBadge2 },
    ],
  },
]

const myTrainingItem: NavItem = { label: 'My Training', to: '/my-training', icon: GraduationCap }
const helpItem: NavItem = { label: 'Help', to: '/help', icon: HelpCircle }

const sidebarWidth = computed(() => (isCollapsed.value ? 'sidebar-collapsed' : 'sidebar-expanded'))

function toggle() {
  isCollapsed.value = !isCollapsed.value
}

function isActive(path: string) {
  return route.path.startsWith(path)
}

// Employee list for switching
const { data: employees } = useEmployees()
const employeeOptions = computed(() => employees?.value ?? [])

// Role label
const roleLabel = computed(() => (auth.isEmployee ? 'Employee' : 'Admin'))
</script>

<template>
  <aside class="sidebar" :class="sidebarWidth">
    <!-- Header -->
    <div class="sidebar-header">
      <span v-if="!isCollapsed" class="sidebar-title">TMS</span>
      <Button variant="ghost" size="icon" class="sidebar-toggle" @click="toggle">
        <Menu class="icon-standard" />
        <span class="sr-only">Toggle sidebar</span>
      </Button>
    </div>

    <!-- Navigation -->
    <nav class="sidebar-nav">
      <template v-for="group in navGroups" :key="group.header">
        <div v-if="!isCollapsed" class="sidebar-group-header">{{ group.header }}</div>
        <RouterLink
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          class="sidebar-nav-item"
          :class="isActive(item.to) ? 'sidebar-nav-item-active' : 'sidebar-nav-item-inactive'"
        >
          <component :is="item.icon" class="icon-standard" />
          <span v-if="!isCollapsed">{{ item.label }}</span>
        </RouterLink>
      </template>
    </nav>

    <!-- Bottom Links (My Training + Help) -->
    <div class="sidebar-nav-bottom">
      <RouterLink
        :to="myTrainingItem.to"
        class="sidebar-nav-item"
        :class="isActive(myTrainingItem.to) ? 'sidebar-nav-item-active' : 'sidebar-nav-item-inactive'"
      >
        <component :is="myTrainingItem.icon" class="icon-standard" />
        <span v-if="!isCollapsed">{{ myTrainingItem.label }}</span>
      </RouterLink>
      <RouterLink
        :to="helpItem.to"
        class="sidebar-nav-item"
        :class="isActive(helpItem.to) ? 'sidebar-nav-item-active' : 'sidebar-nav-item-inactive'"
      >
        <component :is="helpItem.icon" class="icon-standard" />
        <span v-if="!isCollapsed">{{ helpItem.label }}</span>
      </RouterLink>
    </div>

    <!-- User Section -->
    <div class="sidebar-footer">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" class="sidebar-user-trigger" :class="isCollapsed ? 'sidebar-user-trigger-collapsed' : ''">
            <div 
              class="sidebar-avatar"
              :class="auth.isEmployee ? 'sidebar-avatar-employee' : 'sidebar-avatar-admin'"
            >
              {{ auth.userInitials }}
            </div>
            <div v-if="!isCollapsed" class="sidebar-user-info">
              <p class="sidebar-user-name">{{ auth.name }}</p>
              <p class="sidebar-user-role">
                {{ roleLabel }}
                <span v-if="auth.department"> · {{ auth.department }}</span>
              </p>
            </div>
            <ChevronUp v-if="!isCollapsed" class="icon-standard text-slate-400" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent 
          align="end" 
          :side="isCollapsed ? 'right' : 'top'" 
          :side-offset="4"
          class="w-64"
        >
          <!-- Current User Info -->
          <DropdownMenuLabel class="font-normal">
            <div class="btn-group">
              <div 
                class="sidebar-avatar"
                :class="auth.isEmployee ? 'sidebar-avatar-employee' : 'sidebar-avatar-admin'"
              >
                {{ auth.userInitials }}
              </div>
              <div class="sidebar-user-info">
                <p class="sidebar-user-name">{{ auth.name }}</p>
                <p class="sidebar-user-role">{{ roleLabel }}</p>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          
          <!-- Loading indicator -->
          <div v-if="auth.loading" class="btn-group px-2 py-1.5 table-cell-secondary">
            <Loader2 class="icon-standard animate-spin" />
            Switching...
          </div>
          
          <template v-else>
            <!-- Switch Role Section -->
            <DropdownMenuLabel class="page-meta font-normal py-1.5 px-2">
              Switch to
            </DropdownMenuLabel>
            
            <!-- Switch to Admin (for employees) -->
            <DropdownMenuItem 
              v-if="auth.isEmployee"
              @click="auth.loginAsAdmin"
            >
              <User class="icon-standard mr-2" />
              Admin
            </DropdownMenuItem>
            
            <!-- Switch to Employee submenu -->
            <DropdownMenuSub v-if="!auth.isEmployee">
              <DropdownMenuSubTrigger>
                <User class="icon-standard mr-2" />
                Employee
              </DropdownMenuSubTrigger>
              <DropdownMenuSubContent class="w-56 max-h-64 overflow-y-auto">
                <DropdownMenuItem 
                  v-for="emp in employeeOptions" 
                  :key="emp.id"
                  @click="auth.loginAsEmployee(emp.id)"
                >
                  <div class="flex flex-col">
                    <span class="truncate">{{ emp.name }}</span>
                    <span class="page-meta">{{ emp.id }} · {{ emp.department }}</span>
                  </div>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuSub>
            
            <!-- Switch Employee (for employees already) -->
            <DropdownMenuSub v-else>
              <DropdownMenuSubTrigger>
                <RefreshCw class="icon-standard mr-2" />
                Switch Employee
              </DropdownMenuSubTrigger>
              <DropdownMenuSubContent class="w-56 max-h-64 overflow-y-auto">
                <DropdownMenuItem 
                  v-for="emp in employeeOptions" 
                  :key="emp.id"
                  :disabled="emp.id === auth.employeeId"
                  @click="auth.loginAsEmployee(emp.id)"
                >
                  <div class="flex flex-col">
                    <span class="truncate">{{ emp.name }}</span>
                    <span class="page-meta">{{ emp.id }} · {{ emp.department }}</span>
                  </div>
                  <span v-if="emp.id === auth.employeeId" class="ml-auto page-meta">Current</span>
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuSub>
          </template>
          
          <DropdownMenuSeparator />
          <DropdownMenuItem class="table-cell-secondary" @click="auth.clear">
            <LogOut class="icon-standard mr-2" />
            Clear session
          </DropdownMenuItem>
          
          <p v-if="auth.error" class="px-2 py-1.5 form-error">{{ auth.error }}</p>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </aside>
</template>
