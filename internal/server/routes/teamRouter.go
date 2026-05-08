package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Handlers/teamHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func TeamRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /team - Get all team members
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		users, status := teamHandler.GetAll(tenantID)
		responseHelper.WriteResponse(w, status, users)
	})

	// GET /team/{userId} - Get single user
	router.Get("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		userData, status := teamHandler.GetByID(tenantID, uint(userID))
		responseHelper.WriteResponse(w, status, userData)
	})

	// POST /team - Create new team member
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newUser user.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newUser.TenantID = tenantID
		userData, status := teamHandler.Create(&newUser)
		responseHelper.WriteResponse(w, status, userData)
	})

	// PUT /team/{userId} - Update user
	router.Put("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := teamHandler.Update(tenantID, uint(userID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
