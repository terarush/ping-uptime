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
          path: 'users',
          name: 'Users',
          component: AppUsers,
          meta: { requiresAdmin: true },
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
