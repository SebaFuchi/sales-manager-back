package client

import "gorm.io/gorm"

type CategoriaCliente string

const (
	CategoriaSupermercado CategoriaCliente = "Supermercado"
	CategoriaMayorista    CategoriaCliente = "Mayorista"
	CategoriaMinorista    CategoriaCliente = "Minorista"
	CategoriaDistribuidor CategoriaCliente = "Distribuidor"
	CategoriaGeneral      CategoriaCliente = "General"
)

type EstadoCliente string

const (
	EstadoActivo     EstadoCliente = "Activo"
	EstadoDeudor     EstadoCliente = "Deudor"
	EstadoInactivo   EstadoCliente = "Inactivo"
	EstadoSuspendido EstadoCliente = "Suspendido"
)

type EstadoCobranzaMovimiento string

const (
	EstadoAlDia          EstadoCobranzaMovimiento = "al_dia"
	EstadoVenceProxto    EstadoCobranzaMovimiento = "vence_pronto"
	EstadoAtrasado       EstadoCobranzaMovimiento = "atrasado"
	EstadoCobrado        EstadoCobranzaMovimiento = "cobrado"
	EstadoSinVencimiento EstadoCobranzaMovimiento = "sin_vencimiento"
)

type Client struct {
	gorm.Model
	TenantID        uint             `json:"tenantId" gorm:"not null;index"`
	LegalName       string           `json:"legalName" gorm:"not null"`
	TradeName       string           `json:"tradeName"`
	TaxID           string           `json:"taxId" gorm:"not null;index"`
	DeliveryAddress string           `json:"deliveryAddress"`
	City            string           `json:"city" gorm:"not null"`
	Province        string           `json:"province" gorm:"not null"`
	Zone            string           `json:"zone"`
	Category        CategoriaCliente `json:"category" gorm:"type:varchar(30);not null;default:'General'"`
	Status          EstadoCliente    `json:"status" gorm:"type:varchar(20);not null;default:'Activo'"`
	AgentID         uint             `json:"agentId" gorm:"not null;index"`
	AgentName       string           `json:"agentName"`
	Phone           string           `json:"phone"`
	Email           string           `json:"email"`
	Notes           string           `json:"notes" gorm:"type:text"`

	// Relations
	Conditions       []ClientPrincipalCondition `json:"conditions" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
	AccountMovements []AccountMovement          `json:"accountMovements" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE"`
}

type ClientPrincipalCondition struct {
	gorm.Model
	ClientID        uint     `json:"clientId" gorm:"not null;index"`
	PrincipalID     uint     `json:"principalId" gorm:"not null;index"`
	PrincipalName   string   `json:"principalName"`
	PriceList       string   `json:"priceList"`
	PaymentTerm     string   `json:"paymentTerm"`
	PaymentMethod   string   `json:"paymentMethod"`
	CreditLimit     *float64 `json:"creditLimit"`
	SpecialDiscount string   `json:"specialDiscount"`
	ServiceChannel  string   `json:"serviceChannel"`
	Notes           string   `json:"notes" gorm:"type:text"`
}

type TipoMovimiento string

const (
	TipoFactura     TipoMovimiento = "factura"
	TipoPago        TipoMovimiento = "pago"
	TipoNotaCredito TipoMovimiento = "nota_credito"
	TipoNotaDebito  TipoMovimiento = "nota_debito"
)

type AccountMovement struct {
	gorm.Model
	ClientID         uint                     `json:"clientId" gorm:"not null;index"`
	Date             string                   `json:"date" gorm:"not null"`
	VoucherNumber    string                   `json:"voucherNumber" gorm:"not null"`
	Type             TipoMovimiento           `json:"type" gorm:"type:varchar(20);not null"`
	Detail           string                   `json:"detail"`
	Debit            float64                  `json:"debit" gorm:"not null;default:0"`
	Credit           float64                  `json:"credit" gorm:"not null;default:0"`
	Balance          float64                  `json:"balance" gorm:"not null;default:0"`
	DueDate          *string                  `json:"dueDate"`
	CollectionStatus EstadoCobranzaMovimiento `json:"collectionStatus" gorm:"type:varchar(30);not null;default:'sin_vencimiento'"`
}
