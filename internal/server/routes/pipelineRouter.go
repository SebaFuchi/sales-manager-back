package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/pipeline"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/pipelineHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func PipelineRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /pipeline - Get all deals (with optional ?vendedorId filter)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		vendedorIDStr := r.URL.Query().Get("vendedorId")

		var deals interface{}
		var status response.Status

		if vendedorIDStr != "" {
			vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 64)
			if err != nil {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			deals, status = pipelineHandler.GetByVendedor(tenantID, uint(vendedorID))
		} else {
			deals, status = pipelineHandler.GetAll(tenantID)
		}

		responseHelper.WriteResponse(w, status, deals)
	})

	// POST /pipeline - Create new deal
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newDeal pipeline.Deal
		if err := json.NewDecoder(r.Body).Decode(&newDeal); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newDeal.TenantID = tenantID
		deal, status := pipelineHandler.Create(&newDeal)
		responseHelper.WriteResponse(w, status, deal)
	})

	// PATCH /pipeline/{dealId} - Update deal (stage, amount, etc.)
	router.Patch("/{dealId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		dealID, err := strconv.ParseUint(chi.URLParam(r, "dealId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := pipelineHandler.Update(tenantID, uint(dealID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	// DELETE /pipeline/{dealId} - Delete deal
	router.Delete("/{dealId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		dealID, err := strconv.ParseUint(chi.URLParam(r, "dealId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := pipelineHandler.Delete(tenantID, uint(dealID))
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
