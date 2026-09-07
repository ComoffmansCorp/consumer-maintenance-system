package middleware

import (
	"context"
	"net/http"
	"strings"

	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

type AuthContext struct {
	UserID int64
	Role   string
}

type authContextKey struct{}

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey{}, auth)
}

func AuthFromContext(ctx context.Context) (AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey{}).(AuthContext)
	return auth, ok
}

// JWTAuth validates the bearer access token and, if cacheClient is non-nil,
// rejects it with 401 when its jti has been revoked (see auth.Service.Logout
// in internal/auth/service.go, which is what populates "revoked:<jti>" in
// Redis). cacheClient may be nil -- e.g. Redis briefly unavailable at
// startup, or a caller that doesn't wire one up -- in which case the
// revocation check is skipped entirely: fail open on the extra safety net
// rather than locking every user out of an otherwise-healthy API over a
// cache outage, the access token's own short TTL is still the backstop.
func JWTAuth(tokens *platformauth.Service, cacheClient *platformcache.Client, publicPaths map[string]struct{}) func(http.Handler) http.Handler {
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

			if cacheClient != nil && claims.ID != "" {
				if revoked, err := cacheClient.Exists(r.Context(), "revoked:"+claims.ID); err == nil && revoked {
					httpx.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "Token has been revoked")
					return
				}
			}

			ctx := WithAuth(r.Context(), AuthContext{
				UserID: claims.UserID,
				Role:   claims.Role,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
