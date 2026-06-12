<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { Button, Modal } from '@/components/ui'
import { useAuditMeta } from '@/composables/useAuditMeta'
import { useLocaleFormat } from '@/composables/useLocaleFormat'
import type { AuditLogEntry } from '@/types/admin'
import { isValidObjectId } from '@/utils/objectId'

const props = defineProps<{
  open: boolean
  log: AuditLogEntry | null
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const { t, te } = useI18n()
const { formatDateTime } = useLocaleFormat()
const { auditMetaRows, rowLabel } = useAuditMeta()
const router = useRouter()

const metaRows = computed(() => auditMetaRows(props.log?.meta))

const bookingId = computed(() => {
  if (props.log?.entity === 'booking') return props.log.entityId
  return null
})

const bookingRef = computed(() => {
  const ref = props.log?.meta?.bookingRef
  return typeof ref === 'string' ? ref : null
})

const actorUserId = computed(() =>
  isValidObjectId(props.log?.actorId) ? props.log!.actorId! : null,
)

const canViewCustomerHistory = computed(() =>
  Boolean(actorUserId.value && (bookingId.value || bookingRef.value)),
)

function formatAuditAction(action: string) {
  const key = `admin.logs.auditActions.${action}`
  return te(key) ? t(key) : action
}

function viewCustomerHistory() {
  if (!actorUserId.value) return

  const query: Record<string, string> = {}
  if (bookingId.value) query.bookingId = bookingId.value
  else if (bookingRef.value) query.bookingRef = bookingRef.value

  emit('update:open', false)
  router.push({
    name: 'admin-user-bookings',
    params: { userId: actorUserId.value },
    query,
  })
}
</script>

<template>
  <Modal
    :open="open"
    size="xl"
    :title="log ? formatAuditAction(log.action) : t('admin.logs.detailTitle')"
    @update:open="emit('update:open', $event)"
  >
    <div v-if="log" class="space-y-4 text-sm">
      <dl class="grid gap-3 sm:grid-cols-2">
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('common.when') }}</dt>
          <dd>{{ formatDateTime(log.createdAt) }}</dd>
        </div>
        <div class="space-y-1">
          <dt class="text-copy-secondary">{{ t('admin.logs.entity') }}</dt>
          <dd>{{ log.entity }}</dd>
        </div>
        <div class="space-y-1 sm:col-span-2">
          <dt class="text-copy-secondary">{{ t('admin.logs.entityId') }}</dt>
          <dd class="break-all font-mono text-xs text-copy-muted">{{ log.entityId }}</dd>
        </div>
        <div v-if="actorUserId" class="space-y-1 sm:col-span-2">
          <dt class="text-copy-secondary">{{ t('admin.logs.actorId') }}</dt>
          <dd class="break-all font-mono text-xs text-copy-muted">{{ actorUserId }}</dd>
        </div>
      </dl>

      <div v-if="metaRows.length" class="space-y-2">
        <p class="text-copy-secondary">{{ t('admin.logs.details') }}</p>
        <dl class="grid gap-3 rounded-lg bg-subtle p-3 sm:grid-cols-2">
          <div
            v-for="row in metaRows"
            :key="row.key"
            class="space-y-1"
            :class="row.mono ? 'sm:col-span-2' : undefined"
          >
            <dt class="text-xs text-copy-secondary">{{ rowLabel(row) }}</dt>
            <dd
              class="text-sm text-copy-primary"
              :class="row.mono ? 'break-all font-mono text-xs text-copy-muted' : undefined"
            >
              {{ row.value }}
            </dd>
          </div>
        </dl>
      </div>

      <Button v-if="canViewCustomerHistory" variant="secondary" @click="viewCustomerHistory">
        {{ t('admin.logs.viewCustomerHistory') }}
      </Button>
    </div>
  </Modal>
</template>
