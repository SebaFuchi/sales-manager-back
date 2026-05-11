package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Handlers/tenantHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

type TenantRouter struct{}

func (tr *TenantRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(authHelper.UserRoleKey).(string)
	if role != string(user.RoleSuperAdmin) {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	tenants, status := tenantHandler.GetAll()
	responseHelper.WriteResponse(w, status, tenants)
}

func (tr *TenantRouter) GetByID(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(authHelper.UserRoleKey).(string)
	if role != string(user.RoleSuperAdmin) {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	idParam, err := strconv.ParseUint(chi.URLParam(r, "tenantId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	tenantItem, status := tenantHandler.GetByID(uint(idParam))
	responseHelper.WriteResponse(w, status, tenantItem)
}

func (tr *TenantRouter) Create(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(authHelper.UserRoleKey).(string)
	if role != string(user.RoleSuperAdmin) {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	var req tenant.Tenant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	// Set defaults
	if req.Plan == "" {
		req.Plan = tenant.PlanStarter
	}
	if req.Status == "" {
		req.Status = tenant.EstadoTrial
	}
	req.RegistrationDate = time.Now().Format("2006-01-02")
	req.LastActivity = time.Now().Format("2006-01-02")

	created, status := tenantHandler.Create(&req)
	responseHelper.WriteResponse(w, status, created)
}

func (tr *TenantRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(authHelper.RequireAuthMiddleware)

	r.Get("/", tr.GetAll)
	r.Post("/", tr.Create)
	r.Get("/{tenantId}", tr.GetByID)

	return r
}
