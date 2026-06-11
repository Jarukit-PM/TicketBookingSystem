<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchBookingTicket } from '@/api/tickets'
import TicketCard from '@/components/TicketCard.vue'
import { Button, Card, CardContent, CardHeader, CardTitle } from '@/components/ui'
import type { TicketDetail } from '@/types/ticket'

const route = useRoute()
const router = useRouter()
const ticket = ref<TicketDetail | null>(null)
const error = ref<string | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    ticket.value = await fetchBookingTicket(route.params.bookingId as string)
  } catch {
    error.value = 'Unable to load your ticket.'
  } finally {
    loading.value = false
  }
})

function goHome(): void {
  router.push({ name: 'home' })
}
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-8 md:px-6">
    <div class="mx-auto max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>Your ticket</CardTitle>
          <p class="text-sm text-copy-secondary">Present this QR code at the cinema entrance.</p>
        </CardHeader>
        <CardContent>
          <p v-if="loading" class="text-copy-secondary">Loading ticket…</p>
          <p v-else-if="error" class="text-state-error">{{ error }}</p>
          <TicketCard v-else-if="ticket" :ticket="ticket" />
          <div class="pt-6">
            <Button type="button" variant="ghost" @click="goHome">Back to home</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
