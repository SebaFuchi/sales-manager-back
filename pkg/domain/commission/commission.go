package commission

import "gorm.io/gorm"

type EstadoComisionFabrica string

const (
	EstadoPagada    EstadoComisionFabrica = "Pagada"
	EstadoPendCobro EstadoComisionFabrica = "Pend. Cobro"
	EstadoRechazada EstadoComisionFabrica = "Rechazada"
)

type EstadoComisionSub string

const (
	EstadoSubPagada   EstadoComisionSub = "Pagada"
	EstadoSubRetenida EstadoComisionSub = "Retenida (Falta Cobro)"
	EstadoSubNoAplica EstadoComisionSub = "No Aplica"
)

type Commission struct {
	gorm.Model
	TenantID          uint                  `json:"tenantId" gorm:"not null;index"`
	SaleID            uint                  `json:"saleId" gorm:"not null;index"`
	PrincipalID       uint                  `json:"principalId" gorm:"not null;index"`
	PrincipalName     string                `json:"principalName"`
	ResponsibleID     uint                  `json:"responsibleId" gorm:"not null;index"`
	ResponsibleName   string                `json:"responsibleName"`
	CalculationBase   float64               `json:"calculationBase" gorm:"not null;default:0"`
	Percentage        float64               `json:"percentage" gorm:"not null;default:0"`
	Gross             float64               `json:"gross" gorm:"not null;default:0"`
	Expenses          float64               `json:"expenses" gorm:"not null;default:0"`
	Net               float64               `json:"net" gorm:"not null;default:0"`
	OwnerDistribution float64               `json:"ownerDistribution" gorm:"not null;default:0"`
	SubDistribution   float64               `json:"subDistribution" gorm:"not null;default:0"`
	FactoryStatus     EstadoComisionFabrica `json:"factoryStatus" gorm:"type:varchar(20);not null;default:'Pend. Cobro'"`
	SubStatus         EstadoComisionSub     `json:"subStatus" gorm:"type:varchar(30);not null;default:'No Aplica'"`
}
