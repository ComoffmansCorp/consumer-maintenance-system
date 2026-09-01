import { createRouter, createWebHistory } from 'vue-router'
import type { Role } from '@/types'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    marketplaceRoles?: Role[]
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'marketplace',
      component: () => import('@/views/marketplace/MarketplaceLandingView.vue'),
      meta: { public: true },
    },
    {
      path: '/login',
      name: 'marketplace-login',
      component: () => import('@/views/marketplace/MarketplaceLoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/register',
      name: 'marketplace-register',
      component: () => import('@/views/marketplace/MarketplaceRegisterView.vue'),
      meta: { public: true },
    },
    {
      path: '/new-request',
      name: 'marketplace-new-request',
      component: () => import('@/views/marketplace/MarketplaceNewRequestView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
    },
    {
      path: '/my-requests',
      name: 'marketplace-my-requests',
      component: () => import('@/views/marketplace/MarketplaceMyRequestsView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
    },
    {
      path: '/requests/:id',
      name: 'marketplace-request-detail',
      component: () => import('@/views/marketplace/MarketplaceRequestDetailView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
      props: true,
    },
    {
      path: '/favorites',
      name: 'marketplace-favorites',
      component: () => import('@/views/marketplace/MarketplaceFavoritesView.vue'),
      meta: { marketplaceRoles: ['CLIENT'] },
    },
    {
      path: '/admin/categories',
      name: 'admin-categories',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminCategoriesView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
    },
    {
      path: '/admin/services',
      name: 'admin-services',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminServicesView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
    },
    {
      path: '/admin/requests',
      name: 'admin-requests',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminRequestsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
    },
    {
      path: '/admin/masters',
      name: 'admin-masters',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminMastersView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
    },
    {
      path: '/admin/reviews',
      name: 'admin-reviews',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminReviewsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
    },
    {
      path: '/admin/payments',
      name: 'admin-payments',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminPaymentsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'] },
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
    return { name: 'marketplace-login', query: { redirect: to.fullPath } }
  }
  if (to.meta.public && auth.isAuthenticated) {
    if (to.name === 'marketplace-login' || to.name === 'marketplace-register') {
      return { name: 'marketplace-my-requests' }
    }
  }
  if (to.meta.marketplaceRoles && auth.role && !to.meta.marketplaceRoles.includes(auth.role)) {
    return { name: 'marketplace' }
  }
  return true
})

export default router
