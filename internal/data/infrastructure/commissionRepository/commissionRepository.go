package commissionRepository

import (
	"sales-manager-back/pkg/domain/commission"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
)

func GetAll(tenantID uint) ([]commission.Commission, response.Status) {
	var commissions []commission.Commission
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).Find(&commissions)

	if err := result.Error; err != nil {
		return commissions, response.StatusInternalServerError
	}

	return commissions, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]commission.Commission, response.Status) {
	var commissions []commission.Commission
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND responsible_id = ?", tenantID, vendedorID).
		Find(&commissions)

	if err := result.Error; err != nil {
		return commissions, response.StatusInternalServerError
	}

	return commissions, response.StatusOk
}

func GetBySale(tenantID, saleID uint) ([]commission.Commission, response.Status) {
	var commissions []commission.Commission
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND sale_id = ?", tenantID, saleID).
		Find(&commissions)

	if err := result.Error; err != nil {
		return commissions, response.StatusInternalServerError
	}

	return commissions, response.StatusOk
}

func Create(newCommission *commission.Commission) (*commission.Commission, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newCommission)
	if err := result.Error; err != nil {
		return newCommission, response.StatusInternalServerError
	}

	return newCommission, response.StatusCreated
}

func Update(tenantID, commissionID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&commission.Commission{}).
		Where("tenant_id = ? AND id = ?", tenantID, commissionID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
