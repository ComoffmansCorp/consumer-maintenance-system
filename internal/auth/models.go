package auth

import "time"

type Role string

const (
	RoleSuperAdmin  Role = "SUPER_ADMIN"
	RoleTenantAdmin Role = "TENANT_ADMIN"
	RoleDispatcher  Role = "DISPATCHER"
	RoleElectrician Role = "ELECTRICIAN"

	// Marketplace roles: platform-level (TenantID == nil), same as
	// RoleSuperAdmin — not tied to a single tenant/company.
	RoleMaster Role = "MASTER"
	RoleClient Role = "CLIENT"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	FullName  string
	Role      Role
	TenantID  *int64
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}
