import { api } from './client'
import type { AddressDTO, Page } from '@/types'

export interface AddressListParams {
  page?: number
  pageSize?: number
  search?: string
  consumerId?: number
}

export interface AddressRequest {
  street: string
  house: string
  building?: string
  apartment?: string
  consumerId?: number | null
}

export const addressesApi = {
  list: (params: AddressListParams = {}) =>
    api.get<Page<AddressDTO>>('/addresses', { params }).then((r) => r.data),
  get: (id: number) => api.get<AddressDTO>(`/addresses/${id}`).then((r) => r.data),
  create: (req: AddressRequest) => api.post<AddressDTO>('/addresses', req).then((r) => r.data),
  update: (id: number, req: AddressRequest) =>
    api.patch<AddressDTO>(`/addresses/${id}`, req).then((r) => r.data),
}
