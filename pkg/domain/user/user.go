package user

import "gorm.io/gorm"

type UserRole string

const (
	RoleAgency      UserRole = "agency"
	RoleSubvendedor UserRole = "subvendedor"
	RoleSuperAdmin  UserRole = "superadmin"
)

type EquipoRol string

const (
	EquipoRolTitular     EquipoRol = "Titular / Director"
	EquipoRolSubvendedor EquipoRol = "Sub-vendedor"
	EquipoRolBackoffice  EquipoRol = "Backoffice / Admin"
)

type EstadoGeneral string

const (
	EstadoActivo     EstadoGeneral = "Activo"
	EstadoInactivo   EstadoGeneral = "Inactivo"
	EstadoSuspendido EstadoGeneral = "Suspendido"
	EstadoTrial      EstadoGeneral = "Trial"
)

type User struct {
	gorm.Model
	TenantID           uint          `json:"tenantId" gorm:"not null;index"`
	Name               string        `json:"name" gorm:"not null"`
	Email              string        `json:"email" gorm:"uniqueIndex"`
	Phone              string        `json:"phone"`
	TeamRole           EquipoRol     `json:"teamRole" gorm:"type:varchar(50);not null"`
	Role               UserRole      `json:"role" gorm:"type:varchar(20);not null;default:'subvendedor'"`
	Clients            string        `json:"clients"`
	Split              string        `json:"split"`
	SplitPercentageSub float64       `json:"splitPercentageSub" gorm:"not null;default:0"`
	Base               string        `json:"base"`
	Status             EstadoGeneral `json:"status" gorm:"type:varchar(20);not null;default:'Activo'"`
	Initials           string        `json:"initials" gorm:"size:3"`
	CalendarLinked     bool          `json:"calendarLinked" gorm:"default:false"`
	FirebaseUID        string        `json:"firebaseUid" gorm:"uniqueIndex"`
}
