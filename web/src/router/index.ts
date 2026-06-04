import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'
import AppIndex from '@/views/app/Index.vue'
import AppLayout from '@/layouts/AppLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Index,
    },
    {
      path: '/app',
      component: AppLayout,
      children: [
        {
          path: '',
          name: 'app',
          component: AppIndex,
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
