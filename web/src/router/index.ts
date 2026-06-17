import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { siteConfig } from '@/content/config'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/Index.vue'),
      meta: { guestOnly: true },
    },
    {
      path: siteConfig.appPath,
      component: () => import('@/layouts/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'app',
          component: () => import('@/views/App.Index.vue'),
        },
        {
          path: 'monitors',
          name: 'Monitors',
          component: () => import('@/views/App.Monitors.vue'),
        },
        {
          path: 'status-pages',
          name: 'StatusPages',
          component: () => import('@/views/App.StatusPages.vue'),
        },
        {
          path: 'analytics',
          name: 'Analytics',
          component: () => import('@/views/App.Analytics.vue'),
        },
        {
          path: 'incidents',
          name: 'Incidents',
          component: () => import('@/views/App.Incidents.vue'),
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/App.Users.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'settings/notifications',
          name: 'Notifications',
          component: () => import('@/views/App.Notifications.vue'),
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('@/views/App.Settings.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue'),
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const { isAuthenticated, currentUser, verifyToken } = useAuth()

  if (!isAuthenticated.value) {
    await verifyToken()
  }

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/')
  } else if (to.meta.guestOnly && isAuthenticated.value) {
    next(siteConfig.appPath)
  } else if (to.meta.requiresAdmin && (!currentUser.value || currentUser.value.role !== 'admin')) {
    next(siteConfig.appPath)
  } else {
    next()
  }
})

export default router
