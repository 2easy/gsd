import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import InboxView from '../views/InboxView.vue'
import ProcessItemsView from '@/views/ProcessItemsView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/next-actions',
      name: 'Next Actions',
      component: () => import('../views/NextActionsView.vue'),
    },
    {
      path: '/projects',
      name: 'Projects',
      component: () => import('../views/ProjectsView.vue'),
    },
    {
      path: '/inbox',
      name: 'Inbox',
      component: InboxView
    },
    {
      path: '/process-items',
      name: 'ProcessItems',
      component: ProcessItemsView
    }
  ],
})

export default router
