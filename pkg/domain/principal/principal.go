package principal

import "gorm.io/gorm"

type EstadoRepresentada string

const (
	EstadoActiva     EstadoRepresentada = "Activa"
	EstadoInactiva   EstadoRepresentada = "Inactiva"
	EstadoSuspendida EstadoRepresentada = "Suspendida"
)

type Principal struct {
	gorm.Model
	TenantID       uint               `json:"tenantId" gorm:"not null;index"`
	Name           string             `json:"name" gorm:"not null"`
	Brands         int                `json:"brands" gorm:"not null;default:0"`
	Category       string             `json:"category"`
	BaseCommission float64            `json:"baseCommission" gorm:"not null;default:0"`
	Status         EstadoRepresentada `json:"status" gorm:"type:varchar(20);not null;default:'Activa'"`
	Initial        string             `json:"initial" gorm:"size:1"`
	ColorClass     string             `json:"colorClass"`
	Contact        string             `json:"contact"`
	Email          string             `json:"email"`

	// Relations
	PriceLists []PriceList `json:"priceLists" gorm:"foreignKey:PrincipalID;constraint:OnDelete:CASCADE"`
	Catalogs   []Catalog   `json:"catalogs" gorm:"foreignKey:PrincipalID;constraint:OnDelete:CASCADE"`
	Promotions []Promotion `json:"promotions" gorm:"foreignKey:PrincipalID;constraint:OnDelete:CASCADE"`
}

type PriceList struct {
	gorm.Model
	PrincipalID uint   `json:"principalId" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null"`
	ValidFrom   string `json:"validFrom"`
	ValidUntil  string `json:"validUntil"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

type Catalog struct {
	gorm.Model
	PrincipalID uint    `json:"principalId" gorm:"not null;index"`
	Name        string  `json:"name" gorm:"not null"`
	SizeKb      float64 `json:"sizeKb"`
	URL         string  `json:"url"`
}

type Promotion struct {
	gorm.Model
	PrincipalID uint   `json:"principalId" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null"`
	ValidUntil  string `json:"validUntil"`
	Application string `json:"application"`
	Active      bool   `json:"active" gorm:"default:true"`
}
