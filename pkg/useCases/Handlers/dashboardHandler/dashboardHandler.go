package dashboardHandler

import (
	"sales-manager-back/internal/data/infrastructure/dashboardRepository"
	"sales-manager-back/pkg/domain/dashboard"
	"sales-manager-back/pkg/domain/response"
)

// ============ ALERTS ============

func GetAlerts(tenantID, vendedorID uint) (interface{}, response.Status) {
	return dashboardRepository.GetAlerts(tenantID, vendedorID)
}

func CreateAlert(newAlert *dashboard.Alert) (interface{}, response.Status) {
	return dashboardRepository.CreateAlert(newAlert)
}

func MarkAlertAsRead(tenantID, alertID uint) response.Status {
	return dashboardRepository.MarkAlertAsRead(tenantID, alertID)
}

// ============ QUICK NOTES ============

func GetQuickNotes(tenantID, vendedorID uint) (interface{}, response.Status) {
	return dashboardRepository.GetQuickNotes(tenantID, vendedorID)
}

func CreateQuickNote(newNote *dashboard.QuickNote) (interface{}, response.Status) {
	return dashboardRepository.CreateQuickNote(newNote)
}

func UpdateQuickNote(tenantID, noteID uint, updates map[string]interface{}) response.Status {
	return dashboardRepository.UpdateQuickNote(tenantID, noteID, updates)
}

func DeleteQuickNote(tenantID, noteID uint) response.Status {
	return dashboardRepository.DeleteQuickNote(tenantID, noteID)
}

// ============ GOALS ============

func GetGoals(tenantID, vendedorID uint) (interface{}, response.Status) {
	return dashboardRepository.GetGoals(tenantID, vendedorID)
}

func CreateGoal(newGoal *dashboard.Goal) (interface{}, response.Status) {
	return dashboardRepository.CreateGoal(newGoal)
}

func UpdateGoal(tenantID, goalID uint, updates map[string]interface{}) response.Status {
	return dashboardRepository.UpdateGoal(tenantID, goalID, updates)
}
