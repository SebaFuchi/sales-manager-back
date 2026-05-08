package dashboard

import "gorm.io/gorm"

type TipoAlerta string

const (
	TipoVencimiento TipoAlerta = "vencimiento"
	TipoMeta        TipoAlerta = "meta"
	TipoStock       TipoAlerta = "stock"
	TipoSeguimiento TipoAlerta = "seguimiento"
	TipoComision    TipoAlerta = "comision"
)

type NivelAlerta string

const (
	NivelInfo    NivelAlerta = "info"
	NivelWarning NivelAlerta = "warning"
	NivelError   NivelAlerta = "error"
)

type Alert struct {
	gorm.Model
	TenantID    uint        `json:"tenantId" gorm:"not null;index"`
	AgentID     uint        `json:"agentId" gorm:"not null;index"`
	Type        TipoAlerta  `json:"type" gorm:"type:varchar(20);not null"`
	Level       NivelAlerta `json:"level" gorm:"type:varchar(10);not null;default:'info'"`
	Title       string      `json:"title" gorm:"not null"`
	Description string      `json:"description" gorm:"type:text"`
	Seen        bool        `json:"seen" gorm:"default:false"`
	EntityID    *uint       `json:"entityId" gorm:"index"`
	EntityType  string      `json:"entityType"`
}

type QuickNote struct {
	gorm.Model
	TenantID uint   `json:"tenantId" gorm:"not null;index"`
	AgentID  uint   `json:"agentId" gorm:"not null;index"`
	Content  string `json:"content" gorm:"type:text;not null"`
	Color    string `json:"color" gorm:"default:'#fef3c7'"`
	Order    int    `json:"order" gorm:"default:0"`
}

type Goal struct {
	gorm.Model
	TenantID    uint    `json:"tenantId" gorm:"not null;index"`
	AgentID     uint    `json:"agentId" gorm:"not null;index"`
	Title       string  `json:"title" gorm:"not null"`
	MonthlyGoal float64 `json:"monthlyGoal" gorm:"not null;default:0"`
	Achieved    float64 `json:"achieved" gorm:"not null;default:0"`
	Period      string  `json:"period" gorm:"not null"`
	GoalType    string  `json:"goalType" gorm:"type:varchar(20);not null;default:'ventas'"`
}
