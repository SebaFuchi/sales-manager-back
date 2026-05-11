package tenantRepository

import (
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

// GetAll returns all tenants in the system
func GetAll() ([]tenant.Tenant, response.Status) {
	var tenants []tenant.Tenant
	db := databaseHelper.Db

	result := db.Find(&tenants)

	if err := result.Error; err != nil {
		return tenants, response.StatusInternalServerError
	}

	return tenants, response.StatusOk
}

// GetByID returns a specific tenant by ID
func GetByID(id uint) (tenant.Tenant, response.Status) {
	var t tenant.Tenant
	db := databaseHelper.Db

	result := db.First(&t, id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return t, response.StatusNotFound
		}
		return t, response.StatusInternalServerError
	}

	return t, response.StatusOk
}

// Create adds a new tenant to the database
func Create(newTenant *tenant.Tenant) (*tenant.Tenant, response.Status) {
	db := databaseHelper.Db
	if err := db.Create(newTenant).Error; err != nil {
		return nil, response.StatusInternalServerError
	}
	return newTenant, response.StatusCreated
}
