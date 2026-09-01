import { api } from './client'
import type { Page, PaymentDTO, PaymentStatus } from '@/types'

export const paymentsApi = {
  getForRequest: (requestId: number) => api.get<PaymentDTO>(`/requests/${requestId}/payment`).then((r) => r.data),
  listAdmin: (page = 1, pageSize = 20) =>
    api.get<Page<PaymentDTO>>('/admin/payments', { params: { page, pageSize } }).then((r) => r.data),
}

export const paymentStatusLabels: Record<PaymentStatus, string> = {
  HELD: 'Удержано',
  RELEASED: 'Переведено',
  REFUNDED: 'Возвращено',
}

export const paymentStatusStyles: Record<PaymentStatus, { bg: string; fg: string }> = {
  HELD: { bg: 'oklch(0.94 0.04 75)', fg: 'oklch(0.42 0.1 75)' },
  RELEASED: { bg: 'oklch(0.94 0.04 145)', fg: 'oklch(0.42 0.1 145)' },
  REFUNDED: { bg: 'oklch(0.95 0.03 25)', fg: '#B3261E' },
}
