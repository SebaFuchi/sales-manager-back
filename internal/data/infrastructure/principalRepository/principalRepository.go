package principalRepository

import (
	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetAll(tenantID uint) ([]principal.Principal, response.Status) {
	var principals []principal.Principal
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Preload("PriceLists").
		Preload("Catalogs").
		Preload("Promotions").
		Find(&principals)

	if err := result.Error; err != nil {
		return principals, response.StatusInternalServerError
	}

	return principals, response.StatusOk
}

func GetByID(tenantID, principalID uint) (principal.Principal, response.Status) {
	var principalItem principal.Principal
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, principalID).
		Preload("PriceLists").
		Preload("Catalogs").
		Preload("Promotions").
		First(&principalItem)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return principalItem, response.StatusNotFound
		}
		return principalItem, response.StatusInternalServerError
	}

	return principalItem, response.StatusOk
}

func Create(newPrincipal *principal.Principal) (*principal.Principal, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newPrincipal)
	if err := result.Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return newPrincipal, response.StatusConflict
		}
		return newPrincipal, response.StatusInternalServerError
	}

	return newPrincipal, response.StatusCreated
}

func BulkCreate(principals []principal.Principal) response.Status {
	db := databaseHelper.Db

	result := db.CreateInBatches(principals, 100)
	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	return response.StatusCreated
}

func Update(tenantID, principalID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&principal.Principal{}).
		Where("tenant_id = ? AND id = ?", tenantID, principalID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
