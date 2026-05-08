package collection

import "gorm.io/gorm"

// Collection represents the record of payments collected
type Collection struct {
	gorm.Model
	TenantID          uint    `json:"tenantId" gorm:"not null;index"`
	AccountMovementID uint    `json:"accountMovementId" gorm:"not null;uniqueIndex"`
	VoucherID         string  `json:"voucherId" gorm:"not null;index"`
	ClientID          uint    `json:"clientId" gorm:"not null;index"`
	AgentID           uint    `json:"agentId" gorm:"not null;index"`
	Amount            float64 `json:"amount" gorm:"not null"`
	CollectionDate    string  `json:"collectionDate" gorm:"not null;index"`
	PaymentMethod     string  `json:"paymentMethod"`
	Notes             string  `json:"notes" gorm:"type:text"`
}
