package authHelper

import (
	"context"
	"net/http"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"
	"strconv"
	"strings"
)

// ContextKey es el tipo para las claves del contexto
type ContextKey string

const (
	TenantIDKey ContextKey = "tenantID"
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
	FirebaseUID ContextKey = "firebaseUID"
)

// ExtractTenantID extrae el tenantID del header (legacy)
func ExtractTenantID(r *http.Request) (uint, error) {
	tenantIDStr := r.Header.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return 0, http.ErrNoCookie
	}

	tenantID, err := strconv.ParseUint(tenantIDStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(tenantID), nil
}

// RequireAuthMiddleware is a chi-compatible middleware wrapper that verifies Firebase JWT
func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// Fallback legacy (si no hay Authorization, intenta usar los headers X-Tenant-ID para dev local si Firebase no está activo)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			if firebaseHelper.AuthClient == nil {
				// MODO MOCK / DEV FALLBACK
				tenantID, err := ExtractTenantID(r)
				if err != nil {
					responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
					return
				}
				ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify Token
		token, err := firebaseHelper.VerifyToken(idToken)
		if err != nil || token == nil {
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		// Extract custom claims (may be missing for new users)
		tenantIDFloat, ok := token.Claims["tenantId"].(float64)
		var tenantID uint = 0
		if ok {
			tenantID = uint(tenantIDFloat)
		}

		userIDFloat, ok := token.Claims["userId"].(float64)
		var userID uint = 0
		if ok {
			userID = uint(userIDFloat)
		}

		role, _ := token.Claims["role"].(string)
		if role == "" {
			role = "new_user" // default role for newly registered via Google without claims
		}

		// Set context
		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)
		ctx = context.WithValue(ctx, FirebaseUID, token.UID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTenantIDFromContext retrieves tenant ID from context
func GetTenantIDFromContext(ctx context.Context) uint {
	tenantID, _ := ctx.Value(TenantIDKey).(uint)
	return tenantID
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) uint {
	userID, _ := ctx.Value(UserIDKey).(uint)
	return userID
}
