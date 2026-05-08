package agendaHandler

import (
	"sales-manager-back/internal/data/infrastructure/agendaRepository"
	"sales-manager-back/pkg/domain/agenda"
	"sales-manager-back/pkg/domain/response"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return agendaRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return agendaRepository.GetByVendedor(tenantID, vendedorID)
}

func GetByID(tenantID, eventID uint) (interface{}, response.Status) {
	return agendaRepository.GetByID(tenantID, eventID)
}

func Create(newEvent *agenda.AgendaEvent) (interface{}, response.Status) {
	return agendaRepository.Create(newEvent)
}

func Update(tenantID, eventID uint, updates map[string]interface{}) response.Status {
	return agendaRepository.Update(tenantID, eventID, updates)
}

func Complete(tenantID, eventID uint) response.Status {
	return agendaRepository.Complete(tenantID, eventID)
}

func Delete(tenantID, eventID uint) response.Status {
	return agendaRepository.Delete(tenantID, eventID)
}
