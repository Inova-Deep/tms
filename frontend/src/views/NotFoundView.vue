<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useRouter } from 'vue-router'
import { Home } from 'lucide-vue-next'

const router = useRouter()
</script>

<template>
  <section class="page-section">
    <div class="card-surface text-center">
      <h1 class="page-title">Page not found</h1>
      <p class="page-description mt-1">The page you are looking for does not exist.</p>
      <Button variant="outline" size="sm" class="mt-4" @click="router.push('/dashboard')">
        <Home class="mr-2 h-4 w-4" />
        Go to Dashboard
      </Button>
    </div>
  </section>
</template>
