package routes

import (
	"net/http"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/tenantHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TenantRouter struct{}

func (tr *TenantRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	// SuperAdmin role check: tenantID must be 0
	if tenantID != 0 {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	tenants, status := tenantHandler.GetAll()
	responseHelper.WriteResponse(w, status, tenants)
}

func (tr *TenantRouter) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	// SuperAdmin role check: tenantID must be 0
	if tenantID != 0 {
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

func (tr *TenantRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", tr.GetAll)
	r.Get("/{tenantId}", tr.GetByID)

	return r
}
