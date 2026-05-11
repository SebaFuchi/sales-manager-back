package routes

import (
	"encoding/json"
	"net/http"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/authHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

type AuthRouter struct{}

func (ar *AuthRouter) Register(w http.ResponseWriter, r *http.Request) {
	// RequireAuthMiddleware should have populated FirebaseUID in the context
	uid, ok := r.Context().Value(authHelper.FirebaseUID).(string)
	if !ok || uid == "" {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	role, _ := r.Context().Value(authHelper.UserRoleKey).(string)
	if role != "new_user" && role != "" {
		// User already has a role, they don't need to register
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	var req authHandler.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	tenantData, status := authHandler.Register(uid, req)
	responseHelper.WriteResponse(w, status, tenantData)
}

func (ar *AuthRouter) Routes() http.Handler {
	r := chi.NewRouter()

	// Apply auth middleware to ensure they have a valid Firebase Token
	r.Use(authHelper.RequireAuthMiddleware)

	r.Post("/register", ar.Register)

	return r
}
