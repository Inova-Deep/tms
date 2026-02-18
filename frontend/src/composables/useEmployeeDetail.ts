import { useQuery } from '@tanstack/vue-query'
import type { Ref } from 'vue'
import { api } from '@/lib/api'
import type { EmployeeDetail, EvidenceRecord } from '@/types/api'
import employeeFixture from '@/data/employee-E1005.json'
import recordsFixture from '@/data/employee-E1005-records.json'

const useFixtures = import.meta.env.VITE_USE_FIXTURES !== 'false'

type IdInput = string | undefined | Ref<string | undefined> | (() => string | undefined)

function resolveId(id: IdInput): string | undefined {
  if (typeof id === 'function') return id()
  if (typeof id === 'object' && 'value' in id) return id.value
  return id
}

export function useEmployeeDetail(id: IdInput) {
  return useQuery({
    queryKey: ['employee', typeof id === 'function' ? id : typeof id === 'object' && 'value' in id ? id : id],
    queryFn: () => {
      const idVal = resolveId(id)
      return useFixtures && idVal === 'E1005'
        ? Promise.resolve(employeeFixture as EmployeeDetail)
        : api<EmployeeDetail>(`/employees/${idVal}`)
    },
    enabled: () => !!resolveId(id),
    staleTime: 30_000,
  })
}

export function useEmployeeRecords(id: IdInput) {
  return useQuery({
    queryKey: ['employee-records', typeof id === 'function' ? id : typeof id === 'object' && 'value' in id ? id : id],
    queryFn: () => {
      const idVal = resolveId(id)
      return useFixtures && idVal === 'E1005'
        ? Promise.resolve(recordsFixture as EvidenceRecord[])
        : api<EvidenceRecord[]>(`/employees/${idVal}/records`)
    },
    enabled: () => !!resolveId(id),
    staleTime: 30_000,
  })
}
