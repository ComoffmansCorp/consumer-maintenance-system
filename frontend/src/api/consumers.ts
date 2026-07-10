import { api } from './client'
import type { ConsumerDTO, ConsumerType, Page } from '@/types'

export interface ConsumerListParams {
  page?: number
  pageSize?: number
  search?: string
}

export interface ConsumerRequest {
  name: string
  type: ConsumerType
  description?: string
}

export const consumersApi = {
  list: (params: ConsumerListParams = {}) =>
    api.get<Page<ConsumerDTO>>('/consumers', { params }).then((r) => r.data),
  get: (id: number) => api.get<ConsumerDTO>(`/consumers/${id}`).then((r) => r.data),
  create: (req: ConsumerRequest) => api.post<ConsumerDTO>('/consumers', req).then((r) => r.data),
  update: (id: number, req: ConsumerRequest) =>
    api.patch<ConsumerDTO>(`/consumers/${id}`, req).then((r) => r.data),
  deactivate: (id: number) => api.delete<ConsumerDTO>(`/consumers/${id}`).then((r) => r.data),
}
