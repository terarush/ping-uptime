import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import DashboardIndex from '@/views/dashboard/Index.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      component: DashboardLayout,
      children: [
        {
          path: '',
          name: 'dashboard',
          component: DashboardIndex,
        },

      ],
    },
    // {
    //   path: '/:pathMatch(.*)*',
    //   name: 'not-found',
    //   component: () => import('@/views/NotFound.vue'),
    // },
  ],
})

export default router
