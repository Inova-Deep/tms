import { defineStore } from 'pinia'
import { api } from '@/lib/api'

type Role = 'admin' | 'employee'

interface LoginResponse {
  token: string
  role: Role
  employeeId: string
}

interface EmployeeDetail {
  id: string
  name: string
  department: string
  site: string
}

let initialized = false

export const useAuthStore = defineStore('auth', {
  state: () => ({
    role: 'admin' as Role,
    employeeId: '' as string,
    name: 'Training Admin' as string,
    department: '' as string,
    site: '' as string,
    token: '' as string,
    loading: false,
    error: '' as string,
    initialized: false,
  }),
  getters: {
    isEmployee: (state) => state.role === 'employee' && !!state.employeeId,
    is_admin: (state) => state.role === 'admin',
    isAuthenticated: (state) => !!state.token,
    displayLabel: (state) => (state.role === 'employee' && state.employeeId ? state.employeeId : 'Admin'),
    userInitials: (state) => {
      const name = state.name || 'U'
      return name
        .split(' ')
        .map((part) => part.charAt(0))
        .join('')
        .slice(0, 2)
        .toUpperCase()
    },
  },
  actions: {
    async init() {
      if (this.initialized || this.token) return
      
      this.loading = true
      try {
        const res = await api<LoginResponse>('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ role: 'admin' }),
        })
        this.role = res.role
        this.employeeId = res.employeeId
        this.name = 'Training Admin'
        this.department = ''
        this.site = ''
        this.token = res.token
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Init failed'
      } finally {
        this.loading = false
        this.initialized = true
      }
    },

    async loginAsAdmin() {
      this.loading = true
      this.error = ''
      try {
        const res = await api<LoginResponse>('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ role: 'admin' }),
        })
        this.role = 'admin'
        this.employeeId = ''
        this.name = 'Training Admin'
        this.department = ''
        this.site = ''
        this.token = res.token
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Login failed'
        throw err
      } finally {
        this.loading = false
      }
    },

    async loginAsEmployee(employeeId: string) {
      this.loading = true
      this.error = ''
      try {
        const res = await api<LoginResponse>('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ role: 'employee', employeeId }),
        })
        this.role = 'employee'
        this.employeeId = res.employeeId
        this.token = res.token
        
        const empDetail = await api<EmployeeDetail>(`/employees/${res.employeeId}`)
        this.name = empDetail.name
        this.department = empDetail.department
        this.site = empDetail.site
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'Login failed'
        throw err
      } finally {
        this.loading = false
      }
    },

    clear() {
      this.role = 'admin'
      this.employeeId = ''
      this.name = 'Training Admin'
      this.department = ''
      this.site = ''
      this.token = ''
      this.error = ''
    },

    setSession(payload: { role: Role; employeeId?: string; name?: string; token: string }) {
      this.role = payload.role
      this.employeeId = payload.employeeId ?? ''
      this.name = payload.name ?? (payload.role === 'admin' ? 'Training Admin' : payload.employeeId ?? 'Employee')
      this.token = payload.token
    },
  },
})
