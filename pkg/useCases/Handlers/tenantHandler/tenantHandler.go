package tenantHandler

import (
	"sales-manager-back/internal/data/infrastructure/tenantRepository"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
)

func GetAll() (interface{}, response.Status) {
	return tenantRepository.GetAll()
}

func GetByID(id uint) (interface{}, response.Status) {
	return tenantRepository.GetByID(id)
}

func Create(newTenant *tenant.Tenant) (interface{}, response.Status) {
	return tenantRepository.Create(newTenant)
}

func Update(tenantID uint, updates map[string]interface{}) response.Status {
	return tenantRepository.Update(tenantID, updates)
}

func Delete(tenantID uint) response.Status {
	return tenantRepository.Delete(tenantID)
}
