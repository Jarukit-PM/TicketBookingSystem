<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import { ApiError, api } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import type { AuditLogEntry, EmailLogEntry } from '@/types/admin'

type LogTab = 'audit' | 'email'

const activeTab = ref<LogTab>('audit')
const bookingId = ref('')
const auditLogs = ref<AuditLogEntry[]>([])
const emailLogs = ref<EmailLogEntry[]>([])
const loading = ref(false)
const errorMessage = ref('')
const page = ref(1)
const limit = 50

const emptyMessage = computed(() =>
  activeTab.value === 'audit'
    ? 'No audit log entries yet.'
    : 'No email log entries match your filter.',
)

async function loadLogs() {
  loading.value = true
  errorMessage.value = ''
  try {
    if (activeTab.value === 'audit') {
      const data = await api.get<{ logs: AuditLogEntry[] }>(
        `/admin/audit-logs?page=${page.value}&limit=${limit}`,
      )
      auditLogs.value = data.logs
    } else {
      const params = new URLSearchParams({
        page: String(page.value),
        limit: String(limit),
      })
      if (bookingId.value.trim()) params.set('bookingId', bookingId.value.trim())
      const data = await api.get<{ logs: EmailLogEntry[] }>(`/admin/email-logs?${params}`)
      emailLogs.value = data.logs
    }
  } catch (error) {
    auditLogs.value = []
    emailLogs.value = []
    errorMessage.value = error instanceof ApiError ? error.message : 'Failed to load logs'
  } finally {
    loading.value = false
  }
}

function setTab(tab: LogTab) {
  activeTab.value = tab
  page.value = 1
}

watch(activeTab, loadLogs, { immediate: true })
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">Logs</h1>
      <p class="mt-1 text-sm text-copy-secondary">
        Operational audit trail and email delivery status (newest first).
      </p>
    </div>

    <div class="flex flex-wrap gap-2">
      <Button :variant="activeTab === 'audit' ? 'primary' : 'secondary'" @click="setTab('audit')">
        Audit log
      </Button>
      <Button :variant="activeTab === 'email' ? 'primary' : 'secondary'" @click="setTab('email')">
        Email log
      </Button>
    </div>

    <Card v-if="activeTab === 'email'">
      <CardHeader>
        <CardTitle>Filter</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-wrap items-end gap-3">
        <label class="block min-w-[16rem] flex-1 space-y-1.5">
          <span class="text-sm text-copy-secondary">Booking ID</span>
          <Input v-model="bookingId" placeholder="MongoDB ObjectId (optional)" autocomplete="off" />
        </label>
        <Button :disabled="loading" @click="loadLogs">Apply</Button>
      </CardContent>
    </Card>

    <p v-if="errorMessage" class="text-sm text-state-error" role="alert">{{ errorMessage }}</p>

    <Card>
      <CardHeader>
        <CardTitle>{{ activeTab === 'audit' ? 'Audit entries' : 'Email deliveries' }}</CardTitle>
      </CardHeader>
      <CardContent>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="sticky top-0 bg-surface text-copy-muted">
              <tr v-if="activeTab === 'audit'">
                <th class="pb-3 pr-4 font-medium">When</th>
                <th class="pb-3 pr-4 font-medium">Action</th>
                <th class="pb-3 pr-4 font-medium">Entity</th>
                <th class="pb-3 font-medium">Entity ID</th>
              </tr>
              <tr v-else>
                <th class="pb-3 pr-4 font-medium">When</th>
                <th class="pb-3 pr-4 font-medium">To</th>
                <th class="pb-3 pr-4 font-medium">Type</th>
                <th class="pb-3 pr-4 font-medium">Status</th>
                <th class="pb-3 font-medium">Booking</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading">
                <td :colspan="activeTab === 'audit' ? 4 : 5" class="py-6 text-copy-muted">Loading…</td>
              </tr>
              <tr v-else-if="activeTab === 'audit' && !auditLogs.length">
                <td colspan="4" class="py-6 text-copy-muted">{{ emptyMessage }}</td>
              </tr>
              <tr v-else-if="activeTab === 'email' && !emailLogs.length">
                <td colspan="5" class="py-6 text-copy-muted">{{ emptyMessage }}</td>
              </tr>
              <template v-else-if="activeTab === 'audit'">
                <tr
                  v-for="log in auditLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ new Date(log.createdAt).toLocaleString() }}
                  </td>
                  <td class="py-3 pr-4 font-medium text-copy-primary">{{ log.action }}</td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.entity }}</td>
                  <td class="py-3 font-mono text-xs text-copy-muted">{{ log.entityId }}</td>
                </tr>
              </template>
              <template v-else>
                <tr
                  v-for="log in emailLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ new Date(log.createdAt).toLocaleString() }}
                  </td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.to }}</td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.type }}</td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.status }}</td>
                  <td class="py-3 font-mono text-xs text-copy-muted">{{ log.bookingId }}</td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
