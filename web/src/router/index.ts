import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'
import AppIndex from '@/views/App.vue'
import AppLayout from '@/layouts/AppLayout.vue'
import { useAuth } from '@/composables/useAuth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Index,
      meta: { guestOnly: true }
    },
    {
      path: '/app',
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'app',
          component: AppIndex,
        },
      ],
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const { isAuthenticated, verifyToken } = useAuth()

  if (!isAuthenticated.value) {
    await verifyToken()
  }

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/')
  } else if (to.meta.guestOnly && isAuthenticated.value) {
    next('/app')
  } else {
    next()
  }
})

export default router
