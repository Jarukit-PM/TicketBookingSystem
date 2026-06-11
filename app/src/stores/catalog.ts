import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import type { CatalogTab } from '@/types/catalog'

const STORAGE_KEY = 'tbs.selectedCinemaId'

function readStoredCinemaId(): string | null {
  try {
    return localStorage.getItem(STORAGE_KEY)
  } catch {
    return null
  }
}

function writeStoredCinemaId(id: string | null): void {
  try {
    if (id) {
      localStorage.setItem(STORAGE_KEY, id)
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  } catch {
    // ignore storage failures
  }
}

export const useCatalogStore = defineStore('catalog', () => {
  const selectedCinemaId = ref<string | null>(readStoredCinemaId())
  const activeTab = ref<CatalogTab>('now_showing')

  watch(selectedCinemaId, (id) => {
    writeStoredCinemaId(id)
  })

  function setCinema(id: string): void {
    selectedCinemaId.value = id
  }

  function setTab(tab: CatalogTab): void {
    activeTab.value = tab
  }

  return {
    selectedCinemaId,
    activeTab,
    setCinema,
    setTab,
  }
})
