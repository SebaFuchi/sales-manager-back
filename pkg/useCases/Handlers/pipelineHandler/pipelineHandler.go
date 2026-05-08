package pipelineHandler

import (
	"sales-manager-back/internal/data/infrastructure/pipelineRepository"
	"sales-manager-back/pkg/domain/pipeline"
	"sales-manager-back/pkg/domain/response"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return pipelineRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return pipelineRepository.GetByVendedor(tenantID, vendedorID)
}

func Create(newDeal *pipeline.Deal) (interface{}, response.Status) {
	return pipelineRepository.Create(newDeal)
}

func Update(tenantID, dealID uint, updates map[string]interface{}) response.Status {
	return pipelineRepository.Update(tenantID, dealID, updates)
}

func Delete(tenantID, dealID uint) response.Status {
	return pipelineRepository.Delete(tenantID, dealID)
}
