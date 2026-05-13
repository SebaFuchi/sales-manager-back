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
	r.Mount("/sales-manager/sales", sr.Routes())
	r.Mount("/sales-manager/clients", cr.Routes())
	r.Mount("/sales-manager/events", AgendaRouter())
	r.Mount("/sales-manager/principals", PrincipalRouter())
	r.Mount("/sales-manager/commissions", CommissionRouter())
	r.Mount("/sales-manager/team", TeamRouter())
	r.Mount("/sales-manager/pipeline", PipelineRouter())
	r.Mount("/sales-manager/dashboard", DashboardRouter())
	r.Mount("/sales-manager/collections", CollectionRouter())

	// My Agency: returns the tenant data for the currently authenticated user,
	// enriched with live counts from related tables
	r.Route("/sales-manager/my-agency", func(sub chi.Router) {
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
