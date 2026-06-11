<script setup lang="ts">
import { Film, LogIn, Mail } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, ErrorAlert, Input } from '@/components/ui'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
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
      errorMessage.value = translateApiError(error.code, error.message)
    } else {
      errorMessage.value = t('auth.login.error')
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
          class="inline-flex items-center justify-center gap-2 transition-opacity hover:opacity-90"
        >
          <Film class="h-6 w-6 text-brand" aria-hidden="true" />
          <span
            class="bg-gradient-brand bg-clip-text text-2xl font-semibold tracking-tight text-transparent"
          >
            {{ t('common.appName') }}
          </span>
        </RouterLink>
        <p class="mt-2 text-sm text-copy-secondary">{{ t('auth.login.subtitle') }}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <LogIn class="h-5 w-5 text-brand" aria-hidden="true" />
            {{ t('auth.login.title') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form class="space-y-4" @submit.prevent="onSubmit">
            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="email">{{ t('auth.login.email') }}</label>
              <Input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                :placeholder="t('auth.login.emailPlaceholder')"
                required
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="password">
                {{ t('auth.login.password') }}
              </label>
              <Input
                id="password"
                v-model="password"
                type="password"
                autocomplete="current-password"
                placeholder="••••••••"
                required
              />
            </div>

            <ErrorAlert v-if="errorMessage" :message="errorMessage" />

            <Button type="submit" class="w-full gap-1.5" :disabled="submitting || auth.loading">
              <LogIn class="h-4 w-4" aria-hidden="true" />
              {{ submitting ? t('auth.login.submitting') : t('auth.login.submit') }}
            </Button>
          </form>

          <div class="relative my-6">
            <div class="absolute inset-0 flex items-center">
              <span class="w-full border-t border-surface-border" />
            </div>
            <div class="relative flex justify-center text-xs uppercase">
              <span class="bg-surface px-2 text-copy-secondary">{{ t('common.or') }}</span>
            </div>
          </div>

          <a
            href="/api/auth/google"
            class="flex w-full items-center justify-center gap-2 rounded-full border border-surface-border bg-surface px-4 py-2.5 text-sm font-medium text-copy-primary transition-colors hover:bg-subtle"
          >
            <Mail class="h-4 w-4" aria-hidden="true" />
            {{ t('auth.login.google') }}
          </a>

          <p class="mt-6 text-center text-sm text-copy-secondary">
            {{ t('auth.login.newHere') }}
            <RouterLink
              class="text-brand hover:underline"
              :to="{ path: '/register', query: route.query }"
            >
              {{ t('auth.login.createAccount') }}
            </RouterLink>
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
