<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { AlertTriangle, X } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { useErrorStore } from '@/stores/errors'

const errorStore = useErrorStore()
const { errors } = storeToRefs(errorStore)
</script>

<template>
  <div v-if="errors.length" class="error-bar">
    <div class="error-bar-content">
      <div
        v-for="err in errors"
        :key="err.id"
        class="error-bar-item"
      >
        <div class="error-bar-message">
          <AlertTriangle class="h-4 w-4 mt-0.5" />
          <div>
            <p class="error-bar-text">{{ err.message }}</p>
            <p v-if="err.detail" class="error-bar-detail">{{ err.detail }}</p>
          </div>
        </div>
        <Button variant="ghost" size="icon" class="error-bar-dismiss" @click="errorStore.dismiss(err.id)">
          <X class="h-4 w-4" />
          <span class="sr-only">Dismiss</span>
        </Button>
      </div>
    </div>
  </div>
</template>
