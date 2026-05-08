package routes

import (
	"net/http"
	templateHandler "template/pkg/useCases/Handlers/templateHandler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type TemplateRouter struct {
	Handler templateHandler.Handler
}

func (tr *TemplateRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowOriginFunc:    func(r *http.Request, origin string) bool { return true },
		AllowCredentials:   true,
		OptionsPassthrough: true,
		Debug:              true,
		MaxAge:             300,
	}))
	return r
}
