package collectionRepository

import (
	"sales-manager-back/pkg/domain/collection"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetPaid(tenantID uint) ([]collection.Collection, response.Status) {
	var collections []collection.Collection
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Order("collection_date DESC").
		Find(&collections)

	if err := result.Error; err != nil {
		return collections, response.StatusInternalServerError
	}

	return collections, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]collection.Collection, response.Status) {
	var collections []collection.Collection
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Order("collection_date DESC").
		Find(&collections)

	if err := result.Error; err != nil {
		return collections, response.StatusInternalServerError
	}

	return collections, response.StatusOk
}

func MarkPaid(newCollection *collection.Collection) (*collection.Collection, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newCollection)
	if err := result.Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return newCollection, response.StatusConflict
		}
		return newCollection, response.StatusInternalServerError
	}

	return newCollection, response.StatusCreated
}

func UnmarkPaid(tenantID, accountMovementID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND account_movement_id = ?", tenantID, accountMovementID).
		Delete(&collection.Collection{})

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusNoContent
}
