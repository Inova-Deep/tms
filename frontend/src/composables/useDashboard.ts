import { computed, type Ref } from 'vue'
import { useQuery } from '@tanstack/vue-query' 
import { api } from '@/lib/api'
import type { RequirementStatusValue, ProfileStatusValue } from '@/types/api'

// Types based on API documentation
export type DashboardKpis = {
  totalEmployees: number
  overallCompliance: number
  valid: number
  expiringSoon: number
  expired: number
  missing: number
}

export type GridRow = {
  groupValue: string
  itemId: string
  itemName: string
  total: number
  valid: number
  expiring: number
  expired: number
  missing: number
  percent: number
}

export type DrilldownRow = {
  id: string
  name: string
  department: string
  site: string
  status: RequirementStatusValue | ProfileStatusValue
  expiryDate: string
}

// Fetch dashboard KPIs
export function useDashboardKpis(
  department?: Ref<string>,
  site?: Ref<string>,
  profile?: Ref<string>,
) {
  return useQuery({
    queryKey: ['dashboard', 'kpis', department, site, profile],
    queryFn: async () => {
      const params = new URLSearchParams()
      if (department?.value) params.append('department', department.value)
      if (site?.value) params.append('site', site.value)
      if (profile?.value) params.append('profile', profile.value)

      const url = `/dashboard/kpis${params.toString() ? `?${params.toString()}` : ''}`
      return api<DashboardKpis>(url)
    },
  })
}

// Fetch coverage grid
export function useDashboardGrid(
  mode: Ref<'course' | 'profile'>,
  groupBy: Ref<'department' | 'site' | 'costCenter'>,
  department?: Ref<string>,
  site?: Ref<string>,
) {
  return useQuery({
    queryKey: ['dashboard', 'grid', mode, groupBy, department, site],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('mode', mode.value)
      params.append('groupBy', groupBy.value)
      if (department?.value) params.append('department', department.value)
      if (site?.value) params.append('site', site.value)

      const url = `/dashboard/grid?${params.toString()}`
      return api<GridRow[]>(url)
    },
  })
}

// Fetch drilldown data
export function useDashboardDrilldown(
  mode: Ref<'course' | 'profile'>,
  groupBy: Ref<'department' | 'site' | 'costCenter'>,
  groupValue: Ref<string | null>,
  itemId: Ref<string | null>,
  department?: Ref<string>,
  site?: Ref<string>,
) {
  const enabled = computed(() => !!groupValue.value && !!itemId.value)

  return useQuery({
    queryKey: ['dashboard', 'drilldown', mode, groupBy, groupValue, itemId, department, site],
    queryFn: async () => {
      if (!groupValue.value || !itemId.value) return []

      const params = new URLSearchParams()
      params.append('mode', mode.value)
      params.append('groupBy', groupBy.value)
      params.append('groupValue', groupValue.value)
      params.append('itemId', itemId.value)
      if (department?.value) params.append('department', department.value)
      if (site?.value) params.append('site', site.value)

      const url = `/dashboard/drilldown?${params.toString()}`
      return api<DrilldownRow[]>(url)
    },
    enabled,
    refetchOnMount: true,
    refetchOnWindowFocus: false,
  })
}
