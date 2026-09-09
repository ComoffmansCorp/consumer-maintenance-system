import { createRouter, createWebHistory } from 'vue-router'
import type { Role } from '@/types'
import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    marketplaceRoles?: Role[]
    title?: string
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
      meta: { public: true, title: 'Вход' },
    },
    {
      path: '/register',
      name: 'marketplace-register',
      component: () => import('@/views/marketplace/MarketplaceRegisterView.vue'),
      meta: { public: true, title: 'Регистрация' },
    },
    {
      path: '/catalog',
      name: 'marketplace-catalog',
      component: () => import('@/views/marketplace/MarketplaceCatalogView.vue'),
      meta: { public: true, title: 'Каталог услуг' },
    },
    {
      path: '/about',
      name: 'marketplace-about',
      component: () => import('@/views/marketplace/MarketplaceAboutView.vue'),
      meta: { public: true, title: 'О платформе' },
    },
    {
      path: '/for-masters',
      name: 'marketplace-for-masters',
      component: () => import('@/views/marketplace/MarketplaceForMastersView.vue'),
      meta: { public: true, title: 'Для мастеров' },
    },
    {
      path: '/masters',
      name: 'marketplace-masters',
      component: () => import('@/views/marketplace/MarketplaceMastersView.vue'),
      meta: { public: true, title: 'Мастера' },
    },
    {
      path: '/masters/:id',
      name: 'marketplace-master-profile',
      component: () => import('@/views/marketplace/MarketplaceMasterProfileView.vue'),
      meta: { public: true, title: 'Профиль мастера' },
      props: true,
    },
    {
      path: '/new-request',
      name: 'marketplace-new-request',
      component: () => import('@/views/marketplace/MarketplaceNewRequestView.vue'),
      meta: { marketplaceRoles: ['CLIENT'], title: 'Новая заявка' },
    },
    {
      path: '/my-requests',
      name: 'marketplace-my-requests',
      component: () => import('@/views/marketplace/MarketplaceMyRequestsView.vue'),
      meta: { marketplaceRoles: ['CLIENT'], title: 'Мои заявки' },
    },
    {
      path: '/requests/:id',
      name: 'marketplace-request-detail',
      component: () => import('@/views/marketplace/MarketplaceRequestDetailView.vue'),
      meta: { marketplaceRoles: ['CLIENT'], title: 'Заявка' },
      props: true,
    },
    {
      path: '/favorites',
      name: 'marketplace-favorites',
      component: () => import('@/views/marketplace/MarketplaceFavoritesView.vue'),
      meta: { marketplaceRoles: ['CLIENT'], title: 'Избранное' },
    },
    {
      path: '/profile',
      name: 'marketplace-profile',
      component: () => import('@/views/marketplace/MarketplaceProfileView.vue'),
      meta: { marketplaceRoles: ['CLIENT', 'MASTER', 'SUPER_ADMIN'], title: 'Личный кабинет' },
    },
    {
      path: '/admin/categories',
      name: 'admin-categories',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminCategoriesView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Категории' },
    },
    {
      path: '/admin/services',
      name: 'admin-services',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminServicesView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Услуги' },
    },
    {
      path: '/admin/requests',
      name: 'admin-requests',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminRequestsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Заявки' },
    },
    {
      path: '/admin/masters',
      name: 'admin-masters',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminMastersView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Мастера' },
    },
    {
      path: '/admin/reviews',
      name: 'admin-reviews',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminReviewsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Отзывы' },
    },
    {
      path: '/admin/payments',
      name: 'admin-payments',
      component: () => import('@/views/marketplace/admin/MarketplaceAdminPaymentsView.vue'),
      meta: { marketplaceRoles: ['SUPER_ADMIN'], title: 'Админ · Платежи' },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
      meta: { public: true, title: 'Страница не найдена' },
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
      // Same role-aware fallback as MarketplaceLoginView's post-submit
      // redirect: 'marketplace-my-requests' is CLIENT-only, so sending
      // every role there bounced MASTER/SUPER_ADMIN back to the landing
      // page via the marketplaceRoles check below.
      const roleHome: Record<string, string> = {
        CLIENT: 'marketplace-my-requests',
        MASTER: 'marketplace-profile',
        SUPER_ADMIN: 'admin-categories',
      }
      return { name: roleHome[auth.role ?? ''] ?? 'marketplace' }
    }
  }
  if (to.meta.marketplaceRoles && auth.role && !to.meta.marketplaceRoles.includes(auth.role)) {
    return { name: 'marketplace' }
  }
  return true
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} — Мастерская` : 'Мастерская'
})

export default router
