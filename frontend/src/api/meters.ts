import { api } from './client'
import type { MeterDTO, MeterType, SealState } from '@/types'

export interface MeterRequest {
  type: MeterType
  serialNumber: string
  manufactureYear?: number | null
  verificationDate?: string | null
  sealState?: SealState
  transformationRatio?: number | null
}

export const metersApi = {
  list: (actId: number) => api.get<MeterDTO[]>(`/acts/inspection/${actId}/meters`).then((r) => r.data),
  create: (actId: number, req: MeterRequest) =>
    api.post<MeterDTO>(`/acts/inspection/${actId}/meters`, req).then((r) => r.data),
  update: (actId: number, meterId: number, req: MeterRequest) =>
    api.patch<MeterDTO>(`/acts/inspection/${actId}/meters/${meterId}`, req).then((r) => r.data),
  remove: (actId: number, meterId: number) =>
    api.delete(`/acts/inspection/${actId}/meters/${meterId}`),
}
