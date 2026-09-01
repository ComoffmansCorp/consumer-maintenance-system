import { api } from './client'
import type { Page, ProfileDTO } from '@/types'

export interface UpdateProfileRequest {
  city: string
  bio: string
  specializationIds: number[]
}

export const masterApi = {
  getProfile: () => api.get<ProfileDTO>('/master/profile').then((r) => r.data),
  updateProfile: (req: UpdateProfileRequest) => api.put<ProfileDTO>('/master/profile', req).then((r) => r.data),

  // Admin: GET /api/admin/masters -- paginated (see internal/master/handler.go).
  listAdmin: (page = 1, pageSize = 20) =>
    api.get<Page<ProfileDTO>>('/admin/masters', { params: { page, pageSize } }).then((r) => r.data),
}
