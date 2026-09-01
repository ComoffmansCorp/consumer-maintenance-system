import { api } from './client'
import type { CategoryDTO, ServiceDTO } from '@/types'

export interface CreateCategoryRequest {
  name: string
  parentCategoryId?: number
}

export interface UpdateCategoryRequest {
  name: string
  active: boolean
}

export interface CreateServiceRequest {
  categoryId: number
  name: string
  description?: string
  priceFrom?: number
  priceTo?: number
  unit?: string
}

export interface UpdateServiceRequest {
  name: string
  description?: string
  priceFrom?: number
  priceTo?: number
  unit?: string
  active: boolean
}

// Public browsing (GET /api/catalog/*) is unauthenticated -- the landing
// page needs it before login. Mutations live under /api/admin/* per the
// backend's router.go note: JWTAuth's publicPaths only matches exact paths,
// so admin writes can't share the same "/categories"/"/services" path.
export const catalogApi = {
  listCategories: () => api.get<CategoryDTO[]>('/catalog/categories').then((r) => r.data),
  listServices: (categoryId?: number) =>
    api.get<ServiceDTO[]>('/catalog/services', { params: categoryId ? { categoryId } : {} }).then((r) => r.data),

  createCategory: (req: CreateCategoryRequest) =>
    api.post<CategoryDTO>('/admin/categories', req).then((r) => r.data),
  updateCategory: (id: number, req: UpdateCategoryRequest) =>
    api.put<CategoryDTO>(`/admin/categories/${id}`, req).then((r) => r.data),
  createService: (req: CreateServiceRequest) =>
    api.post<ServiceDTO>('/admin/services', req).then((r) => r.data),
  updateService: (id: number, req: UpdateServiceRequest) =>
    api.put<ServiceDTO>(`/admin/services/${id}`, req).then((r) => r.data),
}

// Flattens a (possibly one-level-nested) category tree into a plain list --
// handy anywhere a flat picker is needed (service create/edit forms, etc).
export function flattenCategories(categories: CategoryDTO[]): CategoryDTO[] {
  const out: CategoryDTO[] = []
  for (const c of categories) {
    out.push(c)
    if (c.subcategories) out.push(...c.subcategories)
  }
  return out
}
