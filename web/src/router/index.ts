import { createRouter, createWebHistory } from 'vue-router';
import DashboardLayout from '@/layouts/DashboardLayout.vue';
import DashboardIndex from '@/views/dashboard/Index.vue';

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
        {
          path: '/monitors',
          name: 'monitors',
          component: DashboardIndex, // temporary stub
        },
        {
          path: '/incidents',
          name: 'incidents',
          component: DashboardIndex, // temporary stub
        },
        {
          path: '/status-pages',
          name: 'status-pages',
          component: DashboardIndex, // temporary stub
        },
        {
          path: '/users',
          name: 'users',
          component: DashboardIndex, // temporary stub
        },
        {
          path: '/settings',
          name: 'settings',
          component: DashboardIndex, // temporary stub
        },
      ],
    },
    // {
    //   path: '/:pathMatch(.*)*',
    //   name: 'not-found',
    //   component: () => import('@/views/NotFound.vue'),
    // },
  ],
});

export default router;
