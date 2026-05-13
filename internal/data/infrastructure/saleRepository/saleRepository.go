package saleRepository

import (
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/sale"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetAll(tenantID uint) ([]sale.Sale, response.Status) {
	var sales []sale.Sale
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Preload("Items").
		Order("date DESC").
		Find(&sales)

	if err := result.Error; err != nil {
		return sales, response.StatusInternalServerError
	}

	return sales, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]sale.Sale, response.Status) {
	var sales []sale.Sale
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Preload("Items").
		Order("date DESC").
		Find(&sales)

	if err := result.Error; err != nil {
		return sales, response.StatusInternalServerError
	}

	return sales, response.StatusOk
}

func GetByID(tenantID, saleID uint) (sale.Sale, response.Status) {
	var saleItem sale.Sale
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, saleID).
		Preload("Items").
		First(&saleItem)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return saleItem, response.StatusNotFound
		}
		return saleItem, response.StatusInternalServerError
	}

	return saleItem, response.StatusOk
}

func Create(newSale *sale.Sale) (*sale.Sale, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newSale)
	if err := result.Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return newSale, response.StatusConflict
		}
		return newSale, response.StatusInternalServerError
	}

	return newSale, response.StatusCreated
}

func Update(tenantID, saleID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&sale.Sale{}).
		Where("tenant_id = ? AND id = ?", tenantID, saleID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.StatusNotFound
		}
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

func Delete(tenantID, saleID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, saleID).Delete(&sale.Sale{})
	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}
	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
