<script setup lang="ts">
import { Film, UserPlus } from 'lucide-vue-next'
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

const name = ref('')
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
    await auth.register(name.value, email.value, password.value)
    await router.replace(redirectTarget.value)
  } catch (error) {
    if (error instanceof ApiError) {
      errorMessage.value = translateApiError(error.code, error.message)
    } else {
      errorMessage.value = t('auth.register.error')
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
        <p class="mt-2 text-sm text-copy-secondary">{{ t('auth.register.subtitle') }}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle class="flex items-center gap-2">
            <UserPlus class="h-5 w-5 text-brand" aria-hidden="true" />
            {{ t('auth.register.title') }}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form class="space-y-4" @submit.prevent="onSubmit">
            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="name">{{ t('auth.register.name') }}</label>
              <Input
                id="name"
                v-model="name"
                autocomplete="name"
                :placeholder="t('auth.register.namePlaceholder')"
                required
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="email">{{ t('auth.register.email') }}</label>
              <Input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                :placeholder="t('auth.register.emailPlaceholder')"
                required
              />
            </div>

            <div class="space-y-2">
              <label class="text-sm text-copy-secondary" for="password">
                {{ t('auth.register.password') }}
              </label>
              <Input
                id="password"
                v-model="password"
                type="password"
                autocomplete="new-password"
                minlength="8"
                :placeholder="t('auth.register.passwordPlaceholder')"
                required
              />
            </div>

            <ErrorAlert v-if="errorMessage" :message="errorMessage" />

            <Button type="submit" class="w-full gap-1.5" :disabled="submitting || auth.loading">
              <UserPlus class="h-4 w-4" aria-hidden="true" />
              {{ submitting ? t('auth.register.submitting') : t('auth.register.submit') }}
            </Button>
          </form>

          <p class="mt-6 text-center text-sm text-copy-secondary">
            {{ t('auth.register.hasAccount') }}
            <RouterLink
              class="text-brand hover:underline"
              :to="{ path: '/login', query: route.query }"
            >
              {{ t('auth.register.signIn') }}
            </RouterLink>
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
