export type Role = 'SUPER_ADMIN' | 'CLIENT' | 'MASTER'

export interface AuthResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  userId: number
  fullName: string
  role: Role
}

export interface UserDTO {
  id: number
  username: string
  fullName: string
  role: Role
  createdAt: string
}

export interface Page<T> {
  items: T[]
  page: number
  pageSize: number
  totalItems: number
  totalPages: number
}

export interface Problem {
  type?: string
  title: string
  status: number
  detail?: string
  errors?: Record<string, unknown>
}

// --- catalog ---

export interface CategoryDTO {
  id: number
  name: string
  active: boolean
  subcategories?: CategoryDTO[]
}

export interface ServiceDTO {
  id: number
  categoryId: number
  name: string
  description?: string
  priceFrom?: number
  priceTo?: number
  unit?: string
  imageUrl?: string
  active: boolean
}

// --- master ---

export interface ProfileDTO {
  userId: number
  city?: string
  bio?: string
  avatarUrl?: string
  ratingAvg: number
  ratingCount: number
  specializationIds: number[]
}

// --- request / offers / favorites ---

export type RequestStatus = 'OPEN' | 'ASSIGNED' | 'COMPLETED' | 'CANCELED'
export type OfferStatus = 'PENDING' | 'ACCEPTED' | 'REJECTED' | 'WITHDRAWN'

export interface StatusHistoryEntryDTO {
  fromStatus?: string
  toStatus: string
  changedBy: number
  comment?: string
  createdAt: string
}

export interface RequestDTO {
  id: number
  serviceId: number
  serviceName?: string
  description: string
  addressText: string
  latitude?: number
  longitude?: number
  status: RequestStatus
  clientId: number
  masterId?: number
  agreedPrice?: number
  cancelReason?: string
  createdAt: string
  updatedAt: string
  history?: StatusHistoryEntryDTO[]
}

export interface OfferDTO {
  id: number
  requestId: number
  masterId: number
  price: number
  comment?: string
  status: OfferStatus
  createdAt: string
  updatedAt: string
  masterAvatarUrl?: string
}

export interface FavoriteDTO {
  masterId: number
  createdAt: string
}

// --- reviews ---

export interface ReviewDTO {
  id: number
  requestId: number
  clientId: number
  masterId: number
  rating: number
  comment?: string
  createdAt: string
}

// --- payments ---

export type PaymentStatus = 'HELD' | 'RELEASED' | 'REFUNDED'

export interface PaymentDTO {
  id: number
  requestId: number
  amount: number
  platformFee: number
  status: PaymentStatus
  createdAt: string
  updatedAt: string
}

// --- chat ---

export interface MessageDTO {
  id: number
  requestId: number
  senderId: number
  text: string
  createdAt: string
  readAt?: string
}
