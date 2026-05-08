package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/collection"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/collectionHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func CollectionRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /collections/paid - Get paid collections (with optional ?vendedorId filter)
	router.Get("/paid", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		vendedorIDStr := r.URL.Query().Get("vendedorId")

		var collections interface{}
		var status response.Status

		if vendedorIDStr != "" {
			vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 64)
			if err != nil {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			collections, status = collectionHandler.GetByVendedor(tenantID, uint(vendedorID))
		} else {
			collections, status = collectionHandler.GetPaid(tenantID)
		}

		responseHelper.WriteResponse(w, status, collections)
	})

	// POST /collections/{movementId}/pay - Mark account movement as paid
	router.Post("/{movementId}/pay", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		movementID, err := strconv.ParseUint(chi.URLParam(r, "movementId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var newCollection collection.Collection
		if err := json.NewDecoder(r.Body).Decode(&newCollection); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newCollection.TenantID = tenantID
		newCollection.AccountMovementID = uint(movementID)

		collectionData, status := collectionHandler.MarkPaid(&newCollection)
		responseHelper.WriteResponse(w, status, collectionData)
	})

	// POST /collections/{movementId}/unpay - Remove payment record
	router.Post("/{movementId}/unpay", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		movementID, err := strconv.ParseUint(chi.URLParam(r, "movementId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := collectionHandler.UnmarkPaid(tenantID, uint(movementID))
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
