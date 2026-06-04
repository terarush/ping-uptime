import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'
import AppIndex from '@/views/App.vue'
import AppUsers from '@/views/App.Users.vue'
import AppLayout from '@/layouts/AppLayout.vue'
import { useAuth } from '@/composables/useAuth'
import { siteConfig } from '@/content/config'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Index,
      meta: { guestOnly: true },
    },
    {
      path: '/status/:slug',
      name: 'StatusPagePublic',
      component: () => import('@/views/StatusPagePublic.vue'),
    },
    {
      path: siteConfig.appPath,
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'app',
          component: AppIndex,
        },
        {
          path: 'monitors',
          name: 'Monitors',
          component: () => import('@/views/App.Monitors.vue'),
        },
        {
          path: 'incidents',
          name: 'Incidents',
          component: () => import('@/views/App.Incidents.vue'),
        },
        {
          path: 'status-pages',
          name: 'StatusPages',
          component: () => import('@/views/App.StatusPages.vue'),
        },
        {
          path: 'users',
          name: 'Users',
          component: AppUsers,
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
