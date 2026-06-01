import { createRouter, createWebHashHistory } from 'vue-router'
import Workspace from '../views/Workspace.vue'

const routes = [
  { path: '/', redirect: '/workspace' },
  { path: '/workspace', name: 'Workspace', component: Workspace },
  { path: '/history', name: 'History', component: () => import('../views/History.vue') },
  { path: '/docs', name: 'Docs', component: () => import('../views/Docs.vue') },
  { path: '/test-runner', name: 'TestRunner', component: () => import('../views/TestRunner.vue') },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
