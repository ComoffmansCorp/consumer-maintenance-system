import { createRouter, createWebHistory } from 'vue-router'
import type { Role } from '@/types'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    roles?: Role[]
    // Separate from `roles`: the admin app (tenant-scoped) and the public
    // marketplace (platform-level, CLIENT/MASTER) are two different auth
    // realms that happen to share one Vue app -- an unauthenticated visit
    // must bounce to the matching login, not the other app's.
    marketplaceRoles?: Role[]
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
    // --- marketplace (public, platform-level: CLIENT posts requests, ---
    // --- MASTER claims them -- independent of the tenant admin app above) ---
    {
      path: '/marketplace',
      name: 'marketplace',
      component: () => import('@/views/marketplace/MarketplaceLandingView.vue'),
      meta: { public: true },
    },
    {
      path: '/marketplace/login',
      name: 'marketplace-login',
      component: () => import('@/views/marketplace/MarketplaceLoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/marketplace/register',
      name: 'marketplace-register',
      component: () => import('@/views/marketplace/MarketplaceRegisterView.vue'),
      meta: { public: true },
    },
    {
      path: '/marketplace/new-request',
      name: 'marketplace-new-request',
      component: () => import('@/views/marketplace/MarketplaceNewRequestView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
    },
    {
      path: '/marketplace/my-requests',
      name: 'marketplace-my-requests',
      component: () => import('@/views/marketplace/MarketplaceMyRequestsView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
    },
    {
      path: '/marketplace/requests/:id',
      name: 'marketplace-request-detail',
      component: () => import('@/views/marketplace/MarketplaceRequestDetailView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
      props: true,
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
  const isMarketplaceRoute = to.path.startsWith('/marketplace')

  if (!to.meta.public && !auth.isAuthenticated) {
    return isMarketplaceRoute
      ? { name: 'marketplace-login', query: { redirect: to.fullPath } }
      : { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.public && auth.isAuthenticated) {
    if (to.name === 'login' || to.name === 'register-company') return { name: 'dashboard' }
    if (to.name === 'marketplace-login' || to.name === 'marketplace-register') {
      return { name: 'marketplace-my-requests' }
    }
  }
  if (to.meta.roles && auth.role && !to.meta.roles.includes(auth.role)) {
    return { name: 'dashboard' }
  }
  if (to.meta.marketplaceRoles && auth.role && !to.meta.marketplaceRoles.includes(auth.role)) {
    return { name: 'marketplace' }
  }
  return true
})

export default router
