package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/dashboard"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/dashboardHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func DashboardRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// ============ ALERTS ============
	// GET /dashboard/alerts - Get alerts for a user
	router.Get("/alerts", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID := authHelper.GetUserIDFromContext(r.Context())

		alerts, status := dashboardHandler.GetAlerts(tenantID, userID)
		responseHelper.WriteResponse(w, status, alerts)
	})

	// POST /dashboard/alerts - Create new alert
	router.Post("/alerts", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newAlert dashboard.Alert
		if err := json.NewDecoder(r.Body).Decode(&newAlert); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newAlert.TenantID = tenantID
		alert, status := dashboardHandler.CreateAlert(&newAlert)
		responseHelper.WriteResponse(w, status, alert)
	})

	// PATCH /dashboard/alerts/{alertId}/read - Mark alert as read
	router.Patch("/alerts/{alertId}/read", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		alertID, err := strconv.ParseUint(chi.URLParam(r, "alertId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := dashboardHandler.MarkAlertAsRead(tenantID, uint(alertID))
		responseHelper.WriteResponse(w, status, nil)
	})

	// ============ QUICK NOTES ============
	// GET /dashboard/notes - Get quick notes for a user
	router.Get("/notes", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID := authHelper.GetUserIDFromContext(r.Context())

		notes, status := dashboardHandler.GetQuickNotes(tenantID, userID)
		responseHelper.WriteResponse(w, status, notes)
	})

	// POST /dashboard/notes - Create new note
	router.Post("/notes", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID := authHelper.GetUserIDFromContext(r.Context())

		var newNote dashboard.QuickNote
		if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newNote.TenantID = tenantID
		newNote.AgentID = userID
		note, status := dashboardHandler.CreateQuickNote(&newNote)
		responseHelper.WriteResponse(w, status, note)
	})

	// PATCH /dashboard/notes/{noteId} - Update note
	router.Patch("/notes/{noteId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		noteID, err := strconv.ParseUint(chi.URLParam(r, "noteId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := dashboardHandler.UpdateQuickNote(tenantID, uint(noteID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	// DELETE /dashboard/notes/{noteId} - Delete note
	router.Delete("/notes/{noteId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		noteID, err := strconv.ParseUint(chi.URLParam(r, "noteId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := dashboardHandler.DeleteQuickNote(tenantID, uint(noteID))
		responseHelper.WriteResponse(w, status, nil)
	})

	// ============ GOALS ============
	// GET /dashboard/goals - Get goals for a user
	router.Get("/goals", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID := authHelper.GetUserIDFromContext(r.Context())

		goals, status := dashboardHandler.GetGoals(tenantID, userID)
		responseHelper.WriteResponse(w, status, goals)
	})

	// POST /dashboard/goals - Create new goal
	router.Post("/goals", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		userID := authHelper.GetUserIDFromContext(r.Context())

		var newGoal dashboard.Goal
		if err := json.NewDecoder(r.Body).Decode(&newGoal); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newGoal.TenantID = tenantID
		newGoal.AgentID = userID
		goal, status := dashboardHandler.CreateGoal(&newGoal)
		responseHelper.WriteResponse(w, status, goal)
	})

	// PATCH /dashboard/goals/{goalId} - Update goal
	router.Patch("/goals/{goalId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		goalID, err := strconv.ParseUint(chi.URLParam(r, "goalId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := dashboardHandler.UpdateGoal(tenantID, uint(goalID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	return router
}
