<script setup lang="ts">
import { RouterView } from 'vue-router'
import { ref, computed, provide, onMounted } from 'vue'
import AppSidebar from '@/components/common/AppSidebar.vue'
import GlobalErrorBar from '@/components/common/GlobalErrorBar.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const isAuthReady = ref(false)

onMounted(async () => {
  await auth.init()
  isAuthReady.value = true
})

const sessionLabel = computed(() => (auth.isEmployee ? `Employee · ${auth.employeeId || '—'}` : 'Admin'))

const isSidebarCollapsed = ref(false)
provide('sidebarCollapsed', isSidebarCollapsed)

const contentClass = computed(() =>
  isSidebarCollapsed.value
    ? 'max-w-full px-4'
    : 'max-w-7xl px-4 md:px-6'
)
</script>

<template>
  <div v-if="!isAuthReady" class="loading-container h-screen border-none">
    <div class="text-center">
      <div class="loading-spinner loading-spinner-lg mb-4"></div>
      <p class="table-cell-muted">Initializing...</p>
    </div>
  </div>
  <div v-else class="app-container">
    <AppSidebar @collapse-change="(val: boolean) => isSidebarCollapsed = val" />
    <div class="app-main">
      <header class="app-header">
        <div :class="[contentClass, 'app-header-content']">
          <div class="app-header-title">TMS</div>
          <div class="app-header-meta">
            <p class="page-meta">UK time · 24h format</p>
            <p class="form-label">{{ sessionLabel }}</p>
          </div>
        </div>
      </header>

      <GlobalErrorBar />

      <main class="app-content">
        <div :class="[contentClass, 'app-content-wrapper']">
          <RouterView />
        </div>
      </main>
    </div>
  </div>
</template>
