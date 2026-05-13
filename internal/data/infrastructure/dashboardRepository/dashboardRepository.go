package dashboardRepository

import (
	"sales-manager-back/pkg/domain/dashboard"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
)

// ============ ALERTS ============

func GetAlerts(tenantID, vendedorID uint) ([]dashboard.Alert, response.Status) {
	var alerts []dashboard.Alert
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Order("created_at DESC").
		Find(&alerts)

	if err := result.Error; err != nil {
		return alerts, response.StatusInternalServerError
	}

	return alerts, response.StatusOk
}

func CreateAlert(newAlert *dashboard.Alert) (*dashboard.Alert, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newAlert)
	if err := result.Error; err != nil {
		return newAlert, response.StatusInternalServerError
	}

	return newAlert, response.StatusCreated
}

func MarkAlertAsRead(tenantID, alertID uint) response.Status {
	db := databaseHelper.Db

	result := db.Model(&dashboard.Alert{}).
		Where("tenant_id = ? AND id = ?", tenantID, alertID).
		Update("seen", true)

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

// ============ QUICK NOTES ============

func GetQuickNotes(tenantID, vendedorID uint) ([]dashboard.QuickNote, response.Status) {
	var notes []dashboard.QuickNote
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Order("order ASC").
		Find(&notes)

	if err := result.Error; err != nil {
		return notes, response.StatusInternalServerError
	}

	return notes, response.StatusOk
}

func CreateQuickNote(newNote *dashboard.QuickNote) (*dashboard.QuickNote, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newNote)
	if err := result.Error; err != nil {
		return newNote, response.StatusInternalServerError
	}

	return newNote, response.StatusCreated
}

func UpdateQuickNote(tenantID, noteID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&dashboard.QuickNote{}).
		Where("tenant_id = ? AND id = ?", tenantID, noteID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}

func DeleteQuickNote(tenantID, noteID uint) response.Status {
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, noteID).
		Delete(&dashboard.QuickNote{})

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusNoContent
}

// ============ GOALS ============

func GetGoals(tenantID, vendedorID uint) ([]dashboard.Goal, response.Status) {
	var goals []dashboard.Goal
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND agent_id = ?", tenantID, vendedorID).
		Find(&goals)

	if err := result.Error; err != nil {
		return goals, response.StatusInternalServerError
	}

	return goals, response.StatusOk
}

func CreateGoal(newGoal *dashboard.Goal) (*dashboard.Goal, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newGoal)
	if err := result.Error; err != nil {
		return newGoal, response.StatusInternalServerError
	}

	return newGoal, response.StatusCreated
}

func UpdateGoal(tenantID, goalID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&dashboard.Goal{}).
		Where("tenant_id = ? AND id = ?", tenantID, goalID).
		Updates(databaseHelper.CamelToSnakeMap(updates))

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
