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
      path: '/forgot-password',
      component: () => import('@/views/ForgotPassword.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/reset-password',
      component: () => import('@/views/ResetPassword.vue'),
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
          path: 'maintenances',
          name: 'Maintenances',
          component: () => import('@/views/App.Maintenances.vue'),
        },
        {
          path: 'audit-logs',
          name: 'AuditLogs',
          component: () => import('@/views/App.AuditLogs.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/App.Users.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'tags',
          name: 'Tags',
          component: () => import('@/views/App.Tags.vue'),
        },
        {
          path: 'api-tokens',
          name: 'ApiTokens',
          component: () => import('@/views/App.ApiTokens.vue'),
        },
        {
          path: 'ssl-monitors',
          name: 'SSLMonitors',
          component: () => import('@/views/App.SSLMonitors.vue'),
        },
        {
          path: 'subscribers',
          name: 'Subscribers',
          component: () => import('@/views/App.Subscribers.vue'),
        },
        {
          path: 'notification-logs',
          name: 'NotificationLogs',
          component: () => import('@/views/App.NotificationLogs.vue'),
        },
        {
          path: 'reports',
          name: 'Reports',
          component: () => import('@/views/App.Reports.vue'),
        },
        {
          path: 'integrations',
          name: 'Integrations',
          component: () => import('@/views/App.Integrations.vue'),
        },
        {
          path: 'teams',
          name: 'Teams',
          component: () => import('@/views/App.Teams.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'backup',
          name: 'Backup',
          component: () => import('@/views/App.Backup.vue'),
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
      path: '/status/:slug',
      name: 'StatusPagePublic',
      component: () => import('@/views/StatusPagePublic.vue'),
    },
    {
      path: '/status/:slug/monitor/:monitorId',
      name: 'StatusPageMonitorDetail',
      component: () => import('@/views/StatusPageMonitorDetail.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue'),
    },
  ],
})

router.beforeEach(async (to, from) => {
  const { isAuthenticated, currentUser, verifyToken } = useAuth()

  if (!isAuthenticated.value) {
    await verifyToken()
  }

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    return '/'
  }
  if (to.meta.guestOnly && isAuthenticated.value) {
    return siteConfig.appPath
  }
  if (to.meta.requiresAdmin && (!currentUser.value || currentUser.value.role !== 'admin')) {
    return siteConfig.appPath
  }
})

export default router
