import { defineStore } from 'pinia'

type AppError = {
  id: number
  message: string
  detail?: string
}

let counter = 1

export const useErrorStore = defineStore('errors', {
  state: () => ({
    errors: [] as AppError[],
  }),
  actions: {
    push(message: string, detail?: string) {
      this.errors.push({ id: counter++, message, detail })
    },
    dismiss(id: number) {
      this.errors = this.errors.filter((err) => err.id !== id)
    },
    clear() {
      this.errors = []
    },
  },
})
