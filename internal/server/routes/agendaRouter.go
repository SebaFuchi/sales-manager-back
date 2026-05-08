package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/agenda"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/agendaHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func AgendaRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /agenda - Get all events (with optional ?vendedorId filter)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		vendedorIDStr := r.URL.Query().Get("vendedorId")

		var events interface{}
		var status response.Status

		if vendedorIDStr != "" {
			vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 64)
			if err != nil {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			events, status = agendaHandler.GetByVendedor(tenantID, uint(vendedorID))
		} else {
			events, status = agendaHandler.GetAll(tenantID)
		}

		responseHelper.WriteResponse(w, status, events)
	})

	// GET /agenda/{eventId} - Get single event
	router.Get("/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		eventID, err := strconv.ParseUint(chi.URLParam(r, "eventId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		event, status := agendaHandler.GetByID(tenantID, uint(eventID))
		responseHelper.WriteResponse(w, status, event)
	})

	// POST /agenda - Create new event
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newEvent agenda.AgendaEvent
		if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newEvent.TenantID = tenantID
		event, status := agendaHandler.Create(&newEvent)
		responseHelper.WriteResponse(w, status, event)
	})

	// PATCH /agenda/{eventId} - Update event
	router.Patch("/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		eventID, err := strconv.ParseUint(chi.URLParam(r, "eventId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := agendaHandler.Update(tenantID, uint(eventID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	// POST /agenda/{eventId}/complete - Mark event as completed
	router.Post("/{eventId}/complete", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		eventID, err := strconv.ParseUint(chi.URLParam(r, "eventId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		updates := map[string]interface{}{
			"status": "completado",
		}
		status := agendaHandler.Update(tenantID, uint(eventID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	// DELETE /agenda/{eventId} - Delete event
	router.Delete("/{eventId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		eventID, err := strconv.ParseUint(chi.URLParam(r, "eventId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := agendaHandler.Delete(tenantID, uint(eventID))
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
