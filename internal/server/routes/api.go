package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Initialize struct-based routers
	sr := &SaleRouter{}
	cr := &ClientRouter{}

	// Mount all routers under /sales-manager prefix (already under /api from server.go)
	r.Mount("/sales-manager/ventas", sr.Routes())
	r.Mount("/sales-manager/clientes", cr.Routes())
	r.Mount("/sales-manager/agenda", AgendaRouter())
	r.Mount("/sales-manager/representadas", PrincipalRouter())
	r.Mount("/sales-manager/comisiones", CommissionRouter())
	r.Mount("/sales-manager/equipo", TeamRouter())
	r.Mount("/sales-manager/pipeline", PipelineRouter())
	r.Mount("/sales-manager/dashboard", DashboardRouter())
	r.Mount("/sales-manager/cobranzas", CollectionRouter())

	return r
}
