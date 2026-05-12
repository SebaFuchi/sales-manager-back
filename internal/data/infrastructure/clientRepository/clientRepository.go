package clientRepository

import (
	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetAll(tenantID uint) ([]client.Client, response.Status) {
	var clients []client.Client
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Preload("Conditions").
		Preload("AccountMovements").
		Find(&clients)

	if err := result.Error; err != nil {
		return clients, response.StatusInternalServerError
	}

	return clients, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]client.Client, response.Status) {
	var clients []client.Client
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Preload("Conditions").
		Preload("AccountMovements").
		Find(&clients)

	if err := result.Error; err != nil {
		return clients, response.StatusInternalServerError
	}

	return clients, response.StatusOk
}

func GetByID(tenantID, clientID uint) (client.Client, response.Status) {
	var clientItem client.Client
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, clientID).
		Preload("Conditions").
		Preload("AccountMovements", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).
		First(&clientItem)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return clientItem, response.StatusNotFound
		}
		return clientItem, response.StatusInternalServerError
	}

	return clientItem, response.StatusOk
}

func Create(newClient *client.Client) (*client.Client, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newClient)
	if err := result.Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return newClient, response.StatusConflict
		}
		return newClient, response.StatusInternalServerError
	}

	return newClient, response.StatusCreated
}

func Update(tenantID, clientID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&client.Client{}).
		Where("tenant_id = ? AND id = ?", tenantID, clientID).
		Updates(updates)

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

func Search(tenantID uint, query string) ([]client.Client, response.Status) {
	var clients []client.Client
	db := databaseHelper.Db

	searchPattern := "%" + query + "%"
	result := db.Where("tenant_id = ? AND (razon_social LIKE ? OR nombre_fantasia LIKE ? OR cuit LIKE ?)",
		tenantID, searchPattern, searchPattern, searchPattern).
		Preload("Conditions").
		Find(&clients)

	if err := result.Error; err != nil {
		return clients, response.StatusInternalServerError
	}

	return clients, response.StatusOk
}

func Delete(tenantID, clientID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, clientID).Delete(&client.Client{})
	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}
	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
