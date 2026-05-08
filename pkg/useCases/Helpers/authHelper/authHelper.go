package authHelper

import (
	"context"
	"net/http"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"
	"strconv"
)

// ContextKey es el tipo para las claves del contexto
type ContextKey string

const (
	TenantIDKey ContextKey = "tenantID"
	UserIDKey   ContextKey = "userID"
	UserRoleKey ContextKey = "userRole"
)

// ExtractTenantID extrae el tenantID del header
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

// ExtractUserID extrae el userID del header
func ExtractUserID(r *http.Request) (uint, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, http.ErrNoCookie
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}

// RequireAuth middleware validates both tenant and user authentication headers
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate Tenant ID
		tenantID, err := ExtractTenantID(r)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		// Validate User ID
		userID, err := ExtractUserID(r)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		// Add both to context
		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		ctx = context.WithValue(ctx, UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireAuthMiddleware is a chi-compatible middleware wrapper
func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate Tenant ID
		tenantID, err := ExtractTenantID(r)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		// Validate User ID
		userID, err := ExtractUserID(r)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
			return
		}

		// Add both to context
		ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
		ctx = context.WithValue(ctx, UserIDKey, userID)

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
