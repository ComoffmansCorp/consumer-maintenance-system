import { api } from './client'
import type { NotificationDTO, Page } from '@/types'

export const notificationsApi = {
  list: (params: { page?: number; pageSize?: number; unread?: boolean } = {}) =>
    api.get<Page<NotificationDTO>>('/notifications', { params }).then((r) => r.data),
  unreadCount: () => api.get<{ unreadCount: number }>('/notifications/unread-count').then((r) => r.data),
  markRead: (id: number) => api.post<NotificationDTO>(`/notifications/${id}/read`).then((r) => r.data),
  markAllRead: () => api.post('/notifications/read-all'),
}
