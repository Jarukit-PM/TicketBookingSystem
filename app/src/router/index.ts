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
      path: '/movies/:id',
      name: 'movie-detail',
      component: () => import('../views/MovieDetailView.vue'),
    },
    {
      path: '/book/:showtimeId',
      name: 'book',
      component: () => import('../views/SeatMapView.vue'),
    },
    {
      path: '/book/:showtimeId/checkout',
      name: 'checkout',
      component: () => import('../views/CheckoutView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/book/:showtimeId/confirmation/:bookingId',
      name: 'booking-confirmation',
      component: () => import('../views/BookingConfirmationView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/bookings/:bookingId/ticket',
      name: 'booking-ticket',
      component: () => import('../views/TicketView.vue'),
      meta: { requiresAuth: true },
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
      component: () => import('../views/MyBookingsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/my-bookings/:id',
      name: 'booking-detail',
      component: () => import('../views/BookingDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/ticket/:bookingRef',
      name: 'ticket',
      component: () => import('../views/PublicTicketView.vue'),
    },
    {
      path: '/admin',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/AdminDashboardView.vue'),
        },
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
        {
          path: 'bookings',
          name: 'admin-bookings',
          component: () => import('../views/admin/AdminBookingsView.vue'),
        },
        {
          path: 'users/:userId/bookings',
          name: 'admin-user-bookings',
          component: () => import('../views/admin/AdminUserBookingsView.vue'),
        },
        {
          path: 'scan',
          name: 'admin-scan',
          component: () => import('../views/admin/AdminScanView.vue'),
        },
        {
          path: 'logs',
          name: 'admin-logs',
          component: () => import('../views/admin/AdminLogsView.vue'),
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
