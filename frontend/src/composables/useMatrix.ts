import { useQuery } from '@tanstack/vue-query'
import { api } from '@/lib/api'
import type { MatrixResponse } from '@/types/api'

export function useMatrix() {
  return useQuery({
    queryKey: ['matrix'],
    queryFn: () => api<MatrixResponse>('/matrix'),
  })
}
