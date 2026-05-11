package tenantHandler

import (
	"sales-manager-back/internal/data/infrastructure/tenantRepository"
	"sales-manager-back/pkg/domain/response"
)

func GetAll() (interface{}, response.Status) {
	return tenantRepository.GetAll()
}

func GetByID(id uint) (interface{}, response.Status) {
	return tenantRepository.GetByID(id)
}
