package authHandler

import (
	"log"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
	"time"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	CompanyName string `json:"companyName"`
}

func Register(uid string, req RegisterRequest) (interface{}, response.Status) {
	db := databaseHelper.Db

	// Start a transaction since we are creating multiple related records
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create the Tenant (Agency)
	newTenant := tenant.Tenant{
		Name:             req.CompanyName,
		Owner:            req.DisplayName,
		Email:            req.Email,
		Plan:             tenant.PlanStarter,
		Status:           tenant.EstadoTrial,
		RegistrationDate: time.Now().Format("2006-01-02"),
		LastActivity:     time.Now().Format("2006-01-02"),
	}

	if err := tx.Create(&newTenant).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating tenant: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 2. Create the User (Agency Owner)
	newUser := user.User{
		TenantID:           newTenant.ID,
		Name:               req.DisplayName,
		Email:              req.Email,
		TeamRole:           user.RoleAgency,
		Status:             user.StatusActivo,
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

	// 3. Update Firebase Custom Claims
	claims := map[string]interface{}{
		"tenantId": float64(newTenant.ID),
		"userId":   float64(newUser.ID),
		"role":     string(user.RoleAgency),
	}

	err := firebaseHelper.SetCustomUserClaims(uid, claims)
	if err != nil {
		// Log the error but don't fail the request, the user was created in DB.
		// In a robust system, we would have a retry mechanism.
		log.Printf("CRITICAL: Failed to set Firebase claims for UID %s: %v", uid, err)
	}

	return newTenant, response.StatusCreated
}
