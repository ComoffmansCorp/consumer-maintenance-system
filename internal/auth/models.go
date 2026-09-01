package auth

import "time"

type Role string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleClient     Role = "CLIENT"
	RoleMaster     Role = "MASTER"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	FullName     string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}
