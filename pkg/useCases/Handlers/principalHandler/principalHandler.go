package principalHandler

import (
	"sales-manager-back/internal/data/infrastructure/principalRepository"
	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/response"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return principalRepository.GetAll(tenantID)
}

func GetByID(tenantID, principalID uint) (interface{}, response.Status) {
	return principalRepository.GetByID(tenantID, principalID)
}

func Create(newPrincipal *principal.Principal) (interface{}, response.Status) {
	return principalRepository.Create(newPrincipal)
}

func Update(tenantID, principalID uint, updates map[string]interface{}) response.Status {
	return principalRepository.Update(tenantID, principalID, updates)
}
