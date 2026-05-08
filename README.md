# Sales Manager Backend

Backend Go para la aplicación Sales Manager. Arquitectura multi-tenant con patrón Repository -> Handler -> Router.

## Stack Tecnológico

- **Language:** Go 1.25
- **HTTP Router:** chi/v5 v5.2.3
- **ORM:** GORM v1.25.2
- **Database:** MySQL
- **Auth:** Firebase UID + Custom Headers

## Arquitectura

### Multi-Tenancy
- **Tenant Isolation:** Todas las queries incluyen `tenant_id`
- **Headers requeridos:**
  - `X-Tenant-ID`: ID del tenant
  - `X-User-ID`: ID del usuario autenticado
  - `Authorization`: Token de Firebase (opcional, para validación)

### Estructura del Proyecto

```
sales-manager-back/
├── cmd/
│   └── main.go                 # Punto de entrada
├── internal/
│   ├── server/
│   │   ├── server.go          # Configuración del servidor
│   │   └── routes/
│   │       ├── api.go         # Montaje de routers
│   │       ├── saleRouter.go
│   │       ├── clientRouter.go
│   │       ├── agendaRouter.go
│   │       ├── principalRouter.go
│   │       ├── commissionRouter.go
│   │       ├── teamRouter.go
│   │       ├── pipelineRouter.go
│   │       ├── dashboardRouter.go
│   │       └── collectionRouter.go
│   └── data/infrastructure/
│       ├── saleRepository/
│       ├── clientRepository/
│       ├── agendaRepository/
│       ├── principalRepository/
│       ├── commissionRepository/
│       ├── teamRepository/
│       ├── pipelineRepository/
│       ├── dashboardRepository/
│       └── collectionRepository/
└── pkg/
    ├── domain/
    │   ├── tenant/
    │   ├── user/
    │   ├── sale/
    │   ├── client/
    │   ├── agenda/
    │   ├── principal/
    │   ├── commission/
    │   ├── pipeline/
    │   ├── dashboard/
    │   ├── collection/
    │   └── response/
    └── useCases/
        ├── Handlers/
        │   ├── saleHandler/
        │   ├── clientHandler/
        │   ├── agendaHandler/
        │   ├── principalHandler/
        │   ├── commissionHandler/
        │   ├── teamHandler/
        │   ├── pipelineHandler/
        │   ├── dashboardHandler/
        │   └── collectionHandler/
        └── Helpers/
            ├── databaseHelper/    # DB connection + migrations
            ├── responseHelper/    # HTTP response formatting
            └── authHelper/        # Auth middleware + context

```

## Modelos de Dominio

### Tenant
- Multi-tenancy root entity
- Plans: Starter, Pro, Enterprise
- Tracks: Users, Clients, Principals, Operations

### User (Team Member)
- TeamRole: Titular/Director, Sub-vendedor, Backoffice/Admin
- UserRole: agency, subvendedor, superadmin
- Commission split configuration
- Firebase UID linkage

### Client
- Legal & trade names, tax ID, delivery address
- Category: A, B, C, General
- Linked to sales agent (AgentID)
- **Relations:**
  - ClientPrincipalCondition (M:M with Principal for commercial terms)
  - AccountMovement (current account ledger)

### Principal
- Manufacturers/brands represented
- Base commission, categories, contact info
- **Sub-entities:**
  - PriceList (name, validity dates, URL)
  - Catalog (name, size, URL)
  - Promotion (name, expiry, active status)

### Sale
- Date, Client, Principal, Agent
- NetAmount, SaleStatus (Pedido/Confirmado/Entregado/Facturado/Cancelado)
- CollectionStatus (Pendiente/Cobrado/Atrasado/En Gestión)
- **Sub-entity:**
  - SaleItem (product name, quantity, unit price, subtotal)

### Commission
- Linked to Sale, Principal, Responsible (User)
- Calculation fields: CalculationBase, Percentage, Gross, Expenses, Net
- Distribution: OwnerDistribution, SubDistribution
- Status: FactoryStatus, SubStatus

### AgendaEvent
- Type: reunion, visita, llamada, recordatorio, oportunidad, compromiso
- Status: pendiente, completado, cancelado
- Links to Client, Principal, Sale (optional)
- Date, Time, Notes

### Deal (Pipeline)
- Stage: prospecto, contactado, propuesta, negociacion, cerrado, perdido
- EstimatedAmount, ClientName, PrincipalName, AgentID
- Date, Notes

### Dashboard Entities
- **Alert:** Type (vencimiento, meta, stock, seguimiento, comision), Level (info, warning, error), Seen status
- **QuickNote:** Content, Color, Order
- **Goal:** MonthlyGoal, Achieved, Period, GoalType (ventas, etc.)

### Collection
- Links AccountMovement to collection record
- VoucherID, ClientID, AgentID, Amount, CollectionDate, PaymentMethod
- M:1 with AccountMovement (uniqueIndex)

## API Endpoints

Base path: `/api/sales-manager`

### Ventas (Sales)
- `GET /ventas` - Get all sales (filter: `?vendedorId=X`)
- `GET /ventas/{saleId}` - Get single sale with items
- `POST /ventas` - Create sale
- `PUT /ventas/{saleId}` - Update sale

