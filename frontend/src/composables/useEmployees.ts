import { useQuery } from '@tanstack/vue-query'
import { api } from '@/lib/api'
import type { Employee } from '@/types/api'
import employeesFixture from '@/data/employees.json'

const useFixtures = import.meta.env.VITE_USE_FIXTURES === 'true'

export function useEmployees() {
  return useQuery({
    queryKey: ['employees'],
    queryFn: async () => {
      if (useFixtures) {
        return employeesFixture as Employee[]
      }
      return api<Employee[]>('/employees')
    },
    staleTime: 60_000,
  })
}
