import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/views/DashboardLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/OverviewView.vue'),
        },
        {
          path: 'routes',
          name: 'routes',
          component: () => import('@/views/RoutesView.vue'),
        },
        {
          path: 'blocklist',
          name: 'blocklist',
          component: () => import('@/views/BlocklistView.vue'),
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('@/views/LogsView.vue'),
        },
      ],
    },
  ],
})

export default router
