package agenda

import "gorm.io/gorm"

type AgendaEventType string

const (
	TipoReunion      AgendaEventType = "reunion"
	TipoVisita       AgendaEventType = "visita"
	TipoLlamada      AgendaEventType = "llamada"
	TipoRecordatorio AgendaEventType = "recordatorio"
	TipoOportunidad  AgendaEventType = "oportunidad"
	TipoCompromiso   AgendaEventType = "compromiso"
)

type AgendaEventStatus string

const (
	StatusPendiente  AgendaEventStatus = "pendiente"
	StatusCompletado AgendaEventStatus = "completado"
	StatusCancelado  AgendaEventStatus = "cancelado"
)

type AgendaEvent struct {
	gorm.Model
	TenantID    uint              `json:"tenantId" gorm:"not null;index"`
	AgentID     uint              `json:"agentId" gorm:"not null;index"`
	Type        AgendaEventType   `json:"type" gorm:"type:varchar(20);not null"`
	Title       string            `json:"title" gorm:"not null"`
	Date        string            `json:"date" gorm:"not null;index"`
	Time        string            `json:"time"`
	ClientID    *uint             `json:"clientId" gorm:"index"`
	ClientName  string            `json:"clientName"`
	PrincipalID *uint             `json:"principalId" gorm:"index"`
	SaleID      *uint             `json:"saleId" gorm:"index"`
	Notes       string            `json:"notes" gorm:"type:text"`
	Status      AgendaEventStatus `json:"status" gorm:"type:varchar(20);not null;default:'pendiente'"`
}
