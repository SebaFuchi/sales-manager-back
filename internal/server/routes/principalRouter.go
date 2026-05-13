package routes

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Handlers/principalHandler"
	"sales-manager-back/pkg/useCases/Helpers/authHelper"
	"sales-manager-back/pkg/useCases/Helpers/responseHelper"

	"github.com/go-chi/chi/v5"
)

func PrincipalRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authHelper.RequireAuthMiddleware)

	// GET /principals - Get all principals
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principals, status := principalHandler.GetAll(tenantID)
		responseHelper.WriteResponse(w, status, principals)
	})

	// GET /principals/{principalId} - Get single principal
	router.Get("/{principalId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principalID, err := strconv.ParseUint(chi.URLParam(r, "principalId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		principalData, status := principalHandler.GetByID(tenantID, uint(principalID))
		responseHelper.WriteResponse(w, status, principalData)
	})

	// POST /principals - Create new principal
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

		var newPrincipal principal.Principal
		if err := json.NewDecoder(r.Body).Decode(&newPrincipal); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		newPrincipal.TenantID = tenantID
		principalData, status := principalHandler.Create(&newPrincipal)
		responseHelper.WriteResponse(w, status, principalData)
	})

	// PUT /principals/{principalId} - Update principal
	router.Put("/{principalId}", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())
		principalID, err := strconv.ParseUint(chi.URLParam(r, "principalId"), 10, 64)
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		status := principalHandler.Update(tenantID, uint(principalID), updates)
		responseHelper.WriteResponse(w, status, nil)
	})

	router.Post("/import", func(w http.ResponseWriter, r *http.Request) {
		tenantID := authHelper.GetTenantIDFromContext(r.Context())

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
		headers, err := reader.Read()
		if err != nil {
			responseHelper.WriteResponse(w, response.StatusBadRequest, nil)
			return
		}

		headerMap := make(map[string]int)
		for i, h := range headers {
			headerMap[strings.ToLower(strings.TrimSpace(h))] = i
		}

		var principals []principal.Principal
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

			getVal := func(keys ...string) string {
				for _, k := range keys {
					if idx, ok := headerMap[strings.ToLower(k)]; ok && idx < len(record) {
						return strings.TrimSpace(record[idx])
					}
				}
				return ""
			}

			name := getVal("nombre", "name", "razon social")
			if name == "" {
				errorsList = append(errorsList, "Fila "+strconv.Itoa(rowCount)+": Falta nombre de representada")
				continue
			}

			brandsStr := getVal("marcas", "brands")
			brands := 0
			if brandsStr != "" {
				b, err := strconv.Atoi(brandsStr)
				if err == nil {
					brands = b
				}
			}

			commissionStr := getVal("comision base", "comisión base", "commission")
			commission := 0.0
			if commissionStr != "" {
				c, err := strconv.ParseFloat(commissionStr, 64)
				if err == nil {
					commission = c
				}
			}

			initial := ""
			if len(name) > 0 {
				initial = string(name[0])
			}

			p := principal.Principal{
				TenantID:       tenantID,
				Name:           name,
				Brands:         brands,
				Category:       getVal("categoria", "categoría", "category"),
				BaseCommission: commission,
				Status:         principal.EstadoRepresentada(getVal("estado", "status")),
				Initial:        initial,
				ColorClass:     "bg-slate-100 text-slate-700", // Default color
				Contact:        getVal("contacto", "contact"),
				Email:          getVal("email", "correo"),
			}

			if p.Status == "" {
				p.Status = principal.EstadoActiva
			}
			if p.Category == "" {
				p.Category = "General"
			}

			principals = append(principals, p)
		}

		if len(principals) > 0 {
			status := principalHandler.BulkCreate(principals)
			if status != response.StatusCreated {
				errorsList = append(errorsList, "Error guardando en base de datos")
			}
		}

		responseHelper.WriteResponse(w, response.StatusCreated, map[string]interface{}{
			"imported": len(principals),
			"errors":   errorsList,
		})
	})

	return router
}
