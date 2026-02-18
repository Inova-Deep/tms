import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { api } from '@/lib/api'
import type { TrainingEvent, CreateEventRequest, UpdateEventRequest, ConfirmAttendanceRequest } from '@/types/api'

export function useEvents() {
  return useQuery({
    queryKey: ['events'],
    queryFn: async () => {
      return api<TrainingEvent[]>('/events')
    },
  })
}

export function useCreateEvent() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (event: CreateEventRequest) => {
      return api<TrainingEvent>('/events', {
        method: 'POST',
        body: JSON.stringify(event),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
    },
  })
}

export function useUpdateEvent() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateEventRequest }) => {
      return api<TrainingEvent>(`/events/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
    },
  })
}

export function useConfirmAttendance() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ eventId, data }: { eventId: string; data: ConfirmAttendanceRequest }) => {
      return api<{ status: string }>(`/events/${eventId}/attendance/confirm`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['events'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    },
  })
}
