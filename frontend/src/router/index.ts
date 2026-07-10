import { createRouter, createWebHistory } from 'vue-router'
import type { Role } from '@/types'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    roles?: Role[]
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/register-company',
      name: 'register-company',
      component: () => import('@/views/RegisterCompanyView.vue'),
      meta: { public: true },
    },
    {
      path: '/bootstrap',
      name: 'bootstrap',
      component: () => import('@/views/BootstrapView.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/tenants',
      name: 'tenants',
      component: () => import('@/views/TenantsView.vue'),
      meta: { roles: ['SUPER_ADMIN'] },
    },
    {
      path: '/consumers',
      name: 'consumers',
      component: () => import('@/views/ConsumersView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
    },
    {
      path: '/addresses',
      name: 'addresses',
      component: () => import('@/views/AddressesView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import('@/views/TasksView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
    },
    {
      path: '/tasks/:id',
      name: 'task-detail',
      component: () => import('@/views/TaskDetailView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
      props: true,
    },
    {
      path: '/acts/inspection/:id',
      name: 'inspection-act',
      component: () => import('@/views/InspectionActView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
      props: true,
    },
    {
      path: '/acts/replacement/:id',
      name: 'replacement-act',
      component: () => import('@/views/ReplacementActView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
      props: true,
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('@/views/UsersView.vue'),
      meta: { roles: ['TENANT_ADMIN', 'DISPATCHER'] },
    },
    {
      path: '/notifications',
      name: 'notifications',
      component: () => import('@/views/NotificationsView.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { public: true },
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (!to.meta.public && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.public && auth.isAuthenticated && (to.name === 'login' || to.name === 'register-company')) {
    return { name: 'dashboard' }
  }
  if (to.meta.roles && auth.role && !to.meta.roles.includes(auth.role)) {
    return { name: 'dashboard' }
  }
  return true
})

export default router
