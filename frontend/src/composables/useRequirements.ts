import { useQuery } from '@tanstack/vue-query'
import { api } from '@/lib/api'
import type { RequirementWithStatus } from '@/types/api'

export function useRequirements() {
  return useQuery({
    queryKey: ['requirements'],
    queryFn: async () => {
      // For now, we only fetch P1 (Safety Base) requirements as they are the primary courses
      // In a real app, we might have a dedicated /requirements endpoint
      return api<RequirementWithStatus[]>('/profiles/P1/requirements')
    },
  })
}
