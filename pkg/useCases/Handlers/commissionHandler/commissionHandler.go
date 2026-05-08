package commissionHandler

import (
	"sales-manager-back/internal/data/infrastructure/commissionRepository"
	"sales-manager-back/pkg/domain/commission"
	"sales-manager-back/pkg/domain/response"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return commissionRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return commissionRepository.GetByVendedor(tenantID, vendedorID)
}

func GetBySale(tenantID, saleID uint) (interface{}, response.Status) {
	return commissionRepository.GetBySale(tenantID, saleID)
}

func Create(newCommission *commission.Commission) (interface{}, response.Status) {
	return commissionRepository.Create(newCommission)
}

func Update(tenantID, commissionID uint, updates map[string]interface{}) response.Status {
	return commissionRepository.Update(tenantID, commissionID, updates)
}
