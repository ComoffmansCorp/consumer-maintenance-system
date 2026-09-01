package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    int64  `json:"userId"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(secret string, accessTTL, refreshTTL time.Duration) (*Service, error) {
	if secret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}
	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (s *Service) AccessTTL() time.Duration {
	return s.accessTTL
}

func (s *Service) RefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *Service) IssueAccessToken(userID int64, role string) (string, int64, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTTL)
	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, fmt.Errorf("sign access token: %w", err)
	}
	return signed, int64(s.accessTTL.Seconds()), nil
}

func (s *Service) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	if claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}
