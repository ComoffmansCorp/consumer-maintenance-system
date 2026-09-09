package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	platformauth "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/auth"
	platformcache "github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/cache"
	"github.com/myurbondarchuk/consumer-maintenance-system/internal/platform/httpx"
)

var ErrTokenRevoked = errors.New("token has been revoked")

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

// AuthenticateToken validates a bearer token string and applies the same
// revocation check JWTAuth does, factored out so a non-HTTP-header caller
// (the chat WebSocket handshake, which authenticates via a query param
// instead -- browsers can't set a custom Authorization header on a
// WebSocket upgrade request) doesn't have to duplicate this logic.
func AuthenticateToken(ctx context.Context, tokens *platformauth.Service, cacheClient *platformcache.Client, token string) (AuthContext, error) {
	claims, err := tokens.ParseAccessToken(token)
	if err != nil {
		return AuthContext{}, err
	}
	if cacheClient != nil && claims.ID != "" {
		if revoked, err := cacheClient.Exists(ctx, "revoked:"+claims.ID); err == nil && revoked {
			return AuthContext{}, ErrTokenRevoked
		}
	}
	return AuthContext{UserID: claims.UserID, Role: claims.Role}, nil
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
			auth, err := AuthenticateToken(r.Context(), tokens, cacheClient, token)
			if err != nil {
				status, detail := http.StatusUnauthorized, "Invalid or expired token"
				if errors.Is(err, ErrTokenRevoked) {
					detail = "Token has been revoked"
				}
				httpx.WriteProblem(w, status, "Unauthorized", detail)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithAuth(r.Context(), auth)))
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
