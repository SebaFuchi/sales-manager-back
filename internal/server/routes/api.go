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

	// My Agency: returns the tenant data for the currently authenticated user,
	// enriched with live counts from related tables
	r.Route("/sales-manager/mi-agencia", func(sub chi.Router) {
		sub.Use(authHelper.RequireAuthMiddleware)
		sub.Get("/", func(w http.ResponseWriter, rr *http.Request) {
			tenantID := authHelper.GetTenantIDFromContext(rr.Context())
			if tenantID == 0 {
				responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
				return
			}
			enriched, status := tenantHandler.GetMyAgency(tenantID)
			responseHelper.WriteResponse(w, status, enriched)
		})
	})

	return r
}
