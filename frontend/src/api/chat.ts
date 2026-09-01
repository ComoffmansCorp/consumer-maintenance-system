import { api } from './client'
import type { MessageDTO } from '@/types'

export const chatApi = {
  listMessages: (requestId: number, sinceId?: number) =>
    api
      .get<MessageDTO[]>(`/requests/${requestId}/messages`, { params: sinceId ? { sinceId } : {} })
      .then((r) => r.data),
  sendMessage: (requestId: number, text: string) =>
    api.post<MessageDTO>(`/requests/${requestId}/messages`, { text }).then((r) => r.data),
}
