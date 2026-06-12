import { useI18n } from 'vue-i18n'

import { useLocaleFormat } from '@/composables/useLocaleFormat'
import {
  auditMetaRowLabel,
  buildAuditMetaRows,
  type AuditMetaFieldKey,
  type AuditMetaRow,
} from '@/lib/auditMeta'

export function useAuditMeta() {
  const { t, te } = useI18n()
  const { formatTHB } = useLocaleFormat()

  function labelFor(key: AuditMetaFieldKey): string {
    return t(`admin.logs.metaFields.${key}`)
  }

  function formatReason(reason: string): string {
    const key = `admin.logs.metaReasons.${reason}`
    return te(key) ? t(key) : reason
  }

  function auditMetaRows(meta?: Record<string, unknown>): AuditMetaRow[] {
    return buildAuditMetaRows(meta, {
      formatTotal: formatTHB,
      formatReason,
    })
  }

  function rowLabel(row: AuditMetaRow): string {
    return auditMetaRowLabel(row, labelFor)
  }

  function formatAuditSummary(meta?: Record<string, unknown>): string {
    const rows = auditMetaRows(meta)
    if (!rows.length) return t('common.dash')
    return rows.map((row) => `${rowLabel(row)}: ${row.value}`).join(' · ')
  }

  return {
    auditMetaRows,
    rowLabel,
    formatAuditSummary,
  }
}
