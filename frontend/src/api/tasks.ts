import { api } from './client'
import type { Page, TaskDTO, TaskStatus, TaskType } from '@/types'

export interface TaskListParams {
  page?: number
  pageSize?: number
  status?: TaskStatus
  type?: TaskType
  assigneeId?: number
}

export interface CreateTaskRequest {
  type: TaskType
  addressId: number
  dueDate?: string | null
  assigneeId?: number | null
}

export const tasksApi = {
  list: (params: TaskListParams = {}) => api.get<Page<TaskDTO>>('/tasks', { params }).then((r) => r.data),
  get: (id: number) => api.get<TaskDTO>(`/tasks/${id}`).then((r) => r.data),
  create: (req: CreateTaskRequest) => api.post<TaskDTO>('/tasks', req).then((r) => r.data),
  assign: (id: number, assigneeId: number) =>
    api.post<TaskDTO>(`/tasks/${id}/assign`, { assigneeId }).then((r) => r.data),
  start: (id: number) => api.post<TaskDTO>(`/tasks/${id}/start`).then((r) => r.data),
  complete: (id: number) => api.post<TaskDTO>(`/tasks/${id}/complete`).then((r) => r.data),
  cancel: (id: number, reason: string) =>
    api.post<TaskDTO>(`/tasks/${id}/cancel`, { reason }).then((r) => r.data),
}
