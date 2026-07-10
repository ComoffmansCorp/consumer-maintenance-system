package middleware

import (
	"context"
	"net/http"
	"strings"

	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

type AuthContext struct {
	UserID     int64
	Role       string
	TenantID   *int64
	TenantCode *string
}

type authContextKey struct{}

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey{}, auth)
}

func AuthFromContext(ctx context.Context) (AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey{}).(AuthContext)
	return auth, ok
}

func JWTAuth(tokens *platformauth.Service, publicPaths map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := publicPaths[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Missing bearer token")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			claims, err := tokens.ParseAccessToken(token)
			if err != nil {
				httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
				return
			}

			ctx := WithAuth(r.Context(), AuthContext{
				UserID:     claims.UserID,
				Role:       claims.Role,
				TenantID:   claims.TenantID,
				TenantCode: claims.TenantCode,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TenantIDFromContext returns the tenant id carried by the authenticated
// request. Platform-level users (SUPER_ADMIN) have no tenant and ok is false.
func TenantIDFromContext(ctx context.Context) (int64, bool) {
	auth, ok := AuthFromContext(ctx)
	if !ok || auth.TenantID == nil {
		return 0, false
	}
	return *auth.TenantID, true
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, ok := AuthFromContext(r.Context())
			if !ok {
				httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
				return
			}
			if _, ok := allowed[auth.Role]; !ok {
				httpx.WriteProblem(w, http.StatusForbidden, "Forbidden", "Insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
