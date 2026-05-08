package collectionHandler

import (
	"sales-manager-back/internal/data/infrastructure/collectionRepository"
	"sales-manager-back/pkg/domain/collection"
	"sales-manager-back/pkg/domain/response"
)

func GetPaid(tenantID uint) (interface{}, response.Status) {
	return collectionRepository.GetPaid(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return collectionRepository.GetByVendedor(tenantID, vendedorID)
}

func MarkPaid(newCollection *collection.Collection) (interface{}, response.Status) {
	return collectionRepository.MarkPaid(newCollection)
}

func UnmarkPaid(tenantID, accountMovementID uint) response.Status {
	return collectionRepository.UnmarkPaid(tenantID, accountMovementID)
}
