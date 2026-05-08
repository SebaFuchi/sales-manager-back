package routes

import (
	"encoding/json"
	"net/http"
	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/clientHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ClientRouter struct{}

func (cr *ClientRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	// Check if there's a vendedorId query param
	vendedorIDStr := r.URL.Query().Get("vendedorId")
	if vendedorIDStr != "" {
		vendedorID, err := strconv.ParseUint(vendedorIDStr, 10, 32)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}
		clients, status := clientHandler.GetByVendedor(tenantID, uint(vendedorID))
		responseHelper.WriteResponse(w, status, clients)
		return
	}

	// Check if there's a search query
	q := r.URL.Query().Get("q")
	if q != "" {
		clients, status := clientHandler.Search(tenantID, q)
		responseHelper.WriteResponse(w, status, clients)
		return
	}

	clients, status := clientHandler.GetAll(tenantID)
	responseHelper.WriteResponse(w, status, clients)
}

func (cr *ClientRouter) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	clientID, err := strconv.ParseUint(chi.URLParam(r, "clientId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	clientItem, status := clientHandler.GetByID(tenantID, uint(clientID))
	responseHelper.WriteResponse(w, status, clientItem)
}

func (cr *ClientRouter) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	var newClient client.Client
	if err := json.NewDecoder(r.Body).Decode(&newClient); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	newClient.TenantID = tenantID
	createdClient, status := clientHandler.Create(&newClient)
	responseHelper.WriteResponse(w, status, createdClient)
}

func (cr *ClientRouter) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, err := authHelper.ExtractTenantID(r)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusUnauthorized, nil)
		return
	}

	clientID, err := strconv.ParseUint(chi.URLParam(r, "clientId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	status := clientHandler.Update(tenantID, uint(clientID), updates)
	responseHelper.WriteResponse(w, status, nil)
}

func (cr *ClientRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", cr.GetAll)
	r.Get("/{clientId}", cr.GetByID)
	r.Post("/", cr.Create)
	r.Put("/{clientId}", cr.Update)

	return r
}
