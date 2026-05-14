package sale

import "gorm.io/gorm"

type EstadoVenta string

const (
	EstadoPedido     EstadoVenta = "Pedido"
	EstadoConfirmado EstadoVenta = "Confirmado"
	EstadoEntregado  EstadoVenta = "Entregado"
	EstadoFacturado  EstadoVenta = "Facturado"
	EstadoCancelado  EstadoVenta = "Cancelado"
)

type EstadoCobranza string

const (
	EstadoPendiente EstadoCobranza = "Pendiente"
	EstadoCobrado   EstadoCobranza = "Cobrado"
	EstadoAtrasado  EstadoCobranza = "Atrasado"
	EstadoEnGestion EstadoCobranza = "En Gestión"
)

type Sale struct {
	gorm.Model
	TenantID         uint           `json:"tenantId" gorm:"not null;index"`
	Date             string         `json:"date" gorm:"not null;index"`
	ClientID         uint           `json:"clientId" gorm:"not null;index"`
	ClientName       string         `json:"clientName"`
	PrincipalID      uint           `json:"principalId" gorm:"not null;index"`
	PrincipalName    string         `json:"principalName"`
	AgentID          uint           `json:"agentId" gorm:"not null;index"`
	AgentName        string         `json:"agentName"`
	NetAmount        float64        `json:"netAmount" gorm:"not null;default:0"`
	SaleStatus       EstadoVenta    `json:"saleStatus" gorm:"type:varchar(20);not null;default:'Pedido'"`
	CollectionStatus   EstadoCobranza `json:"collectionStatus" gorm:"type:varchar(20);not null;default:'Pendiente'"`
	CommissionOverride *float64       `json:"commissionOverride" gorm:"type:decimal(5,2)"`
	Notes              string         `json:"notes" gorm:"type:text"`

	// Relations
	Items []SaleItem `json:"items" gorm:"foreignKey:SaleID;constraint:OnDelete:CASCADE"`
}

type SaleItem struct {
	gorm.Model
	SaleID      uint    `json:"saleId" gorm:"not null;index"`
	ProductName string  `json:"productName" gorm:"not null"`
	Quantity    float64 `json:"quantity" gorm:"not null"`
	UnitPrice   float64 `json:"unitPrice" gorm:"not null"`
	Subtotal    float64 `json:"subtotal" gorm:"not null"`
}
