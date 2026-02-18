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

    <Card class="card-surface">
      <CardContent class="p-4">
        <div class="space-y-3">
          <div>
            <h3 class="form-label mb-1">What is this feature?</h3>
            <p class="page-description">Record and track external certifications obtained outside your company.</p>
          </div>
          <div>
            <h3 class="form-label mb-1">What problem does it solve?</h3>
            <p class="page-description">Keep track of professional qualifications, licenses, and certifications with expiry dates.</p>
          </div>
          <div>
            <h3 class="form-label mb-2">How it works</h3>
            <ol class="page-description space-y-1 list-decimal list-inside">
              <li>Add a certification record when an employee obtains one</li>
              <li>Enter the issuer, issue date, and expiry date</li>
              <li>TMS tracks when it expires and alerts you</li>
            </ol>
          </div>
        </div>
      </CardContent>
    </Card>

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
