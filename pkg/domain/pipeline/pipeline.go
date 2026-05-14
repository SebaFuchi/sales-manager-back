package pipeline

import "gorm.io/gorm"

type PipelineStage string

const (
	StageProspecto   PipelineStage = "prospecto"
	StageContactado  PipelineStage = "contactado"
	StagePropuesta   PipelineStage = "propuesta"
	StageNegociacion PipelineStage = "negociacion"
	StageCerrado     PipelineStage = "cerrado"
	StagePerdido     PipelineStage = "perdido"
)

type Deal struct {
	gorm.Model
	TenantID        uint          `json:"tenantId" gorm:"not null;index"`
	ClientName      string        `json:"clientName" gorm:"not null"`
	PrincipalName   string        `json:"principalName" gorm:"not null"`
	EstimatedAmount float64       `json:"estimatedAmount" gorm:"not null;default:0"`
	Stage           PipelineStage `json:"stage" gorm:"type:varchar(20);not null;default:'prospecto'"`
	AgentName       string        `json:"agentName" gorm:"not null"`
	AgentID         uint          `json:"agentId" gorm:"not null;index"`
	Date               string        `json:"date" gorm:"not null;index"`
	CommissionOverride *float64      `json:"commissionOverride" gorm:"type:decimal(5,2)"`
	Notes              string        `json:"notes" gorm:"type:text"`
}
