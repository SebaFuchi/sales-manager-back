package routes

import (
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/commissionHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func CommissionRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /commissions - Get all commissions (with optional ?vendedorId filter)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		vendedorIDStr := r.URL.Query().Get("vendedorId")

		var commissions interface{}
		var status response.Status

		if vendedorIDStr != "" {
			vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 64)
			if err != nil {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			commissions, status = commissionHandler.GetByVendedor(tenantID, uint(vendedorID))
		} else {
			commissions, status = commissionHandler.GetAll(tenantID)
		}

		responseHelper.WriteResponse(w, status, commissions)
	})

	return router
}
