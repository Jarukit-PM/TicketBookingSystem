<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { translateApiError } from '@/api/errors'
import { ApiError, api } from '@/api/client'
import AuditLogDetailModal from '@/components/admin/AuditLogDetailModal.vue'
import TableSkeleton from '@/components/skeletons/TableSkeleton.vue'
import { Button, Card, CardContent, CardHeader, CardTitle, EmptyState, ErrorAlert, Input } from '@/components/ui'
import { CheckCircle2, Eye, FileText, Loader2, Mail } from 'lucide-vue-next'
import { useAuditMeta } from '@/composables/useAuditMeta'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { AuditLogEntry, EmailLogEntry } from '@/types/admin'

type LogTab = 'audit' | 'email'

const AUDIT_ACTIONS = [
  'booking_success',
  'booking_timeout',
  'seat_released',
  'booking_failed',
  'system_error',
  'create',
  'update',
  'delete',
] as const

const AUDIT_ENTITIES = ['booking', 'showtime', 'movie', 'cinema', 'screen'] as const

const EMAIL_TYPES = ['CONFIRMATION'] as const
const EMAIL_STATUSES = ['SENT', 'FAILED'] as const

const selectClass =
  'w-full rounded-lg border border-surface-border bg-surface px-3 py-2 text-sm text-copy-primary focus:outline-none focus:ring-2 focus:ring-accent-glow focus:border-brand/50'

const { t, te } = useI18n()
const { formatDateTime } = useLocaleFormat()
const { formatAuditSummary } = useAuditMeta()

const activeTab = ref<LogTab>('audit')
const bookingId = ref('')
const emailTo = ref('')
const emailType = ref('')
const emailStatus = ref('')
const emailSentFrom = ref('')
const emailSentTo = ref('')
const auditAction = ref('')
const auditEntity = ref('')
const auditEntityId = ref('')
const auditActorId = ref('')
const auditBookingRef = ref('')
const auditLogs = ref<AuditLogEntry[]>([])
const emailLogs = ref<EmailLogEntry[]>([])
const loading = ref(false)
const errorMessage = ref('')
const toastMessage = ref('')
const resendingBookingId = ref('')

let toastTimer: ReturnType<typeof setTimeout> | null = null

function showToast(message: string): void {
  toastMessage.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 4000)
}

onUnmounted(() => {
  if (toastTimer) clearTimeout(toastTimer)
})
const page = ref(1)
const limit = 50
const auditDetailOpen = ref(false)
const selectedAuditLog = ref<AuditLogEntry | null>(null)

const currentLogs = computed(() =>
  activeTab.value === 'audit' ? auditLogs.value : emailLogs.value,
)
const canGoPrev = computed(() => page.value > 1)
const canGoNext = computed(() => currentLogs.value.length === limit)
const showPagination = computed(
  () => !loading.value && currentLogs.value.length > 0 && (canGoPrev.value || canGoNext.value),
)

const hasAuditFilters = computed(() =>
  Boolean(
    auditAction.value ||
      auditEntity.value ||
      auditEntityId.value.trim() ||
      auditActorId.value.trim() ||
      auditBookingRef.value.trim(),
  ),
)

const hasEmailFilters = computed(() =>
  Boolean(
    bookingId.value.trim() ||
      emailTo.value.trim() ||
      emailType.value ||
      emailStatus.value ||
      emailSentFrom.value.trim() ||
      emailSentTo.value.trim(),
  ),
)

const emptyMessage = computed(() => {
  if (activeTab.value === 'email') {
    return hasEmailFilters.value ? t('admin.logs.emptyEmailFiltered') : t('admin.logs.emptyEmail')
  }
  return hasAuditFilters.value ? t('admin.logs.emptyAuditFiltered') : t('admin.logs.emptyAudit')
})

function formatAuditAction(action: string) {
  const key = `admin.logs.auditActions.${action}`
  return te(key) ? t(key) : action
}

function showAuditDetail(log: AuditLogEntry) {
  selectedAuditLog.value = log
  auditDetailOpen.value = true
}

