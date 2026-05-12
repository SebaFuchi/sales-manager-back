package saleHandler

import (
	"sales-manager-back/internal/data/infrastructure/saleRepository"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/sale"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return saleRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return saleRepository.GetByVendedor(tenantID, vendedorID)
}

func GetByID(tenantID, saleID uint) (interface{}, response.Status) {
	return saleRepository.GetByID(tenantID, saleID)
}

func Create(newSale *sale.Sale) (interface{}, response.Status) {
	return saleRepository.Create(newSale)
}

func Update(tenantID, saleID uint, updates map[string]interface{}) response.Status {
	return saleRepository.Update(tenantID, saleID, updates)
}

func Delete(tenantID, saleID uint) response.Status {
	return saleRepository.Delete(tenantID, saleID)
}
