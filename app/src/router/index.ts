import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/auth/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/auth/RegisterView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/my-bookings',
      name: 'my-bookings',
      component: () => import('../views/MyBookingsPlaceholderView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { requiresAdmin: true },
      children: [
        { path: '', redirect: { name: 'admin-movies' } },
        {
          path: 'movies',
          name: 'admin-movies',
          component: () => import('../views/admin/AdminMoviesView.vue'),
        },
        {
          path: 'cinemas',
          name: 'admin-cinemas',
          component: () => import('../views/admin/AdminCinemasView.vue'),
        },
        {
          path: 'screens',
          name: 'admin-screens',
          component: () => import('../views/admin/AdminScreensView.vue'),
        },
        {
          path: 'showtimes',
          name: 'admin-showtimes',
          component: () => import('../views/admin/AdminShowtimesView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.ensureSession()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAdmin) {
    if (!auth.isAuthenticated) {
      return { name: 'login', query: { redirect: to.fullPath } }
    }
    if (!auth.isAdmin) {
      return { name: 'home' }
    }
  }

  if (to.meta.guestOnly && auth.isAuthenticated) {
    const redirect = to.query.redirect
    if (typeof redirect === 'string' && redirect.startsWith('/')) {
      return redirect
    }
    return { name: 'home' }
  }
})

export default router
