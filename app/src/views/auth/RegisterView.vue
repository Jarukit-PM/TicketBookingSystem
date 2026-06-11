<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { translateApiError } from '@/api/errors'
import { ApiError } from '@/api/client'
import { Button, Card, CardContent, CardHeader, CardTitle, Input } from '@/components/ui'
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
          class="bg-gradient-brand bg-clip-text text-2xl font-semibold tracking-tight text-transparent"
        >
          {{ t('common.appName') }}
        </RouterLink>
        <p class="mt-2 text-sm text-copy-secondary">{{ t('auth.register.subtitle') }}</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{{ t('auth.register.title') }}</CardTitle>
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

            <p v-if="errorMessage" class="text-sm text-state-error" role="alert">
              {{ errorMessage }}
            </p>

            <Button type="submit" class="w-full" :disabled="submitting || auth.loading">
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
