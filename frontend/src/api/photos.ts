import { api } from './client'
import type { PhotoDTO } from '@/types'

export const photosApi = {
  listInspection: (actId: number) => api.get<PhotoDTO[]>(`/photos/inspection/${actId}`).then((r) => r.data),
  listReplacement: (actId: number) => api.get<PhotoDTO[]>(`/photos/replacement/${actId}`).then((r) => r.data),
  downloadPath: (id: number) => `/photos/${id}`,
}
