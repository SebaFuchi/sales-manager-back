package agendaRepository

import (
	"sales-manager-back/pkg/domain/agenda"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetAll(tenantID uint) ([]agenda.AgendaEvent, response.Status) {
	var events []agenda.AgendaEvent
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).
		Order("date DESC, time DESC").
		Find(&events)

	if err := result.Error; err != nil {
		return events, response.StatusInternalServerError
	}

	return events, response.StatusOk
}

func GetByVendedor(tenantID, vendedorID uint) ([]agenda.AgendaEvent, response.Status) {
	var events []agenda.AgendaEvent
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Order("date DESC, time DESC").
		Find(&events)

	if err := result.Error; err != nil {
		return events, response.StatusInternalServerError
	}

	return events, response.StatusOk
}

func GetByID(tenantID, eventID uint) (agenda.AgendaEvent, response.Status) {
	var event agenda.AgendaEvent
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, eventID).
		First(&event)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return event, response.StatusNotFound
		}
		return event, response.StatusInternalServerError
	}

	return event, response.StatusOk
}

func Create(newEvent *agenda.AgendaEvent) (*agenda.AgendaEvent, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newEvent)
	if err := result.Error; err != nil {
		return newEvent, response.StatusInternalServerError
	}

	return newEvent, response.StatusCreated
}

func Update(tenantID, eventID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&agenda.AgendaEvent{}).
		Where("tenant_id = ? AND id = ?", tenantID, eventID).
		Updates(updates)

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

func Complete(tenantID, eventID uint) response.Status {
	db := databaseHelper.Db

	result := db.Model(&agenda.AgendaEvent{}).
		Where("tenant_id = ? AND id = ?", tenantID, eventID).
		Update("status", agenda.StatusCompletado)

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

func Delete(tenantID, eventID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, eventID).
		Delete(&agenda.AgendaEvent{})

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusNoContent
}
