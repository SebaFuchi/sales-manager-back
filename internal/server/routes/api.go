package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/tenantHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Initialize struct-based routers
	sr := &SaleRouter{}
	cr := &ClientRouter{}
	tr := &TenantRouter{}
	ar := &AuthRouter{}

	// Mount all routers under /sales-manager prefix (already under /api from server.go)
	r.Mount("/sales-manager/auth", ar.Routes())
	r.Mount("/sales-manager/admin/tenants", tr.Routes())
	r.Mount("/sales-manager/ventas", sr.Routes())
	r.Mount("/sales-manager/clientes", cr.Routes())
	r.Mount("/sales-manager/agenda", AgendaRouter())
	r.Mount("/sales-manager/representadas", PrincipalRouter())
	r.Mount("/sales-manager/comisiones", CommissionRouter())
	r.Mount("/sales-manager/equipo", TeamRouter())
	r.Mount("/sales-manager/pipeline", PipelineRouter())
	r.Mount("/sales-manager/dashboard", DashboardRouter())
	r.Mount("/sales-manager/cobranzas", CollectionRouter())

	// My Agency: returns the tenant data for the currently authenticated user
	r.Route("/sales-manager/mi-agencia", func(sub chi.Router) {
		sub.Use(authHelper.RequireAuthMiddleware)
		sub.Get("/", func(w http.ResponseWriter, r *http.Request) {
			tenantID := authHelper.GetTenantIDFromContext(r.Context())
			if tenantID == 0 {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			tenantData, status := tenantHandler.GetByID(tenantID)
			responseHelper.WriteResponse(w, status, tenantData)
		})
	})

	return r
}
