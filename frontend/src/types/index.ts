export type Role = 'SUPER_ADMIN' | 'TENANT_ADMIN' | 'DISPATCHER' | 'ELECTRICIAN' | 'MASTER' | 'CLIENT'

export interface AuthResponse {
  accessToken: string
  refreshToken: string
  expiresIn: number
  userId: number
  fullName: string
  role: Role
  tenantId?: number
  tenantCode?: string
  tenantName?: string
  tenantPlan?: string
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

// --- organization (platform tenants) ---

export type TenantPlan = 'FREE' | 'BUSINESS' | 'ENTERPRISE'

export interface TenantDTO {
  id: number
  name: string
  code: string
  plan: TenantPlan
  active: boolean
}

// --- consumer ---

export type ConsumerType = 'COMMERCIAL' | 'GOVERNMENT' | 'RESIDENTIAL'

export interface ConsumerDTO {
  id: number
  name: string
  type: ConsumerType
  description?: string
  active: boolean
  createdAt: string
  updatedAt: string
}

// --- address ---

export interface AddressDTO {
  id: number
  street: string
  house: string
  building?: string
  apartment?: string
  consumerId?: number
  consumerName?: string
  createdAt: string
  updatedAt: string
}

// --- meter ---

export type MeterType = 'SINGLE_PHASE' | 'THREE_PHASE_DIRECT' | 'THREE_PHASE_TRANSFORMER'
export type SealState = 'INTACT' | 'BROKEN' | 'MISSING' | ''

export interface MeterDTO {
  id: number
  type: MeterType
  serialNumber: string
  manufactureYear?: number
  verificationDate?: string
  sealState?: SealState
  transformationRatio?: number
  createdAt: string
}

// --- task ---

export type TaskType = 'INSPECTION' | 'REPLACEMENT'
export type TaskStatus = 'PENDING' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELED'

export interface TaskDTO {
  id: number
  type: TaskType
  addressId: number
  addressLabel?: string
  status: TaskStatus
  dueDate?: string
  assigneeId?: number
  assigneeName?: string
  createdAt: string
  updatedAt: string
  completedAt?: string
  canceledAt?: string
  cancelReason?: string
}

// --- act ---

export type InspectionType = 'SCHEDULED' | 'UNSCHEDULED'

export interface InspectionActDTO {
  id: number
  taskId: number
  addressId: number
  addressLabel?: string
  inspectionDate?: string
  consumerId?: number
  consumerName?: string
  inspectionType: InspectionType
  notes?: string
  meterCount: number
  photoCount: number
  createdAt: string
  updatedAt: string
}

export interface ReplacementActDTO {
  id: number
  taskId: number
  addressId: number
  addressLabel?: string
  accountNumber: string
  installationDate?: string
  oldBrand?: string
  oldSerialNumber?: string
  oldReadings?: number
  newBrand?: string
  newSerialNumber?: string
  newReadings?: number
  photoCount: number
  createdAt: string
  updatedAt: string
}

// --- photo ---

export interface PhotoDTO {
  id: number
  note?: string
  originalFilename: string
  contentType: string
  sizeBytes: number
  createdAt: string
}

// --- notification ---

export interface NotificationDTO {
  id: number
  type: string
  title: string
  message: string
  payload?: Record<string, unknown>
  read: boolean
  createdAt: string
}

// --- marketplace (заявки клиентов, независимо от tenant) ---

export interface CategoryDTO {
  id: number
  name: string
  active: boolean
}

export interface ServiceDTO {
  id: number
  categoryId: number
  name: string
  description?: string
  active: boolean
}

export type RequestStatus = 'OPEN' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELED'

export interface RequestDTO {
  id: number
  serviceId: number
  serviceName?: string
  categoryName?: string
  description: string
  addressText: string
  latitude?: number
  longitude?: number
  status: RequestStatus
  clientId: number
  clientName?: string
  masterId?: number
  masterName?: string
  createdAt: string
  updatedAt: string
  claimedAt?: string
  completedAt?: string
  canceledAt?: string
  cancelReason?: string
}

export interface MasterProfileDTO {
  city?: string
  bio?: string
  specializationIds: number[]
}