async function resendEmail(bookingId: string) {
  if (resendingBookingId.value) return
  resendingBookingId.value = bookingId
  errorMessage.value = ''
  try {
    await api.post(`/admin/bookings/${bookingId}/resend-email`, {})
    await loadLogs()
    showToast(t('admin.logs.resendQueued'))
  } catch (error) {
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.logs.resendFailed')
  } finally {
    resendingBookingId.value = ''
  }
}

async function loadLogs() {
  loading.value = true
  errorMessage.value = ''
  try {
    if (activeTab.value === 'audit') {
      const params = new URLSearchParams({
        page: String(page.value),
        limit: String(limit),
      })
      if (auditAction.value) params.set('action', auditAction.value)
      if (auditEntity.value) params.set('entity', auditEntity.value)
      if (auditEntityId.value.trim()) params.set('entityId', auditEntityId.value.trim())
      if (auditActorId.value.trim()) params.set('actorId', auditActorId.value.trim())
      if (auditBookingRef.value.trim()) params.set('bookingRef', auditBookingRef.value.trim())
      const data = await api.get<{ logs: AuditLogEntry[] }>(`/admin/audit-logs?${params}`)
      auditLogs.value = data.logs ?? []
    } else {
      const params = new URLSearchParams({
        page: String(page.value),
        limit: String(limit),
      })
      if (bookingId.value.trim()) params.set('bookingId', bookingId.value.trim())
      if (emailTo.value.trim()) params.set('to', emailTo.value.trim())
      if (emailType.value) params.set('type', emailType.value)
      if (emailStatus.value) params.set('status', emailStatus.value)
      if (emailSentFrom.value.trim()) params.set('sentFrom', emailSentFrom.value.trim())
      if (emailSentTo.value.trim()) params.set('sentTo', emailSentTo.value.trim())
      const data = await api.get<{ logs: EmailLogEntry[] }>(`/admin/email-logs?${params}`)
      emailLogs.value = data.logs ?? []
    }
  } catch (error) {
    auditLogs.value = []
    emailLogs.value = []
    errorMessage.value =
      error instanceof ApiError
        ? translateApiError(error.code, error.message)
        : t('admin.logs.loadFailed')
  } finally {
    loading.value = false
  }
}

function setTab(tab: LogTab) {
  activeTab.value = tab
  page.value = 1
  auditDetailOpen.value = false
  selectedAuditLog.value = null
}

function applyEmailFilter() {
  page.value = 1
  void loadLogs()
}

function clearEmailFilters() {
  bookingId.value = ''
  emailTo.value = ''
  emailType.value = ''
  emailStatus.value = ''
  emailSentFrom.value = ''
  emailSentTo.value = ''
  page.value = 1
  void loadLogs()
}

function applyAuditFilter() {
  page.value = 1
  void loadLogs()
}

function clearAuditFilters() {
  auditAction.value = ''
  auditEntity.value = ''
  auditEntityId.value = ''
  auditActorId.value = ''
  auditBookingRef.value = ''
  page.value = 1
  void loadLogs()
}

function goPrev() {
  if (!canGoPrev.value) return
  page.value -= 1
  void loadLogs()
}

function goNext() {
  if (!canGoNext.value) return
  page.value += 1
  void loadLogs()
}

watch(activeTab, loadLogs, { immediate: true })
</script>

