import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import InboxView from '../views/InboxView.vue'
import ProcessInboxView from '@/views/ProcessInboxView.vue'

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
      path: '/process-inbox',
      name: 'ProcessIinbox',
      component: ProcessInboxView
    }
  ],
})

export default router
