import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient, QueryCache, MutationCache } from '@tanstack/vue-query'

import App from './App.vue'
import router from './router'
import '@/assets/main.css'
import { useErrorStore } from '@/stores/errors'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => {
      const store = useErrorStore()
      store.push('Request failed', error instanceof Error ? error.message : String(error))
    },
  }),
  mutationCache: new MutationCache({
    onError: (error) => {
      const store = useErrorStore()
      store.push('Action failed', error instanceof Error ? error.message : String(error))
    },
  }),
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
    },
  },
})

app.use(VueQueryPlugin, { queryClient })
app.use(router)

app.config.errorHandler = (err, _instance, info) => {
  const store = useErrorStore()
  store.push('Unexpected error', info ?? (err instanceof Error ? err.message : String(err)))
  console.error(err)
}

app.mount('#app')
