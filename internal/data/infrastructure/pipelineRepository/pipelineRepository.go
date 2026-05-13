package pipelineRepository

import (
	"sales-manager-back/pkg/domain/pipeline"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
)

func GetAll(tenantID uint) ([]pipeline.Deal, response.Status) {
	var deals []pipeline.Deal
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Order("date DESC").
		Find(&deals)

	if err := result.Error; err != nil {
		return deals, response.StatusInternalServerError
	}

	return deals, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]pipeline.Deal, response.Status) {
	var deals []pipeline.Deal
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Order("date DESC").
		Find(&deals)

	if err := result.Error; err != nil {
		return deals, response.StatusInternalServerError
	}

	return deals, response.StatusOk
}

func Create(newDeal *pipeline.Deal) (*pipeline.Deal, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newDeal)
	if err := result.Error; err != nil {
		return newDeal, response.StatusInternalServerError
	}

	return newDeal, response.StatusCreated
}

func Update(tenantID, dealID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&pipeline.Deal{}).
		Where("tenant_id = ? AND id = ?", tenantID, dealID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

func Delete(tenantID, dealID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, dealID).
		Delete(&pipeline.Deal{})

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusNoContent
}
