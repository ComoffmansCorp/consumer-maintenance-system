package auth

import "time"

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	UserID       int64  `json:"userId"`
	FullName     string `json:"fullName"`
	Role         Role   `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type BootstrapSuperAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

// RegisterRequest is shared by client and master self-registration.
type RegisterRequest struct {
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

func ToUserDTO(u User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
