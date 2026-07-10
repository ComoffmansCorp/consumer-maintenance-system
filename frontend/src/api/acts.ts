import { api } from './client'
import type { InspectionActDTO, InspectionType, ReplacementActDTO } from '@/types'

export interface CreateInspectionRequest {
  taskId: number
  inspectionDate?: string | null
  consumerId?: number | null
  inspectionType: InspectionType
  notes?: string
}

export type UpdateInspectionRequest = Omit<CreateInspectionRequest, 'taskId'>

export interface CreateReplacementRequest {
  taskId: number
  accountNumber: string
  installationDate?: string | null
  oldBrand?: string
  oldSerialNumber?: string
  oldReadings?: number | null
  newBrand?: string
  newSerialNumber?: string
  newReadings?: number | null
}

export type UpdateReplacementRequest = Omit<CreateReplacementRequest, 'taskId'>

export const actsApi = {
  createInspection: (req: CreateInspectionRequest) =>
    api.post<InspectionActDTO>('/acts/inspection', req).then((r) => r.data),
  getInspection: (id: number) => api.get<InspectionActDTO>(`/acts/inspection/${id}`).then((r) => r.data),
  getInspectionByTask: (taskId: number) =>
    api.get<InspectionActDTO>(`/acts/inspection/by-task/${taskId}`).then((r) => r.data),
  updateInspection: (id: number, req: UpdateInspectionRequest) =>
    api.patch<InspectionActDTO>(`/acts/inspection/${id}`, req).then((r) => r.data),
  inspectionPdfPath: (id: number) => `/acts/inspection/${id}/pdf`,

  createReplacement: (req: CreateReplacementRequest) =>
    api.post<ReplacementActDTO>('/acts/replacement', req).then((r) => r.data),
  getReplacement: (id: number) => api.get<ReplacementActDTO>(`/acts/replacement/${id}`).then((r) => r.data),
  getReplacementByTask: (taskId: number) =>
    api.get<ReplacementActDTO>(`/acts/replacement/by-task/${taskId}`).then((r) => r.data),
  updateReplacement: (id: number, req: UpdateReplacementRequest) =>
    api.patch<ReplacementActDTO>(`/acts/replacement/${id}`, req).then((r) => r.data),
  replacementPdfPath: (id: number) => `/acts/replacement/${id}/pdf`,
}
