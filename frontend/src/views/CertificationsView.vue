<script setup lang="ts">
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { FileBadge2 } from 'lucide-vue-next'
</script>

<template>
  <section class="page-section">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header-left">
        <h1 class="page-title">Certifications</h1>
        <p class="page-description">Manage employee certifications and uploaded documents.</p>
      </div>
    </div>

    <Card>
      <CardHeader>
        <CardTitle class="form-label">Certification Management</CardTitle>
        <CardDescription>Upload and manage certification documents for employees</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="empty-state h-[300px]">
          <div class="empty-state-content">
            <div class="empty-state-icon">
              <FileBadge2 class="empty-state-icon-inner" />
            </div>
            <h3 class="empty-state-title">Coming Soon</h3>
            <p class="empty-state-description">
              Certification upload and management features are being developed.
            </p>
            <Badge variant="secondary">In Development</Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
