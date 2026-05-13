package clientHandler

import (
	"sales-manager-back/internal/data/infrastructure/clientRepository"
	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/response"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return clientRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return clientRepository.GetByVendedor(tenantID, vendedorID)
}

func GetByID(tenantID, clientID uint) (interface{}, response.Status) {
	return clientRepository.GetByID(tenantID, clientID)
}

func Create(newClient *client.Client) (interface{}, response.Status) {
	return clientRepository.Create(newClient)
}

func BulkCreate(clients []client.Client) response.Status {
	return clientRepository.BulkCreate(clients)
}

func Update(tenantID, clientID uint, updates map[string]interface{}) response.Status {
	return clientRepository.Update(tenantID, clientID, updates)
}

func Search(tenantID uint, query string) (interface{}, response.Status) {
	return clientRepository.Search(tenantID, query)
}

func Delete(tenantID, clientID uint) response.Status {
	return clientRepository.Delete(tenantID, clientID)
}
