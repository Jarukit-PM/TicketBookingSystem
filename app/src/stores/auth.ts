import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { api } from '@/api/client'
import type { AuthResponse, User } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const initialized = ref(false)
  const loading = ref(false)

  const isAuthenticated = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function fetchMe() {
    loading.value = true
    try {
      const response = await api.get<{ user: User }>('/auth/me')
      user.value = response.user
    } catch {
      user.value = null
    } finally {
      initialized.value = true
      loading.value = false
    }
  }

  async function ensureSession() {
    if (!initialized.value) {
      await fetchMe()
    }
  }

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const response = await api.post<AuthResponse>('/auth/login', { email, password })
      user.value = response.user
      initialized.value = true
    } finally {
      loading.value = false
    }
  }

  async function register(name: string, email: string, password: string) {
    loading.value = true
    try {
      const response = await api.post<AuthResponse>('/auth/register', {
        name,
        email,
        password,
      })
      user.value = response.user
      initialized.value = true
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    try {
      await api.post('/auth/logout')
      user.value = null
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    initialized,
    loading,
    isAuthenticated,
    isAdmin,
    fetchMe,
    ensureSession,
    login,
    register,
    logout,
  }
})
