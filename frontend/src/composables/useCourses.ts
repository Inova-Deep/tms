import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import type { Ref, MaybeRef } from 'vue'
import { unref } from 'vue'
import { api } from '@/lib/api'
import type { Course, CreateCourseRequest, UpdateCourseRequest } from '@/types/api'

export function useCourses(type?: MaybeRef<'course' | 'certification' | undefined>) {
  return useQuery({
    queryKey: ['courses', type],
    queryFn: async () => {
      const courses = await api<Course[]>('/courses')
      const typeVal = unref(type)
      if (typeVal) {
        return courses.filter(c => c.type === typeVal)
      }
      return courses
    },
    staleTime: 60_000,
  })
}

export function useCourse(id: MaybeRef<string | null | undefined>) {
  return useQuery({
    queryKey: ['course', id],
    queryFn: () => {
      const idVal = unref(id)
      if (!idVal) return Promise.resolve({} as Course)
      return api<Course>(`/courses/${idVal}`)
    },
    enabled: () => {
      const idVal = unref(id)
      return !!idVal
    },
    staleTime: 30_000,
  })
}

export function useCreateCourse() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateCourseRequest) =>
      api<Course>('/courses', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })
}

export function useUpdateCourse() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCourseRequest }) =>
      api<Course>(`/courses/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })
}

export function useDeleteCourse() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) =>
      api<void>(`/courses/${id}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['courses'] })
    },
  })
}
