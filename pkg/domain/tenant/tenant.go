package tenant

import "gorm.io/gorm"

type PlanSaaS string

const (
	PlanStarter    PlanSaaS = "Starter"
	PlanPro        PlanSaaS = "Pro"
	PlanEnterprise PlanSaaS = "Enterprise"
)

type EstadoTenant string

const (
	EstadoActivo     EstadoTenant = "Activo"
	EstadoTrial      EstadoTenant = "Trial"
	EstadoSuspendido EstadoTenant = "Suspendido"
)

type Tenant struct {
	gorm.Model
	Name               string       `json:"name" gorm:"not null"`
	Owner              string       `json:"owner" gorm:"not null"`
	Email              string       `json:"email" gorm:"not null;uniqueIndex"`
	Plan               PlanSaaS     `json:"plan" gorm:"type:varchar(20);not null;default:'Starter'"`
	MonthlyFee         float64      `json:"monthlyFee" gorm:"not null;default:0"`
	MonthVolume        float64      `json:"monthVolume" gorm:"not null;default:0"`
	GrowthMoM          string       `json:"growthMoM"`
	Users              int          `json:"users" gorm:"not null;default:1"`
	UserLimit          int          `json:"userLimit" gorm:"not null;default:5"`
	OwnersCount        int          `json:"ownersCount" gorm:"not null;default:0"`
	SubAgentsCount     int          `json:"subAgentsCount" gorm:"not null;default:0"`
	Status             EstadoTenant `json:"status" gorm:"type:varchar(20);not null;default:'Trial'"`
	RegistrationDate   string       `json:"registrationDate"`
	LastActivity       string       `json:"lastActivity"`
	ClientsLoaded      int          `json:"clientsLoaded" gorm:"not null;default:0"`
	PrincipalsLoaded   int          `json:"principalsLoaded" gorm:"not null;default:0"`
	OperationsRecorded int          `json:"operationsRecorded" gorm:"not null;default:0"`
}