### Clientes (Clients)
- `GET /clientes` - Get all clients (filters: `?vendedorId=X`, `?q=search`)
- `GET /clientes/{clientId}` - Get client with conditions & account movements
- `POST /clientes` - Create client
- `PUT /clientes/{clientId}` - Update client

### Agenda
- `GET /agenda` - Get all events (filter: `?vendedorId=X`)
- `GET /agenda/{eventId}` - Get single event
- `POST /agenda` - Create event
- `PATCH /agenda/{eventId}` - Update event
- `POST /agenda/{eventId}/complete` - Mark as completed
- `DELETE /agenda/{eventId}` - Delete event

### Representadas (Principals)
- `GET /representadas` - Get all principals
- `GET /representadas/{principalId}` - Get principal with pricelists, catalogs, promotions
- `POST /representadas` - Create principal
- `PUT /representadas/{principalId}` - Update principal

### Comisiones (Commissions)
- `GET /comisiones` - Get all commissions (filter: `?vendedorId=X`)

### Equipo (Team)
- `GET /equipo` - Get all team members
- `GET /equipo/{userId}` - Get user
- `POST /equipo` - Create user
- `PUT /equipo/{userId}` - Update user

### Pipeline
- `GET /pipeline` - Get all deals (filter: `?vendedorId=X`)
- `POST /pipeline` - Create deal
- `PATCH /pipeline/{dealId}` - Update deal (stage, amount, etc.)
- `DELETE /pipeline/{dealId}` - Delete deal

### Dashboard
- `GET /dashboard/alerts` - Get alerts for current user
- `POST /dashboard/alerts` - Create alert
- `PATCH /dashboard/alerts/{alertId}/read` - Mark alert as read
- `GET /dashboard/notes` - Get quick notes for current user
- `POST /dashboard/notes` - Create note
- `PATCH /dashboard/notes/{noteId}` - Update note
- `DELETE /dashboard/notes/{noteId}` - Delete note
- `GET /dashboard/goals` - Get goals for current user
- `POST /dashboard/goals` - Create goal
- `PATCH /dashboard/goals/{goalId}` - Update goal

### Cobranzas (Collections)
- `GET /cobranzas/paid` - Get paid collections (filter: `?vendedorId=X`)
- `POST /cobranzas/{movementId}/pay` - Mark account movement as paid
- `POST /cobranzas/{movementId}/unpay` - Remove payment record

## Response Format

```json
{
  "message": "Success message or error",
  "data": { ... }
}
```

**Status Codes:**
- 200 OK
- 201 Created
- 204 No Content
- 400 Bad Request
- 401 Unauthorized
- 404 Not Found
- 409 Conflict
- 500 Internal Server Error

## Configuration

### Environment Variables

```bash
# Database
DB=user:password@tcp(host:3306)

# Server
PORT=8080
```

### Database
- **Name:** `sales_manager`
- **Timezone:** UTC-3
- **Naming:** SingularTable (GORM)
- **Migrations:** Auto-migrate on startup via `databaseHelper.InitDB()`

### Tables Created (19 total)
1. tenant
2. user
3. principal
4. price_list
5. catalog
6. promotion
7. client
8. client_principal_condition
9. account_movement
10. sale
11. sale_item
12. commission
13. agenda_event
14. deal
15. alert
16. quick_note
17. goal
18. collection

## Setup & Run

### Quick Start

1. **Run the setup script:**
```bash
chmod +x setup.sh
./setup.sh
```

### Manual Setup

1. **Navigate to backend directory:**
```bash
cd sales-manager-back
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Set environment variables:**
```bash
# Copy example file
cp .env.example .env

# Edit .env with your database credentials
export DB="user:password@tcp(localhost:3306)/sales_manager"
export PORT=8080
```

4. **Run the server:**
```bash
# Development mode
go run cmd/main.go

# Or build and run
go build -o bin/sales-manager-api cmd/main.go
./bin/sales-manager-api
```

5. **Verify it's running:**
```bash
curl http://localhost:8080/health
# Should return: OK
```

## Development Notes

- **Language Policy:** ALL backend code (field names, variables, comments) is in English. Frontend remains in Spanish for UX.
- **CORS:** Configured for all origins, methods: GET, POST, PUT, PATCH, DELETE
- **Auth Middleware:** `authHelper.RequireAuth` validates headers and populates context with TenantID & UserID
- **Tenant Isolation:** Every query filters by `tenant_id` from context
- **GORM Preload:** All repositories use Preload for related entities
- **Cascade Deletes:** Foreign keys have `constraint:OnDelete:CASCADE`

## Testing

```bash
# Example: Get sales for tenant
curl -H "X-Tenant-ID: 1" -H "X-User-ID: 1" http://localhost:8080/api/sales-manager/ventas

# Example: Create client
curl -X POST -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 1" -H "X-User-ID: 1" \
  -d '{"legalName":"ACME Corp","taxId":"30-12345678-9","city":"CABA","province":"Buenos Aires","agentId":1}' \
  http://localhost:8080/api/sales-manager/clientes
```

## Next Steps

1. Add integration tests
2. Implement Firebase token validation in authHelper
3. Add pagination support for large datasets
4. Implement WebSocket for real-time alerts
5. Add bulk import/export endpoints
6. Document Swagger/OpenAPI specs
