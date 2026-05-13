package tenantHandler

import (
	"log"
	"sales-manager-back/internal/data/infrastructure/tenantRepository"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
)

func GetAll() (interface{}, response.Status) {
	return tenantRepository.GetAll()
}

func GetByID(id uint) (interface{}, response.Status) {
	return tenantRepository.GetByID(id)
}

func Create(newTenant *tenant.Tenant) (interface{}, response.Status) {
	db := databaseHelper.Db

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create the Tenant
	if err := tx.Create(newTenant).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating tenant: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 2. Create the User (Agency Owner)
	newUser := user.User{
		TenantID:           newTenant.ID,
		Name:               newTenant.Owner,
		Email:              newTenant.Email,
		TeamRole:           user.EquipoRolTitular,
		Role:               user.RoleAgency,
		Status:             user.EstadoActivo,
		SplitPercentageSub: 0,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating user: %v", err)
		return nil, response.StatusInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 3. Attempt to create in Firebase Auth
	// Check if user already exists
	firebaseUser, err := firebaseHelper.GetUserByEmail(newUser.Email)
	if err != nil || firebaseUser == nil {
		// User doesn't exist, create them with a temp password
		tempPassword := "CommerciaHub2026!"
		firebaseUser, err = firebaseHelper.CreateUser(newUser.Email, tempPassword, newUser.Name)
	}

	if firebaseUser != nil {
		claims := map[string]interface{}{
			"tenantId": float64(newTenant.ID),
			"userId":   float64(newUser.ID),
			"role":     string(user.RoleAgency),
		}
		errClaims := firebaseHelper.SetCustomUserClaims(firebaseUser.UID, claims)
		if errClaims != nil {
			log.Printf("Failed to set Firebase claims for UID %s: %v", firebaseUser.UID, errClaims)
		}

		// Save FirebaseUID back to the user record so GetByFirebaseUID works on login
		db := databaseHelper.Db
		db.Model(&user.User{}).Where("id = ?", newUser.ID).Update("firebase_uid", firebaseUser.UID)
	} else {
		log.Printf("Error creating/fetching Firebase user for email %s: %v", newUser.Email, err)
	}

	return newTenant, response.StatusCreated
}

func Update(tenantID uint, updates map[string]interface{}) response.Status {
	return tenantRepository.Update(tenantID, updates)
}

func Delete(tenantID uint) response.Status {
	return tenantRepository.Delete(tenantID)
}

// GetMyAgency returns the tenant data enriched with live counts from related tables
func GetMyAgency(tenantID uint) (interface{}, response.Status) {
	t, status := tenantRepository.GetByID(tenantID)
	if status != response.StatusOk {
		return nil, status
	}

	db := databaseHelper.Db

	// Count real data from related tables
	var userCount int64
	db.Model(&user.User{}).Where("tenant_id = ?", tenantID).Count(&userCount)

	var clientCount int64
	db.Table("client").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Count(&clientCount)

	var principalCount int64
	db.Table("principal").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Count(&principalCount)

	var saleCount int64
	db.Table("sale").Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Count(&saleCount)

	// Build enriched response
	result := map[string]interface{}{
		"id":                    t.ID,
		"name":                  t.Name,
		"owner":                 t.Owner,
		"email":                 t.Email,
		"plan":                  t.Plan,
		"status":                t.Status,
		"monthlyFee":            t.MonthlyFee,
		"users":                 userCount,
		"userLimit":             t.UserLimit,
		"clientsLoaded":         clientCount,
		"principalsLoaded":      principalCount,
		"operationsRecorded":    saleCount,
		"registrationDate":      t.RegistrationDate,
	}

	return result, response.StatusOk
}
