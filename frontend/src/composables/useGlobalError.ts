import { useErrorStore } from '@/stores/errors'

export function useGlobalError() {
  const store = useErrorStore()

  function notify(message: string, detail?: string) {
    store.push(message, detail)
  }

  function clear() {
    store.clear()
  }

  return {
    errors: store.errors,
    notify,
    clear,
    dismiss: store.dismiss,
  }
}