<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-copy-primary">{{ t('admin.logs.title') }}</h1>
      <p class="mt-1 text-sm text-copy-secondary">{{ t('admin.logs.subtitle') }}</p>
    </div>

    <div class="flex flex-wrap gap-2">
      <Button :variant="activeTab === 'audit' ? 'primary' : 'secondary'" @click="setTab('audit')">
        {{ t('admin.logs.auditTab') }}
      </Button>
      <Button :variant="activeTab === 'email' ? 'primary' : 'secondary'" @click="setTab('email')">
        {{ t('admin.logs.emailTab') }}
      </Button>
    </div>

    <Card v-if="activeTab === 'audit'">
      <CardHeader>
        <CardTitle>{{ t('admin.logs.auditFilterTitle') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 sm:grid-cols-2">
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.auditAction') }}</span>
            <select v-model="auditAction" :class="selectClass">
              <option value="">{{ t('admin.logs.auditActionAll') }}</option>
              <option v-for="action in AUDIT_ACTIONS" :key="action" :value="action">
                {{ formatAuditAction(action) }}
              </option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.auditEntity') }}</span>
            <select v-model="auditEntity" :class="selectClass">
              <option value="">{{ t('admin.logs.auditEntityAll') }}</option>
              <option v-for="entity in AUDIT_ENTITIES" :key="entity" :value="entity">
                {{ entity }}
              </option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.entityId') }}</span>
            <Input
              v-model="auditEntityId"
              :placeholder="t('admin.logs.bookingIdPlaceholder')"
              autocomplete="off"
            />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.auditBookingRef') }}</span>
            <Input
              v-model="auditBookingRef"
              :placeholder="t('admin.logs.auditBookingRefPlaceholder')"
              autocomplete="off"
            />
          </label>
          <label class="block space-y-1.5 sm:col-span-2">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.auditActorId') }}</span>
            <Input
              v-model="auditActorId"
              :placeholder="t('admin.logs.auditActorIdPlaceholder')"
              autocomplete="off"
            />
          </label>
        </div>
        <div class="flex flex-wrap gap-3">
          <Button :disabled="loading" @click="applyAuditFilter">{{ t('common.apply') }}</Button>
          <Button variant="secondary" :disabled="loading" @click="clearAuditFilters">
            {{ t('common.clear') }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <Card v-if="activeTab === 'email'">
      <CardHeader>
        <CardTitle>{{ t('admin.logs.emailFilterTitle') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 sm:grid-cols-2">
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.bookingId') }}</span>
            <Input v-model="bookingId" :placeholder="t('admin.logs.bookingIdPlaceholder')" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.emailTo') }}</span>
            <Input
              v-model="emailTo"
              type="email"
              :placeholder="t('admin.logs.emailToPlaceholder')"
              autocomplete="off"
            />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.type') }}</span>
            <select v-model="emailType" :class="selectClass">
              <option value="">{{ t('admin.logs.emailTypeAll') }}</option>
              <option v-for="type in EMAIL_TYPES" :key="type" :value="type">
                {{ type }}
              </option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('common.status') }}</span>
            <select v-model="emailStatus" :class="selectClass">
              <option value="">{{ t('admin.logs.emailStatusAll') }}</option>
              <option v-for="status in EMAIL_STATUSES" :key="status" :value="status">
                {{ status }}
              </option>
            </select>
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.sentFrom') }}</span>
            <Input v-model="emailSentFrom" type="date" autocomplete="off" />
          </label>
          <label class="block space-y-1.5">
            <span class="text-sm text-copy-secondary">{{ t('admin.logs.sentTo') }}</span>
            <Input v-model="emailSentTo" type="date" autocomplete="off" />
          </label>
        </div>
        <div class="flex flex-wrap gap-3">
          <Button :disabled="loading" @click="applyEmailFilter">{{ t('common.apply') }}</Button>
          <Button variant="secondary" :disabled="loading" @click="clearEmailFilters">
            {{ t('common.clear') }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <ErrorAlert v-if="errorMessage" :message="errorMessage" />

    <Teleport to="body">
      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="translate-y-2 opacity-0"
        enter-to-class="translate-y-0 opacity-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="translate-y-0 opacity-100"
        leave-to-class="translate-y-2 opacity-0"
      >
        <div
          v-if="toastMessage"
          class="pointer-events-none fixed right-4 top-4 z-50 flex max-w-sm items-center gap-2 rounded-lg border border-state-success/40 bg-state-success-dim px-4 py-3 text-sm text-state-success shadow-1"
          role="status"
          aria-live="polite"
        >
          <CheckCircle2 class="h-4 w-4 shrink-0" aria-hidden="true" />
          <span>{{ toastMessage }}</span>
        </div>
      </Transition>
    </Teleport>

    <Card>
      <CardHeader class="flex flex-row flex-wrap items-center justify-between gap-3">
        <CardTitle>{{ activeTab === 'audit' ? t('admin.logs.auditEntries') : t('admin.logs.emailDeliveries') }}</CardTitle>
        <p v-if="showPagination" class="text-sm text-copy-muted">
          {{ t('admin.logs.pageInfo', { page }) }}
        </p>
      </CardHeader>
      <CardContent class="space-y-4">
        <TableSkeleton v-if="loading" :columns="activeTab === 'email' ? 6 : 6" :rows="6" />
        <EmptyState
          v-else-if="(activeTab === 'audit' && !auditLogs.length) || (activeTab === 'email' && !emailLogs.length)"
          :icon="FileText"
          :title="emptyMessage"
          class="py-10"
        />
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="sticky top-0 bg-surface text-copy-muted">
              <tr v-if="activeTab === 'audit'">
                <th class="pb-3 pr-4 font-medium">{{ t('common.when') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.action') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.entity') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.entityId') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.details') }}</th>
                <th class="pb-3 font-medium">{{ t('admin.logs.actions') }}</th>
              </tr>
              <tr v-else>
                <th class="pb-3 pr-4 font-medium">{{ t('common.when') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.to') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.type') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('common.status') }}</th>
                <th class="pb-3 pr-4 font-medium">{{ t('admin.logs.booking') }}</th>
                <th class="pb-3 font-medium">{{ t('admin.logs.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-if="activeTab === 'audit'">
                <tr
                  v-for="log in auditLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ formatDateTime(log.createdAt) }}
                  </td>
                  <td class="py-3 pr-4 font-medium text-copy-primary">
                    {{ formatAuditAction(log.action) }}
                  </td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.entity }}</td>
                  <td class="py-3 pr-4 font-mono text-xs text-copy-muted">{{ log.entityId }}</td>
                  <td class="py-3 pr-4 text-xs text-copy-secondary">{{ formatAuditSummary(log.meta) }}</td>
                  <td class="py-3 whitespace-nowrap">
                    <Button
                      variant="secondary"
                      class="h-10 w-10 shrink-0 p-0 sm:h-auto sm:w-auto sm:gap-1.5 sm:px-4 sm:py-2.5"
                      :aria-label="t('admin.logs.viewDetail')"
                      @click="showAuditDetail(log)"
                    >
                      <Eye class="h-4 w-4 shrink-0" aria-hidden="true" />
                      <span class="hidden sm:inline">{{ t('admin.logs.viewDetail') }}</span>
                    </Button>
                  </td>
                </tr>
              </template>
              <template v-else>
                <tr
                  v-for="log in emailLogs"
                  :key="log.id"
                  class="border-t border-surface-border"
                >
                  <td class="py-3 pr-4 text-copy-secondary">
                    {{ formatDateTime(log.createdAt) }}
                  </td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.to }}</td>
                  <td class="py-3 pr-4 text-copy-secondary">{{ log.type }}</td>
                  <td class="py-3 pr-4 text-copy-primary">{{ log.status }}</td>
                  <td class="py-3 pr-4 font-mono text-xs text-copy-muted">{{ log.bookingId }}</td>
                  <td class="py-3 whitespace-nowrap">
                    <Button
                      variant="secondary"
                      class="h-10 w-10 shrink-0 p-0 sm:h-auto sm:w-auto sm:gap-1.5 sm:px-4 sm:py-2.5"
                      :disabled="!!resendingBookingId"
                      :aria-busy="resendingBookingId === log.bookingId"
                      :aria-label="
                        resendingBookingId === log.bookingId
                          ? t('admin.logs.resending')
                          : t('admin.logs.resend')
                      "
                      @click="resendEmail(log.bookingId)"
                    >
                      <Loader2
                        v-if="resendingBookingId === log.bookingId"
                        class="h-4 w-4 shrink-0 animate-spin"
                        aria-hidden="true"
                      />
                      <Mail v-else class="h-4 w-4 shrink-0" aria-hidden="true" />
                      <span class="hidden whitespace-nowrap sm:inline">
                        {{
                          resendingBookingId === log.bookingId
                            ? t('admin.logs.resending')
                            : t('admin.logs.resend')
                        }}
                      </span>
                    </Button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
        <div v-if="showPagination" class="flex flex-wrap items-center justify-end gap-3">
          <Button variant="secondary" :disabled="loading || !canGoPrev" @click="goPrev">
            {{ t('common.previous') }}
          </Button>
          <Button variant="secondary" :disabled="loading || !canGoNext" @click="goNext">
            {{ t('common.next') }}
          </Button>
        </div>
      </CardContent>
    </Card>

    <AuditLogDetailModal v-model:open="auditDetailOpen" :log="selectedAuditLog" />
  </div>
</template>
