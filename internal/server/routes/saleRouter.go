package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/sale"
	"sales-manager-back/pkg/useCases/Handlers/saleHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

type SaleRouter struct{}

func (sr *SaleRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	sales, status := saleHandler.GetAll(tenantID)
	responseHelper.WriteResponse(w, status, sales)
}

func (sr *SaleRouter) GetByVendedor(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	vendedorIDStr := r.URL.Query().Get("vendedorId")
	vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 32)
	if err != nil || vendedorID == 0 {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	sales, status := saleHandler.GetByVendedor(tenantID, uint(vendedorID))
	responseHelper.WriteResponse(w, status, sales)
}

func (sr *SaleRouter) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	saleID, err := strconv.ParseUint(chi.URLParam(r, "saleId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	saleItem, status := saleHandler.GetByID(tenantID, uint(saleID))
	responseHelper.WriteResponse(w, status, saleItem)
}

func (sr *SaleRouter) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	var newSale sale.Sale
	if err := json.NewDecoder(r.Body).Decode(&newSale); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	newSale.TenantID = tenantID
	createdSale, status := saleHandler.Create(&newSale)
	responseHelper.WriteResponse(w, status, createdSale)
}

func (sr *SaleRouter) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	saleID, err := strconv.ParseUint(chi.URLParam(r, "saleId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	status := saleHandler.Update(tenantID, uint(saleID), updates)
	responseHelper.WriteResponse(w, status, nil)
}

func (sr *SaleRouter) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	saleID, err := strconv.ParseUint(chi.URLParam(r, "saleId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	status := saleHandler.Delete(tenantID, uint(saleID))
	responseHelper.WriteResponse(w, status, nil)
}

func (sr *SaleRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(authHelper.RequireAuthMiddleware)

	r.Get("/", sr.GetAll)
	r.Get("/{saleId}", sr.GetByID)
	r.Post("/", sr.Create)
	r.Put("/{saleId}", sr.Update)
	r.Delete("/{saleId}", sr.Delete)

	return r
}
