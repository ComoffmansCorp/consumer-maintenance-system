package auth

import "time"

type AuthResponse struct {
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"refreshToken"`
	ExpiresIn    int64   `json:"expiresIn"`
	FullName     string  `json:"fullName"`
	Role         Role    `json:"role"`
	TenantID     *int64  `json:"tenantId,omitempty"`
	TenantCode   *string `json:"tenantCode,omitempty"`
	TenantName   *string `json:"tenantName,omitempty"`
	TenantPlan   *string `json:"tenantPlan,omitempty"`
}

type LoginRequest struct {
	TenantCode string `json:"tenantCode"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type CompanyRegistrationRequest struct {
	TenantName string `json:"tenantName"`
	TenantCode string `json:"tenantCode"`
	Plan       string `json:"plan"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	FullName   string `json:"fullName"`
}

type BootstrapSuperAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UserDTO struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"fullName"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTenantUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Role     Role   `json:"role"`
}

func ToUserDTO(u User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
