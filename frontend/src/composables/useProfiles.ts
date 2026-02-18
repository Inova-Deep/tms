import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import type { Ref } from 'vue'
import { api } from '@/lib/api'
import type { Profile, Requirement, ProfileMember, ProfileExclusion, AddExclusionRequest, UpdateExclusionRequest } from '@/types/api'

export function useProfiles() {
  return useQuery({
    queryKey: ['profiles'],
    queryFn: () => api<Profile[]>('/profiles'),
    staleTime: 60_000,
  })
}

export function useProfile(id: Ref<string | null | undefined> | (() => string | null | undefined)) {
  return useQuery({
    queryKey: ['profile', id],
    queryFn: () => {
      const idVal = typeof id === 'function' ? id() : id.value
      if (!idVal) return Promise.resolve({} as Profile)
      return api<Profile>(`/profiles/${idVal}`)
    },
    enabled: () => {
      const idVal = typeof id === 'function' ? id() : id.value
      return !!idVal
    },
    staleTime: 30_000,
  })
}

export function useProfileRequirements(id: Ref<string | null | undefined> | (() => string | null | undefined)) {
  return useQuery({
    queryKey: ['profile-requirements', id],
    queryFn: () => {
      const idVal = typeof id === 'function' ? id() : id.value
      if (!idVal) return Promise.resolve([] as Requirement[])
      return api<Requirement[]>(`/profiles/${idVal}/requirements`)
    },
    enabled: () => {
      const idVal = typeof id === 'function' ? id() : id.value
      return !!idVal
    },
    staleTime: 30_000,
  })
}

export function useProfileMembers(profileId: Ref<string | null | undefined> | (() => string | null | undefined)) {
  return useQuery({
    queryKey: ['profile-members', profileId],
    queryFn: () => {
      const idVal = typeof profileId === 'function' ? profileId() : profileId.value
      if (!idVal) return Promise.resolve({ members: [], total: 0 })
      return api<{ members: ProfileMember[]; total: number }>(`/profiles/${idVal}/members`)
    },
    enabled: () => {
      const idVal = typeof profileId === 'function' ? profileId() : profileId.value
      return !!idVal
    },
    staleTime: 30_000,
  })
}

export function useAddProfileMember() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, employeeId }: { profileId: string; employeeId: string }) =>
      api<{ status: string }>(`/profiles/${profileId}/members`, {
        method: 'POST',
        body: JSON.stringify({ employeeId }),
      }),
    onSuccess: (_, { profileId }) => {
      queryClient.invalidateQueries({ queryKey: ['profiles'] })
      queryClient.invalidateQueries({ queryKey: ['profile-members', profileId] })
    },
  })
}

export function useRemoveProfileMember() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, employeeId }: { profileId: string; employeeId: string }) =>
      api<{ status: string }>(`/profiles/${profileId}/members/${employeeId}`, {
        method: 'DELETE',
      }),
    onSuccess: (_, { profileId }) => {
      queryClient.invalidateQueries({ queryKey: ['profiles'] })
      queryClient.invalidateQueries({ queryKey: ['profile-members', profileId] })
    },
  })
}

export function useProfileExclusions(profileId: Ref<string | null | undefined> | (() => string | null | undefined)) {
  return useQuery({
    queryKey: ['profile-exclusions', profileId],
    queryFn: () => {
      const idVal = typeof profileId === 'function' ? profileId() : profileId.value
      if (!idVal) return Promise.resolve([])
      return api<ProfileExclusion[]>(`/profiles/${idVal}/exclusions`)
    },
    enabled: () => {
      const idVal = typeof profileId === 'function' ? profileId() : profileId.value
      return !!idVal
    },
    staleTime: 30_000,
  })
}

export function useAddProfileExclusion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, ...data }: { profileId: string } & AddExclusionRequest) =>
      api<ProfileExclusion>(`/profiles/${profileId}/exclusions`, {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { profileId }) => {
      queryClient.invalidateQueries({ queryKey: ['profiles'] })
      queryClient.invalidateQueries({ queryKey: ['profile-exclusions', profileId] })
    },
  })
}

export function useUpdateProfileExclusion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, exclusionId, ...data }: { profileId: string; exclusionId: string } & UpdateExclusionRequest) =>
      api<ProfileExclusion>(`/profiles/${profileId}/exclusions/${exclusionId}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    onSuccess: (_, { profileId }) => {
      queryClient.invalidateQueries({ queryKey: ['profile-exclusions', profileId] })
    },
  })
}

export function useRemoveProfileExclusion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, exclusionId }: { profileId: string; exclusionId: string }) =>
      api<{ status: string }>(`/profiles/${profileId}/exclusions/${exclusionId}`, {
        method: 'DELETE',
      }),
    onSuccess: (_, { profileId }) => {
      queryClient.invalidateQueries({ queryKey: ['profiles'] })
      queryClient.invalidateQueries({ queryKey: ['profile-exclusions', profileId] })
    },
  })
}

export function useCreateProfile() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: { name: string; description?: string }) => {
      return api<Profile>('/profiles', {
        method: 'POST',
        body: JSON.stringify(data),
      })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profiles'] })
    },
  })
}

export function useAddProfileRequirement() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, courseId }: { profileId: string; courseId: string }) =>
      api<Requirement>(`/profiles/${profileId}/requirements`, {
        method: 'POST',
        body: JSON.stringify({ courseId }),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile-requirements'] })
      queryClient.invalidateQueries({ queryKey: ['profile'] })
    },
  })
}

export function useRemoveProfileRequirement() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ profileId, requirementId }: { profileId: string; requirementId: string }) =>
      api<{ status: string }>(`/profiles/${profileId}/requirements/${requirementId}`, {
        method: 'DELETE',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile-requirements'] })
      queryClient.invalidateQueries({ queryKey: ['profile'] })
    },
  })
}
