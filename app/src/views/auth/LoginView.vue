<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { ApiError } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const submitting = ref(false)

const redirectTarget = computed(() => {
  const redirect = route.query.redirect
  return typeof redirect === 'string' && redirect.startsWith('/') ? redirect : '/'
})

async function onSubmit() {
  errorMessage.value = ''
  submitting.value = true
  try {
    await auth.login(email.value, password.value)
    await router.replace(redirectTarget.value)
  } catch (error) {
    if (error instanceof ApiError) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = 'Unable to sign in. Please try again.'
    }
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-base px-4 py-12">
    <div class="mx-auto w-full max-w-md">
      <div class="mb-8 text-center">
        <RouterLink
          to="/"
          class="bg-gradient-brand bg-clip-text text-2xl font-semibold tracking-tight text-transparent"
        >
          Cinema Tickets
        </RouterLink>
        <p class="mt-2 text-sm text-copy-secondary">Sign in to book seats and view your tickets.</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Sign in</CardTitle>
        </CardHeader>
        <CardContent>
          <form class="space-y-4" @submit.prevent="onSubmit">
            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="email">Email</label>
              <Input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                placeholder="you@example.com"
                required
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="password">Password</label>
              <Input
                id="password"
                v-model="password"
                type="password"
                autocomplete="current-password"
                placeholder="••••••••"
                required
              />
            </div>

            <p v-if="errorMessage" class="text-sm text-state-error" role="alert">
              {{ errorMessage }}
            </p>

            <Button type="submit" class="w-full" :disabled="submitting || auth.loading">
              {{ submitting ? 'Signing in…' : 'Sign in' }}
            </Button>
          </form>

          <div class="relative my-6">
            <div class="absolute inset-0 flex items-center">
              <span class="w-full border-t border-border" />
            </div>
            <div class="relative flex justify-center text-xs uppercase">
              <span class="bg-surface px-2 text-copy-secondary">or</span>
            </div>
          </div>

          <a
            href="/api/auth/google"
            class="flex w-full items-center justify-center gap-2 rounded-md border border-border bg-surface px-4 py-2 text-sm font-medium text-copy-primary transition hover:bg-surface-elevated"
          >
            Sign in with Google
          </a>

          <p class="mt-6 text-center text-sm text-copy-secondary">
            New here?
            <RouterLink
              class="text-brand hover:underline"
              :to="{ path: '/register', query: route.query }"
            >
              Create an account
            </RouterLink>
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
