package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/principalHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func PrincipalRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /principals - Get all principals
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principals, status := principalHandler.GetAll(tenantID)
		responseHelper.WriteResponse(w, status, principals)
	})

	// GET /principals/{principalId} - Get single principal
	router.Get("/{principalId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principalID, err := strconv.ParseUint(chi.URLParam(r, "principalId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		principalData, status := principalHandler.GetByID(tenantID, uint(principalID))
		responseHelper.WriteResponse(w, status, principalData)
	})

	// POST /principals - Create new principal
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newPrincipal principal.Principal
		if err := json.NewDecoder(r.Body).Decode(&newPrincipal); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newPrincipal.TenantID = tenantID
		principalData, status := principalHandler.Create(&newPrincipal)
		responseHelper.WriteResponse(w, status, principalData)
	})

	// PUT /principals/{principalId} - Update principal
	router.Put("/{principalId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principalID, err := strconv.ParseUint(chi.URLParam(r, "principalId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := principalHandler.Update(tenantID, uint(principalID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
