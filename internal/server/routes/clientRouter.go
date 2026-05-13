package routes

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/clientHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

type ClientRouter struct{}

func (cr *ClientRouter) GetAll(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

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
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	clientID, err := strconv.ParseUint(chi.URLParam(r, "clientId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	clientItem, status := clientHandler.GetByID(tenantID, uint(clientID))
	responseHelper.WriteResponse(w, status, clientItem)
}

func (cr *ClientRouter) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

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
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

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

func (cr *ClientRouter) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())

	clientID, err := strconv.ParseUint(chi.URLParam(r, "clientId"), 10, 32)
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	status := clientHandler.Delete(tenantID, uint(clientID))
	responseHelper.WriteResponse(w, status, nil)
}

func (cr *ClientRouter) Import(w http.ResponseWriter, r *http.Request) {
	tenantID := authHelper.GetTenantIDFromContext(r.Context())
	agentID := authHelper.GetUserIDFromContext(r.Context())

	// Limit to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read header
	headers, err := reader.Read()
	if err != nil {
		responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
		return
	}

	// map header name to index
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	var clients []client.Client
	var errorsList []string
	rowCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorsList = append(errorsList, "Error leyendo fila")
			continue
		}
		rowCount++

		// Helper to safely get column by name
		getVal := func(keys ...string) string {
			for _, k := range keys {
				if idx, ok := headerMap[strings.ToLower(k)]; ok && idx < len(record) {
					return strings.TrimSpace(record[idx])
				}
			}
			return ""
		}

		cuit := getVal("cuit", "tax id", "rut")
		legalName := getVal("razon social", "razón social", "nombre", "legal name")

		if cuit == "" || legalName == "" {
			errorsList = append(errorsList, "Fila " + strconv.Itoa(rowCount) + ": Faltan campos obligatorios (CUIT o Razón Social)")
			continue
		}

		c := client.Client{
			TenantID:        tenantID,
			AgentID:         agentID,
			AgentName:       "Vendedor Importado", // En un sistema real se traería el nombre de auth
			TaxID:           cuit,
			LegalName:       legalName,
			TradeName:       getVal("nombre fantasia", "nombre de fantasía", "trade name"),
			DeliveryAddress: getVal("direccion", "dirección", "address"),
			City:            getVal("localidad", "ciudad", "city"),
			Province:        getVal("provincia", "province", "estado"),
			Zone:            getVal("zona", "zone"),
			Category:        client.CategoriaCliente(getVal("categoria", "categoría", "category")),
			Status:          client.EstadoCliente(getVal("estado", "status")),
			Phone:           getVal("telefono", "teléfono", "phone"),
			Email:           getVal("email", "correo"),
			Notes:           getVal("observaciones", "notas", "notes"),
		}

		// Defaults
		if c.Category == "" {
			c.Category = client.CategoriaGeneral
		}
		if c.Status == "" {
			c.Status = client.EstadoActivo
		}
		if c.City == "" {
			c.City = "Sin Definir"
		}
		if c.Province == "" {
			c.Province = "Sin Definir"
		}

		clients = append(clients, c)
	}

	if len(clients) > 0 {
		status := clientHandler.BulkCreate(clients)
		if status != response.StatusCreated {
			errorsList = append(errorsList, "Error guardando en base de datos")
		}
	}

	responseHelper.WriteResponse(w, response.StatusCreated, map[string]interface{}{
		"imported": len(clients),
		"errors":   errorsList,
	})
}

func (cr *ClientRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(authHelper.RequireAuthMiddleware)

	r.Get("/", cr.GetAll)
	r.Get("/{clientId}", cr.GetByID)
	r.Post("/", cr.Create)
	r.Post("/import", cr.Import)
	r.Put("/{clientId}", cr.Update)
	r.Delete("/{clientId}", cr.Delete)

	return r
}
